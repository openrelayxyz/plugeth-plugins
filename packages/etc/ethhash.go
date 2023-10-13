package main 

import (
	"os"
	"fmt"
	"time"
	"bytes"
	"errors"
	"math/big"
	"runtime"
	"path/filepath"

	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted/types"
)

// Config are the configuration parameters of the ethash.
type Config struct {
	CacheDir         string
	CachesInMem      int
	CachesOnDisk     int
	CachesLockMmap   bool
	DatasetDir       string
	DatasetsInMem    int
	DatasetsOnDisk   int
	DatasetsLockMmap bool
	PowMode          Mode

	// When set, notifications sent by the remote sealer will
	// be block header JSON objects instead of work package arrays.
	NotifyFull bool

	Log core.Logger `toml:"-"`
	// ECIP-1099
	ECIP1099Block *uint64 `toml:"-"`
}

// Mode defines the type and amount of PoW verification an ethash engine makes.
type Mode uint

const (
	ModeNormal Mode = iota
	ModeShared
	ModeTest
	ModeFake
	ModePoissonFake
	ModeFullFake
)

// Ethash proof-of-work protocol constants.
var (
	maxUncles              = 2                // Maximum number of uncles allowed in a single block
	allowedFutureBlockTime = 15 * time.Second // Max time from current time allowed for blocks, before they're considered future blocks
)

// Various error messages to mark blocks invalid. These should be private to
// prevent engine specific errors from being referenced in the remainder of the
// codebase, inherently breaking if the engine is swapped out. Please put common
// error types into the consensus package.
var (
	errOlderBlockTime    = errors.New("timestamp older than parent")
	errTooManyUncles     = errors.New("too many uncles")
	errDuplicateUncle    = errors.New("duplicate uncle")
	errUncleIsAncestor   = errors.New("uncle is ancestor")
	errDanglingUncle     = errors.New("uncle's parent is not ancestor")
	errInvalidDifficulty = errors.New("non-positive difficulty")
	errInvalidMixDigest  = errors.New("invalid mix digest")
	errInvalidPoW        = errors.New("invalid proof-of-work")
)

var unixNow int64 = time.Now().Unix()

var ErrInvalidDumpMagic = errors.New("invalid dump magic")

var (
	// two256 is a big integer representing 2^256
	two256 = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0))

	// sharedEthash is a full instance that can be shared between multiple users.
	sharedEthash *Ethash

	// algorithmRevision is the data structure version used for file naming.
	algorithmRevision = 23

	// dumpMagic is a dataset dump header to sanity check a data dump.
	dumpMagic = []uint32{0xbaddcafe, 0xfee1dead}
)

func init() {
	sharedConfig := Config{
		PowMode:       ModeNormal,
		CachesInMem:   3,
		DatasetsInMem: 1,
	}
	sharedEthash = New(sharedConfig, nil, false)
}

// New creates a full sized ethash PoW scheme and starts a background thread for
// remote mining, also optionally notifying a batch of remote services of new work
// packages.
func New(config Config, notify []string, noverify bool) *Ethash {
	if config.CachesInMem <= 0 {
		log.Warn("One ethash cache must always be in memory", "requested", config.CachesInMem)
		config.CachesInMem = 1
	}
	if config.CacheDir != "" && config.CachesOnDisk > 0 {
		log.Info("Disk storage enabled for ethash caches", "dir", config.CacheDir, "count", config.CachesOnDisk)
	}
	if config.DatasetDir != "" && config.DatasetsOnDisk > 0 {
		log.Info("Disk storage enabled for ethash DAGs", "dir", config.DatasetDir, "count", config.DatasetsOnDisk)
	}
	ethash := &Ethash{
		config:   config,
		// caches:   newlru(config.CachesInMem, newCache),
		// datasets: newlru(config.DatasetsInMem, newDataset),
		update:   make(chan struct{}),
		// hashrate: metrics.NewMeterForced(),
	}
	if config.PowMode == ModeShared {
		ethash.shared = sharedEthash
	}
	ethash.remote = startRemoteSealer(ethash, notify, noverify)
	return ethash
}

func (ethash *Ethash) Threads() int {
	ethash.lock.Lock()
	defer ethash.lock.Unlock()

	return ethash.threads
}

// verifyHeader checks whether a header conforms to the consensus rules of the
// stock Ethereum ethash engine.
// See YP section 4.3.4. "Block Header Validity"
func (ethash *Ethash) verifyHeader(chain ChainHeaderReader, header, parent *types.Header, uncle bool, seal bool, unixNow int64) error {
	// Ensure that the header's extra-data section is of a reasonable size
	if uint64(len(header.Extra)) > MaximumExtraDataSize {
		return fmt.Errorf("extra-data too long: %d > %d", len(header.Extra), MaximumExtraDataSize)
	}
	// Verify the header's timestamp
	if !uncle {
		if header.Time > uint64(unixNow+int64(allowedFutureBlockTime.Seconds())) {
			return ErrFutureBlock
		}
	}
	if header.Time <= parent.Time {
		return errOlderBlockTime
	}
	// Verify the block's difficulty based on its timestamp and parent's difficulty
	expected := ethash.CalcDifficulty(chain, header.Time, parent)

	if expected.Cmp(header.Difficulty) != 0 {
		return fmt.Errorf("invalid difficulty: have %v, want %v", header.Difficulty, expected)
	}
	// Verify that the gas limit is <= 2^63-1
	if header.GasLimit > MaxGasLimit {
		return fmt.Errorf("invalid gasLimit: have %v, max %v", header.GasLimit, MaxGasLimit)
	}
	// Verify that the gasUsed is <= gasLimit
	if header.GasUsed > header.GasLimit {
		return fmt.Errorf("invalid gasUsed: have %d, gasLimit %d", header.GasUsed, header.GasLimit)
	}
	// Verify the block's gas usage and (if applicable) verify the base fee.
	if !ethash.pluginConfig.IsEnabled(ethash.pluginConfig.GetEIP1559Transition, header.Number) {
		// Verify BaseFee not present before EIP-1559 fork.
		if header.BaseFee != nil {
			return fmt.Errorf("invalid baseFee before fork: have %d, expected 'nil'", header.BaseFee)
		}
		if err := VerifyGaslimit(parent.GasLimit, header.GasLimit); err != nil {
			return err
		}
	} else if err := VerifyEIP1559Header(*ethash.pluginConfig, parent, header); err != nil {
		// Verify the header's EIP-1559 attributes.
		return err
	}
	// Verify that the block number is parent's +1
	if diff := new(big.Int).Sub(header.Number, parent.Number); diff.Cmp(big.NewInt(1)) != 0 {
		return ErrInvalidNumber
	}
	if ethash.pluginConfig.IsEnabledByTime(ethash.pluginConfig.GetEIP3860TransitionTime, &header.Time) || ethash.pluginConfig.IsEnabled(ethash.pluginConfig.GetEIP3860Transition, header.Number) {
		return fmt.Errorf("ethash does not support shanghai fork")
	}
	if ethash.pluginConfig.IsEnabledByTime(ethash.pluginConfig.GetEIP4844TransitionTime, &header.Time) {
		return fmt.Errorf("ethash does not support cancun fork")
	}
	// Verify the engine specific seal securing the block
	if seal {
		if err := ethash.verifySeal(chain, header, false); err != nil {
			return err
		}
	}
	// If all checks passed, validate any special fields for hard forks
	if err := VerifyDAOHeaderExtraData(*ethash.pluginConfig, header); err != nil {
		return err
	}
	return nil
}

// func (ethash *Ethash) verifyHeaderWorker(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool, index int, unixNow int64) error {
// 	var parent *types.Header
// 	if index == 0 {
// 		parent = chain.GetHeader(headers[0].ParentHash, headers[0].Number.Uint64()-1)
// 	} else if headers[index-1].Hash() == headers[index].ParentHash {
// 		parent = headers[index-1]
// 	}
// 	if parent == nil {
// 		return consensus.ErrUnknownAncestor
// 	}
// 	return ethash.verifyHeader(chain, headers[index], parent, false, seals[index], unixNow)
// }

// verifySeal checks whether a block satisfies the PoW difficulty requirements,
// either using the usual ethash cache for it, or alternatively using a full DAG
// to make remote mining fast.
func (ethash *Ethash) verifySeal(chain ChainHeaderReader, header *types.Header, fulldag bool) error {
	// If we're running a fake PoW, accept any seal as valid
	if ethash.config.PowMode == ModeFake || ethash.config.PowMode == ModePoissonFake || ethash.config.PowMode == ModeFullFake {
		time.Sleep(ethash.fakeDelay)
		if ethash.fakeFail == header.Number.Uint64() {
			return errInvalidPoW
		}
		return nil
	}
	// If we're running a shared PoW, delegate verification to it
	if ethash.shared != nil {
		return ethash.shared.verifySeal(chain, header, fulldag)
	}
	// Ensure that we have a valid difficulty for the block
	if header.Difficulty.Sign() <= 0 {
		return errInvalidDifficulty
	}
	// Recompute the digest and PoW values
	number := header.Number.Uint64()

	var (
		digest []byte
		result []byte
	)
	// If fast-but-heavy PoW verification was requested, use an ethash dataset
	if fulldag {
		dataset := ethash.dataset(number, true)
		if dataset.generated() {
			digest, result = hashimotoFull(dataset.dataset, ethash.SealHash(header).Bytes(), header.Nonce.Uint64())

			// Datasets are unmapped in a finalizer. Ensure that the dataset stays alive
			// until after the call to hashimotoFull so it's not unmapped while being used.
			runtime.KeepAlive(dataset)
		} else {
			// Dataset not yet generated, don't hang, use a cache instead
			fulldag = false
		}
	}
	// If slow-but-light PoW verification was requested (or DAG not yet ready), use an ethash cache
	if !fulldag {
		cache := ethash.cache(number)
		epochLength := calcEpochLength(number, ethash.config.ECIP1099Block)
		epoch := calcEpoch(number, epochLength)
		size := datasetSize(epoch)
		if ethash.config.PowMode == ModeTest {
			size = 32 * 1024
		}
		digest, result = hashimotoLight(size, cache.cache, ethash.SealHash(header).Bytes(), header.Nonce.Uint64())

		// Caches are unmapped in a finalizer. Ensure that the cache stays alive
		// until after the call to hashimotoLight so it's not unmapped while being used.
		runtime.KeepAlive(cache)
	}
	// Verify the calculated values against the ones provided in the header
	if !bytes.Equal(header.MixDigest[:], digest) {
		return errInvalidMixDigest
	}
	target := new(big.Int).Div(two256, header.Difficulty)
	if new(big.Int).SetBytes(result).Cmp(target) > 0 {
		return errInvalidPoW
	}
	return nil
}

// dataset tries to retrieve a mining dataset for the specified block number
// by first checking against a list of in-memory datasets, then against DAGs
// stored on disk, and finally generating one if none can be found.
//
// If async is specified, not only the future but the current DAG is also
// generates on a background thread.
func (ethash *Ethash) dataset(block uint64, async bool) *dataset {
	// Retrieve the requested ethash dataset
	epochLength := calcEpochLength(block, ethash.config.ECIP1099Block)
	epoch := calcEpoch(block, epochLength)
	log.Error("230", "block", block, "epochLength", epochLength, "ECIP1099Block", ethash.config.ECIP1099Block)
	current, future := ethash.datasets.get(epoch, epochLength, ethash.config.ECIP1099Block)

	// set async false if ecip-1099 transition in case of regeneratiion bad DAG on disk
	if epochLength == epochLengthECIP1099 && (epoch == 42 || epoch == 195) {
		async = false
	}

	// If async is specified, generate everything in a background thread
	if async && !current.generated() {
		go func() {
			current.generate(ethash.config.DatasetDir, ethash.config.DatasetsOnDisk, ethash.config.DatasetsLockMmap, ethash.config.PowMode == ModeTest)
			if future != nil {
				future.generate(ethash.config.DatasetDir, ethash.config.DatasetsOnDisk, ethash.config.DatasetsLockMmap, ethash.config.PowMode == ModeTest)
			}
		}()
	} else {
		// Either blocking generation was requested, or already done
		current.generate(ethash.config.DatasetDir, ethash.config.DatasetsOnDisk, ethash.config.DatasetsLockMmap, ethash.config.PowMode == ModeTest)
		if future != nil {
			go future.generate(ethash.config.DatasetDir, ethash.config.DatasetsOnDisk, ethash.config.DatasetsLockMmap, ethash.config.PowMode == ModeTest)
		}
	}
	return current
}

// cache tries to retrieve a verification cache for the specified block number
// by first checking against a list of in-memory caches, then against caches
// stored on disk, and finally generating one if none can be found.
func (ethash *Ethash) cache(block uint64) *cache {
	// var num *uint64
	// bi := big.NewInt(11700000).Uint64()
	// num = &bi
	epochLength := calcEpochLength(block, ethash.config.ECIP1099Block)
	epoch := calcEpoch(block, epochLength)
	log.Error("261", "block", block, "epochLength", epochLength, "ECIP1099Block", ethash.config.ECIP1099Block)
	current, future := ethash.caches.get(epoch, epochLength, ethash.config.ECIP1099Block)

	// Wait for generation finish.
	current.generate(ethash.config.CacheDir, ethash.config.CachesOnDisk, ethash.config.CachesLockMmap, ethash.config.PowMode == ModeTest)

	// If we need a new future cache, now's a good time to regenerate it.
	if future != nil {
		go future.generate(ethash.config.CacheDir, ethash.config.CachesOnDisk, ethash.config.CachesLockMmap, ethash.config.PowMode == ModeTest)
	}
	return current
}

// generated returns whether this particular dataset finished generating already
// or not (it may not have been started at all). This is useful for remote miners
// to default to verification caches instead of blocking on DAG generations.
func (d *dataset) generated() bool {
	return d.done.Load()
}

// get retrieves or creates an item for the given epoch. The first return value is always
// non-nil. The second return value is non-nil if lru thinks that an item will be useful in
// the near future.
func (lru *lru[T]) get(epoch uint64, epochLength uint64, ecip1099FBlock *uint64) (item, future T) {
	log.Error("is it nil", "lru", lru)
	lru.mu.Lock()
	defer lru.mu.Unlock()

	// Use the sum of epoch and epochLength as the cache key.
	// This is not perfectly safe, but it's good enough (at least for the first 30000 epochs, or the first 427 years).
	cacheKey := epochLength + epoch

	// Get or create the item for the requested epoch.
	item, ok := lru.cache.Get(cacheKey)
	if !ok {
		if lru.future > 0 && lru.future == epoch {
			item = lru.futureItem
		} else {
			log.Trace("Requiring new ethash "+lru.what, "epoch", epoch)
			item = lru.new(epoch, epochLength)
		}
		lru.cache.Add(cacheKey, item)
	}

	// Ensure pre-generation handles ecip-1099 changeover correctly
	var nextEpoch = epoch + 1
	var nextEpochLength = epochLength
	if ecip1099FBlock != nil {
		nextEpochBlock := nextEpoch * epochLength
		// Note that == demands that the ECIP1099 activation block is situated
		// at the beginning of an epoch.
		// https://github.com/ethereumclassic/ECIPs/blob/master/_specs/ecip-1099.md#implementation
		if nextEpochBlock == *ecip1099FBlock && epochLength == epochLengthDefault {
			nextEpoch = nextEpoch / 2
			nextEpochLength = epochLengthECIP1099
		}
	}

	// Update the 'future item' if epoch is larger than previously seen.
	// Last conditional clause ('lru.future > nextEpoch') handles the ECIP1099 case where
	// the next epoch is expected to be LESSER THAN that of the previous state's future epoch number.
	if epoch < maxEpoch-1 && lru.future != nextEpoch {
		log.Trace("Requiring new future ethash "+lru.what, "epoch", nextEpoch)
		future = lru.new(nextEpoch, nextEpochLength)
		lru.future = nextEpoch
		lru.futureItem = future
	}
	return item, future
}

// generate ensures that the cache content is generated before use.
func (c *cache) generate(dir string, limit int, lock bool, test bool) {
	c.once.Do(func() {
		size := cacheSize(c.epoch)
		seed := seedHash(c.epoch, c.epochLength)
		if test {
			size = 1024
		}
		// If we don't store anything on disk, generate and return.
		if dir == "" {
			c.cache = make([]uint32, size/4)
			generateCache(c.cache, c.epoch, c.epochLength, seed)
			return
		}
		// Disk storage is needed, this will get fancy
		var endian string
		if !isLittleEndian() {
			endian = ".be"
		}
		// The file path naming scheme was changed to include epoch values in the filename,
		// which enables a filepath glob with scan to identify out-of-bounds caches and remove them.
		// The legacy path declaration is provided below as a comment for reference.
		//
		// path := filepath.Join(dir, fmt.Sprintf("cache-R%d-%x%s", algorithmRevision, seed[:8], endian))                 // LEGACY
		path := filepath.Join(dir, fmt.Sprintf("cache-R%d-%d-%x%s", algorithmRevision, c.epoch, seed[:8], endian)) // CURRENT
		// logger := log.New("epoch", c.epoch, "epochLength", c.epochLength)

		// We're about to mmap the file, ensure that the mapping is cleaned up when the
		// cache becomes unused.
		runtime.SetFinalizer(c, (*cache).finalizer)

		// Try to load the file from disk and memory map it
		var err error
		c.dump, c.mmap, c.cache, err = memoryMap(path, lock)
		if err == nil {
			log.Debug("Loaded old ethash cache from disk")
			return
		}
		log.Debug("Failed to load old ethash cache", "err", err)

		// No usable previous cache available, create a new cache file to fill
		c.dump, c.mmap, c.cache, err = memoryMapAndGenerate(path, size, lock, func(buffer []uint32) { generateCache(buffer, c.epoch, c.epochLength, seed) })
		if err != nil {
			log.Error("Failed to generate mapped ethash cache", "err", err)

			c.cache = make([]uint32, size/4)
			generateCache(c.cache, c.epoch, c.epochLength, seed)
		}

		// Iterate over all cache file instances, deleting any out of bounds (where epoch is below lower limit, or above upper limit).
		matches, _ := filepath.Glob(filepath.Join(dir, fmt.Sprintf("cache-R%d*", algorithmRevision)))
		for _, file := range matches {
			var ar int   // algorithm revision
			var e uint64 // epoch
			var s string // seed
			if _, err := fmt.Sscanf(filepath.Base(file), "cache-R%d-%d-%s"+endian, &ar, &e, &s); err != nil {
				// There is an unrecognized file in this directory.
				// See if the name matches the expected pattern of the legacy naming scheme.
				if _, err := fmt.Sscanf(filepath.Base(file), "cache-R%d-%s"+endian, &ar, &s); err == nil {
					// This file matches the previous generation naming pattern (sans epoch).
					if err := os.Remove(file); err != nil {
						log.Error("Failed to remove legacy ethash cache file", "file", file, "err", err)
					} else {
						log.Warn("Deleted legacy ethash cache file", "path", file)
					}
				}
				// Else the file is unrecognized (unknown name format), leave it alone.
				continue
			}
			if e <= c.epoch-uint64(limit) || e > c.epoch+1 {
				if err := os.Remove(file); err == nil {
					log.Debug("Deleted ethash cache file", "target.epoch", e, "file", file)
				} else {
					log.Error("Failed to delete ethash cache file", "target.epoch", e, "file", file, "err", err)
				}
			}
		}
	})
}

// generate ensures that the dataset content is generated before use.
func (d *dataset) generate(dir string, limit int, lock bool, test bool) {
	d.once.Do(func() {
		// Mark the dataset generated after we're done. This is needed for remote
		defer d.done.Store(true)

		csize := cacheSize(d.epoch)
		dsize := datasetSize(d.epoch)
		seed := seedHash(d.epoch, d.epochLength)
		if test {
			csize = 1024
			dsize = 32 * 1024
		}
		// If we don't store anything on disk, generate and return
		if dir == "" {
			cache := make([]uint32, csize/4)
			generateCache(cache, d.epoch, d.epochLength, seed)

			d.dataset = make([]uint32, dsize/4)
			generateDataset(d.dataset, d.epoch, d.epochLength, cache)

			return
		}
		// Disk storage is needed, this will get fancy
		var endian string
		if !isLittleEndian() {
			endian = ".be"
		}
		path := filepath.Join(dir, fmt.Sprintf("full-R%d-%d-%x%s", algorithmRevision, d.epoch, seed[:8], endian))
		// logger := log.New("epoch", d.epoch)

		// We're about to mmap the file, ensure that the mapping is cleaned up when the
		// cache becomes unused.
		runtime.SetFinalizer(d, (*dataset).finalizer)

		// Try to load the file from disk and memory map it
		var err error
		d.dump, d.mmap, d.dataset, err = memoryMap(path, lock)
		if err == nil {
			log.Debug("Loaded old ethash dataset from disk", "path", path)
			return
		}
		log.Debug("Failed to load old ethash dataset", "err", err)

		// No usable previous dataset available, create a new dataset file to fill
		cache := make([]uint32, csize/4)
		generateCache(cache, d.epoch, d.epochLength, seed)

		d.dump, d.mmap, d.dataset, err = memoryMapAndGenerate(path, dsize, lock, func(buffer []uint32) { generateDataset(buffer, d.epoch, d.epochLength, cache) })
		if err != nil {
			log.Error("Failed to generate mapped ethash dataset", "err", err)

			d.dataset = make([]uint32, dsize/4)
			generateDataset(d.dataset, d.epoch, d.epochLength, cache)
		}

		// Iterate over all full file instances, deleting any out of bounds (where epoch is below lower limit, or above upper limit).
		matches, _ := filepath.Glob(filepath.Join(dir, fmt.Sprintf("full-R%d*", algorithmRevision)))
		for _, file := range matches {
			var ar int   // algorithm revision
			var e uint64 // epoch
			var s string // seed
			if _, err := fmt.Sscanf(filepath.Base(file), "full-R%d-%d-%s"+endian, &ar, &e, &s); err != nil {
				// There is an unrecognized file in this directory.
				// See if the name matches the expected pattern of the legacy naming scheme.
				if _, err := fmt.Sscanf(filepath.Base(file), "full-R%d-%s"+endian, &ar, &s); err == nil {
					// This file matches the previous generation naming pattern (sans epoch).
					if err := os.Remove(file); err != nil {
						log.Error("Failed to remove legacy ethash full file", "file", file, "err", err)
					} else {
						log.Warn("Deleted legacy ethash full file", "path", file)
					}
				}
				// Else the file is unrecognized (unknown name format), leave it alone.
				continue
			}
			if e <= d.epoch-uint64(limit) || e > d.epoch+1 {
				if err := os.Remove(file); err == nil {
					log.Debug("Deleted ethash full file", "target.epoch", e, "file", file)
				} else {
					log.Error("Failed to delete ethash full file", "target.epoch", e, "file", file, "err", err)
				}
			}
		}
	})
}



// func (ethash *Ethash) mine(block *types.Block, id int, seed uint64, abort chan struct{}, found chan *types.Block) {
// 	// Extract some data from the header
// 	var (
// 		header  = block.Header()
// 		hash    = ethash.SealHash(header).Bytes()
// 		target  = new(big.Int).Div(two256, header.Difficulty)
// 		number  = header.Number.Uint64()
// 		dataset = ethash.dataset(number, false)
// 	)
// 	// Start generating random nonces until we abort or find a good one
// 	var (
// 		attempts  = int64(0)
// 		nonce     = seed
// 		powBuffer = new(big.Int)
// 	)
// 	// logger := ethash.config.Log.New("miner", id)
// 	// logger.Trace("Started ethash search for new nonces", "seed", seed)
// 	// TODO PM talk to AR regarding these log.New methods
// search:
// 	for {
// 		select {
// 		case <-abort:
// 			// Mining terminated, update stats and abort
// 			log.Trace("Ethash nonce search aborted", "attempts", nonce-seed)
// 			// ethash.hashrate.Mark(attempts)
// 			break search

// 		default:
// 			// We don't have to update hash rate on every nonce, so update after after 2^X nonces
// 			attempts++
// 			if (attempts % (1 << 15)) == 0 {
// 				// ethash.hashrate.Mark(attempts)
// 				attempts = 0
// 			}
// 			// Compute the PoW value of this nonce
// 			digest, result := hashimotoFull(dataset.dataset, hash, nonce)
// 			if powBuffer.SetBytes(result).Cmp(target) <= 0 {
// 				// Correct nonce found, create a new header with it
// 				header = types.CopyHeader(header)
// 				header.Nonce = types.EncodeNonce(nonce)
// 				header.MixDigest = core.BytesToHash(digest)

// 				// Seal and return a block (if still needed)
// 				select {
// 				case found <- block.WithSeal(header):
// 					log.Trace("Ethash nonce found and reported", "attempts", nonce-seed, "nonce", nonce)
// 				case <-abort:
// 					log.Trace("Ethash nonce found but discarded", "attempts", nonce-seed, "nonce", nonce)
// 				}
// 				break search
// 			}
// 			nonce++
// 		}
// 	}
// 	// Datasets are unmapped in a finalizer. Ensure that the dataset stays live
// 	// during sealing so it's not unmapped while being read.
// 	runtime.KeepAlive(dataset)
// }

// // Threads returns the number of mining threads currently enabled. This doesn't
// // necessarily mean that mining is running!
// func (ethash *Ethash) Threads() int {
// 	ethash.lock.Lock()
// 	defer ethash.lock.Unlock()

// 	return ethash.threads
// }

func (ethash *Ethash) SetThreads(threads int) {
	ethash.lock.Lock()
	defer ethash.lock.Unlock()

	// If we're running a shared PoW, set the thread count on that instead
	if ethash.shared != nil {
		ethash.shared.SetThreads(threads)
		return
	}
	// Update the threads and ping any running seal to pull in any changes
	ethash.threads = threads
	select {
	case ethash.update <- struct{}{}:
	default:
	}
}

// SeedHash is the seed to use for generating a verification cache and the mining
// dataset.
func SeedHash(epoch uint64, epochLength uint64) []byte {
	return seedHash(epoch, epochLength)
}

// CalcEpochLength returns the epoch length for a given block number (ECIP-1099)
func CalcEpochLength(block uint64, ecip1099FBlock *uint64) uint64 {
	return calcEpochLength(block, ecip1099FBlock)
}

// CalcEpoch returns the epoch for a given block number (ECIP-1099)
func CalcEpoch(block uint64, epochLength uint64) uint64 {
	return calcEpoch(block, epochLength)
}
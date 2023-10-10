package main 

import (
	"os"
	"fmt"
	"hash"
	"errors"
	"time"
	"reflect"
	"unsafe"
	"sync/atomic"
	"math/big"
	"math/rand"
	"bytes"
	"encoding/hex"
	"encoding/binary"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"gonum.org/v1/gonum/stat/distuv"

	"golang.org/x/crypto/sha3"
	"github.com/edsrzf/mmap-go"
	exprand "golang.org/x/exp/rand"

	"github.com/openrelayxyz/plugeth-utils/restricted/types"
	"github.com/openrelayxyz/plugeth-utils/restricted/crypto"
)

// BigMax returns the larger of x or y.
func BigMax(x, y *big.Int) *big.Int {
	if x.Cmp(y) < 0 {
		return y
	}
	return x
}

// parent_diff_over_dbd is a  convenience fn for CalcDifficulty
func parent_diff_over_dbd(p *types.Header) *big.Int {
	return new(big.Int).Div(p.Difficulty, DifficultyBoundDivisor)
}

// parent_time_delta is a convenience fn for CalcDifficulty
func parent_time_delta(t uint64, p *types.Header) *big.Int {
	return new(big.Int).Sub(new(big.Int).SetUint64(t), new(big.Int).SetUint64(p.Time))
}

// VerifyDAOHeaderExtraData validates the extra-data field of a block header to
// ensure it conforms to DAO hard-fork rules.
//
// DAO hard-fork extension to the header validity:
//
//   - if the node is no-fork, do not accept blocks in the [fork, fork+10) range
//     with the fork specific extra-data set.
//   - if the node is pro-fork, require blocks in the specific range to have the
//     unique extra-data set.
func VerifyDAOHeaderExtraData(config PluginConfigurator, header *types.Header) error {
	// If the config wants the DAO fork, it should validate the extra data.
	// Otherwise, like any other block or any other config, it should not care.
	daoForkBlock := config.GetEthashEIP779Transition()
	if daoForkBlock == nil {
		return nil
	}
	daoForkBlockB := new(big.Int).SetUint64(*daoForkBlock)
	// Make sure the block is within the fork's modified extra-data range
	limit := new(big.Int).Add(daoForkBlockB, DAOForkExtraRange)
	if header.Number.Cmp(daoForkBlockB) < 0 || header.Number.Cmp(limit) >= 0 {
		return nil
	}
	if !bytes.Equal(header.Extra, DAOForkBlockExtra) {
		return ErrBadProDAOExtra
	}
	return nil

	// Leaving the "old" code in as dead commented code for reference.
	//
	// // Short circuit validation if the node doesn't care about the DAO fork
	// daoForkBlock := config.GetEthashEIP779Transition()
	// // Second clause catches test configs with nil fork blocks (maybe set dynamically or
	// // testing agnostic of chain config).
	// if daoForkBlock == nil && !generic.AsGenericCC(config).DAOSupport() {
	//	return nil
	// }
	//
	// if daoForkBlock == nil {
	//
	// }
	//
	// daoForkBlockB := new(big.Int).SetUint64(*daoForkBlock)
	//
	// // Make sure the block is within the fork's modified extra-data range
	// limit := new(big.Int).Add(daoForkBlockB, DAOForkExtraRange)
	// if header.Number.Cmp(daoForkBlockB) < 0 || header.Number.Cmp(limit) >= 0 {
	//	return nil
	// }
	// // Depending on whether we support or oppose the fork, validate the extra-data contents
	// if generic.AsGenericCC(config).DAOSupport() {
	//	if !bytes.Equal(header.Extra, DAOForkBlockExtra) {
	//		return ErrBadProDAOExtra
	//	}
	// } else {
	//	if bytes.Equal(header.Extra, DAOForkBlockExtra) {
	//		return ErrBadNoDAOExtra
	//	}
	// }
	// // All ok, header has the same extra-data we expect
	// return nil
}

// FromHex returns the bytes represented by the hexadecimal string s.
// s may be prefixed with "0x".
func FromHex(s string) []byte {
	if has0xPrefix(s) {
		s = s[2:]
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return Hex2Bytes(s)
}

// Hex2Bytes returns the bytes represented by the hexadecimal string str.
func Hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}

// has0xPrefix validates str begins with '0x' or '0X'.
func has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns
// the difficulty that a new block should have when created at time
// given the parent block's time and difficulty.
func CalcDifficulty(config *PluginConfigurator, time uint64, parent *types.Header) *big.Int {
	next := new(big.Int).Add(parent.Number, big1)
	out := new(big.Int)

	// TODO (meowbits): do we need this?
	// if config.IsEnabled(config.GetEthashTerminalTotalDifficulty, next) {
	// 	return big.NewInt(1)
	// }

	// ADJUSTMENT algorithms
	if config.IsEnabled(config.GetEthashEIP100BTransition, next) {
		// https://github.com/ethereum/EIPs/issues/100
		// algorithm:
		// diff = (parent_diff +
		//         (parent_diff / 2048 * max((2 if len(parent.uncles) else 1) - ((timestamp - parent.timestamp) // 9), -99))
		//        ) + 2^(periodCount - 2)
		out.Div(parent_time_delta(time, parent), EIP100FDifficultyIncrementDivisor)

		if parent.UncleHash == types.EmptyUncleHash {
			out.Sub(big1, out)
		} else {
			out.Sub(big2, out)
		}
		out.Set(BigMax(out, bigMinus99))
		out.Mul(parent_diff_over_dbd(parent), out)
		out.Add(out, parent.Difficulty)
	} else if config.IsEnabled(config.GetEIP2Transition, next) {
		// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-2.md
		// algorithm:
		// diff = (parent_diff +
		//         (parent_diff / 2048 * max(1 - (block_timestamp - parent_timestamp) // 10, -99))
		//        )
		out.Div(parent_time_delta(time, parent), EIP2DifficultyIncrementDivisor)
		out.Sub(big1, out)
		out.Set(BigMax(out, bigMinus99))
		out.Mul(parent_diff_over_dbd(parent), out)
		out.Add(out, parent.Difficulty)
	} else {
		// FRONTIER
		// algorithm:
		// diff =
		//   if parent_block_time_delta < params.DurationLimit
		//      parent_diff + (parent_diff // 2048)
		//   else
		//      parent_diff - (parent_diff // 2048)
		out.Set(parent.Difficulty)
		if parent_time_delta(time, parent).Cmp(DurationLimit) < 0 {
			out.Add(out, parent_diff_over_dbd(parent))
		} else {
			out.Sub(out, parent_diff_over_dbd(parent))
		}
	}

	// after adjustment and before bomb
	out.Set(BigMax(out, MinimumDifficulty))

	if config.IsEnabled(config.GetEthashECIP1041Transition, next) {
		return out
	}

	// EXPLOSION delays

	// exPeriodRef the explosion clause's reference point
	exPeriodRef := new(big.Int).Add(parent.Number, big1)

	if config.IsEnabled(config.GetEthashECIP1010PauseTransition, next) {
		ecip1010Explosion(*config, next, exPeriodRef)
	} else if len(config.GetEthashDifficultyBombDelaySchedule()) > 0 {
		// This logic varies from the original fork-based logic (below) in that
		// configured delay values are treated as compounding values (-2000000 + -3000000 = -5000000@constantinople)
		// as opposed to hardcoded pre-compounded values (-5000000@constantinople).
		// Thus the Sub-ing.
		fakeBlockNumber := new(big.Int).Set(exPeriodRef)
		for activated, dur := range config.GetEthashDifficultyBombDelaySchedule() {
			if exPeriodRef.Cmp(big.NewInt(int64(activated))) < 0 {
				continue
			}
			fakeBlockNumber.Sub(fakeBlockNumber, dur)
		}
		exPeriodRef.Set(fakeBlockNumber)
	} else if config.IsEnabled(config.GetEthashEIP5133Transition, next) {
		// calcDifficultyEip4345 is the difficulty adjustment algorithm as specified by EIP 4345.
		// It offsets the bomb a total of 10.7M blocks.
		fakeBlockNumber := new(big.Int)
		delayWithOffset := new(big.Int).Sub(EIP5133DifficultyBombDelay, big1)
		if parent.Number.Cmp(delayWithOffset) >= 0 {
			fakeBlockNumber = fakeBlockNumber.Sub(parent.Number, delayWithOffset)
		}
		exPeriodRef.Set(fakeBlockNumber)
	} else if config.IsEnabled(config.GetEthashEIP4345Transition, next) {
		// calcDifficultyEip4345 is the difficulty adjustment algorithm as specified by EIP 4345.
		// It offsets the bomb a total of 10.7M blocks.
		fakeBlockNumber := new(big.Int)
		delayWithOffset := new(big.Int).Sub(EIP4345DifficultyBombDelay, big1)
		if parent.Number.Cmp(delayWithOffset) >= 0 {
			fakeBlockNumber = fakeBlockNumber.Sub(parent.Number, delayWithOffset)
		}
		exPeriodRef.Set(fakeBlockNumber)
	} else if config.IsEnabled(config.GetEthashEIP3554Transition, next) {
		// calcDifficultyEIP3554 is the difficulty adjustment algorithm for London (December 2021).
		// The calculation uses the Byzantium rules, but with bomb offset 9.7M.
		fakeBlockNumber := new(big.Int)
		delayWithOffset := new(big.Int).Sub(EIP3554DifficultyBombDelay, big1)
		if parent.Number.Cmp(delayWithOffset) >= 0 {
			fakeBlockNumber = fakeBlockNumber.Sub(parent.Number, delayWithOffset)
		}
		exPeriodRef.Set(fakeBlockNumber)
	} else if config.IsEnabled(config.GetEthashEIP2384Transition, next) {
		// calcDifficultyEIP2384 is the difficulty adjustment algorithm for Muir Glacier.
		// The calculation uses the Byzantium rules, but with bomb offset 9M.
		fakeBlockNumber := new(big.Int)
		delayWithOffset := new(big.Int).Sub(EIP2384DifficultyBombDelay, big1)
		if parent.Number.Cmp(delayWithOffset) >= 0 {
			fakeBlockNumber = fakeBlockNumber.Sub(parent.Number, delayWithOffset)
		}
		exPeriodRef.Set(fakeBlockNumber)
	} else if config.IsEnabled(config.GetEthashEIP1234Transition, next) {
		// calcDifficultyEIP1234 is the difficulty adjustment algorithm for Constantinople.
		// The calculation uses the Byzantium rules, but with bomb offset 5M.
		// Specification EIP-1234: https://eips.ethereum.org/EIPS/eip-1234
		// Note, the calculations below looks at the parent number, which is 1 below
		// the block number. Thus we remove one from the delay given

		// calculate a fake block number for the ice-age delay
		// Specification: https://eips.ethereum.org/EIPS/eip-1234
		fakeBlockNumber := new(big.Int)
		delayWithOffset := new(big.Int).Sub(EIP1234DifficultyBombDelay, big1)
		if parent.Number.Cmp(delayWithOffset) >= 0 {
			fakeBlockNumber = fakeBlockNumber.Sub(parent.Number, delayWithOffset)
		}
		exPeriodRef.Set(fakeBlockNumber)
	} else if config.IsEnabled(config.GetEthashEIP649Transition, next) {
		// The calculation uses the Byzantium rules, with bomb offset of 3M.
		// Specification EIP-649: https://eips.ethereum.org/EIPS/eip-649
		// Related meta-ish EIP-669: https://github.com/ethereum/EIPs/pull/669
		// Note, the calculations below looks at the parent number, which is 1 below
		// the block number. Thus we remove one from the delay given

		fakeBlockNumber := new(big.Int)
		delayWithOffset := new(big.Int).Sub(EIP649DifficultyBombDelay, big1)
		if parent.Number.Cmp(delayWithOffset) >= 0 {
			fakeBlockNumber = fakeBlockNumber.Sub(parent.Number, delayWithOffset)
		}
		exPeriodRef.Set(fakeBlockNumber)
	}

	// EXPLOSION

	// the 'periodRef' (from above) represents the many ways of hackishly modifying the reference number
	// (ie the 'currentBlock') in order to lie to the function about what time it really is
	//
	//   2^(( periodRef // EDP) - 2)
	//
	x := new(big.Int)
	x.Div(exPeriodRef, ExpDiffPeriod) // (periodRef // EDP)
	if x.Cmp(big1) > 0 {                     // if result large enough (not in algo explicitly)
		x.Sub(x, big2)      // - 2
		x.Exp(big2, x, nil) // 2^
	} else {
		x.SetUint64(0)
	}
	out.Add(out, x)
	return out
}

// VerifyGaslimit verifies the header gas limit according increase/decrease
// in relation to the parent gas limit.
func VerifyGaslimit(parentGasLimit, headerGasLimit uint64) error {
	// Verify that the gas limit remains within allowed bounds
	diff := int64(parentGasLimit) - int64(headerGasLimit)
	if diff < 0 {
		diff *= -1
	}
	limit := parentGasLimit / GasLimitBoundDivisor
	if uint64(diff) >= limit {
		return fmt.Errorf("invalid gas limit: have %d, want %d +-= %d", headerGasLimit, parentGasLimit, limit-1)
	}
	if headerGasLimit < MinGasLimit {
		return errors.New("invalid gas limit below 5000")
	}
	return nil
}

func ecip1010Explosion(config PluginConfigurator, next *big.Int, exPeriodRef *big.Int) {
	// https://github.com/ethereumproject/ECIPs/blob/master/ECIPs/ECIP-1010.md

	if next.Uint64() < *config.GetEthashECIP1010ContinueTransition() {
		exPeriodRef.SetUint64(*config.GetEthashECIP1010PauseTransition())
	} else {
		length := new(big.Int).SetUint64(*config.GetEthashECIP1010ContinueTransition() - *config.GetEthashECIP1010PauseTransition())
		exPeriodRef.Sub(exPeriodRef, length)
	}
}

// hashimoto aggregates data from the full dataset in order to produce our final
// value for a particular header hash and nonce.
func hashimoto(hash []byte, nonce uint64, size uint64, lookup func(index uint32) []uint32) ([]byte, []byte) {
	// Calculate the number of theoretical rows (we use one buffer nonetheless)
	rows := uint32(size / mixBytes)

	// Combine header+nonce into a 40 byte seed
	seed := make([]byte, 40)
	copy(seed, hash)
	binary.LittleEndian.PutUint64(seed[32:], nonce)

	seed = crypto.Keccak512(seed)
	seedHead := binary.LittleEndian.Uint32(seed)

	// Start the mix with replicated seed
	mix := make([]uint32, mixBytes/4)
	for i := 0; i < len(mix); i++ {
		mix[i] = binary.LittleEndian.Uint32(seed[i%16*4:])
	}
	// Mix in random dataset nodes
	temp := make([]uint32, len(mix))

	for i := 0; i < loopAccesses; i++ {
		parent := fnv(uint32(i)^seedHead, mix[i%len(mix)]) % rows
		for j := uint32(0); j < mixBytes/hashBytes; j++ {
			copy(temp[j*hashWords:], lookup(2*parent+j))
		}
		fnvHash(mix, temp)
	}
	// Compress mix
	for i := 0; i < len(mix); i += 4 {
		mix[i/4] = fnv(fnv(fnv(mix[i], mix[i+1]), mix[i+2]), mix[i+3])
	}
	mix = mix[:len(mix)/4]

	digest := make([]byte, HashLength)
	for i, val := range mix {
		binary.LittleEndian.PutUint32(digest[i*4:], val)
	}
	return digest, crypto.Keccak256(append(seed, digest...))
}

// hashimotoLight aggregates data from the full dataset (using only a small
// in-memory cache) in order to produce our final value for a particular header
// hash and nonce.
func hashimotoLight(size uint64, cache []uint32, hash []byte, nonce uint64) ([]byte, []byte) {
	keccak512 := makeHasher(sha3.NewLegacyKeccak512())

	lookup := func(index uint32) []uint32 {
		rawData := generateDatasetItem(cache, index, keccak512)

		data := make([]uint32, len(rawData)/4)
		for i := 0; i < len(data); i++ {
			data[i] = binary.LittleEndian.Uint32(rawData[i*4:])
		}
		return data
	}
	return hashimoto(hash, nonce, size, lookup)
}

// hashimotoFull aggregates data from the full dataset (using the full in-memory
// dataset) in order to produce our final value for a particular header hash and
// nonce.
func hashimotoFull(dataset []uint32, hash []byte, nonce uint64) ([]byte, []byte) {
	lookup := func(index uint32) []uint32 {
		offset := index * hashWords
		return dataset[offset : offset+hashWords]
	}
	return hashimoto(hash, nonce, uint64(len(dataset))*4, lookup)
}

// calcEpochLength returns the epoch length for a given block number (ECIP-1099)
func calcEpochLength(block uint64, ecip1099FBlock *uint64) uint64 {
	if ecip1099FBlock != nil {
		if block >= *ecip1099FBlock {
			return epochLengthECIP1099
		}
	}
	return epochLengthDefault
}

// calcEpoch returns the epoch for a given block number (ECIP-1099)
func calcEpoch(block uint64, epochLength uint64) uint64 {
	epoch := block / epochLength
	return epoch
}

// datasetSize returns the size of the ethash mining dataset that belongs to a certain
// block number.
func datasetSize(epoch uint64) uint64 {
	if epoch < maxEpoch {
		return datasetSizes[int(epoch)]
	}
	return calcDatasetSize(epoch)
}

// fnv is an algorithm inspired by the FNV hash, which in some cases is used as
// a non-associative substitute for XOR. Note that we multiply the prime with
// the full 32-bit input, in contrast with the FNV-1 spec which multiplies the
// prime with one byte (octet) in turn.
func fnv(a, b uint32) uint32 {
	return a*0x01000193 ^ b
}

// fnvHash mixes in data into mix using the ethash fnv method.
func fnvHash(mix []uint32, data []uint32) {
	for i := 0; i < len(mix); i++ {
		mix[i] = mix[i]*0x01000193 ^ data[i]
	}
}

// hasher is a repetitive hasher allowing the same hash data structures to be
// reused between hash runs instead of requiring new ones to be created.
type hasher func(dest []byte, data []byte)

// makeHasher creates a repetitive hasher, allowing the same hash data structures to
// be reused between hash runs instead of requiring new ones to be created. The returned
// function is not thread safe!
func makeHasher(h hash.Hash) hasher {
	// sha3.state supports Read to get the sum, use it to avoid the overhead of Sum.
	// Read alters the state but we reset the hash before every operation.
	type readerHash interface {
		hash.Hash
		Read([]byte) (int, error)
	}
	rh, ok := h.(readerHash)
	if !ok {
		panic("can't find Read method on hash")
	}
	outputLen := rh.Size()
	return func(dest []byte, data []byte) {
		rh.Reset()
		rh.Write(data)
		rh.Read(dest[:outputLen])
	}
}

// generateDatasetItem combines data from 256 pseudorandomly selected cache nodes,
// and hashes that to compute a single dataset node.
func generateDatasetItem(cache []uint32, index uint32, keccak512 hasher) []byte {
	// Calculate the number of theoretical rows (we use one buffer nonetheless)
	rows := uint32(len(cache) / hashWords)

	// Initialize the mix
	mix := make([]byte, hashBytes)

	binary.LittleEndian.PutUint32(mix, cache[(index%rows)*hashWords]^index)
	for i := 1; i < hashWords; i++ {
		binary.LittleEndian.PutUint32(mix[i*4:], cache[(index%rows)*hashWords+uint32(i)])
	}
	keccak512(mix, mix)

	// Convert the mix to uint32s to avoid constant bit shifting
	intMix := make([]uint32, hashWords)
	for i := 0; i < len(intMix); i++ {
		intMix[i] = binary.LittleEndian.Uint32(mix[i*4:])
	}
	// fnv it with a lot of random cache nodes based on index
	for i := uint32(0); i < datasetParents; i++ {
		parent := fnv(index^i, intMix[i%16]) % rows
		fnvHash(intMix, cache[parent*hashWords:])
	}
	// Flatten the uint32 mix into a binary one and return
	for i, val := range intMix {
		binary.LittleEndian.PutUint32(mix[i*4:], val)
	}
	keccak512(mix, mix)
	return mix
}

// calcDatasetSize calculates the dataset size for epoch. The dataset size grows linearly,
// however, we always take the highest prime below the linearly growing threshold in order
// to reduce the risk of accidental regularities leading to cyclic behavior.
func calcDatasetSize(epoch uint64) uint64 {
	size := datasetInitBytes + datasetGrowthBytes*epoch - mixBytes
	for !new(big.Int).SetUint64(size / mixBytes).ProbablyPrime(1) { // Always accurate for n < 2^64
		size -= 2 * mixBytes
	}
	return size
}

// cacheSize returns the size of the ethash verification cache that belongs to a certain
// block number.
func cacheSize(epoch uint64) uint64 {
	if epoch < maxEpoch {
		return cacheSizes[int(epoch)]
	}
	return calcCacheSize(epoch)
}

// seedHash is the seed to use for generating a verification cache and the mining
// dataset. The block number passed should be pre-rounded to an epoch boundary + 1
// e.g: seedHash(calcEpochBlock(epoch, epochLength))
func seedHash(epoch uint64, epochLength uint64) []byte {
	block := calcEpochBlock(epoch, epochLength)

	seed := make([]byte, 32)
	if block < epochLengthDefault {
		return seed
	}

	keccak256 := makeHasher(sha3.NewLegacyKeccak256())
	for i := 0; i < int(block/epochLengthDefault); i++ {
		keccak256(seed, seed)
	}
	return seed
}

// generateCache creates a verification cache of a given size for an input seed.
// The cache production process involves first sequentially filling up 32 MB of
// memory, then performing two passes of Sergio Demian Lerner's RandMemoHash
// algorithm from Strict Memory Hard Hashing Functions (2014). The output is a
// set of 524288 64-byte values.
// This method places the result into dest in machine byte order.
func generateCache(dest []uint32, epoch uint64, epochLength uint64, seed []byte) {
	// Print some debug logs to allow analysis on low end devices
	// logger := log.New("epoch", epoch)

	start := time.Now()
	defer func() {
		elapsed := time.Since(start)

		logFn := log.Debug
		if elapsed > 3*time.Second {
			logFn = log.Info
		}
		logFn("Generated ethash verification cache", "epochLength", epochLength, "elapsed", elapsed)
	}()
	// Convert our destination slice to a byte buffer
	var cache []byte
	cacheHdr := (*reflect.SliceHeader)(unsafe.Pointer(&cache))
	dstHdr := (*reflect.SliceHeader)(unsafe.Pointer(&dest))
	cacheHdr.Data = dstHdr.Data
	cacheHdr.Len = dstHdr.Len * 4
	cacheHdr.Cap = dstHdr.Cap * 4

	// Calculate the number of theoretical rows (we'll store in one buffer nonetheless)
	size := uint64(len(cache))
	rows := int(size) / hashBytes

	// Start a monitoring goroutine to report progress on low end devices
	var progress atomic.Uint32

	done := make(chan struct{})
	defer close(done)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-time.After(3 * time.Second):
				log.Info("Generating ethash verification cache", "epochLength", epochLength, "percentage", progress.Load()*100/uint32(rows)/(cacheRounds+1), "elapsed", time.Since(start))
			}
		}
	}()
	// Create a hasher to reuse between invocations
	keccak512 := makeHasher(sha3.NewLegacyKeccak512())

	// Sequentially produce the initial dataset
	keccak512(cache, seed)
	for offset := uint64(hashBytes); offset < size; offset += hashBytes {
		keccak512(cache[offset:], cache[offset-hashBytes:offset])
		progress.Add(1)
	}
	// Use a low-round version of randmemohash
	temp := make([]byte, hashBytes)

	for i := 0; i < cacheRounds; i++ {
		for j := 0; j < rows; j++ {
			var (
				srcOff = ((j - 1 + rows) % rows) * hashBytes
				dstOff = j * hashBytes
				xorOff = (binary.LittleEndian.Uint32(cache[dstOff:]) % uint32(rows)) * hashBytes
			)
			XORBytes(temp, cache[srcOff:srcOff+hashBytes], cache[xorOff:xorOff+hashBytes])
			keccak512(cache[dstOff:], temp)

			progress.Add(1)
		}
	}
	// Swap the byte order on big endian systems and return
	if !isLittleEndian() {
		swap(cache)
	}
}

// isLittleEndian returns whether the local system is running in little or big
// endian byte order.
func isLittleEndian() bool {
	n := uint32(0x01020304)
	return *(*byte)(unsafe.Pointer(&n)) == 0x04
}

// memoryMap tries to memory map a file of uint32s for read only access.
func memoryMap(path string, lock bool) (*os.File, mmap.MMap, []uint32, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, nil, nil, err
	}
	mem, buffer, err := memoryMapFile(file, false)
	if err != nil {
		file.Close()
		return nil, nil, nil, err
	}
	for i, magic := range dumpMagic {
		if buffer[i] != magic {
			mem.Unmap()
			file.Close()
			return nil, nil, nil, ErrInvalidDumpMagic
		}
	}
	if lock {
		if err := mem.Lock(); err != nil {
			mem.Unmap()
			file.Close()
			return nil, nil, nil, err
		}
	}
	return file, mem, buffer[len(dumpMagic):], err
}

// memoryMapAndGenerate tries to memory map a temporary file of uint32s for write
// access, fill it with the data from a generator and then move it into the final
// path requested.
func memoryMapAndGenerate(path string, size uint64, lock bool, generator func(buffer []uint32)) (*os.File, mmap.MMap, []uint32, error) {
	// Ensure the data folder exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, nil, nil, err
	}
	// Create a huge temporary empty file to fill with data
	temp := path + "." + strconv.Itoa(rand.Int())

	dump, err := os.Create(temp)
	if err != nil {
		return nil, nil, nil, err
	}
	if err = ensureSize(dump, int64(len(dumpMagic))*4+int64(size)); err != nil {
		dump.Close()
		os.Remove(temp)
		return nil, nil, nil, err
	}
	// Memory map the file for writing and fill it with the generator
	mem, buffer, err := memoryMapFile(dump, true)
	if err != nil {
		dump.Close()
		os.Remove(temp)
		return nil, nil, nil, err
	}
	copy(buffer, dumpMagic)

	data := buffer[len(dumpMagic):]
	generator(data)

	if err := mem.Unmap(); err != nil {
		return nil, nil, nil, err
	}
	if err := dump.Close(); err != nil {
		return nil, nil, nil, err
	}
	if err := os.Rename(temp, path); err != nil {
		return nil, nil, nil, err
	}
	return memoryMap(path, lock)
}

// generateDataset generates the entire ethash dataset for mining.
// This method places the result into dest in machine byte order.
func generateDataset(dest []uint32, epoch uint64, epochLength uint64, cache []uint32) {
	// Print some debug logs to allow analysis on low end devices
	// logger := log.New("epoch", epoch)

	start := time.Now()
	defer func() {
		elapsed := time.Since(start)

		logFn := log.Debug
		if elapsed > 3*time.Second {
			logFn = log.Info
		}
		logFn("Generated ethash verification dataset", "epochLength", epochLength, "elapsed", elapsed)
	}()

	// Figure out whether the bytes need to be swapped for the machine
	swapped := !isLittleEndian()

	// Convert our destination slice to a byte buffer
	var dataset []byte
	datasetHdr := (*reflect.SliceHeader)(unsafe.Pointer(&dataset))
	destHdr := (*reflect.SliceHeader)(unsafe.Pointer(&dest))
	datasetHdr.Data = destHdr.Data
	datasetHdr.Len = destHdr.Len * 4
	datasetHdr.Cap = destHdr.Cap * 4

	// Generate the dataset on many goroutines since it takes a while
	threads := runtime.NumCPU()
	size := uint64(len(dataset))

	var pend sync.WaitGroup
	pend.Add(threads)

	var progress atomic.Uint64
	for i := 0; i < threads; i++ {
		go func(id int) {
			defer pend.Done()

			// Create a hasher to reuse between invocations
			keccak512 := makeHasher(sha3.NewLegacyKeccak512())

			// Calculate the data segment this thread should generate
			batch := (size + hashBytes*uint64(threads) - 1) / (hashBytes * uint64(threads))
			first := uint64(id) * batch
			limit := first + batch
			if limit > size/hashBytes {
				limit = size / hashBytes
			}
			// Calculate the dataset segment
			percent := size / hashBytes / 100
			for index := first; index < limit; index++ {
				item := generateDatasetItem(cache, uint32(index), keccak512)
				if swapped {
					swap(item)
				}
				copy(dataset[index*hashBytes:], item)

				if status := progress.Add(1); status%percent == 0 {
					log.Info("Generating DAG in progress", "epochLength", epochLength, "percentage", (status*100)/(size/hashBytes), "elapsed", time.Since(start))
				}
			}
		}(i)
	}
	// Wait for all the generators to finish and return
	pend.Wait()
}

// calcCacheSize calculates the cache size for epoch. The cache size grows linearly,
// however, we always take the highest prime below the linearly growing threshold in order
// to reduce the risk of accidental regularities leading to cyclic behavior.
func calcCacheSize(epoch uint64) uint64 {
	size := cacheInitBytes + cacheGrowthBytes*epoch - hashBytes
	for !new(big.Int).SetUint64(size / hashBytes).ProbablyPrime(1) { // Always accurate for n < 2^64
		size -= 2 * hashBytes
	}
	return size
}

// calcEpochBlock returns the epoch start block for a given epoch (ECIP-1099)
func calcEpochBlock(epoch uint64, epochLength uint64) uint64 {
	return epoch*epochLength + 1
}

// memoryMapFile tries to memory map an already opened file descriptor.
func memoryMapFile(file *os.File, write bool) (mmap.MMap, []uint32, error) {
	// Try to memory map the file
	flag := mmap.RDONLY
	if write {
		flag = mmap.RDWR
	}
	mem, err := mmap.Map(file, flag, 0)
	if err != nil {
		return nil, nil, err
	}
	// The file is now memory-mapped. Create a []uint32 view of the file.
	var view []uint32
	header := (*reflect.SliceHeader)(unsafe.Pointer(&view))
	header.Data = (*reflect.SliceHeader)(unsafe.Pointer(&mem)).Data
	header.Cap = len(mem) / 4
	header.Len = header.Cap
	return mem, view, nil
}

// swap changes the byte order of the buffer assuming a uint32 representation.
func swap(buffer []byte) {
	for i := 0; i < len(buffer); i += 4 {
		binary.BigEndian.PutUint32(buffer[i:], binary.LittleEndian.Uint32(buffer[i:]))
	}
}

// finalizer unmaps the memory and closes the file.
func (c *cache) finalizer() {
	if c.mmap != nil {
		c.mmap.Unmap()
		c.dump.Close()
		c.mmap, c.dump = nil, nil
	}
}

// finalizer closes any file handlers and memory maps open.
func (d *dataset) finalizer() {
	if d.mmap != nil {
		d.mmap.Unmap()
		d.dump.Close()
		d.mmap, d.dump = nil, nil
	}
}

// ensureSize expands the file to the given size. This is to prevent runtime
// errors later on, if the underlying file expands beyond the disk capacity,
// even though it ostensibly is already expanded, but due to being sparse
// does not actually occupy the full declared size on disk.
func ensureSize(f *os.File, size int64) error {
	// On systems which do not support fallocate, we merely truncate it.
	// More robust alternatives  would be to
	// - Use posix_fallocate, or
	// - explicitly fill the file with zeroes.
	return f.Truncate(size)
}

// makePoissonFakeDelay uses the ethash.threads value as a mean time (lambda)
// for a Poisson distribution, returning a random value from
// that discrete function. I think a Poisson distribution probably
// fairly accurately models real world block times.
// Note that this is a hacky way to use ethash.threads since
// lower values will yield faster blocks, but it saves having
// to add or modify any more code than necessary.
func (ethash *Ethash) makePoissonFakeDelay() float64 {
	p := distuv.Poisson{
		Lambda: float64(ethash.Threads()),
		Src:    exprand.NewSource(uint64(time.Now().UnixNano())),
	}
	return p.Rand()
}

// StopRemoteSealer stops the remote sealer
func (ethash *Ethash) StopRemoteSealer() error {
	ethash.closeOnce.Do(func() {
		// Short circuit if the exit channel is not allocated.
		if ethash.remote == nil {
			return
		}
		close(ethash.remote.requestExit)
		<-ethash.remote.exitCh
	})
	return nil
}
package main 

import (
	"fmt"
	"time"
	"errors"
	"math/big"

	"github.com/openrelayxyz/plugeth-utils/core"

	// "github.com/openrelayxyz/plugeth-utils/restricted/consensus"
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

// CalcDifficulty is the difficulty adjustment algorithm. It returns
// the difficulty that a new block should have when created at time
// given the parent block's time and difficulty.
func (ethash *Ethash) CalcDifficulty(chain ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	return CalcDifficulty(chain.Config(), time, parent)
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
	if !chain.Config().IsEnabled(chain.Config().GetEIP1559Transition, header.Number) {
		// Verify BaseFee not present before EIP-1559 fork.
		if header.BaseFee != nil {
			return fmt.Errorf("invalid baseFee before fork: have %d, expected 'nil'", header.BaseFee)
		}
		if err := VerifyGaslimit(parent.GasLimit, header.GasLimit); err != nil {
			return err
		}
	} else if err := VerifyEIP1559Header(chain.Config(), parent, header); err != nil {
		// Verify the header's EIP-1559 attributes.
		return err
	}
	// Verify that the block number is parent's +1
	if diff := new(big.Int).Sub(header.Number, parent.Number); diff.Cmp(big.NewInt(1)) != 0 {
		return ErrInvalidNumber
	}
	if chain.Config().IsEnabledByTime(chain.Config().GetEIP3860TransitionTime, &header.Time) || chain.Config().IsEnabled(chain.Config().GetEIP3860Transition, header.Number) {
		return fmt.Errorf("ethash does not support shanghai fork")
	}
	if chain.Config().IsEnabledByTime(chain.Config().GetEIP4844TransitionTime, &header.Time) {
		return fmt.Errorf("ethash does not support cancun fork")
	}
	// Verify the engine specific seal securing the block
	// if seal {
	// 	if err := ethash.verifySeal(chain, header, false); err != nil {
	// 		return err
	// 	}
	// }
	// If all checks passed, validate any special fields for hard forks
	if err := VerifyDAOHeaderExtraData(chain.Config(), header); err != nil {
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
// func (ethash *Ethash) verifySeal(chain consensus.ChainHeaderReader, header *types.Header, fulldag bool) error {
// 	// If we're running a fake PoW, accept any seal as valid
// 	if ethash.config.PowMode == ModeFake || ethash.config.PowMode == ModePoissonFake || ethash.config.PowMode == ModeFullFake {
// 		time.Sleep(ethash.fakeDelay)
// 		if ethash.fakeFail == header.Number.Uint64() {
// 			return errInvalidPoW
// 		}
// 		return nil
// 	}
// 	// If we're running a shared PoW, delegate verification to it
// 	if ethash.shared != nil {
// 		return ethash.shared.verifySeal(chain, header, fulldag)
// 	}
// 	// Ensure that we have a valid difficulty for the block
// 	if header.Difficulty.Sign() <= 0 {
// 		return errInvalidDifficulty
// 	}
// 	// Recompute the digest and PoW values
// 	number := header.Number.Uint64()

// 	var (
// 		digest []byte
// 		result []byte
// 	)
// 	// If fast-but-heavy PoW verification was requested, use an ethash dataset
// 	if fulldag {
// 		dataset := ethash.dataset(number, true)
// 		if dataset.generated() {
// 			digest, result = hashimotoFull(dataset.dataset, ethash.SealHash(header).Bytes(), header.Nonce.Uint64())

// 			// Datasets are unmapped in a finalizer. Ensure that the dataset stays alive
// 			// until after the call to hashimotoFull so it's not unmapped while being used.
// 			runtime.KeepAlive(dataset)
// 		} else {
// 			// Dataset not yet generated, don't hang, use a cache instead
// 			fulldag = false
// 		}
// 	}
// 	// If slow-but-light PoW verification was requested (or DAG not yet ready), use an ethash cache
// 	if !fulldag {
// 		cache := ethash.cache(number)
// 		epochLength := calcEpochLength(number, ethash.config.ECIP1099Block)
// 		epoch := calcEpoch(number, epochLength)
// 		size := datasetSize(epoch)
// 		if ethash.config.PowMode == ModeTest {
// 			size = 32 * 1024
// 		}
// 		digest, result = hashimotoLight(size, cache.cache, ethash.SealHash(header).Bytes(), header.Nonce.Uint64())

// 		// Caches are unmapped in a finalizer. Ensure that the cache stays alive
// 		// until after the call to hashimotoLight so it's not unmapped while being used.
// 		runtime.KeepAlive(cache)
// 	}
// 	// Verify the calculated values against the ones provided in the header
// 	if !bytes.Equal(header.MixDigest[:], digest) {
// 		return errInvalidMixDigest
// 	}
// 	target := new(big.Int).Div(two256, header.Difficulty)
// 	if new(big.Int).SetBytes(result).Cmp(target) > 0 {
// 		return errInvalidPoW
// 	}
// 	return nil
// }

package main

import (
	"sync"
	"time"
	"math/big"
	"math/rand"
	"runtime"
	"errors"

	"golang.org/x/crypto/sha3"
	mapset "github.com/deckarep/golang-set/v2"

	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted"
	"github.com/openrelayxyz/plugeth-utils/restricted/consensus"
	"github.com/openrelayxyz/plugeth-utils/restricted/params"
	trie "github.com/openrelayxyz/plugeth-utils/restricted/hasher"
	"github.com/openrelayxyz/plugeth-utils/restricted/rlp"
	"github.com/openrelayxyz/plugeth-utils/restricted/types"
)

// type engine struct {
// }

func CreateEngine(chainConfig *params.ChainConfig, db restricted.Database) consensus.Engine {
	
	return &Ethash{

		// config:   config,
		pluginConfig: NewPluginConfig(),
		// caches:   newlru(config.CachesInMem, newCache),
		// datasets: newlru(config.DatasetsInMem, newDataset),
		update:   make(chan struct{}),
		// hashrate: metrics.NewMeterForced(),
	}
}

type Ethash struct {
	config Config

	pluginConfig *PluginConfigurator

	caches   *lru[*cache]   // In memory caches to avoid regenerating too often
	datasets *lru[*dataset] // In memory datasets to avoid regenerating too often

	// Mining related fields
	rand     *rand.Rand    // Properly seeded random source for nonces
	threads  int           // Number of threads to mine on if mining
	update   chan struct{} // Notification channel to update mining parameters
	// hashrate metrics.Meter // Meter tracking the average hashrate TODO PM make conversion to Cardianl metrics library
	remote   *remoteSealer

	// The fields below are hooks for testing
	shared    *Ethash       // Shared PoW verifier to avoid cache regeneration
	fakeFail  uint64        // Block number which fails PoW check even in fake mode
	fakeDelay time.Duration // Time delay to sleep for before returning from verify

	lock      sync.Mutex // Ensures thread safety for the in-memory caches and mining fields
	closeOnce sync.Once  // Ensures exit channel will not be closed twice.

}

func NewPluginConfig() *PluginConfigurator {
	return etc_config
}

// Author implements consensus.Engine, returning the header's coinbase as the
// proof-of-work verified author of the block.
func (ethash *Ethash) Author(header *types.Header) (core.Address, error) {
	return header.Coinbase, nil
}

// VerifyHeader checks whether a header conforms to the consensus rules of the
// stock Ethereum ethash engine.
func (ethash *Ethash) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, seal bool) error {
	// If we're running a full engine faking, accept any input as valid
	if ethash.config.PowMode == ModeFullFake {
		return nil
	}
	// Short circuit if the header is known, or its parent not
	number := header.Number.Uint64()
	if chain.GetHeader(header.Hash(), number) != nil {
		return nil
	}
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return ErrUnknownAncestor
	}
	// Sanity checks passed, do a proper verification
	return ethash.verifyHeader(chain, header, parent, false, seal, time.Now().Unix())
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
// concurrently. The method returns a quit channel to abort the operations and
// a results channel to retrieve the async verifications.
func (ethash *Ethash) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	// If we're running a full engine faking, accept any input as valid
	if ethash.config.PowMode == ModeFullFake || len(headers) == 0 {
		abort, results := make(chan struct{}), make(chan error, len(headers))
		for i := 0; i < len(headers); i++ {
			results <- nil
		}
		return abort, results
	}

	// Spawn as many workers as allowed threads
	workers := runtime.GOMAXPROCS(0)
	if len(headers) < workers {
		workers = len(headers)
	}

	// Create a task channel and spawn the verifiers
	var (
		inputs  = make(chan int)
		done    = make(chan int, workers)
		errors  = make([]error, len(headers))
		abort   = make(chan struct{})
		unixNow = time.Now().Unix()
	)
	for i := 0; i < workers; i++ {
		go func() {
			for index := range inputs {
				errors[index] = ethash.verifyHeaderWorker(chain, headers, seals, index, unixNow)
				done <- index
			}
		}()
	}

	errorsOut := make(chan error, len(headers))
	go func() {
		defer close(inputs)
		var (
			in, out = 0, 0
			checked = make([]bool, len(headers))
			inputs  = inputs
		)
		for {
			select {
			case inputs <- in:
				if in++; in == len(headers) {
					// Reached end of headers. Stop sending to workers.
					inputs = nil
				}
			case index := <-done:
				for checked[index] = true; checked[out]; out++ {
					errorsOut <- errors[out]
					if out == len(headers)-1 {
						return
					}
				}
			case <-abort:
				return
			}
		}
	}()
	return abort, errorsOut
}

func (ethash *Ethash) verifyHeaderWorker(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool, index int, unixNow int64) error {
	var parent *types.Header
	if index == 0 {
		parent = chain.GetHeader(headers[0].ParentHash, headers[0].Number.Uint64()-1)
	} else if headers[index-1].Hash() == headers[index].ParentHash {
		parent = headers[index-1]
	}
	if parent == nil {
		return ErrUnknownAncestor
	}
	return ethash.verifyHeader(chain, headers[index], parent, false, seals[index], unixNow)
}

// VerifyUncles verifies that the given block's uncles conform to the consensus
// rules of the stock Ethereum ethash engine.
func (ethash *Ethash) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	// If we're running a full engine faking, accept any input as valid
	if ethash.config.PowMode == ModeFullFake {
		return nil
	}
	// Verify that there are at most 2 uncles included in this block
	if len(block.Uncles()) > maxUncles {
		return errTooManyUncles
	}
	if len(block.Uncles()) == 0 {
		return nil
	}
	// Gather the set of past uncles and ancestors
	uncles, ancestors := mapset.NewSet[core.Hash](), make(map[core.Hash]*types.Header)

	number, parent := block.NumberU64()-1, block.ParentHash()
	for i := 0; i < 7; i++ {
		ancestorHeader := chain.GetHeader(parent, number)
		if ancestorHeader == nil {
			break
		}
		ancestors[parent] = ancestorHeader
		// If the ancestor doesn't have any uncles, we don't have to iterate them
		if ancestorHeader.UncleHash != types.EmptyUncleHash {
			// Need to add those uncles to the banned list too
			ancestor := chain.GetBlock(parent, number)
			if ancestor == nil {
				break
			}
			for _, uncle := range ancestor.Uncles() {
				uncles.Add(uncle.Hash())
			}
		}
		parent, number = ancestorHeader.ParentHash, number-1
	}
	ancestors[block.Hash()] = block.Header()
	uncles.Add(block.Hash())

	// Verify each of the uncles that it's recent, but not an ancestor
	for _, uncle := range block.Uncles() {
		// Make sure every uncle is rewarded only once
		hash := uncle.Hash()
		if uncles.Contains(hash) {
			return errDuplicateUncle
		}
		uncles.Add(hash)

		// Make sure the uncle has a valid ancestry
		if ancestors[hash] != nil {
			return errUncleIsAncestor
		}
		if ancestors[uncle.ParentHash] == nil || uncle.ParentHash == block.ParentHash() {
			return errDanglingUncle
		}
		if err := ethash.verifyHeader(chain, uncle, ancestors[uncle.ParentHash], true, true, time.Now().Unix()); err != nil {
			return err
		}
	}
	return nil
}

// Prepare implements consensus.Engine, initializing the difficulty field of a
// header to conform to the ethash protocol. The changes are done inline.
func (ethash *Ethash) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	parent := chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
	if parent == nil {
		return ErrUnknownAncestor
	}
	header.Difficulty = ethash.CalcDifficulty(chain, header.Time, parent)
	return nil
}

// Finalize implements consensus.Engine, accumulating the block and uncle rewards.
func (ethash *Ethash) Finalize(chain consensus.ChainHeaderReader, header *types.Header, state core.RWStateDB, txs []*types.Transaction, uncles []*types.Header, withdrawals []*types.Withdrawal) {
	// Accumulate any block and uncle rewards and commit the final state root
	AccumulateRewards(ethash.pluginConfig, state, header, uncles)
}

// FinalizeAndAssemble implements consensus.Engine, accumulating the block and
// uncle rewards, setting the final state and assembling the block.
func (ethash *Ethash) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state core.RWStateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt, withdrawals []*types.Withdrawal) (*types.Block, error) {
	if len(withdrawals) > 0 {
		return nil, errors.New("ethash does not support withdrawals")
	}
	// Finalize block
	ethash.Finalize(chain, header, state, txs, uncles, nil)


	// Assign the final state root to header.
	header.Root = state.IntermediateRoot(ethash.pluginConfig.IsEnabled(ethash.pluginConfig.GetEIP161dTransition, header.Number))

	// Header seems complete, assemble into a block and return
	return types.NewBlock(header, txs, uncles, receipts, trie.NewStackTrie(nil)), nil
}

// Seal implements consensus.Engine, attempting to find a nonce that satisfies
// the block's difficulty requirements.
// func (ethash *Ethash) Seal(chain consensus.ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
// 	// If we're running a fake PoW, simply return a 0 nonce immediately
// 	if ethash.config.PowMode == ModeFake || ethash.config.PowMode == ModeFullFake {
// 		header := block.Header()
// 		header.Nonce, header.MixDigest = types.BlockNonce{}, core.Hash{}
// 		select {
// 		case results <- block.WithSeal(header):
// 		default:
// 			ethash.config.Log.Warn("Sealing result is not read by miner", "mode", "fake", "sealhash", ethash.SealHash(block.Header()))
// 		}
// 		return nil
// 	} else if ethash.config.PowMode == ModePoissonFake {
// 		go func(header *types.Header) {
// 			// Assign random (but non-zero) values to header nonce and mix.
// 			header.Nonce = types.EncodeNonce(uint64(rand.Int63n(math.MaxInt64)))
// 			b, _ := header.Nonce.MarshalText()
// 			header.MixDigest = core.BytesToHash(b)

// 			// Wait some amount of time.
// 			timeout := time.NewTimer(time.Duration(ethash.makePoissonFakeDelay()) * time.Second)
// 			defer timeout.Stop()

// 			select {
// 			case <-stop:
// 				return
// 			case <-ethash.update:
// 				timeout.Stop()
// 				if err := ethash.Seal(chain, block, results, stop); err != nil {
// 					ethash.config.Log.Error("Failed to restart sealing after update", "err", err)
// 				}
// 			case <-timeout.C:
// 				// Send the results when the timeout expires.
// 				select {
// 				case results <- block.WithSeal(header):
// 				default:
// 					ethash.config.Log.Warn("Sealing result is not read by miner", "mode", "fake", "sealhash", ethash.SealHash(block.Header()))
// 				}
// 			}
// 		}(block.Header())
// 		return nil
// 	}
// 	// If we're running a shared PoW, delegate sealing to it
// 	if ethash.shared != nil {
// 		return ethash.shared.Seal(chain, block, results, stop)
// 	}
// 	// Create a runner and the multiple search threads it directs
// 	abort := make(chan struct{})

// 	ethash.lock.Lock()
// 	threads := ethash.threads
// 	if ethash.rand == nil {
// 		seed, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
// 		if err != nil {
// 			ethash.lock.Unlock()
// 			return err
// 		}
// 		ethash.rand = rand.New(rand.NewSource(seed.Int64()))
// 	}
// 	ethash.lock.Unlock()
// 	if threads == 0 {
// 		threads = runtime.NumCPU()
// 	}
// 	if threads < 0 {
// 		threads = 0 // Allows disabling local mining without extra logic around local/remote
// 	}
// 	// Push new work to remote sealer
// 	if ethash.remote != nil {
// 		ethash.remote.workCh <- &sealTask{block: block, results: results}
// 	}
// 	var (
// 		pend   sync.WaitGroup
// 		locals = make(chan *types.Block)
// 	)
// 	for i := 0; i < threads; i++ {
// 		pend.Add(1)
// 		go func(id int, nonce uint64) {
// 			defer pend.Done()
// 			ethash.mine(block, id, nonce, abort, locals)
// 		}(i, uint64(ethash.rand.Int63()))
// 	}
// 	// Wait until sealing is terminated or a nonce is found
// 	go func() {
// 		var result *types.Block
// 		select {
// 		case <-stop:
// 			// Outside abort, stop all miner threads
// 			close(abort)
// 		case result = <-locals:
// 			// One of the threads found a block, abort all others
// 			select {
// 			case results <- result:
// 			default:
// 				ethash.config.Log.Warn("Sealing result is not read by miner", "mode", "local", "sealhash", ethash.SealHash(block.Header()))
// 			}
// 			close(abort)
// 		case <-ethash.update:
// 			// Thread count was changed on user request, restart
// 			close(abort)
// 			if err := ethash.Seal(chain, block, results, stop); err != nil {
// 				ethash.config.Log.Error("Failed to restart sealing after update", "err", err)
// 			}
// 		}
// 		// Wait for all miners to terminate and return the block
// 		pend.Wait()
// 	}()
// 	return nil
// }

// // SealHash returns the hash of a block prior to it being sealed.
func (ethash *Ethash) SealHash(header *types.Header) (hash core.Hash) {
	hasher := sha3.NewLegacyKeccak256()

	enc := []interface{}{
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra,
	}
	if header.BaseFee != nil {
		enc = append(enc, header.BaseFee)
	}
	if header.WithdrawalsHash != nil {
		panic("withdrawal hash set on ethash")
	}
	rlp.Encode(hasher, enc)
	hasher.Sum(hash[:0])
	return hash
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns
// the difficulty that a new block should have when created at time
// given the parent block's time and difficulty.
func (ethash *Ethash) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	return CalcDifficulty(ethash.pluginConfig, time, parent)
}

func (e *Ethash) APIs(chain consensus.ChainHeaderReader) []core.API {
	return []core.API{}
}

// Close closes the exit channel to notify all backend threads exiting.
func (ethash *Ethash) Close() error {
	return ethash.StopRemoteSealer()
}

// func CreateEngine(chainConfig *params.ChainConfig, db restricted.Database) consensus.Engine {
// 	return &engine{}
// }
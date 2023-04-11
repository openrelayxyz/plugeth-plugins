package main

import(
	"errors"
	"math/big"

	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted"
	"github.com/openrelayxyz/plugeth-utils/restricted/types"
	"github.com/openrelayxyz/plugeth-utils/restricted/hasher"
	"github.com/openrelayxyz/plugeth-utils/restricted/consensus"
)

var (
	pl      core.PluginLoader
	backend restricted.Backend
	log     core.Logger
	events  core.Feed
)
var httpApiFlagName = "http.api"

func Initialize(ctx core.Context, loader core.PluginLoader, logger core.Logger) { 
	pl = loader
	events = pl.GetFeed()
	log = logger
	v := ctx.String(httpApiFlagName)
	if v != "" {
		ctx.Set(httpApiFlagName, v+",plugeth")
	} else {
		ctx.Set(httpApiFlagName, "eth,net,web3,plugeth")
		log.Info("Loaded consensus engine plugin")
	}
}

type engine struct {
}

func (e *engine) Author(header *types.Header) (core.Address, error) {
	// log.Error("inside of author")
	// log.Error("Inside of Author()", "header", header)
	return core.Address{}, nil
}
func (e *engine) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, seal bool) error {
	// log.Error("inside of verifyHeader")
	return nil
}
func (e *engine) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	quit := make(chan struct{})
	err := make(chan error)
	go func () {
		for i, h := range headers {
			select {
			case <-quit:
				return 
			case err<- e.VerifyHeader(chain, h, seals[i]):
			}
		} 
	} ()
	// log.Error("inside of verifyHeaders")
	return quit, err
}
func (e *engine) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	// log.Error("inside of verify uncles")
	return nil
}
func (e *engine) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	header.Difficulty = new(big.Int).SetUint64(123456789)
	header.UncleHash = core.HexToHash("1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347")
	// log.Error("inside of prepare")
	// log.Error("inside of prepare", "header", header)
	return nil
}
func (e *engine) Finalize(chain consensus.ChainHeaderReader, header *types.Header, state core.RWStateDB, txs []*types.Transaction,uncles []*types.Header, withdrawals []*types.Withdrawal) {
	// log.Error("inside of Finalize")
}
func (e *engine) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state core.RWStateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt, withdrawals []*types.Withdrawal) (*types.Block, error) {
	log.Error("inside of FinalizeAndAssemble")
	// header.Root = state.IntermediateRoot(false)
	header.Root = core.HexToHash("0x8fb57ddad3e9f83990092976adcec803324fd045988fd60225058f87a1439eab")
	// return types.NewBlockWithHeader(header).WithBody(txs, uncles).WithWithdrawals(withdrawals), nil
	hasher := hasher.NewStackTrie(nil)
	block := types.NewBlockWithWithdrawals(header, txs, uncles, receipts, withdrawals, hasher)
	return block, nil
}
func (e *engine) Seal(chain consensus.ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	// log.Error("inside of Seal")
	if len(block.Transactions()) == 0 {
		return errors.New("sealing paused while waiting for transactions")
	}
	go func () {
		results <- block 
		close(results)
	} ()
	// TO DO: the stop channel will need to be addressed in a non test case scenerio
	return nil
}
func (e *engine) SealHash(header *types.Header) core.Hash {
	// log.Error("inside of SealHash", "header.Hash()", header.Hash())
	return header.Hash()
}
func (e *engine) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	log.Error("inside of CalcDifficulty")
	return new(big.Int).SetUint64(uint64(123456789))
}
func (e *engine) APIs(chain consensus.ChainHeaderReader) []core.API {
	// log.Error("inside of APIs")
	return []core.API{}
}
func (e *engine) Close() error {
	 log.Error("inside of Close")
	return nil
}

func CreateEngine(core.Node, []string, bool, restricted.Database) consensus.Engine {
	return &engine{}
}
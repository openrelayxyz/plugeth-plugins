package main

import (
	"context"
	"math/big"
	"errors"
	
	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted"
	"github.com/openrelayxyz/plugeth-utils/restricted/types"

	"github.com/openrelayxyz/plugeth-utils/restricted/consensus"
)

var log core.Logger

var httpApiFlagName = "http.api"

func Initialize(ctx core.Context, loader core.PluginLoader, logger core.Logger) {
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
	// log.Error("inside of verifyHeadersss")
	return quit, err
}
func (e *engine) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	// log.Error("inside of verify uncles")
	return nil
}
func (e *engine) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	header.Difficulty = new(big.Int).SetUint64(123456789)
	// log.Error("inside of prepare")
	return nil
}
func (e *engine) Finalize(chain consensus.ChainHeaderReader, header *types.Header, state core.RWStateDB, txs []*types.Transaction,uncles []*types.Header, withdrawals []*types.Withdrawal) {
	// log.Error("inside of Finalize")
}
func (e *engine) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state core.RWStateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt, withdrawals []*types.Withdrawal) (*types.Block, error) {
	// log.Error("inside of FinalizeAndAssemble")
			
	return types.NewBlockWithHeader(header).WithBody(txs, uncles).WithWithdrawals(withdrawals), nil
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
	//  log.Error("inside of SealHash")
	return header.Hash()
}
func (e *engine) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	//  log.Error("inside of CalcDifficulty")
	return new(big.Int).SetUint64(uint64(123456789))
}
func (e *engine) APIs(chain consensus.ChainHeaderReader) []core.API {
	//  log.Error("inside of APIs")
	return []core.API{}
}
func (e *engine) Close() error {
	//  log.Error("inside of Close")
	return nil
}

func CreateEngine(core.Node, []string, bool, restricted.Database) consensus.Engine {
	return &engine{}
}

var apis []core.API

func GetAPIs(stack core.Node, backend core.Backend) []core.API {
	apis = []core.API{
		{
			Namespace: "plugeth",
			Version:   "1.0",
			Service:   &engineService{backend, stack},
			Public:    true,
		},
	}

	return apis
}

type engineService struct {
	backend core.Backend
	stack core.Node
}


var errs []error

func HookTester() {
	
	cl := apis[0].Service.(*engineService).stack
	client, err := cl.Attach()
	if err != nil {
		errs = append(errs, err)
		log.Error("Error connecting with client")
	}

	var x interface{}
	err = client.Call(&x, "plugeth_test")
	if err != nil {
		errs = append(errs, err)
		log.Error("failed to call method plugeth_test", "err", err)
	}
	log.Error("this is the return value for isSynced", "test", x, "len", len(errs))

	var y interface{}
	err = client.Call(&y, "eth_coinbase")
	if err != nil {
		errs = append(errs, err)
		log.Error("failed to call method eth_coinbase", "err", err)
	}
	log.Error("this is the return value for eth_coinbase", "test", y, "len", len(errs))

	

// 	if len(errs) > 0 {
// 		for _, err := range errs {
// 		log.Error("Error", "err", err)
// 		}
// 	// os.Exit(1)
// 	}

// 	// os.Exit(0)
	
}

// what should be done next is set up a json object to pass into eth_sendTransactions so we can append one block to the chain from there we can start to 
// experiment on what can be done to excerise the hooks

// the as of now command to start geth is => /geth --dev --http --http.api eth --verbosity=5,  with this plugin loaded 


func (service *engineService) Test(ctx context.Context) string {
	return "this is a placeholder function"
}


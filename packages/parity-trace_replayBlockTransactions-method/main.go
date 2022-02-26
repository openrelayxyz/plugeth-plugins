package main



import (
	"context"

	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted/hexutil"
	"github.com/openrelayxyz/plugeth-utils/restricted/rlp"
"github.com/openrelayxyz/plugeth-utils/restricted/types"

	"gopkg.in/urfave/cli.v1"
)



type FinalResult struct {
	// Output    string          `json:"output"`
	StateDiff map[string]*LayerTwo        `json:"stateDiff"`
	Trace     []*ParityResult `json:"trace"`
	TransactionHash core.Hash       `json:"transactionHash"`
	//VMTrace   interface{}        `json:"vmTrace"`
}

type ParityTrace struct {
	backend core.Backend
	stack   core.Node
}

var Tracers = map[string]func(core.StateDB,  core.BlockContext) core.TracerResult{
	// "plugethVMTracer": func(sdb core.StateDB, bctx core.BlockContext) core.TracerResult {
	// 	return &VMTracerService{StateDB: sdb}
	// },
	"plugethStateDiffTracer": func(sdb core.StateDB, bctx core.BlockContext) core.TracerResult {
		return &SDTracerService{stateDB: sdb, blockContext: bctx}
	},
}

func GetAPIs(stack core.Node, backend core.Backend) []core.API {
	return []core.API{
		{
			Namespace: "trace",
			Version:   "1.0",
			Service:   &ParityTrace{backend, stack},
			Public:    true,
		},
	}
}

var log core.Logger
var httpApiFlagName = "http.api"

func Initialize(ctx *cli.Context, loader core.PluginLoader, logger core.Logger) {
	log = logger
	v := ctx.GlobalString(httpApiFlagName)
	if v != "" {
		ctx.GlobalSet(httpApiFlagName, v+",trace")
	} else {
		ctx.GlobalSet(httpApiFlagName, "eth,net,web3,trace")
		log.Info("Loaded tester plugin")
	}
}

// func (vm *ParityTrace) ReplayTransaction(ctx context.Context, txHash core.Hash, tracerType []string) (interface{}, error) {
// 	result := &FinalResult{}
// 	var err error
//
// 	for _, typ := range tracerType {
// 		if typ == "trace" {
// 				result.Trace, err = vm.TraceVariant(ctx, txHash)
// 				if err != nil {return nil, err}
// 				}
// 		if typ == "vmTrace" {
// 			  result.VMTrace, err = vm.VMTraceVariant(ctx, txHash)
// 					if err != nil {return nil, err}
// 		    }
// 		if typ == "stateDiff" {
// 			result.StateDiff, err = vm.StateDiffVariant(ctx, txHash)
// 				if err != nil {return nil, err}
// 				    }
// 		}
// 	return result, nil
// }

func (pt *ParityTrace) ReplayBlockTransactions(ctx context.Context, bkNum string, tracerType []string) (interface{}, error) {
	// result := &FinalResult{}
	var err error
	// var traces [][]*ParityResult
	var diffs []map[string]*LayerTwo

	for _, typ := range tracerType {
		// if typ == "trace" {
		// 		traces, err = pt.TraceVariant(ctx, bkNum)
		// 		if err != nil {return nil, err}
		// 		}
		// if typ == "vmTrace" {
		// 		result.VMTrace, err = vm.VMTraceVariant(ctx, txHash)
		// 			if err != nil {return nil, err}
		// 		}
		if typ == "stateDiff" {
			diffs, err = pt.StateDiffVariant(ctx, bkNum)
				if err != nil {return nil, err}
						}
		}

	blockNM, err := hexutil.DecodeUint64(bkNum)
	if err != nil {
		return nil, err
	}
	block := &types.Block{}
	bkNB := int64(blockNM)
	rlpBlock, err := pt.backend.BlockByNumber(ctx, bkNB)
	if err != nil {
		return nil, err
	}
	if err := rlp.DecodeBytes(rlpBlock, block); err != nil {
		return nil, err
	}

	// transactions := block.Transactions()

	// tr := [][]*ParityResult{}
	// for _, item := range gr {
	// 	tAddress := make([]int, 0)
	// 	pr = append(pr, GethParity(item.Result, tAddress, strings.ToLower(item.Result.Type)))
	// }
	// results := make([]FinalResult, len(transactions))
	// for i, _ := range results {
	// 	// if gr[i].Result.Output == "" {
	// 	// 	gr[i].Result.Output = "0x"
	// 	// }
	// 	results[i] = FinalResult{
	// 		// Output:          gr[i].Result.Output,
	// 		StateDiff:       diffs[i],
	// 		Trace:           traces[i],
	// 		TransactionHash: transactions[i].Hash(),
	// 		// VMTrace:         nil,
	// 	}
	// }
	return diffs, nil
}

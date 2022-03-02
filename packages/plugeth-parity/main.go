package main



import (
	"context"

	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted"
	"github.com/openrelayxyz/plugeth-utils/restricted/hexutil"
	"github.com/openrelayxyz/plugeth-utils/restricted/types"
	"gopkg.in/urfave/cli.v1"
)

type FinalResult struct {
	Output    string          `json:"output"`
	StateDiff map[string]*LayerTwo        `json:"stateDiff"`
	Trace     []*ParityResult `json:"trace"`
	VMTrace   interface{}        `json:"vmTrace"`
}

type ParityTrace struct {
	backend restricted.Backend
	stack   core.Node
}

var Tracers = map[string]func(core.StateDB,  core.BlockContext) core.TracerResult{
	"plugethVMTracer": func(sdb core.StateDB, bctx core.BlockContext) core.TracerResult {
		return &VMTracerService{StateDB: sdb}
	},
	"plugethStateDiffTracer": func(sdb core.StateDB, bctx core.BlockContext) core.TracerResult {
		return &SDTracerService{stateDB: sdb, blockContext: bctx, log:log}
	},
}

func GetAPIs(stack core.Node, backend restricted.Backend) []core.API {
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

func (pt *ParityTrace) Call(ctx context.Context, txObject map[string]interface{}, tracerType []string) (interface{}, error) {
	result := &FinalResult{}
	var output string
	var err error

	for _, typ := range tracerType {
		if typ == "trace" {
				result.Trace, output, err = pt.TraceVariant(ctx, txObject)
				if err != nil {return nil, err}
			}
		if typ == "vmTrace" {
			  result.VMTrace, output, err = pt.VMTraceVariant(ctx, txObject)
					if err != nil {return nil, err}
			}
		if typ == "stateDiff" {
			result.StateDiff, output, err = pt.StateDiffVariant(ctx, txObject)
				if err != nil {return nil, err}
		}
	}
	result.Output = output
	return result, nil
}


func (pt *ParityTrace) RawTransaction(ctx context.Context, data hexutil.Bytes, tracerType []string) (interface{}, error) {
	result := &FinalResult{}
	var output string
	var err error
	tx := types.Transaction{}
	err = tx.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}
	//config := *params.GoerliChainConfig
	config := pt.backend.ChainConfig()
	hs := types.LatestSigner(config)
	sender, err := hs.Sender(&tx)
	if err != nil {
		return nil, err
	}
	txObject := make(map[string]interface{})
	txObject["from"] = sender
	txObject["to"] = tx.To()
	gas := hexutil.EncodeUint64(tx.Gas())
	txObject["gas"] = gas
	dt := hexutil.Encode(tx.Data())
	txObject["data"] = dt
	price := hexutil.EncodeBig(tx.GasPrice())
	txObject["gasPrice"] = price
	vl := hexutil.EncodeBig(tx.Value())
	txObject["value"] = vl


	for _, typ := range tracerType {
		if typ == "trace" {
				result.Trace, output, err = pt.TraceVariant(ctx, txObject)
				if err != nil {return nil, err}
				}
		if typ == "vmTrace" {
			  result.VMTrace, output, err = pt.VMTraceVariant(ctx, txObject)
					if err != nil {return nil, err}
		    }
		if typ == "stateDiff" {
			result.StateDiff, output, err = pt.StateDiffVariant(ctx, txObject)
				if err != nil {return nil, err}
				    }
		}
	result.Output = output
	return result, nil
}

package main



import (
	"context"

	"github.com/openrelayxyz/plugeth-utils/core"
	"gopkg.in/urfave/cli.v1"
)

type FinalResult struct {
	Output    string          `json:"output"`
	StateDiff map[string]*LayerTwo        `json:"stateDiff"`
	Trace     []*ParityResult `json:"trace"`
	VMTrace   interface{}        `json:"vmTrace"`
}

type ParityTrace struct {
	backend core.Backend
	stack   core.Node
}

var Tracers = map[string]func(core.StateDB,  core.BlockContext) core.TracerResult{
	"plugethVMTracer": func(sdb core.StateDB, bctx core.BlockContext) core.TracerResult {
		return &VMTracerService{StateDB: sdb}
	},
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

func (vm *ParityTrace) Call(ctx context.Context, txObject map[string]string, tracerType []string) (interface{}, error) {
	result := &FinalResult{}
	var err error

	for _, typ := range tracerType {
		if typ == "trace" {
				result.Trace, err = vm.TraceVariant(ctx, txObject, "latest")
				if err != nil {return nil, err}
				}
		if typ == "vmTrace" {
			  result.VMTrace, err = vm.VMTraceVariant(ctx, txObject, "latest")
					if err != nil {return nil, err}
		    }
		if typ == "stateDiff" {
			result.StateDiff, err = vm.StateDiffVariant(ctx, txObject, "latest")
				if err != nil {return nil, err}
				    }
		}
	return result, nil
}

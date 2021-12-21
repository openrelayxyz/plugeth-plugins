package main

import (
	"context"
	"math/big"
	"time"

	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted/hexutil"
	"gopkg.in/urfave/cli.v1"
)

type OuterResult struct {
	Output    hexutil.Bytes `json:"output"`
	StateDiff *string       `json:"stateDiff"`
	Trace     []string      `json:"trace"`
	VMTrace   VMTrace       `json:"vmTrace"`
}

type VMTrace struct {
	Code   hexutil.Bytes `json:"code"`
	Ops    []Ops         `json:"ops"`
	parent *VMTrace
}

type Ops struct {
	Cost uint64   `json:"cost"`
	Ex   string   `json:"ex"`
	PC   uint64   `json:"pc"`
	Sub  *VMTrace `json:"sub"`
}

type ParityVMTrace struct {
	backend core.Backend
	stack   core.Node
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
		log.Info("Loaded Open Ethereum vmTracer plugin")
	}
}

func GetAPIs(stack core.Node, backend core.Backend) []core.API {
	return []core.API{
		{
			Namespace: "trace",
			Version:   "1.0",
			Service:   &ParityVMTrace{backend, stack},
			Public:    true,
		},
	}
}

var Tracers = map[string]func(core.StateDB) core.TracerResult{
	"plugethVMTracer": func(sdb core.StateDB) core.TracerResult {
		return &TracerService{StateDB: sdb}
	},
}

func (vm *ParityVMTrace) ReplayTransaction(ctx context.Context, txHash core.Hash, tracer []string) (interface{}, error) {
	client, err := vm.stack.Attach()
	if err != nil {
		return nil, err
	}
	response := make(map[string]string)
	client.Call(&response, "eth_getTransactionByHash", txHash)
	var code hexutil.Bytes
	err = client.Call(&code, "eth_getCode", response["to"], response["blockNumber"])
	if err != nil {
		return nil, err
	}
	tr := TracerResult{}
	err = client.Call(&tr, "debug_traceTransaction", txHash, map[string]string{"tracer": "plugethVMTracer"})
	//var result []interface{}
	//result = append(result, tr.Count, tr.CountTwo, tr.CountThree, tr.CountFour, tr.CountFive)
	//result = append(result, tr.PCs)
	ops := make([]Ops, tr.Count)
	for i := range ops {
		ops[i].Cost = tr.Costs[i]
		ops[i].PC = tr.PCs[i]
	}
	trace := make([]string, 0)
	result := &OuterResult{
		Output:    tr.Output,
		StateDiff: nil,
		Trace:     trace,
		VMTrace: VMTrace{
			Code: code,
			Ops:  ops,
		},
	}

	return result, nil
}

//Note: If transactions is a contract deployment then the input is the 'code' that we are trying to capture with getCode

type TracerService struct {
	StateDB      core.StateDB
	Output       hexutil.Bytes
	Cost         uint64
	PC           uint64
	Costs        []uint64
	PCs          []uint64
	Greeting     string
	Depth        []int
	OpCodes      []core.OpCode
	Ops          []core.OpCode
	ErrorOps     []core.OpCode
	GasUsed      []uint64
	Count        int
	CountTwo     int
	CountThree   int
	CountFour    int
	CountFive    int
	CurrentTrace *VMTrace
}

func (r *TracerService) CaptureStart(from core.Address, to core.Address, create bool, input []byte, gas uint64, value *big.Int) {
	r.Greeting = "Goodbuy Horses"
	r.OpCodes = []core.OpCode{}
	r.Ops = []core.OpCode{}
	r.Depth = []int{}
	r.ErrorOps = []core.OpCode{}
	r.GasUsed = []uint64{}
	r.Costs = []uint64{}
	r.CurrentTrace = &VMTrace{Code: r.StateDB.GetCode(to), Ops: []Ops{}}
}
func (r *TracerService) CaptureState(pc uint64, op core.OpCode, gas, cost uint64, scope core.ScopeContext, rData []byte, depth int, err error) {
	if depth == 1 {
		r.OpCodes = append(r.OpCodes, op)
		r.Costs = append(r.Costs, cost)
		r.PCs = append(r.PCs, pc)
	}
	r.Depth = append(r.Depth, depth)
	// if depth > 1 {
	// 	r.SecodaryResult.Id =
	// }

	//append to r.CurrentTrace.Ops

}
func (r *TracerService) CaptureFault(pc uint64, op core.OpCode, gas, cost uint64, scope core.ScopeContext, depth int, err error) {
	r.ErrorOps = append(r.ErrorOps, op)
}
func (r *TracerService) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) {
	r.GasUsed = append(r.GasUsed, gasUsed)
	r.Output = output
}
func (r *TracerService) CaptureEnter(typ core.OpCode, from core.Address, to core.Address, input []byte, gas uint64, value *big.Int) {
	//r.Ops = append(r.Ops, typ)
	trace := &VMTrace{Code: r.StateDB.GetCode(to), Ops: []Ops{}, parent: r.CurrentTrace}
	r.CurrentTrace.Ops[len(r.CurrentTrace.Ops)-1].Sub = trace
	r.CurrentTrace = trace

}
func (r *TracerService) CaptureExit(output []byte, gasUsed uint64, err error) {
	r.CurrentTrace = r.CurrentTrace.parent
}
func (r *TracerService) Result() (interface{}, error) {

	r.Count = len(r.OpCodes)
	r.CountTwo = len(r.Ops)
	r.CountThree = len(r.Depth)
	r.CountFour = len(r.ErrorOps)
	r.CountFive = len(r.GasUsed)
	return r, nil
}

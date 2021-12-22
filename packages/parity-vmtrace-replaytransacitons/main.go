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
	tr := TracerService{}
	err = client.Call(&tr, "debug_traceTransaction", txHash, map[string]string{"tracer": "plugethVMTracer"})
	return tr, nil
}

//Note: If transactions is a contract deployment then the input is the 'code' that we are trying to capture with getCode

type TracerService struct {
	StateDB core.StateDB
	CurrentTrace *VMTrace
}

func (r *TracerService) CaptureStart(from core.Address, to core.Address, create bool, input []byte, gas uint64, value *big.Int) {
	r.CurrentTrace = &VMTrace{Code: r.StateDB.GetCode(to), Ops: []Ops{}}
}
func (r *TracerService) CaptureState(pc uint64, op core.OpCode, gas, cost uint64, scope core.ScopeContext, rData []byte, depth int, err error) {
	ops := Ops{Cost: cost,
		PC: pc}
	r.CurrentTrace.Ops = append(r.CurrentTrace.Ops, ops)
}
func (r *TracerService) CaptureFault(pc uint64, op core.OpCode, gas, cost uint64, scope core.ScopeContext, depth int, err error) {
}
func (r *TracerService) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) {
}
func (r *TracerService) CaptureEnter(typ core.OpCode, from core.Address, to core.Address, input []byte, gas uint64, value *big.Int) {
	trace := &VMTrace{Code: r.StateDB.GetCode(to), Ops: []Ops{}, parent: r.CurrentTrace}
	r.CurrentTrace.Ops[len(r.CurrentTrace.Ops)-1].Sub = trace
	r.CurrentTrace = trace
}
func (r *TracerService) CaptureExit(output []byte, gasUsed uint64, err error) {
	r.CurrentTrace = r.CurrentTrace.parent
}
func (r *TracerService) Result() (interface{}, error) {
	return r, nil
}

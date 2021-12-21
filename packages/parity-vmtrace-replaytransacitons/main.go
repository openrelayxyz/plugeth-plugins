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
	Code hexutil.Bytes `json:"code"`
	Ops  []core.OpCode `json:"ops"`
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
	"myTracer": func(core.StateDB) core.TracerResult {
		return &ParityVMTrace{}
	},
}

func (vm *ParityVMTrace) ReplayTransaction(ctx context.Context, txHash core.Hash, tracer []string) (interface{}, error) {
	client, err := vm.stack.Attach()
	if err != nil {
		return nil, err
	}
	response := make(map[string]string)
	client.Call(&response, "eth_getTransactionByHash", txHash)
	var code interface{}
	err = client.Call(&code, "eth_getCode", response["to"], response["blockNumber"])
	if err != nil {
		return nil, err
	}
	return code, nil
}

//Note: If transactions is a contract deployment then the input is the 'code' that we are trying to capture with getCode

func (b *ParityVMTrace) CaptureStart(from core.Address, to core.Address, create bool, input []byte, gas uint64, value *big.Int) {
}
func (b *ParityVMTrace) CaptureState(pc uint64, op core.OpCode, gas, cost uint64, scope core.ScopeContext, rData []byte, depth int, err error) {
}
func (b *ParityVMTrace) CaptureFault(pc uint64, op core.OpCode, gas, cost uint64, scope core.ScopeContext, depth int, err error) {
}
func (b *ParityVMTrace) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) {
}
func (b *ParityVMTrace) CaptureEnter(typ core.OpCode, from core.Address, to core.Address, input []byte, gas uint64, value *big.Int) {
}
func (b *ParityVMTrace) CaptureExit(output []byte, gasUsed uint64, err error) {
}
func (b *ParityVMTrace) Result() (interface{}, error) { return "hello world", nil }

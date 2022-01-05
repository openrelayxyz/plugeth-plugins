package main

import (
	"context"
	"math/big"
	"strconv"
	"time"

	"github.com/holiman/uint256"
	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted"
	"github.com/openrelayxyz/plugeth-utils/restricted/hexutil"
	"gopkg.in/urfave/cli.v1"
)

type OuterResult struct {
	Output    hexutil.Bytes `json:"output"`
	StateDiff *string       `json:"stateDiff"`
	Trace     []string      `json:"trace"`
	VMTrace   *VMTrace      `json:"vmTrace"`
}

type VMTrace struct {
	Code   hexutil.Bytes `json:"code"`
	Ops    []Ops         `json:"ops"`
	parent *VMTrace
}

type Ops struct {
	Op          string
	pushcount   int
	orientation int
	Cost        uint64   `json:"cost"`
	Ex          Ex       `json:"ex"`
	PC          uint64   `json:"pc"`
	Sub         *VMTrace `json:"sub"`
}

type Ex struct {
	Mem   *Mem           `json:"mem"`
	Push  []*uint256.Int `json:"push"`
	Store *Store         `json:"store"`
	Used  uint64         `json:"used"`
}

type Mem struct {
	Data interface{} `json:"data"`
	Off  uint64      `json:"off"`
}

type Store struct {
	Key   *uint256.Int `json:"key"`
	Value *uint256.Int `json:"val"`
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
	tr := TracerService{}
	err = client.Call(&tr, "debug_traceTransaction", txHash, map[string]string{"tracer": "plugethVMTracer"})
	trace := make([]string, 0)
	result := OuterResult{
		Output:    tr.Output,
		StateDiff: nil,
		Trace:     trace,
		VMTrace:   tr.CurrentTrace,
	}
	return result, nil
}

//Note: If transactions is a contract deployment then the input is the 'code' that we are trying to capture with getCode

type TracerService struct {
	StateDB      core.StateDB
	CurrentTrace *VMTrace
	Output       hexutil.Bytes
	Mem          Mem
	Store        Store
}

func (r *TracerService) CaptureStart(from core.Address, to core.Address, create bool, input []byte, gas uint64, value *big.Int) {
	r.CurrentTrace = &VMTrace{Code: r.StateDB.GetCode(to), Ops: []Ops{}}
}
func (r *TracerService) CaptureState(pc uint64, op core.OpCode, gas, cost uint64, scope core.ScopeContext, rData []byte, depth int, err error) {
	count := 0
	direction := 0
	var mem *Mem
	var str *Store
	if size := len(r.CurrentTrace.Ops); size > 0 {
		switch r.CurrentTrace.Ops[size-1].orientation {
		case 0:
			r.CurrentTrace.Ops[size-1].Ex.Push = make([]*uint256.Int, r.CurrentTrace.Ops[size-1].pushcount)
			for i := 0; i < r.CurrentTrace.Ops[size-1].pushcount; i++ {
				r.CurrentTrace.Ops[size-1].Ex.Push[i] = scope.Stack().Back(i).Clone()
			}
		case 1:
			r.CurrentTrace.Ops[size-1].Ex.Push = make([]*uint256.Int, r.CurrentTrace.Ops[size-1].pushcount)
			for i := 0; i < r.CurrentTrace.Ops[size-1].pushcount; i++ {
				r.CurrentTrace.Ops[size-1].Ex.Push[(len(r.CurrentTrace.Ops[size-1].Ex.Push)-1)-i] = scope.Stack().Back(i).Clone()
			}
		}
	}
	pushCode := restricted.OpCode(op).String()
	switch pushCode {
	case "PUSH1", "PUSH2", "PUSH3", "PUSH4", "PUSH5", "PUSH6", "PUSH7", "PUSH8", "PUSH9", "PUSH10", "PUSH11", "PUSH12", "PUSH13", "PUSH14", "PUSH15", "PUSH16", "PUSH17", "PUSH18", "PUSH19", "PUSH20", "PUSH21", "PUSH22", "PUSH23", "PUSH24", "PUSH25", "PUSH26", "PUSH27", "PUSH28", "PUSH29", "PUSH30", "PUSH31", "PUSH32":
		count = 1
	case "SIGNEXTEND", "ISZERO", "CALLDATASIZE", "STATICCALL", "CALLVALUE", "MLOAD", "EQ", "ADDRESS", "DELEGATECALL", "CALLDATALOAD", "ADD", "LT", "SHR", "GT", "SLOAD", "SHL", "AND", "SUB", "EXTCODESIZE", "GAS", "SLT", "CALLER", "SHA3", "CALL", "RETURNDATASIZE", "NOT", "MUL", "OR", "DIV", "EXP", "BYTE", "TIMESTAMP", "SELFBALANCE":
		count = 1
	case "DUP1", "DUP2", "DUP3", "DUP4", "DUP5", "DUP6", "DUP7", "DUP8", "DUP9", "DUP10", "DUP11", "DUP12", "DUP13", "DUP14", "DUP15", "DUP16":
		x, _ := strconv.Atoi(pushCode[3:len(pushCode)])
		count = x + 1
		direction = 1
	case "SWAP1", "SWAP2", "SWAP3", "SWAP4", "SWAP5", "SWAP6", "SWAP7", "SWAP8", "SWAP9", "SWAP10", "SWAP11", "SWAP12", "SWAP13", "SWAP14", "SWAP15", "SWAP16":
		x, _ := strconv.Atoi(pushCode[4:len(pushCode)])
		count = x + 1
		direction = 1
	}
	memCode := restricted.OpCode(op).String()
	switch memCode {
	case "STATICCALL", "CALL":
		mem = &Mem{
			Off: scope.Stack().Back(4).Uint64(),
		}
	case "MSTORE", "MSTORE8":
		mem = &Mem{
			Data: core.BytesToHash(scope.Stack().Back(1).Bytes()),
			Off:  scope.Stack().Back(0).Uint64(),
		}
	case "MLOAD", "RETURNDATACOPY":
		mem = &Mem{
			Data: core.BytesToHash(scope.Memory().GetCopy(int64(scope.Stack().Back(0).Uint64()), 32)),
			Off:  scope.Stack().Back(0).Uint64(),
		}
	case "CALLDATACOPY":
		mem = &Mem{
			Data: scope.Memory().GetCopy(int64(scope.Stack().Back(0).Uint64()), 32),
			Off:  scope.Stack().Back(0).Uint64(),
		}
	}
	storeCode := restricted.OpCode(op).String()
	switch storeCode {
	case "SSTORE":
		str = &Store{
			Key:   scope.Stack().Back(0).Clone(),
			Value: scope.Stack().Back(1).Clone(),
		}
	}
	ops := Ops{
		orientation: direction,
		pushcount:   count,
		Op:          restricted.OpCode(op).String(),
		Cost:        cost,
		Ex: Ex{Mem: mem,
			Push:  make([]*uint256.Int, 0),
			Store: str,
			Used:  gas - cost},
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
	// currentOp := r.CurrentTrace.Ops[len(r.CurrentTrace.Ops)-1].Op
	// if currentOp == "REVERT" || currentOp == "RETURN" {
	// 	r.CurrentTrace.Ops[len(r.CurrentTrace.Ops)-1].Ex.Push = make([]*uint256.Int, 0)
	// }
	r.CurrentTrace = r.CurrentTrace.parent
	lastOpUsed := r.CurrentTrace.Ops[len(r.CurrentTrace.Ops)-2].Ex.Used
	switch r.CurrentTrace.Ops[len(r.CurrentTrace.Ops)-1].Op {
	case "DELEGATECALL":
		r.CurrentTrace.Ops[len(r.CurrentTrace.Ops)-1].Ex.Used = lastOpUsed - (gasUsed + 2600)
	case "CALL", "STATICCALL":
		r.CurrentTrace.Ops[len(r.CurrentTrace.Ops)-1].Ex.Used = lastOpUsed - (gasUsed + 100)
		r.CurrentTrace.Ops[len(r.CurrentTrace.Ops)-1].Ex.Mem.Data = hexutil.Bytes(output)
		if err != nil {
			r.CurrentTrace.Ops[len(r.CurrentTrace.Ops)-1].Ex.Mem = nil
		}
	}
	r.Output = output
}
func (r *TracerService) Result() (interface{}, error) {
	return r, nil
}

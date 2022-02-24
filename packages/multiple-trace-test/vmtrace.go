package main

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/holiman/uint256"
	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted"
	"github.com/openrelayxyz/plugeth-utils/restricted/hexutil"
)


// type OuterResult struct {
// 	Output    hexutil.Bytes `json:"output"`
// 	StateDiff interface{}   `json:"stateDiff"`
// 	Trace     []string      `json:"trace"`
// 	VMTrace   *string       `json:"vmTrace"`
// }

type VMTrace struct {
	Code            hexutil.Bytes `json:"code"`
	Ops             []Ops         `json:"ops"`
	parent          *VMTrace
	lastReturnValue []byte
}

type Ops struct {
	pushcount   int
	orientation int
	warmAccess  bool
	Cost        uint64   `json:"cost"`
	Ex          Ex       `json:"ex"`
	PC          uint64   `json:"pc"`
	Sub         *VMTrace `json:"sub"`
	Op          string   `json:"Op,omitempty"`
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

// var Tracers = map[string]func(core.StateDB) core.TracerResult{
// 	"plugethVMTracer": func(sdb core.StateDB) core.TracerResult {
// 		return &TracerService{StateDB: sdb}
// 	},
// }

func (vm *ParityTrace) VMTraceVarient(ctx context.Context, txHash core.Hash) (interface{}, error) {
	client, err := vm.stack.Attach()
	if err != nil {
		return nil, err
	}
	tr := VMTracerService{}
	err = client.Call(&tr, "debug_traceTransaction", txHash, map[string]string{"tracer": "plugethVMTracer"})
	// output := string(tr.Output)
	result := tr.CurrentTrace
	return result, nil
}

//Note: If transactions is a contract deployment then the input is the 'code' that we are trying to capture with getCode

func getData(data []byte, start uint64, size uint64) []byte {
	length := uint64(len(data))
	if start > length {
		start = length
	}
	end := start + size
	if end > length {
		end = length
	}
	d := make([]byte, size)
	copy(d, data[start:end])
	return d
}

type VMTracerService struct {

	StateDB      core.StateDB
	CurrentTrace *VMTrace
	Output       hexutil.Bytes
	Mem          Mem
	Store        Store
	warmAccess   bool
}

func (r *VMTracerService) CaptureStart(from core.Address, to core.Address, create bool, input []byte, gas uint64, value *big.Int) {
	r.CurrentTrace = &VMTrace{Code: r.StateDB.GetCode(to), Ops: []Ops{}}
}
func (r *VMTracerService) CaptureState(pc uint64, op core.OpCode, gas, cost uint64, scope core.ScopeContext, Data []byte, depth int, err error) {
	warm := false
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
	case "ADD", "ADDMOD", "ADDRESS", "AND", "BYTE", "CALL", "CALLDATALOAD", "CALLDATASIZE", "CALLER", "CALLVALUE", "CHAINID", "CREATE", "CREATE2", "DELEGATECALL", "DIV", "EQ", "EXP", "EXTCODEHASH", "EXTCODESIZE", "GAS", "GASPRICE", "GT", "ISZERO", "LT", "MLOAD", "MOD", "MSIZE", "MUL", "MULMOD", "NOT", "NUMBER", "OR", "RETURNDATASIZE", "SDIV", "SELFBALANCE", "SGT", "SHA3", "SHL", "SHR", "SIGNEXTEND", "SLOAD", "SLT", "STATICCALL", "SUB", "TIMESTAMP":
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
	case "CALL":
		mem = &Mem{
			Off: scope.Stack().Back(5).Uint64(),
		}
	case "STATICCALL":
		mem = &Mem{
			Off: scope.Stack().Back(4).Uint64(),
		}
	case "MSTORE":
		mem = &Mem{
			Data: core.BytesToHash(scope.Stack().Back(1).Bytes()),
			Off:  scope.Stack().Back(0).Uint64(),
		}
	case "MSTORE8":
		mem = &Mem{
			Data: fmt.Sprintf("%#x", scope.Stack().Back(1).Uint64()),
			Off:  scope.Stack().Back(0).Uint64(),
		}
	case "MLOAD":
		mem = &Mem{
			Data: core.BytesToHash(scope.Memory().GetCopy(int64(scope.Stack().Back(0).Uint64()), 32)),
			Off:  scope.Stack().Back(0).Uint64(),
		}

	case "RETURNDATACOPY":
		var (
			memOffset  = scope.Stack().Back(0)
			codeOffset = scope.Stack().Back(1)
			length     = scope.Stack().Back(2)
		)
		uint64CodeOffset, overflow := codeOffset.Uint64WithOverflow()
		if overflow {
			uint64CodeOffset = 0xffffffffffffffff
		}
		if len(hexutil.Bytes(getData(r.CurrentTrace.lastReturnValue, uint64CodeOffset, length.Uint64()))) == 0 {
			mem = nil
		} else {
			mem = &Mem{
				Data: hexutil.Bytes(getData(r.CurrentTrace.lastReturnValue, uint64CodeOffset, length.Uint64())),
				Off:  memOffset.Uint64(),
			}
		}
	case "CODECOPY":
		var (
			memOffset  = scope.Stack().Back(0)
			codeOffset = scope.Stack().Back(1)
			length     = scope.Stack().Back(2)
		)
		uint64CodeOffset, overflow := codeOffset.Uint64WithOverflow()
		if overflow {
			uint64CodeOffset = 0xffffffffffffffff
		}
		code := scope.Contract().Code()
		mem = &Mem{
			Data: hexutil.Bytes(getData(code, uint64CodeOffset, length.Uint64())),
			Off:  memOffset.Uint64(),
		}
	case "CALLDATACOPY":
		var (
			memOffset  = scope.Stack().Back(0)
			codeOffset = scope.Stack().Back(1)
			length     = scope.Stack().Back(2)
		)
		uint64DataOffset, overflow := codeOffset.Uint64WithOverflow()
		if overflow {
			uint64DataOffset = 0xffffffffffffffff
		}
		data := scope.Contract().Input()
		if len(hexutil.Bytes(getData(data, uint64DataOffset, length.Uint64()))) == 0 {
			mem = nil
		}
		mem = &Mem{
			Data: hexutil.Bytes(getData(data, uint64DataOffset, length.Uint64())),
			Off:  memOffset.Uint64(),
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
		warmAccess:  warm,
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
func (r *VMTracerService) CaptureFault(pc uint64, op core.OpCode, gas, cost uint64, scope core.ScopeContext, depth int, err error) {
}
func (r *VMTracerService) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) {
	r.Output = output
}
func (r *VMTracerService) CaptureEnter(typ core.OpCode, from core.Address, to core.Address, input []byte, gas uint64, value *big.Int) {
	// if restricted.OpCode(type).String() == "CALLDATACOPY" {
	//
	// }
	trace := &VMTrace{Code: r.StateDB.GetCode(to), Ops: []Ops{}, parent: r.CurrentTrace}
	r.CurrentTrace.Ops[len(r.CurrentTrace.Ops)-1].Sub = trace
	r.CurrentTrace = trace
}
func (r *VMTracerService) CaptureExit(output []byte, gasUsed uint64, err error) {
	r.CurrentTrace = r.CurrentTrace.parent
	r.CurrentTrace.lastReturnValue = output
	lastOpUsed := r.CurrentTrace.Ops[len(r.CurrentTrace.Ops)-2].Ex.Used
	switch r.CurrentTrace.Ops[len(r.CurrentTrace.Ops)-1].Op {
	case "DELEGATECALL":
		r.CurrentTrace.Ops[len(r.CurrentTrace.Ops)-1].Ex.Used = lastOpUsed - (gasUsed + 100)
	case "CALL", "STATICCALL":
		r.CurrentTrace.Ops[len(r.CurrentTrace.Ops)-1].Ex.Used = lastOpUsed - (gasUsed + 100)
		r.CurrentTrace.Ops[len(r.CurrentTrace.Ops)-1].Ex.Mem.Data = hexutil.Bytes(output)
		if err != nil || len(output) == 0 {
			// if err != nil || len(output) == 0 || core.BytesToHash(output) == core.HexToHash("0x1") {
			r.CurrentTrace.Ops[len(r.CurrentTrace.Ops)-1].Ex.Mem = nil
		}
	}
	// r.Output = output
}
func (r *VMTracerService) Result() (interface{}, error) {
	return r, nil
}

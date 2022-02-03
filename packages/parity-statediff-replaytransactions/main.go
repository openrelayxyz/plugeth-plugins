package main

import (
	"context"
	"math/big"
	"time"

	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted"
	"github.com/openrelayxyz/plugeth-utils/restricted/hexutil"
	"gopkg.in/urfave/cli.v1"
)

type OuterResult struct {
	Output    hexutil.Bytes          `json:"output"`
	StateDiff map[string]interface{} `json:"stateDiff"`
	Trace     []string               `json:"trace"`
	VMTrace   *string                `json:"vmTrace"`
}

type Balance struct {
	Balance interface{} `json:"balance"`
}

type Code struct {
	Code interface{} `json:"Code"`
}

type Nonce struct {
	Nonce interface{} `json:"nonce"`
}

type Storage struct {
	Storage interface{} `json:"storage"`
}

type Star struct {
	Interior Interior `json:"*"`
}

type Interior struct {
	From interface{} `json:"from"`
	To   interface{} `json:"to"`
}

// type StateDiff struct {
// 	Balance *big.Int `json:"balance"`
// 	Nonce   uint64   `json:"nonce"`
// 	Code    []byte   `json:"code"`
// }

type ParityStateDiffTrace struct {
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
		log.Info("Loaded Open Ethereum stateDiff plugin")
	}
}

func GetAPIs(stack core.Node, backend core.Backend) []core.API {
	return []core.API{
		{
			Namespace: "trace",
			Version:   "1.0",
			Service:   &ParityStateDiffTrace{backend, stack},
			Public:    true,
		},
	}
}

var Tracers = map[string]func(core.StateDB, core.BlockContext) core.TracerResult{
	"plugethStateDiffTracer": func(sdb core.StateDB, bctx core.BlockContext) core.TracerResult {
		return &TracerService{stateDB: sdb, blockContext: bctx}
	},
}

func (sd *ParityStateDiffTrace) ReplayTransaction(ctx context.Context, txHash core.Hash, tracer []string) (interface{}, error) {
	client, err := sd.stack.Attach()
	if err != nil {
		return nil, err
	}
	tr := TracerService{}
	err = client.Call(&tr, "debug_traceTransaction", txHash, map[string]string{"tracer": "plugethStateDiffTracer"})

	trace := make([]string, 0)
	result := OuterResult{
		Output:    tr.Output,
		StateDiff: tr.ReturnObj,
		Trace:     trace,
		VMTrace:   nil,
	}

	// result := make(map[string]interface{})
	// result[tr.Miner.String()] = Balance{Balance: Star{Interior{From: tr.MinerStartBalance, To: tr.MinerReturnBalance}}}
	// result[tr.To.String()] = Balance{Balance: Star{Interior{From: tr.ToStartBalance, To: tr.ToReturnBalance}}}
	// result[tr.From.String()] = Balance{Balance: Star{Interior{From: tr.FromStartBalance, To: tr.FromReturnBalance}}}

	// type lists struct {
	// 	To   []core.Address
	// 	From []core.Address
	// }
	// result := lists{To: tr.EnterTo, From: tr.EnterFrom}
	return result, err
}

type TracerService struct {
	stateDB      core.StateDB
	blockContext core.BlockContext
	Output       hexutil.Bytes
	// StateDiff          InnerResult
	// Miner              core.Address
	// To                 core.Address
	// From               core.Address
	// EnterTo            core.Address
	// EnterFrom          core.Address
	// ToStartBalance     string
	// ToEndBalance       string
	// ToReturnBalance    string
	// FromStartBalance   string
	// FromEndBalance     string
	// FromReturnBalance  string
	// MinerStartBalance  string
	// MinerEndBalance    string
	// MinerReturnBalance string
	// Value              string
	// Difficulty         string
	// BaseFee            string
	// EnterToList        []core.Address
	// EnterFromList      []core.Address
	// StorageToList      []core.Address
	// StorageFromList    []core.Address
	// StoreAddr          []string
	// StoreVal           []core.Hash
	// StoreKey           []core.Hash
	// StoreCode          []core.Hash
	ReturnObj map[string]interface{}
	Inner     map[string]interface{}
	// Ops                []string
	// Types              []string
	// EnterToBalance     []string
	// EnterFromBalance   []string
	// StartToNonce      []uint64
	// EnterToNonce      []uint64
	// StartFromNonce    []uint64
	// EnterFromNonce    []uint64
	// Code    []byte
}

func (r *TracerService) CaptureStart(from core.Address, to core.Address, create bool, input []byte, gas uint64, value *big.Int) {
	// toResult := make(map[string]map[string]string)
	// minerResult := make(map[string]map[string]string)
	// fromResult := make(map[string]map[string]string)
	r.ReturnObj = make(map[string]interface{})
	r.Inner = make(map[string]interface{})

	// r.EnterToList = []core.Address{}
	// r.EnterFromList = []core.Address{}
	// r.StorageToList = []core.Address{}
	// r.StorageFromList = []core.Address{}
	// r.StoreAddr = []string{}
	// r.StoreVal = []core.Hash{}
	// r.StoreKey = []core.Hash{}
	// r.StoreCode = []core.Hash{}

	// r.Ops = []string{}
	// r.Types = []string{}
	// r.EnterToBalance = []string{}
	// r.EnterFromBalance = []string{}
	// r.StartToNonce = []uint64{}
	// r.StartFromNonce = []uint64{}
	// r.EnterToNonce = []uint64{}
	// r.EnterFromNonce = []uint64{}
	// r.StartFrom = append(r.StartFrom, from)
	// r.StartTo = append(r.StartTo, to)
	// r.StartToNonce = append(r.StartToNonce, r.stateDB.GetNonce(to))
	// r.StartFromNonce = append(r.StartFromNonce, r.stateDB.GetNonce(from))
	// r.Value = hexutil.EncodeBig(value)
	// r.Difficulty = hexutil.EncodeBig(r.blockContext.Difficulty)
	// r.BaseFee = hexutil.EncodeBig(r.blockContext.BaseFee)
	// r.Miner = r.blockContext.Coinbase
	// r.MinerStartBalance = hexutil.EncodeBig(r.stateDB.GetBalance(r.Miner))
	// r.To = to
	// r.From = from
	// r.ToStartBalance = hexutil.EncodeBig(r.stateDB.GetBalance(to))
	// r.FromStartBalance = hexutil.EncodeBig(r.stateDB.GetBalance(from))
	// r.StateDiff = make(map[string]*StateDiff)
	// if _, ok := r.StateDiff[to.String()]; !ok {
	// 	r.StateDiff[to.String()] = &StateDiff{
	// 		Balance: r.stateDB.GetBalance(to),
	// 		Nonce:   r.stateDB.GetNonce(to),
	// 		Code:    r.stateDB.GetCode(to),
	// 	}
	// }
	// if _, ok := r.StateDiff[from.String()]; !ok {
	// 	r.StateDiff[from.String()] = &StateDiff{
	// 		Balance: r.stateDB.GetBalance(from),
	// 		Nonce:   r.stateDB.GetNonce(from),
	// 		Code:    r.stateDB.GetCode(from),
	// 	}
	// }
}
func (r *TracerService) CaptureState(pc uint64, op core.OpCode, gas, cost uint64, scope core.ScopeContext, rData []byte, depth int, err error) {
	// r.Ops = append(r.Ops, restricted.OpCode(op).String())
	opCode := restricted.OpCode(op).String()
	switch opCode {
	case "SSTORE":
		popVal := scope.Stack().Back(0).Bytes()
		storageFrom := r.stateDB.GetState(scope.Contract().Address(), core.BytesToHash(popVal))
		storageTo := core.BytesToHash(scope.Stack().Back(1).Bytes())
		storageHash := core.BytesToHash(popVal).String()
		//storageAddr := core.HexToAddress(hexutil.EncodeUint64(scope.Stack().Back(0).Uint64()))
		addr := scope.Contract().Address().String()
		if storageTo != storageFrom {
			//if storageTo != storageFrom && r.stateDB.Empty(storageAddr) == false && r.stateDB.Exist(storageAddr) == true {
			r.Inner[storageHash] = Star{Interior{From: storageFrom, To: storageTo}}
			r.ReturnObj[addr] = Storage{Storage: r.Inner}
		}
	}
	// r.StorageToList = append(r.StorageToList, r.EnterTo)
	// r.StorageFromList = append(r.StorageFromList, r.EnterFrom)
	// r.StoreAddr = append(r.StoreAddr, scope.Contract().Address().String())
	// r.StoreVal = append(r.StoreVal, core.BytesToHash(scope.Stack().Back(1).Bytes()))
	// r.StoreKey = append(r.StoreKey, r.stateDB.GetState(scope.Contract().Address(), core.BytesToHash(scope.Stack().Back(0).Bytes())))
	// r.StoreCode = append(r.StoreCode, core.BytesToHash(scope.Stack().Back(0).Bytes()))
}
func (r *TracerService) CaptureFault(pc uint64, op core.OpCode, gas, cost uint64, scope core.ScopeContext, depth int, err error) {
}
func (r *TracerService) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) {
	// r.ToEndBalance = hexutil.EncodeBig(r.stateDB.GetBalance(r.To))
	// r.FromEndBalance = hexutil.EncodeBig(r.stateDB.GetBalance(r.From))
	// r.MinerEndBalance = hexutil.EncodeBig(r.stateDB.GetBalance(r.Miner))
	r.Output = output
}
func (r *TracerService) CaptureEnter(typ core.OpCode, from core.Address, to core.Address, input []byte, gas uint64, value *big.Int) {
	// r.Types = append(r.Types, restricted.OpCode(typ).String())
	// r.EnterTo = to
	// r.EnterFrom = from
	// r.EnterFromList = append(r.EnterFromList, from)
	// r.EnterToList = append(r.EnterToList, to)
	// r.EnterToBalance = append(r.EnterToBalance, hexutil.EncodeBig(r.stateDB.GetBalance(to)))
	// r.EnterFromBalance = append(r.EnterFromBalance, hexutil.EncodeBig(r.stateDB.GetBalance(from)))
	// r.EnterToNonce = append(r.EnterToNonce, r.stateDB.GetNonce(to))
	// r.EnterFromNonce = append(r.EnterFromNonce, r.stateDB.GetNonce(from))
	// if _, ok := r.StateDiff[to.String()]; !ok {
	// 	r.StateDiff[to.String()] = &StateDiff{
	// 		Balance: r.stateDB.GetBalance(to),
	// 		Nonce:   r.stateDB.GetNonce(to),
	// 		Code:    r.stateDB.GetCode(to),
	// 	}
	// }
	// if _, ok := r.StateDiff[from.String()]; !ok {
	// 	r.StateDiff[from.String()] = &StateDiff{
	// 		Balance: r.stateDB.GetBalance(from),
	// 		Nonce:   r.stateDB.GetNonce(from),
	// 		Code:    r.stateDB.GetCode(from),
	// 	}
	// }
}
func (r *TracerService) CaptureExit(output []byte, gasUsed uint64, err error) {
}
func (r *TracerService) Result() (interface{}, error) {
	// r.ToReturnBalance = hexutil.EncodeBig(r.stateDB.GetBalance(r.To))
	// r.FromReturnBalance = hexutil.EncodeBig(r.stateDB.GetBalance(r.From))
	// r.MinerReturnBalance = hexutil.EncodeBig(r.stateDB.GetBalance(r.Miner))
	return r, nil
}

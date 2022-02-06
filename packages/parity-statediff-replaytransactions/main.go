package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted"
	"github.com/openrelayxyz/plugeth-utils/restricted/hexutil"
	"gopkg.in/urfave/cli.v1"
)

type OuterResult struct {
	Output    hexutil.Bytes `json:"output"`
	StateDiff interface{}   `json:"stateDiff"`
	Trace     []string      `json:"trace"`
	VMTrace   *string       `json:"vmTrace"`
}

type Code struct {
	Code interface{} `json:"Code"`
}

// type Nonce struct {
// 	Nonce interface{} `json:"nonce"`
// }

type LayerTwo struct {
	Balance *Star            `json:"balance"`
	Nonce   *Star            `json:"nonce"`
	Storage map[string]*Star `json:"storage"`
}

type Star struct {
	Interior Interior
}

func (s *Star) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte(`"="`), nil
	}
	if s.Interior.From == s.Interior.To {
		return []byte(`"="`), nil
	}
	interior, err := json.Marshal(s.Interior)
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf(`{"*":%v}`, string(interior))), nil
}

func (s *Star) UnmarshalJSON(input []byte) error {
	if string(input) == `"="` {
		return nil
	}
	x := struct {
		Interior Interior `json:"*"`
	}{}
	if err := json.Unmarshal(input, &x); err != nil {
		return err
	}
	s.Interior = x.Interior
	return nil
}

type Interior struct {
	From string `json:"from"`
	To   string `json:"to"`
}

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

	// trace := make([]string, 0)
	// result := OuterResult{
	// 	Output:    tr.Output,
	// 	StateDiff: tr.ReturnObj,
	// 	Trace:     trace,
	// 	VMTrace:   nil,
	// }

	result := tr.NonceCount
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
	ParityMiner  core.Address
	Miner        core.Address
	To           core.Address
	From         core.Address
	ReturnObj    map[string]LayerTwo
	Inner        map[string]interface{}
	NonceList    []string
	NonceCount   map[string]uint64
	Count        uint64
	PMBalance    *big.Int
}

func (r *TracerService) CaptureStart(from core.Address, to core.Address, create bool, input []byte, gas uint64, value *big.Int) {
	r.NonceCount = make(map[string]uint64)
	// count := uint64(0)
	r.NonceList = []string{}
	r.ReturnObj = make(map[string]LayerTwo)
	r.Inner = make(map[string]interface{})
	r.Count = 0
	r.ParityMiner = core.HexToAddress("0x0000000000000000000000000000000000000000")
	r.PMBalance = r.stateDB.GetBalance(r.ParityMiner)
	// r.NonceCount[r.To.String()] = r.Count
	r.NonceCount[r.From.String()] = r.Count
	r.To = to
	r.From = from
	r.Miner = r.blockContext.Coinbase
	// r.NonceList = append(r.NonceList, hexutil.EncodeUint64(r.stateDB.GetNonce(r.To)), hexutil.EncodeUint64(r.stateDB.GetNonce(r.From)))
	if _, ok := r.ReturnObj[r.To.String()]; !ok {
		r.ReturnObj[r.To.String()] = LayerTwo{Storage: make(map[string]*Star), Balance: &Star{Interior{From: hexutil.EncodeBig(r.stateDB.GetBalance(to))}}}
	}

	if _, ok := r.ReturnObj[r.From.String()]; !ok {
		r.ReturnObj[r.From.String()] = LayerTwo{Storage: make(map[string]*Star), Balance: &Star{Interior{From: hexutil.EncodeBig(r.stateDB.GetBalance(from))}}, Nonce: &Star{Interior{To: hexutil.EncodeUint64(r.stateDB.GetNonce(r.From))}}}
	}

	if _, ok := r.ReturnObj[r.Miner.String()]; !ok {
		r.ReturnObj[r.Miner.String()] = LayerTwo{Storage: make(map[string]*Star), Balance: &Star{Interior{From: hexutil.EncodeBig(r.stateDB.GetBalance(r.Miner))}}}
	}

}
func (r *TracerService) CaptureState(pc uint64, op core.OpCode, gas, cost uint64, scope core.ScopeContext, rData []byte, depth int, err error) {
	// r.NonceList = append(r.NonceList, hexutil.EncodeUint64(r.stateDB.GetNonce(r.To)), hexutil.EncodeUint64(r.stateDB.GetNonce(r.From)))
	// r.Ops = append(r.Ops, restricted.OpCode(op).String())
	opCode := restricted.OpCode(op).String()
	switch opCode {
	case "SSTORE":
		popVal := scope.Stack().Back(0).Bytes()
		storageFrom := r.stateDB.GetState(scope.Contract().Address(), core.BytesToHash(popVal)).String()
		storageTo := core.BytesToHash(scope.Stack().Back(1).Bytes()).String()
		storageHash := core.BytesToHash(popVal).String()
		//storageAddr := core.HexToAddress(hexutil.EncodeUint64(scope.Stack().Back(0).Uint64()))
		addr := scope.Contract().Address().String()
		if storageTo != storageFrom {

			//r.Inner[storageHash] = Star{Interior{From: storageFrom, To: storageTo}
			if storage, ok := r.ReturnObj[addr].Storage[storageHash]; ok {
				storage.Interior.To = storageTo
			} else {

				r.ReturnObj[addr].Storage[storageHash] = &Star{Interior{From: storageFrom, To: storageTo}}
			}
		}
	}

}
func (r *TracerService) CaptureFault(pc uint64, op core.OpCode, gas, cost uint64, scope core.ScopeContext, depth int, err error) {
}
func (r *TracerService) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) {
	// if _, ok := r.ReturnObj[r.To.String()]; !ok {
	// 	r.ReturnObj[r.To.String()] = LayerTwo{Nonce: &Star{Interior{To: hexutil.EncodeUint64(r.stateDB.GetNonce(r.To))}}}
	// }
	// r.NonceList = append(r.NonceList, hexutil.EncodeUint64(r.stateDB.GetNonce(r.To)), hexutil.EncodeUint64(r.stateDB.GetNonce(r.From)))
	r.Output = output
}
func (r *TracerService) CaptureEnter(typ core.OpCode, from core.Address, to core.Address, input []byte, gas uint64, value *big.Int) {
	// r.NonceList = append(r.NonceList, hexutil.EncodeUint64(r.stateDB.GetNonce(r.To)), hexutil.EncodeUint64(r.stateDB.GetNonce(r.From)))
	if _, ok := r.ReturnObj[to.String()]; !ok {
		r.ReturnObj[to.String()] = LayerTwo{Storage: make(map[string]*Star), Balance: &Star{Interior{From: hexutil.EncodeBig(r.stateDB.GetBalance(to))}}, Nonce: &Star{Interior{To: hexutil.EncodeUint64(r.stateDB.GetNonce(to))}}}
	}
	// r.NonceCount[to.String()] = r.Count + r.stateDB.GetNonce(to)
	r.Count = r.stateDB.GetNonce(from)
}
func (r *TracerService) CaptureExit(output []byte, gasUsed uint64, err error) {
	// r.NonceList = append(r.NonceList, hexutil.EncodeUint64(r.stateDB.GetNonce(r.To)), hexutil.EncodeUint64(r.stateDB.GetNonce(r.From)))
	// r.NonceCount[r.To.String()] = r.Count + r.stateDB.GetNonce(r.To)
	r.NonceCount[r.From.String()] = r.Count + r.stateDB.GetNonce(r.From)
	if _, ok := r.ReturnObj[r.To.String()]; !ok {
		r.ReturnObj[r.To.String()] = LayerTwo{Nonce: &Star{Interior{From: hexutil.EncodeUint64(r.stateDB.GetNonce(r.To))}}}
	}
}
func (r *TracerService) Result() (interface{}, error) {
	// r.NonceList = append(r.NonceList, hexutil.EncodeUint64(r.stateDB.GetNonce(r.To)), hexutil.EncodeUint64(r.stateDB.GetNonce(r.From)))
	// r.UltimateToBalance = hexutil.EncodeBig(r.stateDB.GetBalance(r.To))
	// r.UltimateFromBalance = hexutil.EncodeBig(r.stateDB.GetBalance(r.From))
	// r.UltimateMinerBalance = hexutil.EncodeBig(r.stateDB.GetBalance(r.Miner))
	// r.ReturnObj[r.To.String()] = LayerTwo{Storage: make(map[string]*Star), Balance: &Star{Interior{From: hexutil.EncodeBig(r.stateDB.GetBalance(to))}}}
	for addrHex, account := range r.ReturnObj {
		addr := core.HexToAddress(addrHex)
		// if addr == r.Miner{
		// 	account
		// }
		account.Balance.Interior.To = hexutil.EncodeBig(r.stateDB.GetBalance(addr))
		// account.Nonce.Interior.To = hexutil.EncodeUint64(r.stateDB.GetNonce(addr))
		// if account.Nonce == nil && account.Balance == nil && len(account.Storage) == 0 {
		// 	delete(r.ReturnObj, addrHex)
		// }
	}
	return r, nil
}

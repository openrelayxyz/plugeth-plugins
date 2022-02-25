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

)

type LayerTwo struct {
	Balance *Star            `json:"balance"`
	Code    *Star            `json:"code"`
	Nonce   *Star            `json:"nonce"`
	Storage map[string]*Star `json:"storage"`
}

type Star struct {
	Interior Interior
	New      bool
}
type Interior struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func (s *Star) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte(`"="`), nil
	}
	if s.New {
		return []byte(fmt.Sprintf(`{"+":"%v"}`, s.Interior.To)), nil
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
	// TODO: we need to distinguish between the return values if the key is star and the return values if the key is plus
	x := make(map[string]json.RawMessage)
	if err := json.Unmarshal(input, &x); err != nil {
		return err
	}
	if v, ok := x["*"]; ok {
		if err := json.Unmarshal(v, &s.Interior); err != nil {
			return err
		}
		return nil
	}
	if v, ok := x["+"]; ok {
		var y string
		if err := json.Unmarshal(v, &y); err != nil {
			return err
		}
		s.Interior.To = y
		s.New = true
		return nil
	}
	return fmt.Errorf("cannot unmarshall json")
}

// var Tracers = map[string]func(core.StateDB, core.BlockContext) core.TracerResult{
// 	"plugethStateDiffTracer": func(sdb core.StateDB, bctx core.BlockContext) core.TracerResult {
// 		return &TracerService{stateDB: sdb, blockContext: bctx}
// 	},
// }

func (sd *ParityTrace) StateDiffVariant(ctx context.Context, txHash core.Hash) (map[string]*LayerTwo, error) {
	client, err := sd.stack.Attach()
	if err != nil {
		return nil, err
	}
	tr := SDTracerService{}
	err = client.Call(&tr, "debug_traceTransaction", txHash, map[string]string{"tracer": "plugethStateDiffTracer"})

	result := tr.ReturnObj

	return result, err
}

type SDTracerService struct {
	stateDB      core.StateDB
	blockContext core.BlockContext
	Output       hexutil.Bytes
	Miner        core.Address
	To           core.Address
	From         core.Address
	ReturnObj    map[string]*LayerTwo
	ParityMiner core.Address
	MinerInitBalance *big.Int
	PMinerInitBalance *big.Int
}

func (r *SDTracerService) CapturePreStart(from core.Address, to *core.Address, input []byte, gas uint64, value *big.Int) {
	r.ReturnObj = make(map[string]*LayerTwo)
	r.Miner = r.blockContext.Coinbase
	r.ParityMiner = core.HexToAddress("0x0000000000000000000000000000000000000000")
	r.MinerInitBalance = r.stateDB.GetBalance(r.Miner)
	r.PMinerInitBalance = r.stateDB.GetBalance(r.ParityMiner)
	// r.To = to
	// r.From = from
	if to != nil {if _, ok := r.ReturnObj[to.String()]; !ok {
		r.ReturnObj[to.String()] = &LayerTwo{Storage: make(map[string]*Star), Balance: &Star{Interior{From: hexutil.EncodeBig(r.stateDB.GetBalance(*to))}, false}, Nonce: &Star{Interior{From: hexutil.EncodeUint64(r.stateDB.GetNonce(*to))}, false}, Code: &Star{Interior{From: hexutil.Encode(r.stateDB.GetCode(*to))}, false}}
	}}

	if _, ok := r.ReturnObj[from.String()]; !ok {
		r.ReturnObj[from.String()] = &LayerTwo{Storage: make(map[string]*Star), Balance: &Star{Interior{From: hexutil.EncodeBig(r.stateDB.GetBalance(from))}, false}, Nonce: &Star{Interior{From: hexutil.EncodeUint64(r.stateDB.GetNonce(from))}, false}, Code: &Star{Interior{From: hexutil.Encode(r.stateDB.GetCode(from))}, false}}
	}

	// if _, ok := r.ReturnObj[r.Miner.String()]; !ok {
	// 	r.ReturnObj[r.Miner.String()] = &LayerTwo{Storage: make(map[string]*Star), Balance: &Star{Interior{From: hexutil.EncodeBig(r.stateDB.GetBalance(r.Miner))}, false}, Nonce: &Star{Interior{From: hexutil.EncodeUint64(r.stateDB.GetNonce(r.Miner))}, false}, Code: &Star{Interior{From: hexutil.Encode(r.stateDB.GetCode(r.Miner))}, false}}
	// }
	if _, ok := r.ReturnObj[r.ParityMiner.String()]; !ok {
		r.ReturnObj[r.ParityMiner.String()] = &LayerTwo{Storage: make(map[string]*Star), Balance: &Star{Interior{From: hexutil.EncodeBig(r.PMinerInitBalance)}, false}, Nonce: &Star{Interior{From: hexutil.EncodeUint64(r.stateDB.GetNonce(r.ParityMiner))}, false}, Code: &Star{Interior{From: hexutil.Encode(r.stateDB.GetCode(r.ParityMiner))}, false}}
	}
}

func (r *SDTracerService) CaptureStart(from core.Address, to core.Address, create bool, input []byte, gas uint64, value *big.Int) {
	// r.ReturnObj = make(map[string]*LayerTwo)
	// r.Miner = r.blockContext.Coinbase
	// r.To = to
	// r.From = from
	// if _, ok := r.ReturnObj[to.String()]; !ok {
	// 	r.ReturnObj[to.String()] = &LayerTwo{Storage: make(map[string]*Star), Balance: &Star{Interior{From: hexutil.EncodeBig(new(big.Int).Add(r.stateDB.GetBalance(to), value))}, false}, Nonce: &Star{Interior{From: hexutil.EncodeUint64(r.stateDB.GetNonce(to))}, false}, Code: &Star{Interior{From: hexutil.Encode(r.stateDB.GetCode(to))}, false}}
	// 	// r.ReturnObj[r.To.String()] = LayerTwo{Storage: make(map[string]*Star), Balance: &Star{Interior{From: hexutil.EncodeBig(new(big.Int).Sub(r.stateDB.GetBalance(to), value))}}}
	// }
	//
	// if _, ok := r.ReturnObj[r.From.String()]; !ok {
	// 	r.ReturnObj[r.From.String()] = &LayerTwo{Storage: make(map[string]*Star), Balance: &Star{Interior{From: hexutil.EncodeBig(r.stateDB.GetBalance(from))}, false}, Nonce: &Star{Interior{From: hexutil.EncodeUint64(r.stateDB.GetNonce(from) - 1), To: hexutil.EncodeUint64(r.stateDB.GetNonce(from))}, false}, Code: &Star{Interior{From: hexutil.Encode(r.stateDB.GetCode(from))}, false}}
	// }
	//
	// if _, ok := r.ReturnObj[r.Miner.String()]; !ok {
	// 	r.ReturnObj[r.Miner.String()] = &LayerTwo{Storage: make(map[string]*Star), Balance: &Star{Interior{From: hexutil.EncodeBig(r.stateDB.GetBalance(r.Miner))}, false}, Nonce: &Star{Interior{From: hexutil.EncodeUint64(r.stateDB.GetNonce(r.Miner))}, false}, Code: &Star{Interior{From: hexutil.Encode(r.stateDB.GetCode(r.Miner))}, false}}
	// 	// r.ReturnObj[r.Miner.String()] = LayerTwo{Storage: make(map[string]*Star), Balance: &Star{Interior{From: hexutil.EncodeBig(r.stateDB.GetBalance(r.Miner))}}}
	// }
}
func (r *SDTracerService) CaptureState(pc uint64, op core.OpCode, gas, cost uint64, scope core.ScopeContext, rData []byte, depth int, err error) {
	opCode := restricted.OpCode(op).String()
	switch opCode {
	case "SSTORE":
		popVal := scope.Stack().Back(0).Bytes()
		storageFrom := r.stateDB.GetState(scope.Contract().Address(), core.BytesToHash(popVal)).String()
		storageTo := core.BytesToHash(scope.Stack().Back(1).Bytes()).String()
		storageHash := core.BytesToHash(popVal).String()
		addr := scope.Contract().Address().String()
		if storageTo != storageFrom {
			if storage, ok := r.ReturnObj[addr].Storage[storageHash]; ok {
				storage.Interior.To = storageTo
			} else {

				r.ReturnObj[addr].Storage[storageHash] = &Star{Interior{From: storageFrom, To: storageTo}, false}
			}
		}
	}

}
func (r *SDTracerService) CaptureFault(pc uint64, op core.OpCode, gas, cost uint64, scope core.ScopeContext, depth int, err error) {
}
func (r *SDTracerService) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) {
	r.Output = output
}
func (r *SDTracerService) CaptureEnter(typ core.OpCode, from core.Address, to core.Address, input []byte, gas uint64, value *big.Int) {
	localValue := value
	if localValue == nil {
		localValue = new(big.Int)
	}
	if _, ok := r.ReturnObj[to.String()]; !ok {
		r.ReturnObj[to.String()] = &LayerTwo{
			Storage: make(map[string]*Star),
			Balance: &Star{Interior{
				From: hexutil.EncodeBig(r.stateDB.GetBalance(to))}, false},
			Nonce: &Star{Interior{From: hexutil.EncodeUint64(r.stateDB.GetNonce(to))}, false},
			Code: &Star{Interior{From: hexutil.Encode(r.stateDB.GetCode(to))}, false}}
	}
}
func (r *SDTracerService) CaptureExit(output []byte, gasUsed uint64, err error) {
}
func (r *SDTracerService) Result() (interface{}, error) {
	 minerDiff := new(big.Int).Sub(r.stateDB.GetBalance(r.Miner), r.MinerInitBalance)

	for addrHex, account := range r.ReturnObj {
		addr := core.HexToAddress(addrHex)
		account.Balance.Interior.To = hexutil.EncodeBig(r.stateDB.GetBalance(addr))
		account.Nonce.Interior.To = hexutil.EncodeUint64(r.stateDB.GetNonce(addr))
		account.Code.Interior.To = hexutil.Encode(r.stateDB.GetCode(addr))

		if addr == r.ParityMiner {
			account.Balance.Interior.To = hexutil.EncodeBig(new(big.Int).Add(r.PMinerInitBalance, minerDiff))
		}

		if account.Nonce.Interior.To == account.Nonce.Interior.From && account.Balance.Interior.To == account.Balance.Interior.From && account.Code.Interior.To == account.Code.Interior.From && len(account.Storage) == 0 {
			delete(r.ReturnObj, addrHex)
		}

		 if account.Balance.Interior.From == "0x0" && hexutil.Encode(r.stateDB.GetCode(addr)) == "0x" && r.stateDB.GetNonce(addr) == 0 {
			account.Balance.New = true
			account.Nonce.New = true
			account.Code.New = true
		}
	}

	return r, nil
}

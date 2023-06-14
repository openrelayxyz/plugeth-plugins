package main

import (
	// "fmt"
	"context"
	"math/big"
	"time"
	"os"
	
	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted/hexutil"
	// "github.com/openrelayxyz/plugeth-utils/restricted/types"
)

// type engineService struct {
// 	backend core.Backend
// 	stack core.Node
// }


var errs chan error = make(chan error)
var hookChan chan map[string]struct{} = make(chan map[string]struct{})
var quit chan string = make(chan string)
var initialized bool

func HookTester() {

	defer txTracer()

	go func () {
		initialized = true
		for {
			select {
				// case <- time.NewTimer(15 * time.Second).C:
				case <- quit:
					if len(plugins) > 0 {
						log.Error("Exit with Error, Plugins map not empty", "Plugins not called", plugins)
						os.Exit(1)
					} else {
						log.Error("Exit without error", "len", len(plugins))
						os.Exit(0)
					}
				case m := <- hookChan:
					// log.Error("this came in off of the hookChan", "m", m)
					var ok bool
					f := func(key string) bool {_, ok = m[key]; return ok}
					switch {
						case f("OnShutdown"):
							delete(plugins, "OnShutdown")
						case f("StateUpdate"):
							delete(plugins, "StateUpdate")
						case f("PreProcessBlock"):
							delete(plugins, "PreProcessBlock")
						case f("PreProcessTransaction"):
							delete(plugins, "PreProcessTransaction")
						case f("PostProcessTransaction"):
							delete(plugins, "PostProcessTransaction")
						case f("PostProcessBlock"):
							delete(plugins, "PostProcessBlock")
						case f("NewHead"):
							delete(plugins, "NewHead")
						case f("LivePreProcessBlock"):
							delete(plugins, "LivePreProcessBlock")
						case f("LivePreProcessTransaction"):
							delete(plugins, "LivePreProcessTransaction")
						case f("LivePostProcessTransaction"):
							delete(plugins, "LivePostProcessTransaction")
						case f("LivePostProcessBlock"):
							delete(plugins, "LivePostProcessBlock")
						case f("GetRPCCalls"):
							delete(plugins, "GetRPCCalls")
						case f("SetTrieFlushIntervalClone"):
							delete(plugins, "SetTrieFlushIntervalClone")
						case f("StandardCaptureStart"):
							delete(plugins, "StandardCaptureStart")
						case f("StandardCaptureState"):
							delete(plugins, "StandardCaptureState")
						case f("StandardCaptureFault"):
							delete(plugins, "StandardCaptureFault")
						case f("StandardCaptureEnter"):
							delete(plugins, "StandardCaptureEnter")
						case f("StandardCaptureExit"):
							delete(plugins, "StandardCaptureExit")
						case f("StandardCaptureEnd"):
							delete(plugins, "StandardCaptureEnd")
						case f("StandardTracerResult"):
							delete(plugins, "StandardTracerResult")
						case f("LivePreProcessBlock"):
							delete(plugins, "LivePreProcessBlock")
						case f("LiveCaptureStart"):
							delete(plugins, "LiveCaptureStart")
						case f("LiveCaptureState"):
							delete(plugins, "LiveCaptureState")
						case f("LiveCaptureFault"):
							delete(plugins, "LiveCaptureFault")
						case f("LiveCaptureEnter"):
							delete(plugins, "LiveCaptureEnter")
						case f("LiveCaptureExit"):
							delete(plugins, "LiveCaptureExit")
						case f("LiveCaptureEnd"):
							delete(plugins, "LiveCaptureEnd")
						case f("LiveTracerResult"):
							delete(plugins, "LiveTracerResult")
						case f("PreTrieCommit"):
							delete(plugins, "PreTrieCommit")
						case f("PostTrieCommit"):
							delete(plugins, "PostTrieCommit")
				}
			}
		}
	}()
	
	blockFactory()
	// time.Sleep(2 * time.Second)
	txContracts()
	log.Error("called block factory")
	// time.Sleep(2 * time.Second)
	// txContracts()
}

type TransactionArgs struct {
	From                 *core.Address `json:"from"`
	To                   *core.Address `json:"to"`
	Gas                  *hexutil.Uint64 `json:"gas"`
	GasPrice             *hexutil.Big    `json:"gasPrice"`
	MaxFeePerGas         *hexutil.Big    `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *hexutil.Big    `json:"maxPriorityFeePerGas"`
	Value                *hexutil.Big    `json:"value"`
	Nonce 				 *hexutil.Big    `json:"nonce"`
}

var t0 core.Hash
var t1 core.Hash
var t2 core.Hash
var t3 core.Hash
var coinBase *core.Address

func blockFactory() {

	cl := apis[0].Service.(*engineService).stack
	client, err := cl.Attach()
	if err != nil {
		errs <- err
		log.Error("Error connecting with client block factory")
	}

	err = client.Call(&coinBase, "eth_coinbase")
	if err != nil {
		errs <- err
		log.Error("failed to call method eth_coinbase", "err", err)
	}
	log.Info("THIS IS THE CB", "coinbase", coinBase)

	var peerCount hexutil.Uint64
	for peerCount == 0 {
		err = client.Call(&peerCount, "net_peerCount")
		if err != nil {
			errs <- err
			log.Error("failed to call method eth_coinbase", "err", err)
		}
		time.Sleep(100 * time.Millisecond)
	} 

	v := (*hexutil.Big)(big.NewInt(1))

	unlockedAccount := core.HexToAddress("4204477bf7fce868e761caaba991ffc607717dbf")

	tx0_params := &TransactionArgs{
		From: coinBase,
		To: &unlockedAccount,
		Value: v,
	}
	
	err = client.Call(&t0, "eth_sendTransaction", tx0_params)
	if err != nil {
		log.Error("miner to miner transfer failed", "err", err)
	}
	log.Error("this is the return value for eth_sendTransaction zero", "tx0", t0)

	// var val interface{}
	// err = client.Call(&val, "plugeth_setTrieFlushInterval", "2h")
	// if err != nil {
	// 	log.Error("miner to trieFlushInterval call failed", "err", err)
	// }
	// log.Error("this is the return value for trieFlushInterval", "val", val)

	// arg0 := map[string]interface{}{
	// 	// "input": "0x60006000fd",
	// 	"input": "0x61520873000000000000000000000000000000000000000060006000600060006000f1",
	// 	"from": coinBase,
	// }

	// time.Sleep(2 * time.Second)
	// log.Error("second client call")
	// err = client.Call(&t1, "eth_sendTransaction", arg0)
	// if err != nil {
	// 	errs <- err
	// 	log.Error("failed to call method eth_sendTransaction", "err", err)
	// }
	// log.Error("this is the return value for eth_sendTransaction one", "tx1", t1)

	// arg1 := map[string]interface{}{
	// 	"input": "0x60006000fd",
	// 	// "input": "0x61520873000000000000000000000000000000000000000060006000600060006000f1",
	// 	"from": coinBase,
	// }

	// time.Sleep(2 * time.Second)
	// log.Error("third client call")
	// err = client.Call(&t2, "eth_sendTransaction", arg1)
	// if err != nil {
	// 	errs <- err
	// 	log.Error("failed to call method eth_sendTransaction", "err", err)
	// }
	// log.Error("this is the return value for eth_sendTransaction one", "tx1", t2)
	
}

func txContracts() {

	cl := apis[0].Service.(*engineService).stack
	client, err := cl.Attach()
	if err != nil {
		log.Error("Error connecting with client block factory")
	}

	arg0 := map[string]interface{}{
		"input": "0x60018080600053f3",
		// "input": "0x608060405234801561001057600080fd5b5061011a806100206000396000f3fe608060405234801561001057600080fd5b50600436106100375760003560e01c806360fe47b11461003c5780636d4ce63c1461005d57610037565b600080fd5b61004561007e565b60405161005291906100c5565b60405180910390f35b61007c6004803603602081101561007a57600080fd5b50356100c2565b6040516020018083838082843780820191505050505b565b005b6100946100c4565b60405161005291906100bf565b6100d1565b60405180910390f35b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb60e11b815260040161010060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16146101e557600080fd5b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663e7ba30df6040518163ffffffff1660e01b8152600401600060405180830381600087803b1580156101ae57600080fd5b505af11580156101c2573d6000803e3d6000fd5b50505050505050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161461029157600080fd5b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663fdacd5766040518163ffffffff1660e01b8152600401600060405180830381600087803b1580156102f957600080fd5b505af115801561030d573d6000803e3d6000fd5b50505050505050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff168156fea2646970667358221220d4f2763f3a0ae2826cc9ef37a65ff0c14d7a3aafe8d1636ff99f72e2f705413d64736f6c634300060c0033",
		// "input": "0x3859818153F3",
		// "input": "0x61520873000000000000000000000000000000000000000060006000600060006000f1",
		"from": coinBase,
	}

	time.Sleep(2 * time.Second)
	err = client.Call(&t1, "eth_sendTransaction", arg0)
	if err != nil {
		log.Error("intial contract deployment failed", "err", err)
	}
	log.Error("this is the return value for eth_sendTransaction T1", "tx1", t1)
	
	// address := core.HexToAddress("0x4fc22727bc09f881566f467207291f101f78b7d6")
	// miner := core.HexToAddress("f2c207111cb6ef761e439e56b25c7c99ac026a01")
	arg1 := map[string]interface{}{
		// "input": "0x3859818153F3",
		"input": "0x608060405234801561001057600080fd5b5061011a806100206000396000f3fe608060405234801561001057600080fd5b50600436106100375760003560e01c806360fe47b11461003c5780636d4ce63c1461005d57610037565b600080fd5b61004561007e565b60405161005291906100c5565b60405180910390f35b61007c6004803603602081101561007a57600080fd5b50356100c2565b6040516020018083838082843780820191505050505b565b005b6100946100c4565b60405161005291906100bf565b6100d1565b60405180910390f35b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb60e11b815260040161010060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16146101e557600080fd5b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663e7ba30df6040518163ffffffff1660e01b8152600401600060405180830381600087803b1580156101ae57600080fd5b505af11580156101c2573d6000803e3d6000fd5b50505050505050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161461029157600080fd5b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663fdacd5766040518163ffffffff1660e01b8152600401600060405180830381600087803b1580156102f957600080fd5b505af115801561030d573d6000803e3d6000fd5b50505050505050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff168156fea2646970667358221220d4f2763f3a0ae2826cc9ef37a65ff0c14d7a3aafe8d1636ff99f72e2f705413d64736f6c634300060c0033",
		// "input": "0x61520873000000000000000000000000000000000000000060006000600060006000f1",
		//the following address is derived from a curl to the node, should be derived from internal call. 
		// "input": fmt.Sprintf("0x62ffffff73%v60006000600060006000f1", address),
		"from": coinBase,
	}

	time.Sleep(2 * time.Second)
	err = client.Call(&t2, "eth_sendTransaction", arg1)
	if err != nil {
		log.Error("contract call failed", "err", err)
	}
	log.Error("this is the return value for eth_sendTransaction T2", "tx1", t2)

	for i := 0; i < 126; i ++ {
		time.Sleep(2 * time.Second)
		err = client.Call(&t2, "eth_sendTransaction", arg1)
		if err != nil {
			log.Error("looped contract call failed on index", "i", i, "err", err)
		}
		log.Error("this is the return value for eth_sendTransaction i", "tx", i)
	}

	// arg1 := map[string]interface{}{
	// 	// "input": "0x60006000fd",
	// 	// "input": "0x61520873000000000000000000000000000000000000000060006000600060006000f1",
	// 	"data": "0xa9059cbb0000000000000000000000001234567890abcdef0000000000000000000000000000000000000000000000000000000000000064",
	// 	"from": coinBase,
	// }

	// time.Sleep(2 * time.Second)
	// log.Error("third client call")
	// err = client.Call(&t2, "eth_sendTransaction", arg1)
	// if err != nil {
	// 	log.Error("failed to call new TXXXXXX", "err", err)
	// }
	// log.Error("this is the return value for the NEW TXXXXXX", "tx1", t2)
}

type TraceConfig struct {
	Tracer  *string
}

func txTracer() {
	cl := apis[0].Service.(*engineService).stack
	client, err := cl.Attach()
	if err != nil {
		errs <- err
		log.Error("Error connecting with client block factory")
	}

	time.Sleep(2 * time.Second)
	tr := "testTracer"
	t := TraceConfig{
		Tracer: &tr,
	}

	var trResult interface{}
	err = client.Call(&trResult, "debug_traceTransaction", t0, t)
	log.Error("tracer result", "result", trResult, "err", err, "hash", t0)

	arg0 := map[string]interface{}{
		"input": "0x60006000fd",
		"from": coinBase,
	}

	var trResult0 interface{}
	err = client.Call(&trResult0, "debug_traceCall", arg0, "latest", t)
	log.Error("tracer result", "result", trResult0, "err", err)

	arg1 := map[string]interface{}{
		"input": "0x61520873000000000000000000000000000000000000000060006000600060006000f1",
		"from": coinBase,
	}

	var trResult1 interface{}
	err = client.Call(&trResult1, "debug_traceCall", arg1, "latest", t)

	// address := core.HexToAddress("0x4fc22727bc09f881566f467207291f101f78b7d6")
	// miner := core.HexToAddress("f2c207111cb6ef761e439e56b25c7c99ac026a01")
	final := map[string]interface{}{
		// "input": "0x3859818153F3",
		// "input": "0x608060405234801561001057600080fd5b5061011a806100206000396000f3fe608060405234801561001057600080fd5b50600436106100375760003560e01c806360fe47b11461003c5780636d4ce63c1461005d57610037565b600080fd5b61004561007e565b60405161005291906100c5565b60405180910390f35b61007c6004803603602081101561007a57600080fd5b50356100c2565b6040516020018083838082843780820191505050505b565b005b6100946100c4565b60405161005291906100bf565b6100d1565b60405180910390f35b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb60e11b815260040161010060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16146101e557600080fd5b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663e7ba30df6040518163ffffffff1660e01b8152600401600060405180830381600087803b1580156101ae57600080fd5b505af11580156101c2573d6000803e3d6000fd5b50505050505050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161461029157600080fd5b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663fdacd5766040518163ffffffff1660e01b8152600401600060405180830381600087803b1580156102f957600080fd5b505af115801561030d573d6000803e3d6000fd5b50505050505050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff168156fea2646970667358221220d4f2763f3a0ae2826cc9ef37a65ff0c14d7a3aafe8d1636ff99f72e2f705413d64736f6c634300060c0033",
		"input": "0x61520873000000000000000000000000000000000000000060006000600060006000f1",
		//the following address is derived from a curl to the node, should be derived from internal call. 
		// "input": fmt.Sprintf("0x62ffffff73%v60006000600060006000f1", address),
		"from": coinBase,
	}

	time.Sleep(2 * time.Second)
	err = client.Call(&t3, "eth_sendTransaction", final)
	if err != nil {
		log.Error("contract call failed", "err", err)
	}
	log.Error("this is the return value for eth_sendTransaction Tfinal", "tx1", t2)

	quit <- "quit"

}

type innerParams struct {
	to string `json:"to"`
}

// type tracerTypeParams struct {
// 	tracer *string `json:"tracer"`
// }

// type tracerParams struct {
// 	innerParams
// 	*hexutil.Uint64 
// 	tracerTypeParams 

// }

// {"to":"0x32Be343B94f860124dC4fEe278FDCBD38C102D88"},"latest",{"tracer":"myTracer"}],"id":0


func (service *engineService) CaptureShutdown(ctx context.Context) {
	m := map[string]struct{}{
		"OnShutdown":struct{}{},
	}
	hookChan <- m
}

func (service *engineService) CapturePreTrieCommit(ctx context.Context) {
	m := map[string]struct{}{
		"PreTrieCommit":struct{}{},
	}
	hookChan <- m
}

func (service *engineService) CapturePostTrieCommit(ctx context.Context) {
	m := map[string]struct{}{
		"PostTrieCommit":struct{}{},
	}
	hookChan <- m
}


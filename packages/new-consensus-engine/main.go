package main

import (
	"context"
	"math/big"
	"time"
	// "os"
	
	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted/hexutil"
	// "github.com/openrelayxyz/plugeth-utils/restricted/types"
)

type engineService struct {
	backend core.Backend
	stack core.Node
}


var errs chan error = make(chan error)
var hookChan chan map[string]struct{} = make(chan map[string]struct{})

func HookTester() {

	defer txTracer()

	log.Error("inside of hooktester")
	
	blockFactory()


	log.Error("Pre loop map", "plugins", plugins)

	go func () {
		for {
			select {
				// case <- time.NewTimer(5 * time.Second).C:
				// 	if len(plugins) > 0 {
				// 		log.Error("Exit with Error, Plugins map not empty", "Plugins not called", plugins)
				// 		os.Exit(1)
				// 	} else {
				// 		log.Error("Exit without error", "len", len(plugins))
				// 		os.Exit(0)
				// 	}
				case m := <- hookChan:
					log.Error("this came in off of the hookChan", "m", m)
					var ok bool
					f := func(key string) bool {_, ok = m[key]; return ok}
					switch {
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
						// case f("GetRPCCalls"):
						// 	delete(plugins, "GetRPCCalls")
						// case f("SetTrieFlushIntervalClone"):
						// 	delete(plugins, "SetTrieFlushIntervalClone")
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
				}
			}
		}
	}()

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
		log.Error("failed to call method eth_sendTransaction", "err", err)
	}
	log.Error("this is the return value for eth_sendTransaction zero", "tx0", t0)

	arg0 := map[string]interface{}{
		// "input": "0x60006000fd",
		"input": "0x61520873000000000000000000000000000000000000000060006000600060006000f1",
		"from": coinBase,
	}

	time.Sleep(2 * time.Second)
	log.Error("second client call")
	err = client.Call(&t1, "eth_sendTransaction", arg0)
	if err != nil {
		errs <- err
		log.Error("failed to call method eth_sendTransaction", "err", err)
	}
	log.Error("this is the return value for eth_sendTransaction one", "tx1", t1)

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
		errs <- err
		log.Error("Error connecting with client block factory")
	}

	arg0 := map[string]interface{}{
		// "input": "0x60006000fd",
		"input": "0x61520873000000000000000000000000000000000000000060006000600060006000f1",
		"from": coinBase,
	}


	time.Sleep(2 * time.Second)
	log.Error("second client call")
	err = client.Call(&t1, "eth_sendTransaction", arg0)
	if err != nil {
		errs <- err
		log.Error("failed to call method eth_sendTransaction", "err", err)
	}
	log.Error("this is the return value for eth_sendTransaction one", "tx1", t1)

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


func (service *engineService) Test(ctx context.Context) string {
	return "this is a placeholder function"
}


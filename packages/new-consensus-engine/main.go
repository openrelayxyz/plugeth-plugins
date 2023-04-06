package main

import (
	"context"
	"math/big"
	"time"
	"os"
	"time"
	// "errors"
	
	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted/hexutil"
)

type engineService struct {
	backend core.Backend
	stack core.Node
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


var errs chan error = make(chan error)
// var errs []error
// var hookChan chan map[string]func(item interface{}) = make(chan map[string]func(item interface{}))
var hookChan chan string = make(chan string)

func HookTester() {

	log.Error("inside hook tester")

	blockFactory()

	start := time.Now()
	go func () {
		for {
			// if m := <- hookChan {
				m := <- hookChan
				log.Error("came in off of hookchan", "m", m)
				// var val interface{}
				// var ok bool
				// f := func(key string) bool {val, ok = m[key]; return ok}
				// switch {
				// case f("PreProcessBlock"):
				// 	if val == func([]byte, core.Hash, [][]byte, *big.Int) {
				// 		log.Error("pre process")
				// 	}
				// case f("PreProcessTransaction"):
				// 	if val == func(core.Hash, core.Hash) {
				// 		log.Error("pp txn")
				// 	}
				// }
			}
		
	}()

	go func () {
		var e error
		for {
			e := <- errs
			log.Error("Plugin returned error", "err", e)
			}
		if e != nil {
			os.Exit(1)
		}
	}()

// 	if len(errs) > 0 { this needs to be a channel
// 		for _, err := range errs {
// 		log.Error("Error", "err", err)
// 		}
// 	// os.Exit(1)
// 	}

	// os.Exit(0)
	
}



// func evaluate() {
// 	m := <- hookChan
// 	log.Error("eval func", "name", name)
// }

func blockFactory() {

	log.Error("inside block factory")

	cl := apis[0].Service.(*engineService).stack
	client, err := cl.Attach()
	if err != nil {
		errs <- err
		log.Error("Error connecting with client")
	}

	var coinBase *core.Address
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

	var t0 interface{}
	err = client.Call(&t0, "eth_sendTransaction", tx0_params)
	if err != nil {
		errs <- err
		log.Error("failed to call method eth_sendTransaction", "err", err)
	}
	log.Error("this is the return value for eth_sendTransaction zero", "tx0", t0)
}

// this is how to attach to the json shell: ./geth attach /tmp/geth.ipc


func (service *engineService) Test(ctx context.Context) string {
	return "this is a placeholder function"
}


package main

import (
	"context"
	"math/big"
	"time"
	
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


var errs []error

func HookTester() {

	cl := apis[0].Service.(*engineService).stack
	client, err := cl.Attach()
	if err != nil {
		errs = append(errs, err)
		log.Error("Error connecting with client")
	}

	var coinBase *core.Address
	err = client.Call(&coinBase, "eth_coinbase")
	if err != nil {
		errs = append(errs, err)
		log.Error("failed to call method eth_coinbase", "err", err)
	}
	// log.Error("this is the return value for eth_coinbase", "test", coinBase, "type", reflect.TypeOf(coinBase), "len", len(errs))

	var peerCount hexutil.Uint64
	for peerCount == 0 {
		err = client.Call(&peerCount, "net_peerCount")
		if err != nil {
			errs = append(errs, err)
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
		errs = append(errs, err)
		log.Error("failed to call method eth_sendTransaction", "err", err)
	}
	log.Error("this is the return value for eth_sendTransaction zero", "tx0", t0)

	// go evaluate()

	go func () {
		for {
			x := <- hookChan
			log.Error("channel returns", "x", x)
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


// what should be done next is set up a json object to pass into eth_sendTransactions so we can append one block to the chain from there we can start to 
// experiment on what can be done to excerise the hooks

// the as of now command to start geth is => /geth --dev --http --http.api eth --verbosity=5,  with this plugin loaded 

// this is how to attch to the json shell: ./geth attach /tmp/geth.ipc


func (service *engineService) Test(ctx context.Context) string {
	return "this is a placeholder function"
}


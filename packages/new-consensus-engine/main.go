package main

import (
	"context"
	"math/big"
	
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

	acctParams := "password"
	// it is our understanding that the personal namespace will be depricated blah blah
	var newAccount *core.Address
	err = client.Call(&newAccount, "personal_newAccount", acctParams)
	if err != nil {
		errs = append(errs, err)
		log.Error("failed to call method personal_newAccount", "err", err)
	}
	// log.Error("this is the return value for personal_newAccount", "test", newAccount, "type", reflect.TypeOf(newAccount), "len", len(errs))


	v := (*hexutil.Big)(big.NewInt(1))

	tx0_params := &TransactionArgs{
		From: coinBase,
		To: newAccount,
		Value: v,
	}

	var t0 interface{}
	err = client.Call(&t0, "eth_sendTransaction", tx0_params)
	if err != nil {
		errs = append(errs, err)
		log.Error("failed to call method eth_sendTransaction", "err", err)
	}
	log.Error("this is the return value for eth_sendTransaction zero", "tx0", t0)

	// tx1_params := &TransactionArgs{
	// 	From: newAccount,
	// 	To: coinBase,
	// 	Value: v,
	// }

	// var t1 interface{}
	// err = client.Call(&t1, "eth_sendTransaction", tx1_params)
	// if err != nil {
	// 	errs = append(errs, err)
	// 	log.Error("failed to call method eth_sendTransaction", "err", err)
	// }
	// log.Error("this is the return value for eth_sendTransaction one", "tx1", t1)

	
	
	

// 	if len(errs) > 0 { this needs to be a channel
// 		for _, err := range errs {
// 		log.Error("Error", "err", err)
// 		}
// 	// os.Exit(1)
// 	}

	// os.Exit(0)
	
}


// what should be done next is set up a json object to pass into eth_sendTransactions so we can append one block to the chain from there we can start to 
// experiment on what can be done to excerise the hooks

// the as of now command to start geth is => /geth --dev --http --http.api eth --verbosity=5,  with this plugin loaded 

// this is how to attch to the json shell: ./geth attach /tmp/geth.ipc


func (service *engineService) Test(ctx context.Context) string {
	return "this is a placeholder function"
}


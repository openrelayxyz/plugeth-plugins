package main

import (
	"context"
	"math/big"
	"time"
	"os"
	
	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted/hexutil"
)

type engineService struct {
	backend core.Backend
	stack core.Node
}


var errs chan error = make(chan error)
var hookChan chan map[string]struct{} = make(chan map[string]struct{})

func HookTester() {

	blockFactory()

	log.Error("Pre loop map", "plugins", plugins)

	go func () {
		for {
			select {
				case <- time.NewTimer(5 * time.Second).C:
					if len(plugins) > 0 {
						log.Error("told you so", "len", len(plugins))
						os.Exit(1)
					} else {
						log.Error("Exit without error", "len", len(plugins))
						os.Exit(0)
					}
				case m := <- hookChan:
					var ok bool
					f := func(key string) bool {_, ok = m[key]; return ok}
					switch {
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
						case f("GetRPCCalls"):
							delete(plugins, "GetRPCCalls")
						case f("SetTrieFlushIntervalClone"):
							delete(plugins, "SetTrieFlushIntervalClone")
				}
			}
		}
	}()
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

func blockFactory() {

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


func (service *engineService) Test(ctx context.Context) string {
	return "this is a placeholder function"
}


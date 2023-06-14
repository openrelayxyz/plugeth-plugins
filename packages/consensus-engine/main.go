package main

import (
	"context"
	"math/big"
	"time"
	"os"
	
	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted/hexutil"
)

var apis []core.API

type engineService struct {
	backend core.Backend
	stack core.Node
}

func GetAPIs(stack core.Node, backend core.Backend) []core.API {
	apis = []core.API{
		{
			Namespace: "plugeth",
			Version:   "1.0",
			Service:   &engineService{backend, stack},
			Public:    true,
		},
	}
	return apis
}

var coinBase *core.Address
var tx0_hash core.Hash
var tx1_hash core.Hash

func BlockChain() {

	cl := apis[0].Service.(*engineService).stack
	client, err := cl.Attach()
	if err != nil {
		log.Error("Error connecting with BlockChain client")
	}

	err = client.Call(&coinBase, "eth_coinbase")
	if err != nil {
		log.Error("Failed to call method eth_coinbase", "err", err)
	}

	var peerCount hexutil.Uint64
	for peerCount == 0 {
		err = client.Call(&peerCount, "net_peerCount")
		if err != nil {
			log.Error("failed to call method net_peerCount", "err", err)
		}
		time.Sleep(100 * time.Millisecond)
	} 

	tx_params := map[string]interface{}{
		"from": coinBase,
		"to": coinBase,
		"value": (*hexutil.Big)(big.NewInt(1)),
	}

	
	err = client.Call(&tx0_hash, "eth_sendTransaction", tx_params)
	if err != nil {
		log.Error("Initial tx deployment failed", "err", err)
	}

	contract_params := map[string]interface{}{
		"input": "0x608060405234801561001057600080fd5b5061011a806100206000396000f3fe608060405234801561001057600080fd5b50600436106100375760003560e01c806360fe47b11461003c5780636d4ce63c1461005d57610037565b600080fd5b61004561007e565b60405161005291906100c5565b60405180910390f35b61007c6004803603602081101561007a57600080fd5b50356100c2565b6040516020018083838082843780820191505050505b565b005b6100946100c4565b60405161005291906100bf565b6100d1565b60405180910390f35b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb60e11b815260040161010060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16146101e557600080fd5b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663e7ba30df6040518163ffffffff1660e01b8152600401600060405180830381600087803b1580156101ae57600080fd5b505af11580156101c2573d6000803e3d6000fd5b50505050505050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161461029157600080fd5b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663fdacd5766040518163ffffffff1660e01b8152600401600060405180830381600087803b1580156102f957600080fd5b505af115801561030d573d6000803e3d6000fd5b50505050505050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff168156fea2646970667358221220d4f2763f3a0ae2826cc9ef37a65ff0c14d7a3aafe8d1636ff99f72e2f705413d64736f6c634300060c0033",
		"from": coinBase,
	}

	time.Sleep(2 * time.Second)
	err = client.Call(&tx1_hash, "eth_sendTransaction", contract_params)
	if err != nil {
		log.Error("Contract deployment failed", "err", err)
	}

	time.Sleep(2 * time.Second)
	var blockString string
	err = client.Call(&blockString, "eth_blockNumber")
	if err != nil {
		log.Error("blockNumber call failed", "err", err)
	}
	blockNumber, _ := hexutil.DecodeUint64(blockString)
	if blockNumber != uint64(2) {
		log.Error("blockChain malformed, want blocknumber == 2, have blocknumber", "number", blockNumber)
		os.Exit(1)
	} else {
		log.Info("blockChain is good order!")
		os.Exit(0)
	}
}

func (e *engineService) HelloEngine(context.Context) string {
	//this is a paceholder method
	return "hello user"
}


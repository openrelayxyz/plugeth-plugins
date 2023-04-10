package main

import (
	"time"
	"math/big"

	"github.com/openrelayxyz/plugeth-utils/core"

)

// cmd/geth/

var apis []core.API

func GetAPIs(stack core.Node, backend core.Backend) []core.API {
	apis = []core.API{
		{
			Namespace: "plugeth",
			Version:   "1.0",
			Service:   &engineService{backend, stack},
			Public:    true,
		},
		// {
		// 	Namespace: "plugeth",
		// 	Version:   "1.0",
		// 	Service:   &TracerResult{},
		// 	Public:    true,
		// },
		// {
		// 	Namespace: "plugeth",
		// 	Version:   "1.0",
		// 	Service:   &TracerService{},
		// 	Public:    true,
		// },
	}
	// name := "GetAPIs"
	// m := map[string]interface{}{
	// 	name:func(stack core.Node, backend core.Backend) []core.API {}
	// }
	// hookChan <- name
	return apis
}

func OnShutdown(){
	// name := "OnShutdown"
	// m := map[string]interface{}{
	// 	name: func(),
	// }
	// hookChan <- m
}

// core/


func PreProcessBlock(hash core.Hash, number uint64, encoded []byte) {
	name := "PreProcessBlock"
	m := map[string]interface{}{
		name:PreProcessBlock,
	}
	hookChan <- m
}

func PreProcessTransaction(txBytes []byte, txHash, blockHash core.Hash, i int) {
	name := "PreProcessTransaction"
	m := map[string]interface{}{
		name:PreProcessTransaction,
	}
	hookChan <- m
}

func BlockProcessingError(tx core.Hash, block core.Hash, err error) {
	// name := "BlockProcessingError"
	// name := map[string]func(item interface{}){
	// 	name:func(core.Hash, core.Hash, error)
	// }
	// hookChan <- name
}

func PostProcessTransaction(tx core.Hash, block core.Hash, i int, receipt []byte) {
	name := "PostProcessTransaction"
	m := map[string]interface{}{
		name:PostProcessTransaction,
	}
	hookChan <- m
}

func PostProcessBlock(block core.Hash) {
	name := "PostProcessBlock"
	m := map[string]interface{}{
		name:PostProcessBlock,
	}
	hookChan <- m
}

func NewHead(block []byte, hash core.Hash, logs [][]byte, td *big.Int) {
	name := "NewHead"
	m := map[string]interface{}{
		name:NewHead,
	}
	hookChan <- m
}

func NewSideBlock(block []byte, hash core.Hash, logs [][]byte) {
	// name := "NewSideBlock"
	// hookChan <- name
}

func Reorg(commonBlock core.Hash, oldChain, newChain []core.Hash) {
	// name := "Reorg"
	// hookChan <- name
}

func SetTrieFlushIntervalClone(t time.Duration) time.Duration {
	// name := "SetTrieFlushIntervalClone"
	// hookChan<- name
	return t
}

// var Interval time.Duration 

// type TrieIntervalService struct {
// }

// func (service *TrieIntervalService) SetTrieFlushInterval(ctx context.Context, interval string) error {
// 	log.Error("true flush interval", "interval", interval)
// 	return nil
// }

// core/rawdb/

func ModifyAncients(index uint64, freezerUpdate map[string]interface{}) {
	// name := "ModifyAncients"
	// hookChan<- name
}

// func AppendAncient(number uint64, hash, header, body, receipts, td []byte) {
// 	name := "AppendAncient"
// 	hookChan<- name
// }

// core/state/

// func StateUpdate(blockRoot core.Hash, parentRoot core.Hash, coreDestructs map[core.Hash]struct{}, coreAccounts map[core.Hash][]byte, coreStorage map[core.Hash]map[core.Hash][]byte, coreCode map[core.Hash][]byte) {
// 	name := "StateUpdate"
// 	hookChan<- name
// }

// core/vm we have code in core/vm but not hooks

// rpc/

func GetRPCCalls(method string, id string, params string) {
	// name := "GetRPCCalls"
	// m := map[string]interface{}{
	// 	name:GetRPCCalls,
	// }
	// hookChan <- m
}

var plugins map[string]interface{} = map[string]interface{}{
	"PreProcessBlock": PreProcessBlock,
	"PreProcessTransaction": PreProcessTransaction,
	"PostProcessTransaction": PostProcessTransaction,
	"PostProcessBlock": PostProcessBlock,
	"NewHead": NewHead,
	"GetRPCCalls": GetRPCCalls,
} 

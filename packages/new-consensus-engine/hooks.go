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
		{
			Namespace: "plugeth",
			Version:   "1.0",
			Service:   &LiveTracerResult{},
			Public:    true,
		},
	}
	// m := map[string]struct{}{
	// 	"GetAPIs":struct{}{},
	// }
	// hookChan <- m
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
	m := map[string]struct{}{
		"PreProcessBlock":struct{}{},
	}
	hookChan <- m
}

func PreProcessTransaction(txBytes []byte, txHash, blockHash core.Hash, i int) {
	m := map[string]struct{}{
		"PreProcessTransaction":struct{}{},
	}
	hookChan <- m
}

func BlockProcessingError(tx core.Hash, block core.Hash, err error) { //also may be beyond scope
	// m := map[string]struct{}{
	// 	"BlockProcessingError":struct{}{},
	// }
	// hookChan <- m
}

func PostProcessTransaction(tx core.Hash, block core.Hash, i int, receipt []byte) {
	m := map[string]struct{}{
		"PostProcessTransaction":struct{}{},
	}
	hookChan <- m
}

func PostProcessBlock(block core.Hash) {
	m := map[string]struct{}{
		"PostProcessBlock":struct{}{},
	}
	hookChan <- m
}

func NewHead(block []byte, hash core.Hash, logs [][]byte, td *big.Int) {
	m := map[string]struct{}{
		"NewHead":struct{}{},
	}
	hookChan <- m
}

func NewSideBlock(block []byte, hash core.Hash, logs [][]byte) { // beyond the scope of the test at this time
	// m := map[string]struct{}{
	// 	"NewSideBlock":struct{}{},
	// }
	// hookChan <- m
}

func Reorg(commonBlock core.Hash, oldChain, newChain []core.Hash) { // beyond the scope of the test at this time
	// m := map[string]struct{}{
	// 	"Reorg":struct{}{},
	// }
	// hookChan <- m
}

func SetTrieFlushIntervalClone(t time.Duration) time.Duration {
	// m := map[string]struct{}{
	// 	"SetTrieFlushIntervalClone":struct{}{},
	// }
	// hookChan <- m
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

func ModifyAncients(index uint64, freezerUpdate map[string]struct{}) {
	// name := "ModifyAncients"
	// hookChan<- name
}

func AppendAncient(number uint64, hash, header, body, receipts, td []byte) {
	// name := "AppendAncient"
	// hookChan<- name
}

// core/state/

func StateUpdate(blockRoot core.Hash, parentRoot core.Hash, coreDestructs map[core.Hash]struct{}, coreAccounts map[core.Hash][]byte, coreStorage map[core.Hash]map[core.Hash][]byte, coreCode map[core.Hash][]byte) {
	m := map[string]struct{}{
		"StateUpdate":struct{}{},
	}
	hookChan <- m
}

// core/vm we have code in core/vm but not hooks

// rpc/

// func GetRPCCalls(method string, id string, params string) {
// 	m := map[string]struct{}{
// 		"GetRPCCalls":struct{}{},
// 	}
// 	hookChan <- m
// }


var plugins map[string]struct{} = map[string]struct{}{
	"StateUpdate": struct{}{},
	"PreProcessBlock": struct{}{},
	"PreProcessTransaction": struct{}{},
	"PostProcessTransaction": struct{}{},
	"PostProcessBlock": struct{}{},
	"NewHead": struct{}{},
	"StandardCaptureStart": struct{}{},
	"StandardCaptureState": struct{}{},
	"StandardCaptureFault": struct{}{},
	"StandardCaptureEnter": struct{}{},
	"StandardCaptureExit": struct{}{},
	"StandardCaptureEnd": struct{}{},
	"StandardTracerResult": struct{}{},
	// "GetRPCCalls": struct{}{},
	"LivePreProcessBlock": struct{}{},
	"LivePreProcessTransaction": struct{}{},
	"LivePostProcessTransaction": struct{}{},
	"LivePostProcessBlock": struct{}{},
	"LiveCaptureStart": struct{}{},
	"LiveCaptureState": struct{}{},
	"LiveCaptureFault": struct{}{},
	"LiveCaptureEnter": struct{}{},
	"LiveCaptureExit": struct{}{},
	"LiveCaptureEnd": struct{}{},
	"LiveTracerResult": struct{}{},
} 


package main

import (
	"time"
	"math/big"
	
	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted/hexutil"
	
)

// cmd/geth/

var apis []core.API

// var hookChan chan map[string]struct{} = make(chan map[string]struct{})

type LiveTracerResult struct {
	// backend core.Backend
	// stack core.Node
	CallStack []CallStack
	Results   []CallStack
}

type engineService struct {
	backend core.Backend
	stack core.Node
}

func (e *engineService) Test() {
	log.Info("")
}

func (e *LiveTracerResult) Test() {
	log.Info("")
}


type CallStack struct {
	Type    string         `json:"type"`
	From    core.Address   `json:"from"`
	To      core.Address   `json:"to"`
	Value   *big.Int       `json:"value,omitempty"`
	Gas     hexutil.Uint64 `json:"gas"`
	GasUsed hexutil.Uint64 `json:"gasUsed"`
	Input   hexutil.Bytes  `json:"input"`
	Output  hexutil.Bytes  `json:"output"`
	Time    string         `json:"time,omitempty"`
	Calls   []CallStack    `json:"calls,omitempty"`
	Results []CallStack    `json:"results,omitempty"`
	Error   string         `json:"error,omitempty"`
}

func GetAPIs(stack core.Node, backend core.Backend) []core.API {
	// GetAPIs is covered by virtue of the plugeth_captureShutdown method functioning.
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
	return apis
}

func OnShutdown(){
	// this injection is covered by its own test in this package. See documentation for details. 
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

func BlockProcessingError(tx core.Hash, block core.Hash, err error) { 
	// this injection is covered by a stand alone test: plugeth_injection_test.go in the core/ package. 
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
	// this injection is covered by a stand alone test: plugeth_injection_test.go in the core/ package.
}

func Reorg(commonBlock core.Hash, oldChain, newChain []core.Hash) { // beyond the scope of the test at this time
	// this injection is covered by a stand alone test: plugeth_injection_test.go in the core/ package.
}

func SetTrieFlushIntervalClone(t time.Duration) time.Duration {
	m := map[string]struct{}{
		"SetTrieFlushIntervalClone":struct{}{},
	}
	log.Error("INSIDE OF TIFC TIFC TIFC TIFC")
	hookChan <- m
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
	// this injection is covered by a stand alone test: plugeth_injection_test.go in the core/rawdb package. 
}

func AppendAncient(number uint64, hash, header, body, receipts, td []byte) {
	// this injection is covered by a stand alone test: plugeth_injection_test.go in the core/rawdb package.
}

// core/state/

// func StateUpdate(blockRoot core.Hash, parentRoot core.Hash, coreDestructs map[core.Hash]struct{}, coreAccounts map[core.Hash][]byte, coreStorage map[core.Hash]map[core.Hash][]byte, coreCode map[core.Hash][]byte) {
// 	log.Warn("StatueUpdate", "blockRoot", blockRoot, "parentRoot", parentRoot, "coreDestructs", coreDestructs, "coreAccounts", coreAccounts, "coreStorage", coreStorage, "coreCode", coreCode)
// 	m := map[string]struct{}{
// 		"StateUpdate":struct{}{},
// 	}
// 	hookChan <- m
// }

// core/vm we have code in core/vm but not hooks

// rpc/

// func GetRPCCalls(method string, id string, params string) {
// 	m := map[string]struct{}{
// 		"GetRPCCalls":struct{}{},
// 	}
// 	log.Error("inside of get rpccallssss")
// 	hookChan <- m
// }


var plugins map[string]struct{} = map[string]struct{}{
	"OnShutdown": struct{}{},
	"SetTrieFlushIntervalClone":struct{}{},
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
	"GetRPCCalls": struct{}{},
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


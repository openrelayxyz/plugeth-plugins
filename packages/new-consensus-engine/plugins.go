package main

import (
	"time"
	"math/big"

	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted"

)

var (
	pl      core.PluginLoader
	backend restricted.Backend
	log     core.Logger
	events  core.Feed
)

var httpApiFlagName = "http.api"

// cmd/geth/

var hookChan chan interface{} = make(chan interface{})

func Initialize(ctx core.Context, loader core.PluginLoader, logger core.Logger) { 
	pl = loader
	events = pl.GetFeed()
	log = logger
	v := ctx.String(httpApiFlagName)
	if v != "" {
		ctx.Set(httpApiFlagName, v+",plugeth")
	} else {
		ctx.Set(httpApiFlagName, "eth,net,web3,plugeth")
		log.Info("Loaded consensus engine plugin")
	}
}

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
			Service:   &TracerResult{},
			Public:    true,
		},
		{
			Namespace: "plugeth",
			Version:   "1.0",
			Service:   &TracerService{},
			Public:    true,
		},
	}

	return apis
}

func OnShutdown(){
	name := "OnShutdown"
	hookChan <- name
}

// core/


func PreProcessBlock(hash core.Hash, number uint64, encoded []byte) {
	name := "PreProcessBlock"
	hookChan<- name
}

func PreProcessTransaction(tx core.Hash, block core.Hash, i int) {
	name := "PreProcessTransaction"
	hookChan<- name
}

func BlockProcessingError(tx core.Hash, block core.Hash, err error) {
	name := "BlockProcessingError"
	hookChan<- name
}

func PostProcessTransaction(tx core.Hash, block core.Hash, i int, receipt []byte) {
	name := "PostProcessTransaction"
	hookChan<- name
}

func PostProcessBlock(block core.Hash) {
	name := "PostProcessBlock"
	hookChan<- name
}

func NewHead(block []byte, hash core.Hash, logs [][]byte, td *big.Int) {
	log.Error("inside custom newhead function")
	name := "NewHead"
	hookChan <- name
}

func NewSideBlock(block []byte, hash core.Hash, logs [][]byte) {
	name := "NewSideBlock"
	hookChan <- name
}

func Reorg(commonBlock core.Hash, oldChain, newChain []core.Hash) {
	name := "Reord"
	hookChan <- name
}

func SetTrieFlushIntervalClone(t time.Duration) time.Duration {
	name := "SetTrieFlushIntervalClone"
	hookChan<- name
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
	name := "ModifyAncients"
	hookChan<- name
}

func AppendAncient(number uint64, hash, header, body, receipts, td []byte) {
	name := "AppendAncient"
	hookChan<- name
}

// core/state/

func StateUpdate(blockRoot core.Hash, parentRoot core.Hash, coreDestructs map[core.Hash]struct{}, coreAccounts map[core.Hash][]byte, coreStorage map[core.Hash]map[core.Hash][]byte, coreCode map[core.Hash][]byte) {
	name := "StateUpdate"
	hookChan<- name
}

// core/vm we have code in core/vm but not hooks

// rpc/

func GetRPCCalls(s0 string, s1 string, s2 string) {
	log.Error("inside custom newhead GetRPCCalls fucntion", "arg0", s0, "arg1", s1, "arg2", s2)
	name := "GetRPCCalls"
	hookChan <- name
}

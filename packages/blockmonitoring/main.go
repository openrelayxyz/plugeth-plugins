package main

import (
	"context"
	"time"
	"math"
	"math/big"
	"sync"
	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted/hexutil"
	"github.com/openrelayxyz/plugeth-utils/restricted/rlp"
	"github.com/openrelayxyz/plugeth-utils/restricted/types"
)

type eventType byte

const (
	preProcessEvent eventType = iota
	postProcessEvent
	newHeadEvent
)

type blockEvent struct {
	T eventType `json:"eventType"`
	Time time.Time `json:"time"`
	Hash core.Hash `json:"hash"`
	Age time.Duration `json:"age,omitempty"`
}

var (
	events map[int][]blockEvent
	blocks map[core.Hash]int
	log core.Logger
	lock sync.RWMutex
	lastNewHead time.Time
)

func Initialize(ctx core.Context, loader core.PluginLoader, logger core.Logger) {
	events = make(map[int][]blockEvent)
	blocks = make(map[core.Hash]int)
	log = logger
}

func NewHead(blockBytes []byte, hash core.Hash, logsBytes [][]byte, td *big.Int) {
	lock.RLock()
	blockNo := blocks[hash]
	lock.RUnlock()
	events[blockNo] = append(events[blockNo], blockEvent{
		T: newHeadEvent,
		Time: time.Now(),
		Hash: hash,
	})
	
}

func PreProcessBlock(hash core.Hash, num uint64, blockrlp []byte) {
	lock.Lock()
	blocks[hash] = int(num)
	lock.Unlock()
	events[int(num)] = append(events[int(num)], blockEvent{
		T: preProcessEvent,
		Time: time.Now(),
		Hash: hash,
	})
}
func PostProcessBlock(hash core.Hash) {
	lock.RLock()
	blockNo := blocks[hash]
	lock.RUnlock()
	events[blockNo] = append(events[blockNo], blockEvent{
		T: postProcessEvent,
		Time: time.Now(),
		Hash: hash,
	})
}

type blockMonitor struct {
	backend core.Backend
}

func (b *blockMonitor) GetEventsByBlockNumber(ctx context.Context, num hexutil.Uint64) []blockEvent {
	result := make([]blockEvent, len(events[int(num)]))
	for i, e := range events[int(num)] {
		result[i] = blockEvent{
			T: e.T,
			Time: e.Time,
			Hash: e.Hash,
		}
		rlpblock, err := b.backend.BlockByHash(ctx, e.Hash)
		if err != nil {
			log.Warn("Error getting block", "num", num, "hash", e.Hash)
			continue
		}
		var block types.Block
		if err := rlp.DecodeBytes(rlpblock, &block); err != nil {
			log.Error("Failed to decode block", "hash", e.Hash, "err", err)
			continue
		}
		result[i].Age = e.Time.Sub(time.Unix(int64(block.Time()), 0))
	}
	return events[int(num)]
}

func (*blockMonitor) MinMonitoredBlock() hexutil.Uint64{
	val := math.MaxInt
	for i, _ := range events{
		if i < val {
			val = i
		}
	}
	return hexutil.Uint64(val)
}


// GetAPIs exposes the BlockUpdates service under the cardinal namespace.
func GetAPIs(stack core.Node, backend core.Backend) []core.API {
	return []core.API{
	 {
		 Namespace: "plugeth",
		 Version:	 "1.0",
		 Service:	 &blockMonitor{backend},
		 Public:		true,
	 },
 }
}
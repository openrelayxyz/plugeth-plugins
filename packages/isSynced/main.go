package main

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted/rlp"
	"github.com/openrelayxyz/plugeth-utils/restricted/types"
	"gopkg.in/urfave/cli.v1"
)

type IsSyncedService struct {
	backend core.Backend
	stack   core.Node
}

type peerinfo struct {
	Protocols struct {
		Eth struct {
			Difficulty *big.Int `json:"difficulty"`
		} `json:"eth"`
	} `json:"protocols"`
}

var log core.Logger

var httpApiFlagName = "http.api"

func Initialize(ctx *cli.Context, loader core.PluginLoader, logger core.Logger) {
	log = logger
	v := ctx.GlobalString(httpApiFlagName)
	if v != "" {
		ctx.GlobalSet(httpApiFlagName, v+",plugeth")
	} else {
		ctx.GlobalSet(httpApiFlagName, "eth,net,web3,plugeth")
		log.Info("Loaded isSynced plugin")
	}
}

func GetAPIs(stack core.Node, backend core.Backend) []core.API {
	return []core.API{
		{
			Namespace: "plugeth",
			Version:   "1.0",
			Service:   &IsSyncedService{backend, stack},
			Public:    true,
		},
	}
}

func (service *IsSyncedService) IsSynced(ctx context.Context) (interface{}, error) {
	client, err := service.stack.Attach()
	if err != nil {
		return nil, err
	}
	block := &types.Block{}
	var x []peerinfo
	err = client.Call(&x, "admin_peers")
	peers := false
	if err := rlp.DecodeBytes(service.backend.CurrentBlock(), block); err != nil {
		return nil, err
	}
	totalDifficulty := service.backend.GetTd(ctx, block.Hash())
	y := len(x)
	if y > 0 && totalDifficulty != nil {
		for i := range x {
			if totalDifficulty.Cmp(x[i].Protocols.Eth.Difficulty) <= 0 {
				peers = true
				break
			}
		}
	}
	progress := service.backend.Downloader().Progress()
	if progress.HighestBlock() == 0 {
		peers = false
	}
	if time.Now().Unix()-int64(block.Time()) < 60 {
		peers = true
	}
	highest := progress.HighestBlock()
	if highest < block.NumberU64() {
		highest = block.NumberU64()
	}
	return map[string]interface{}{
		"startingBlock":       fmt.Sprintf("%#x", (progress.StartingBlock())),
		"currentBlock":        fmt.Sprintf("%#x", (progress.CurrentBlock())),
		"highestBlock":        fmt.Sprintf("%#x", (highest)),
		"syncedAccounts":      fmt.Sprintf("%#x", (progress.SyncedAccounts())),
		"syncedAccountBytes":  fmt.Sprintf("%#x", (progress.SyncedAccountBytes())),
		"syncedBytecodes":     fmt.Sprintf("%#x", (progress.SyncedBytecodes())),
		"syncedBytecodeBytes": fmt.Sprintf("%#x", (progress.SyncedBytecodeBytes())),
		"syncedStorage":       fmt.Sprintf("%#x", (progress.SyncedStorage())),
		"syncedStorageBytes":  fmt.Sprintf("%#x", (progress.SyncedStorageBytes())),
		"healedTrienodes":     fmt.Sprintf("%#x", (progress.HealedTrienodes())),
		"healedTrienodeBytes": fmt.Sprintf("%#x", (progress.HealedTrienodeBytes())),
		"healedBytecodes":     fmt.Sprintf("%#x", (progress.HealedBytecodes())),
		"healedBytecodeBytes": fmt.Sprintf("%#x", (progress.HealedBytecodeBytes())),
		"healingTrienodes":    fmt.Sprintf("%#x", (progress.HealingTrienodes())),
		"healingBytecode":     fmt.Sprintf("%#x", (progress.HealingBytecode())),
		"activePeers":         peers,
		"nodeIsSynced":        peers && progress.CurrentBlock() >= progress.HighestBlock(),
	}, err
}

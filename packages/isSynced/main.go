package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted/crypto"
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
	var x []peerinfo
	err = client.Call(&x, "admin_peers")
	peers := false
	hash := crypto.Keccak256Hash(service.backend.CurrentHeader())
	totalDifficulty := service.backend.GetTd(ctx, hash)
	y := len(x)
	if y > 0 {
		peers = true
		for i := range x {
			if totalDifficulty.Cmp(x[i].Protocols.Eth.Difficulty) < 0 {
				peers = false
				break
			}
		}
	}
	progress := service.backend.Downloader().Progress()
	return map[string]interface{}{
		"startingBlock": fmt.Sprintf("%#x", (progress.StartingBlock())),
		"currentBlock":  fmt.Sprintf("%#x", (progress.CurrentBlock())),
		"highestBlock":  fmt.Sprintf("%#x", (progress.HighestBlock())),
		"pulledStates":  fmt.Sprintf("%#x", (progress.PulledStates())),
		"knownStates":   fmt.Sprintf("%#x", (progress.KnownStates())),
		"activePeers":   peers,
		"nodeIsSynced":  peers && progress.CurrentBlock() >= progress.HighestBlock(),
	}, err
}

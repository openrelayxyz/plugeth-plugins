package main

import (
	"context"

	"github.com/openrelayxyz/plugeth-utils/core"
)

type BlockByTxHashService struct {
	backend core.Backend
	stack   core.Node
}

var log core.Logger

var httpApiFlagName = "http.api"

func Initialize(ctx core.Context, loader core.PluginLoader, logger core.Logger) {
	log = logger
	v := ctx.String(httpApiFlagName)
	if v != "" {
		ctx.Set(httpApiFlagName, v+",plugeth")
	} else {
		ctx.Set(httpApiFlagName, "eth,net,web3,plugeth")
		log.Info("Loaded getBlockByTransactionHash plugin")
	}
}

func GetAPIs(stack core.Node, backend core.Backend) []core.API {
	return []core.API{
		{
			Namespace: "plugeth",
			Version:   "1.0",
			Service:   &BlockByTxHashService{backend, stack},
			Public:    true,
		},
	}
}

func (service *BlockByTxHashService) GetBlockByTransactionHash(ctx context.Context, txHash core.Hash) (interface{}, error) {
	
	_, blockHash, _, _, err := service.backend.GetTransaction(ctx, txHash)
	if err != nil {
		return nil, err
	}

	client, err := service.stack.Attach()
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := client.Call(&result, "eth_getBlockByHash", blockHash, true); err != nil {
		return nil, err
	}

	return result, nil
	
}
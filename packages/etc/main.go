package main

import (
	"context"

	"github.com/openrelayxyz/plugeth-utils/core"
)

type ClassicService struct {
	backend core.Backend
	stack   core.Node
}

var log core.Logger

var httpApiFlagName = "http.api"

func Initialize(ctx core.Context, loader core.PluginLoader, logger core.Logger) {
	log = logger
	log.Error("inside initalize")
	v := ctx.String(httpApiFlagName)
	if v != "" {
		ctx.Set(httpApiFlagName, v+",plugeth")
	} else {
		ctx.Set(httpApiFlagName, "eth,net,web3,plugeth")
	}
	log.Info("Loaded Ethereum classic plugin")
}

func GetAPIs(stack core.Node, backend core.Backend) []core.API {
	return []core.API{
		{
			Namespace: "plugeth",
			Version:   "1.0",
			Service:   &ClassicService{backend, stack},
			Public:    true,
		},
	}
}

func (service *ClassicService) Test(ctx context.Context) string {
	return "oh me oh me oh my"
}

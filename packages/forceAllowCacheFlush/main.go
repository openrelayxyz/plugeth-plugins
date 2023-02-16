package main

import (
	"context"
	"github.com/openrelayxyz/plugeth-utils/core"
)

var log core.Logger

type FlushCacheService struct {
	stack   core.Node
}

func Initialize(ctx core.Context, loader core.PluginLoader, logger core.Logger) {
	log = logger
	log.Info("loaded Flush-Dirty-Cache plugin")
}

func GetAPIs(stack core.Node, backend core.Backend) []core.API {
	return []core.API{
		{
			Namespace: "plugeth",
			Version:   "1.0",
			Service:   &FlushCacheService{stack},
			Public:    true,
		},
	}
}

var val bool 

func reSetVal() {
	val = false
}

func FlushCache() bool {
	defer reSetVal()
	return val
}

func (service *FlushCacheService) CallFlushCache(ctx context.Context) string {
	val = true
	return "Dirty Cache has been triggered to flush"
}
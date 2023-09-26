package main

import (
	"context"
	"path/filepath"

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

func DefaultDataDir(path string) string {
	return filepath.Join(path, "classic")
}

func SetBootstrapNodes() []string {
	var classicBootnodes = []string{} // this will need to be modified?
	return classicBootnodes
}

func SetNetworkId() *uint64 {
	var networkId *uint64
	classicNetworkId := uint64(1)
	networkId = &classicNetworkId
	return networkId 
}

func SetETHDiscoveryURLs() []string {
	var result []string
	return result
}

func SetSnapDiscoveryURLs() []string {
	var result []string
	return result
}


func (service *ClassicService) Test(ctx context.Context) string {
	return "total classic"
}

package main

import (
	"github.com/openrelayxyz/plugeth-utils/core"
	"gopkg.in/urfave/cli.v1"
)

var (
	log core.Logger
)

type myservice struct{}

func (*myservice) Hello() string {
	return "Hello world"
}

func Initialize(ctx *cli.Context, loader core.PluginLoader, logger core.Logger) {
	log = logger
	log.Info("Initialized hello")
}

func GetAPIs(node core.Node, backend core.Backend) []core.API {
	defer log.Info("APIs Initialized")
	return []core.API{
		{
			Namespace: "mynamespace",
			Version:   "1.0",
			Service:   &myservice{},
			Public:    true,
		},
	}
}

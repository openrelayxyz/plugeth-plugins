package main

import (
	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/urfave/cli/v2"
)

var log core.Logger

func Initialize(ctx *cli.Context, loader core.PluginLoader, logger core.Logger) {
	log = logger
	log.Info("loaded Get Rpc Calls plugin")
}

func GetRPCCalls(id, method, params string) {

	log.Info("Received RPC Call", "id", id, "method", method, "params", params)

}

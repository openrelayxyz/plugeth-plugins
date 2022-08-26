package main

import (
	"context"
	"time"
	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/urfave/cli/v2"
)

var (
	log core.Logger
)

type myservice struct{}

func (*myservice) Hello() string {
	return "Hello world"
}


func (*myservice) Timer(ctx context.Context) (<-chan int64, error) {
	ticker := time.NewTicker(time.Second)
	ch := make(chan int64)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				close(ch)
				return
			case t := <-ticker.C:
				ch <- t.UnixNano()
			}
		}
	}()
	return ch, nil
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

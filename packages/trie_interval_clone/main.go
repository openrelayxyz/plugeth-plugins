package main

import (
	"fmt"
	"time"
	"context"
	"github.com/openrelayxyz/plugeth-utils/core"
)

var log core.Logger

type TrieIntervalService struct {
}

func Initialize(ctx core.Context, loader core.PluginLoader, logger core.Logger) {
	log = logger
	log.Info("loaded setTrieFlushInterval plugin")
}

func GetAPIs(stack core.Node, backend core.Backend) []core.API {
	return []core.API{
		{
			Namespace: "debug",
			Version:   "1.0",
			Service:   &TrieIntervalService{},
			Public:    true,
		},
	}
}

var nodeInterval time.Duration

var ModifiedInterval time.Duration 

func SetTrieFlushIntervalClone(duration time.Duration) time.Duration {
	nodeInterval = duration
	if ModifiedInterval > 0 {
		duration = ModifiedInterval
	}
	return duration
}

func (service *TrieIntervalService) SetTrieFlushInterval(ctx context.Context, arg string) (string, error) {
	newInterval, err := time.ParseDuration(arg)
	if err != nil {
		return "", err
	}
	ModifiedInterval = newInterval

	return fmt.Sprintf("flushInterval has been modified from %v to %v", nodeInterval, ModifiedInterval), nil
}
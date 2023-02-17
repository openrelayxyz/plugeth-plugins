package main

import (
	"time"
	"context"
	"github.com/openrelayxyz/plugeth-utils/core"
)

var log core.Logger

type FlushCacheService struct {
}

func Initialize(ctx core.Context, loader core.PluginLoader, logger core.Logger) {
	log = logger
	log.Info("loaded Allow Force trei flush plugin")
}

func GetAPIs(stack core.Node, backend core.Backend) []core.API {
	return []core.API{
		{
			Namespace: "plugeth",
			Version:   "1.0",
			Service:   &FlushCacheService{},
			Public:    true,
		},
	}
}

var nodeInterval time.Duration 

func PauseTreiCommit(gcproc, duration time.Duration) time.Duration {
	log.Error("these are the values", "gcproc", gcproc, "duration", duration)
	if nodeInterval.Minutes() > 0 {
		duration = nodeInterval
	}
	return duration
}

// or something like where we inspect the gcproc and if we are in puase mode we repeatedly call the function and incredment by some t
// the tickher inspector and incrementor would need to be started and stopped via rpc. Or have a switch that is always false until it is true. 

func (service *FlushCacheService) InspectTrieInterval(ctx context.Context) string {
	return nodeInterval.String()
}

func (service *FlushCacheService) ModifyTrieInterval(ctx context.Context, arg string) (string, error) {
	newInterval, err := time.ParseDuration(arg)
	if err != nil {
		return "", err
	}
	nodeInterval = newInterval

	return "flushInterval has been modified", nil
}

var val bool 

var allowVal bool = true

func AllowTreiCommit() bool {
	defer func() {
		allowVal = true
	}()
	return allowVal
}

func ForceTreiCommit() bool {
	defer func() {
		val = false
	}()
	return val
}

func FlushCache() bool {
	defer func() {
		val = false
	}()
	return val
}

func (service *FlushCacheService) CallAllowTreiCommit(ctx context.Context) string {
	AllowTreiCommit()
	return "Trie is allowed to flush"
}

func (service *FlushCacheService) CallForceTreiCommit(ctx context.Context) string {
	val = true
	return "Trie has been forced to flush"
}

func (service *FlushCacheService) CallFlushCache(ctx context.Context) string {
	val = true
	return "Dirty Cache has been triggered to flush"
}
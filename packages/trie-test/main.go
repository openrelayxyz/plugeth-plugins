package main

import (
	"context"
	// "encoding/json"

	"github.com/openrelayxyz/plugeth-utils/core"
)


var (
	log core.Logger
)

type TrieTestService struct {
	backend core.Backend
	stack   core.Node
}

func Initialize(ctx core.Context, loader core.PluginLoader, logger core.Logger) {
	log = logger
	log.Info("Initialized State Trie test plugin")
}

func GetAPIs(stack core.Node, backend core.Backend) []core.API {
	return []core.API{
		{
			Namespace: "plugeth",
			Version:   "1.0",
			Service:   &TrieTestService{backend, stack},
			Public:    true,
		},
	}
}

type Trie struct {
	Root interface{} `json:"root"` // node
	Owner core.Hash `json:"owner"`
	Unhashed int `json:"unhashed"`
	Reader interface{} `json:"reader"` // *trieReader
	Tracer interface{} `json:"tracer"` // *tracer 
	Preimages map[core.Hash][]byte `json:"preimages"`
	HashKeyBuf []byte `json:"hashKeyBuf"`
	SecKeyCache map[string][]byte `json:"secKeyCache"`
	SecKeyCacheOwner interface{} `json:"secKeyCacheOwner"` // *StateTrie
}

func (t *TrieTestService) GetTrie(ctx context.Context, hash core.Hash) (interface{}, error) {
	x, err := t.backend.GetTrie(hash)
	if err != nil {
		return nil, err
	}
	return x, nil
}

func (t *TrieTestService) GetAccountTrie(ctx context.Context, stateRoot core.Hash, account core.Address) (interface{}, error) {
	x, err := t.backend.GetAccountTrie(stateRoot, account)
	if err != nil {
		return nil, err
	}
	return x, nil
}



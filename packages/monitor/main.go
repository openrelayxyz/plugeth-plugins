package main

import (
	"math/big"
	"bytes"
	"strings"
	"encoding/json"
	lru "github.com/hashicorp/golang-lru"
	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted"
	"github.com/openrelayxyz/plugeth-utils/restricted/crypto"
	"github.com/openrelayxyz/plugeth-utils/restricted/hexutil"
	"github.com/openrelayxyz/plugeth-utils/restricted/rlp"
	"io"
)

var (
	pl core.PluginLoader
	backend restricted.Backend
	lastBlock core.Hash
	cache *lru.Cache
	recentEmits *lru.Cache
	snapshotFlagName = "snapshot"
	log core.Logger
	blockEvents core.Feed
)

var (
	httpApiFlagName = "http.api"
	wsApiFlagName = "ws.api"
)


type stateUpdate struct {
	Destructs map[core.Hash]struct{}
	Accounts map[core.Hash][]byte
	Storage map[core.Hash]map[core.Hash][]byte
	Code map[core.Hash][]byte
}


type kvpair struct {
	Key core.Hash
	Value []byte
}

type storage struct {
	Account core.Hash
	Data []kvpair
}

type storedStateUpdate struct {
	Destructs []core.Hash
	Accounts	[]kvpair
	Storage	 []storage
	Code	[]kvpair
}

func InitializeNode(stack core.Node, b restricted.Backend) {
	backend = b
	log.Info("Initialized state-update monitor plugin")
}

func Initialize(ctx core.Context, loader core.PluginLoader, logger core.Logger) {
	log = logger
	pl = loader
	blockEvents = pl.GetFeed()
	cache, _ = lru.New(128)
	recentEmits, _ = lru.New(128)
	if !ctx.Bool(snapshotFlagName) {
		log.Warn("Snapshots are required for StateUpdate plugins, but are currently disabled. State Updates will be unavailable")
	}
	v := ctx.String(httpApiFlagName)
	if v == "" {
		ctx.Set(httpApiFlagName, "eth,net,web3,plugeth")
	} else if !strings.Contains(v, "plugeth") {
		ctx.Set(httpApiFlagName, v+",plugeth")
	}
	v = ctx.String(wsApiFlagName)
	if v == "" {
		ctx.Set(wsApiFlagName, "eth,net,web3,plugeth")
	} else if !strings.Contains(v, "plugeth") {
		ctx.Set(wsApiFlagName, v+",plugeth")
	}
	log.Info("Loaded state-update monitor plugin")
}


func (su *stateUpdate) MarshalJSON() ([]byte, error) {
	result := make(map[string]interface{})
	destructs := make([]core.Hash, 0, len(su.Destructs))
	for k := range su.Destructs {
		destructs = append(destructs, k)
	}
	result["destructs"] = destructs
	accounts := make(map[string]hexutil.Bytes)
	for k, v := range su.Accounts {
		accounts[k.String()] = hexutil.Bytes(v)
	}
	result["accounts"] = accounts
	storage := make(map[string]map[string]hexutil.Bytes)
	for m, s := range su.Storage {
		storage[m.String()] = make(map[string]hexutil.Bytes)
		for k, v := range s {
			storage[m.String()][k.String()] = hexutil.Bytes(v)
		}
	}
	result["storage"] = storage
	code := make(map[string]hexutil.Bytes)
	for k, v := range su.Code {
		code[k.String()] = hexutil.Bytes(v)
	}
	result["code"] = code
	return json.Marshal(result)
}

func (su *stateUpdate) EncodeRLP(w io.Writer) error {
	destructs := make([]core.Hash, 0, len(su.Destructs))
	for k := range su.Destructs {
		destructs = append(destructs, k)
	}
	accounts := make([]kvpair, 0, len(su.Accounts))
	for k, v := range su.Accounts {
		accounts = append(accounts, kvpair{k, v})
	}
	s := make([]storage, 0, len(su.Storage))
	for a, m := range su.Storage {
		accountStorage := storage{a, make([]kvpair, 0, len(m))}
		for k, v := range m {
			accountStorage.Data = append(accountStorage.Data, kvpair{k, v})
		}
		s = append(s, accountStorage)
	}
	code := make([]kvpair, 0, len(su.Code))
	for k, v := range su.Code {
		code = append(code, kvpair{k, v})
	}
	return rlp.Encode(w, storedStateUpdate{destructs, accounts, s, code})
}

func (su *stateUpdate) DecodeRLP(s *rlp.Stream) error {
	ssu := storedStateUpdate{}
	if err := s.Decode(&ssu); err != nil { return err }
	su.Destructs = make(map[core.Hash]struct{})
	for _, s := range ssu.Destructs {
		su.Destructs[s] = struct{}{}
	}
	su.Accounts = make(map[core.Hash][]byte)
	for _, kv := range ssu.Accounts {
		su.Accounts[kv.Key] = kv.Value
	}
	su.Storage = make(map[core.Hash]map[core.Hash][]byte)
	for _, s := range ssu.Storage {
		su.Storage[s.Account] = make(map[core.Hash][]byte)
		for _, kv := range s.Data {
			su.Storage[s.Account][kv.Key] = kv.Value
		}
	}
	su.Code = make(map[core.Hash][]byte)
	for _, kv := range ssu.Code {
		su.Code[kv.Key] = kv.Value
	}
	return nil
}

var (
	EmptyRootHash  = core.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
	EmptyCodeHash = crypto.Keccak256Hash(nil)
)

type Account struct {
	Nonce    uint64
	Balance  *big.Int
	Root     []byte
	CodeHash []byte
}

func normalizeAccount(acct *Account) {
	if acct.Balance == nil {
		acct.Balance = new(big.Int)
	}
	if bytes.Equal(acct.Root, EmptyRootHash.Bytes()) {
		acct.Root = []byte{}
	}
	if bytes.Equal(acct.CodeHash, EmptyCodeHash.Bytes()) {
		acct.CodeHash = []byte{}
	}
}

func StateUpdate(blockRoot core.Hash, parentRoot core.Hash, destructs map[core.Hash]struct{}, accounts map[core.Hash][]byte, storage map[core.Hash]map[core.Hash][]byte, codeUpdates map[core.Hash][]byte) {
	if backend == nil {
		log.Warn("State update called before InitializeNode", "root", blockRoot)
		return
	}
	t, err := backend.GetTrie(parentRoot)
	if err != nil {
		log.Error("Error getting trie", "root", parentRoot, "err", err)
		return
	}

	for hashedAddr, v := range accounts {
		parentV := t.GetKey(hashedAddr.Bytes())
		var parentAcct, acct Account
		if err := rlp.DecodeBytes(v, &acct); err != nil {
			log.Error("Error decoding acct", "err", err, "acctData", v)
			continue
		}
		if err := rlp.DecodeBytes(parentV, &parentAcct); err != nil {
			log.Error("Error decoding parentAcct", "err", err, "acctData", parentV)
			continue
		}
		normalizeAccount(&acct)
		normalizeAccount(&parentAcct)
		if acct.Nonce == parentAcct.Nonce && bytes.Equal(acct.Root, parentAcct.Root) && bytes.Equal(acct.CodeHash, parentAcct.CodeHash) && acct.Balance.Cmp(parentAcct.Balance) == 0 {
			log.Error("StateUpdate account equal to parent", "block", blockRoot, "parent", parentRoot, "acctHash", hashedAddr)
		}
		// Nonce    uint64
		// Balance  *big.Int
		// Root     Hash // merkle root of the storage trie
		// CodeHash []byte

	}
	// su := &stateUpdate{
	// 	Destructs: destructs,
	// 	Accounts: accounts,
	// 	Storage: storage,
	// 	Code: codeUpdates,
	// }

	// log.Error("envoking state update monitor", "accounts", su.Accounts)
}

// Copyright 2019 The multi-geth Authors
// This file is part of the multi-geth library.
//
// The multi-geth library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The multi-geth library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the multi-geth library. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"context"
	"math/big"
	"errors"
	"sync"
	"sync/atomic"
	"time"
	"os"

	"github.com/edsrzf/mmap-go"

	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted/params"
	"github.com/openrelayxyz/plugeth-utils/restricted/types"
)

var (
	// algorithmRevision is the data structure version used for file naming.
	algorithmRevision = 23

	// dumpMagic is a dataset dump header to sanity check a data dump.
	dumpMagic = []uint32{0xbaddcafe, 0xfee1dead}

	ErrInvalidDumpMagic = errors.New("invalid dump magic")
)

// Lengths of hashes and addresses in bytes.
const (
	// HashLength is the expected length of the hash
	HashLength = 32
	// AddressLength is the expected length of the address
	AddressLength = 20
)

// This file holds the Configurator interfaces.
// Interface methods follow distinct naming and signature patterns
// to enable abstracted logic.
//
// All methods are in pairs; Get'ers and Set'ers.
// Some Set methods are prefixed with Must, ie MustSet. These methods
// are allowed to return errors for debugging and logging, but
// any non-nil errors returned should stop program execution.
//
// All Forking methods (getting and setting Hard-fork requiring protocol changes)
// are suffixed with "Transition", and use *uint64 as in- and out-put variable types.
// A pointer is used because it's important that the fields can be nil, to signal
// being unset (as opposed to zero value uint64 == 0, ie from Genesis).

type ChainHeaderReader interface {
	// Config retrieves the blockchain's chain configuration.
	Config() *params.ChainConfig

	// Config retrieves the blockchain's chain configuration.
	// PluginConfig() ChainConfigurator

	// CurrentHeader retrieves the current header from the local chain.
	CurrentHeader() *types.Header

	// GetHeader retrieves a block header from the database by hash and number.
	GetHeader(hash core.Hash, number uint64) *types.Header

	// GetHeaderByNumber retrieves a block header from the database by number.
	GetHeaderByNumber(number uint64) *types.Header

	// GetHeaderByHash retrieves a block header from the database by its hash.
	GetHeaderByHash(hash core.Hash) *types.Header

	// GetTd retrieves the total difficulty from the database by hash and number.
	GetTd(hash core.Hash, number uint64) *big.Int
}

type remoteSealer struct {
	works        map[core.Hash]*types.Block
	rates        map[core.Hash]hashrate
	currentBlock *types.Block
	currentWork  [4]string
	notifyCtx    context.Context
	cancelNotify context.CancelFunc // cancels all notification requests
	reqWG        sync.WaitGroup     // tracks notification request goroutines

	ethash       *Ethash
	noverify     bool
	notifyURLs   []string
	results      chan<- *types.Block
	workCh       chan *sealTask   // Notification channel to push new work and relative result channel to remote sealer
	fetchWorkCh  chan *sealWork   // Channel used for remote sealer to fetch mining work
	submitWorkCh chan *mineResult // Channel used for remote sealer to submit their mining result
	fetchRateCh  chan chan uint64 // Channel used to gather submitted hash rate for local or remote sealer.
	submitRateCh chan *hashrate   // Channel used for remote sealer to submit their mining hashrate
	requestExit  chan struct{}
	exitCh       chan struct{}
}

// hashrate wraps the hash rate submitted by the remote sealer.
type hashrate struct {
	id   core.Hash
	ping time.Time
	rate uint64

	done chan struct{}
}

// sealTask wraps a seal block with relative result channel for remote sealer thread.
type sealTask struct {
	block   *types.Block
	results chan<- *types.Block
}

// sealWork wraps a seal work package for remote sealer.
type sealWork struct {
	errc chan error
	res  chan [4]string
}

// mineResult wraps the pow solution parameters for the specified block.
type mineResult struct {
	nonce     types.BlockNonce
	mixDigest core.Hash
	hash      core.Hash

	errc chan error
}

// ChainReader defines a small collection of methods needed to access the local
// blockchain during header and/or uncle verification.
type ChainReader interface {
	ChainHeaderReader

	// GetBlock retrieves a block from the database by hash and number.
	GetBlock(hash core.Hash, number uint64) *types.Block
}

type Configurator interface {
	ChainConfigurator
	GenesisBlocker
}

type ChainConfigurator interface {
	String() string

	ProtocolSpecifier
	Forker
	ConsensusEnginator // Consensus Engine
	// CHTer
}

// ProtocolSpecifier defines protocol interfaces that are agnostic of consensus engine.
type ProtocolSpecifier interface {
	GetAccountStartNonce() *uint64
	SetAccountStartNonce(n *uint64) error
	GetMaximumExtraDataSize() *uint64
	SetMaximumExtraDataSize(n *uint64) error
	GetMinGasLimit() *uint64
	SetMinGasLimit(n *uint64) error
	GetGasLimitBoundDivisor() *uint64
	SetGasLimitBoundDivisor(n *uint64) error
	GetNetworkID() *uint64
	SetNetworkID(n *uint64) error
	GetChainID() *big.Int
	SetChainID(i *big.Int) error
	GetSupportedProtocolVersions() []uint
	SetSupportedProtocolVersions(p []uint) error
	GetMaxCodeSize() *uint64
	SetMaxCodeSize(n *uint64) error

	GetElasticityMultiplier() uint64
	SetElasticityMultiplier(n uint64) error
	GetBaseFeeChangeDenominator() uint64
	SetBaseFeeChangeDenominator(n uint64) error

	// Be careful with EIP2.
	// It is a messy EIP, specifying diverse changes, like difficulty, intrinsic gas costs for contract creation,
	// txpool management, and contract OoG handling.
	// It is both Ethash-specific and _not_.
	GetEIP2Transition() *uint64
	SetEIP2Transition(n *uint64) error

	GetEIP7Transition() *uint64
	SetEIP7Transition(n *uint64) error
	GetEIP150Transition() *uint64
	SetEIP150Transition(n *uint64) error
	GetEIP152Transition() *uint64
	SetEIP152Transition(n *uint64) error
	GetEIP160Transition() *uint64
	SetEIP160Transition(n *uint64) error
	GetEIP161abcTransition() *uint64
	SetEIP161abcTransition(n *uint64) error
	GetEIP161dTransition() *uint64
	SetEIP161dTransition(n *uint64) error
	GetEIP170Transition() *uint64
	SetEIP170Transition(n *uint64) error
	GetEIP155Transition() *uint64
	SetEIP155Transition(n *uint64) error
	GetEIP140Transition() *uint64
	SetEIP140Transition(n *uint64) error
	GetEIP198Transition() *uint64
	SetEIP198Transition(n *uint64) error
	GetEIP211Transition() *uint64
	SetEIP211Transition(n *uint64) error
	GetEIP212Transition() *uint64
	SetEIP212Transition(n *uint64) error
	GetEIP213Transition() *uint64
	SetEIP213Transition(n *uint64) error
	GetEIP214Transition() *uint64
	SetEIP214Transition(n *uint64) error
	GetEIP658Transition() *uint64
	SetEIP658Transition(n *uint64) error
	GetEIP145Transition() *uint64
	SetEIP145Transition(n *uint64) error
	GetEIP1014Transition() *uint64
	SetEIP1014Transition(n *uint64) error
	GetEIP1052Transition() *uint64
	SetEIP1052Transition(n *uint64) error
	GetEIP1283Transition() *uint64
	SetEIP1283Transition(n *uint64) error
	GetEIP1283DisableTransition() *uint64
	SetEIP1283DisableTransition(n *uint64) error
	GetEIP1108Transition() *uint64
	SetEIP1108Transition(n *uint64) error
	GetEIP2200Transition() *uint64
	SetEIP2200Transition(n *uint64) error
	GetEIP2200DisableTransition() *uint64
	SetEIP2200DisableTransition(n *uint64) error
	GetEIP1344Transition() *uint64
	SetEIP1344Transition(n *uint64) error
	GetEIP1884Transition() *uint64
	SetEIP1884Transition(n *uint64) error
	GetEIP2028Transition() *uint64
	SetEIP2028Transition(n *uint64) error
	GetECIP1080Transition() *uint64
	SetECIP1080Transition(n *uint64) error
	GetEIP1706Transition() *uint64
	SetEIP1706Transition(n *uint64) error
	GetEIP2537Transition() *uint64
	SetEIP2537Transition(n *uint64) error

	GetECBP1100Transition() *uint64
	SetECBP1100Transition(n *uint64) error
	GetEIP2315Transition() *uint64
	SetEIP2315Transition(n *uint64) error

	// ModExp gas cost
	GetEIP2565Transition() *uint64
	SetEIP2565Transition(n *uint64) error

	// Gas cost increases for state access opcodes
	GetEIP2929Transition() *uint64
	SetEIP2929Transition(n *uint64) error

	// Optional access lists
	GetEIP2930Transition() *uint64
	SetEIP2930Transition(n *uint64) error

	// Typed transaction envelope
	GetEIP2718Transition() *uint64
	SetEIP2718Transition(n *uint64) error

	GetEIP1559Transition() *uint64
	SetEIP1559Transition(n *uint64) error

	GetEIP3541Transition() *uint64
	SetEIP3541Transition(n *uint64) error

	GetEIP3529Transition() *uint64
	SetEIP3529Transition(n *uint64) error

	GetEIP3198Transition() *uint64
	SetEIP3198Transition(n *uint64) error

	// EIP4399 is the RANDOM opcode.
	// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-4399.md
	GetEIP4399Transition() *uint64
	SetEIP4399Transition(n *uint64) error

	// Shanghai:
	//
	// EIP3651: Warm COINBASE
	GetEIP3651TransitionTime() *uint64
	SetEIP3651TransitionTime(n *uint64) error
	// EIP3855: PUSH0 instruction
	GetEIP3855TransitionTime() *uint64
	SetEIP3855TransitionTime(n *uint64) error
	// EIP3860: Limit and meter initcode
	GetEIP3860TransitionTime() *uint64
	SetEIP3860TransitionTime(n *uint64) error
	// EIP4895: Beacon chain push withdrawals as operations
	GetEIP4895TransitionTime() *uint64
	SetEIP4895TransitionTime(n *uint64) error
	// EIP6049: Deprecate SELFDESTRUCT
	GetEIP6049TransitionTime() *uint64
	SetEIP6049TransitionTime(n *uint64) error

	// Shanghai expressed as block activation numbers:
	GetEIP3651Transition() *uint64
	SetEIP3651Transition(n *uint64) error
	GetEIP3855Transition() *uint64
	SetEIP3855Transition(n *uint64) error
	GetEIP3860Transition() *uint64
	SetEIP3860Transition(n *uint64) error
	GetEIP4895Transition() *uint64
	SetEIP4895Transition(n *uint64) error
	GetEIP6049Transition() *uint64
	SetEIP6049Transition(n *uint64) error

	// GetMergeVirtualTransition is a Virtual fork after The Merge to use as a network splitter
	GetMergeVirtualTransition() *uint64
	SetMergeVirtualTransition(n *uint64) error

	// Cancun:
	// EIP4844 - Shard Blob Transactions - https://eips.ethereum.org/EIPS/eip-4844
	GetEIP4844TransitionTime() *uint64
	SetEIP4844TransitionTime(n *uint64) error

	// EIP1153 - Transient Storage opcodes - https://eips.ethereum.org/EIPS/eip-1153
	GetEIP1153TransitionTime() *uint64
	SetEIP1153TransitionTime(n *uint64) error

	// EIP5656 - MCOPY - Memory copying instruction - https://eips.ethereum.org/EIPS/eip-5656
	GetEIP5656TransitionTime() *uint64
	SetEIP5656TransitionTime(n *uint64) error

	// EIP6780 - SELFDESTRUCT only in same transaction - https://eips.ethereum.org/EIPS/eip-6780
	GetEIP6780TransitionTime() *uint64
	SetEIP6780TransitionTime(n *uint64) error
}

type Forker interface {
	// IsEnabled tells if interface has met or exceeded a fork block number.
	// eg. IsEnabled(c.GetEIP1108Transition, big.NewInt(42)))
	IsEnabled(fn func() *uint64, n *big.Int) bool
	IsEnabledByTime(fn func() *uint64, n *uint64) bool

	// ForkCanonHash yields arbitrary number/hash pairs.
	// This is an abstraction derived from the original EIP150 implementation.
	GetForkCanonHash(n uint64) core.Hash
	SetForkCanonHash(n uint64, h core.Hash) error
	GetForkCanonHashes() map[uint64]core.Hash
}

type ConsensusEnginator interface {
	GetConsensusEngineType() ConsensusEngineT
	MustSetConsensusEngineType(t ConsensusEngineT) error
	GetIsDevMode() bool
	SetDevMode(devMode bool) error

	EthashConfigurator
	CliqueConfigurator
	Lyra2Configurator
}

type EthashConfigurator interface {
	GetEthashMinimumDifficulty() *big.Int
	SetEthashMinimumDifficulty(i *big.Int) error
	GetEthashDifficultyBoundDivisor() *big.Int
	SetEthashDifficultyBoundDivisor(i *big.Int) error
	GetEthashDurationLimit() *big.Int
	SetEthashDurationLimit(i *big.Int) error
	GetEthashHomesteadTransition() *uint64
	SetEthashHomesteadTransition(n *uint64) error

	// GetEthashEIP779Transition should return the block if the node wants the fork.
	// Otherwise, nil should be returned.
	GetEthashEIP779Transition() *uint64 // DAO

	// SetEthashEIP779Transition should turn DAO support on (nonnil) or off (nil).
	SetEthashEIP779Transition(n *uint64) error
	GetEthashEIP649Transition() *uint64
	SetEthashEIP649Transition(n *uint64) error
	GetEthashEIP1234Transition() *uint64
	SetEthashEIP1234Transition(n *uint64) error
	GetEthashEIP2384Transition() *uint64
	SetEthashEIP2384Transition(n *uint64) error
	GetEthashEIP3554Transition() *uint64
	SetEthashEIP3554Transition(n *uint64) error
	GetEthashEIP4345Transition() *uint64
	SetEthashEIP4345Transition(n *uint64) error
	GetEthashECIP1010PauseTransition() *uint64
	SetEthashECIP1010PauseTransition(n *uint64) error
	GetEthashECIP1010ContinueTransition() *uint64
	SetEthashECIP1010ContinueTransition(n *uint64) error
	GetEthashECIP1017Transition() *uint64
	SetEthashECIP1017Transition(n *uint64) error
	GetEthashECIP1017EraRounds() *uint64
	SetEthashECIP1017EraRounds(n *uint64) error
	GetEthashEIP100BTransition() *uint64
	SetEthashEIP100BTransition(n *uint64) error
	GetEthashECIP1041Transition() *uint64
	SetEthashECIP1041Transition(n *uint64) error
	GetEthashECIP1099Transition() *uint64
	SetEthashECIP1099Transition(n *uint64) error
	GetEthashEIP5133Transition() *uint64 // Gray Glacier difficulty bomb delay
	SetEthashEIP5133Transition(n *uint64) error

	GetEthashTerminalTotalDifficulty() *big.Int
	SetEthashTerminalTotalDifficulty(n *big.Int) error

	GetEthashTerminalTotalDifficultyPassed() bool
	SetEthashTerminalTotalDifficultyPassed(t bool) error

	IsTerminalPoWBlock(parentTotalDiff *big.Int, totalDiff *big.Int) bool

	GetEthashDifficultyBombDelaySchedule() Uint64BigMapEncodesHex
	SetEthashDifficultyBombDelaySchedule(m Uint64BigMapEncodesHex) error
	GetEthashBlockRewardSchedule() Uint64BigMapEncodesHex
	SetEthashBlockRewardSchedule(m Uint64BigMapEncodesHex) error
}

type CliqueConfigurator interface {
	GetCliquePeriod() uint64
	SetCliquePeriod(n uint64) error
	GetCliqueEpoch() uint64
	SetCliqueEpoch(n uint64) error
}

type Lyra2Configurator interface {
	GetLyra2NonceTransition() *uint64
	SetLyra2NonceTransition(n *uint64) error
}

type BlockSealer interface {
	GetSealingType() BlockSealingT
	SetSealingType(t BlockSealingT) error
	BlockSealerEthereum
}

type BlockSealerEthereum interface {
	GetGenesisSealerEthereumNonce() uint64
	SetGenesisSealerEthereumNonce(n uint64) error
	GetGenesisSealerEthereumMixHash() core.Hash
	SetGenesisSealerEthereumMixHash(h core.Hash) error
}

type GenesisBlocker interface {
	BlockSealer
	Accounter
	GetGenesisDifficulty() *big.Int
	SetGenesisDifficulty(i *big.Int) error
	GetGenesisAuthor() core.Address
	SetGenesisAuthor(a core.Address) error
	GetGenesisTimestamp() uint64
	SetGenesisTimestamp(u uint64) error
	GetGenesisParentHash() core.Hash
	SetGenesisParentHash(h core.Hash) error
	GetGenesisExtraData() []byte
	SetGenesisExtraData(b []byte) error
	GetGenesisGasLimit() uint64
	SetGenesisGasLimit(u uint64) error
}

type Accounter interface {
	ForEachAccount(fn func(address core.Address, bal *big.Int, nonce uint64, code []byte, storage map[core.Hash]core.Hash) error) error
	UpdateAccount(address core.Address, bal *big.Int, nonce uint64, code []byte, storage map[core.Hash]core.Hash) error
}

type ConsensusEngineT int

const (
	ConsensusEngineT_Unknown = iota
	ConsensusEngineT_Ethash
	ConsensusEngineT_Clique
	ConsensusEngineT_Lyra2
)

func (c ConsensusEngineT) String() string {
	switch c {
	case ConsensusEngineT_Ethash:
		return "ethash"
	case ConsensusEngineT_Clique:
		return "clique"
	case ConsensusEngineT_Lyra2:
		return "lyra2"
	default:
		return "unknown"
	}
}

func (c ConsensusEngineT) IsEthash() bool {
	return c == ConsensusEngineT_Ethash
}

func (c ConsensusEngineT) IsClique() bool {
	return c == ConsensusEngineT_Clique
}

func (c ConsensusEngineT) IsLyra2() bool {
	return c == ConsensusEngineT_Lyra2
}

func (c ConsensusEngineT) IsUnknown() bool {
	return c == ConsensusEngineT_Unknown
}

// Uint64BigMapEncodesHex is a map that encodes and decodes w/ JSON hex format.
type Uint64BigMapEncodesHex map[uint64]*big.Int

type BlockSealingT int

const (
	BlockSealing_Unknown = iota
	BlockSealing_Ethereum
)

func (b BlockSealingT) String() string {
	switch b {
	case BlockSealing_Ethereum:
		return "ethereum"
	default:
		return "unknown"
	}
}

var big0 = big.NewInt(0)

var big1 = new(big.Int).SetInt64(1)

var big2 = new(big.Int).SetInt64(2)

var bigMinus99 = big.NewInt(-99)

var DisinflationRateQuotient = big.NewInt(4)

var DisinflationRateDivisor  = big.NewInt(5)

// two256 is a big integer representing 2^256
var two256 = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0))

// DAOForkBlockExtra is the block header extra-data field to set for the DAO fork
// point and a number of consecutive blocks to allow fast/light syncers to correctly
// pick the side they want  ("dao-hard-fork").
var DAOForkBlockExtra = FromHex("0x64616f2d686172642d666f726b")

// DAOForkExtraRange is the number of consecutive blocks from the DAO fork point
// to override the extra-data in to prevent no-fork attacks.
var DAOForkExtraRange = big.NewInt(10)

// DAORefundContract is the address of the refund contract to send DAO balances to.
var DAORefundContract = core.HexToAddress("0xbf4ed7b27f1d666546e30d74d50d173d20bca754")

var ErrBadProDAOExtra = errors.New("bad DAO pro-fork extra-data")

var ExpDiffPeriod *big.Int = big.NewInt(100000)

const (
	datasetInitBytes    = 1 << 30 // Bytes in dataset at genesis
	datasetGrowthBytes  = 1 << 23 // Dataset growth per epoch
	cacheInitBytes      = 1 << 24 // Bytes in cache at genesis
	cacheGrowthBytes    = 1 << 17 // Cache growth per epoch
	epochLengthDefault  = 30000   // Default epoch length (blocks per epoch)
	epochLengthECIP1099 = 60000   // Blocks per epoch if ECIP-1099 is activated
	mixBytes            = 128     // Width of mix
	hashBytes           = 64      // Hash length in bytes
	hashWords           = 16      // Number of 32 bit ints in a hash
	datasetParents      = 256     // Number of parents of each dataset element
	cacheRounds         = 3       // Number of rounds in cache production
	loopAccesses        = 64      // Number of accesses in hashimoto loop
	maxEpoch            = 2048    // Max Epoch for included tables
)

// lru tracks caches or datasets by their last use time, keeping at most N of them.
type lru[T cacheOrDataset] struct {
	what string
	new  func(epoch uint64, epochLength uint64) T
	mu   sync.Mutex
	// Items are kept in a LRU cache, but there is a special case:
	// We always keep an item for (highest seen epoch) + 1 as the 'future item'.
	cache      BasicLRU[uint64, T]
	future     uint64
	futureItem T
}

type cacheOrDataset interface {
	*cache | *dataset
}

// cache wraps an ethash cache with some metadata to allow easier concurrent use.
type cache struct {
	epoch       uint64    // Epoch for which this cache is relevant
	epochLength uint64    // Epoch length (ECIP-1099)
	dump        *os.File  // File descriptor of the memory mapped cache
	mmap        mmap.MMap // Memory map itself to unmap before releasing
	cache       []uint32  // The actual cache data content (may be memory mapped)
	once        sync.Once // Ensures the cache is generated only once
}

// dataset wraps an ethash dataset with some metadata to allow easier concurrent use.
type dataset struct {
	epoch       uint64      // Epoch for which this cache is relevant
	epochLength uint64      // Epoch length (ECIP-1099)
	dump        *os.File    // File descriptor of the memory mapped cache
	mmap        mmap.MMap   // Memory map itself to unmap before releasing
	dataset     []uint32    // The actual cache data content
	once        sync.Once   // Ensures the cache is generated only once
	done        atomic.Bool // Atomic flag to determine generation status
}
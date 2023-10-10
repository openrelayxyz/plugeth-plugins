package main 

import (
	"fmt"
	"sort"
	"math/big"
	"errors"

	"github.com/openrelayxyz/plugeth-utils/core"
)

type PluginConfigurator struct {
	NetworkID                 uint64   `json:"networkId"`
	ChainID                   *big.Int `json:"chainId"`                             // chainId identifies the current chain and is used for replay protection
	// SupportedProtocolVersions []uint   `json:"supportedProtocolVersions,omitempty"` // supportedProtocolVersions identifies the supported eth protocol versions for the current chain

	// HF: Homestead
	// HomesteadBlock *big.Int `json:"homesteadBlock,omitempty"` // Homestead switch block (nil = no fork, 0 = already homestead)
	// "Homestead Hard-fork Changes"
	// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-2.md
	EIP2FBlock *big.Int `json:"eip2FBlock,omitempty"`
	// DELEGATECALL
	// https://eips.ethereum.org/EIPS/eip-7
	EIP7FBlock *big.Int `json:"eip7FBlock,omitempty"`
	// Note: EIP 8 was also included in this fork, but was not backwards-incompatible

	// HF: DAO
	DAOForkBlock *big.Int `json:"daoForkBlock,omitempty"` // TheDAO hard-fork switch block (nil = no fork)
	// DAOForkSupport bool     `json:"daoForkSupport,omitempty"` // Whether the nodes supports or opposes the DAO hard-fork

	// HF: Tangerine Whistle
	// EIP150 implements the Gas price changes (https://github.com/ethereum/EIPs/issues/150)
	EIP150Block *big.Int `json:"eip150Block,omitempty"` // EIP150 HF block (nil = no fork)
	// EIP150Hash  common.Hash `json:"eip150Hash,omitempty"`  // EIP150 HF hash (needed for header only clients as only gas pricing changed)

	// HF: Spurious Dragon
	EIP155Block *big.Int `json:"eip155Block,omitempty"` // EIP155 HF block
	// EIP158Block *big.Int `json:"eip158Block,omitempty"` // EIP158 HF block, includes implementations of 158/161, 160, and 170
	//
	// EXP cost increase
	// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-160.md
	// NOTE: this json tag:
	// (a.) varies from it's 'siblings', which have 'F's in them
	// (b.) without the 'F' will vary from ETH implementations if they choose to accept the proposed changes
	// with corresponding refactoring (https://github.com/ethereum/go-ethereum/pull/18401)
	EIP160FBlock *big.Int `json:"eip160Block,omitempty"`
	// State trie clearing (== EIP158 proper)
	// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-161.md
	EIP161FBlock *big.Int `json:"eip161FBlock,omitempty"`
	// Contract code size limit
	// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-170.md
	EIP170FBlock *big.Int `json:"eip170FBlock,omitempty"`

	// HF: Byzantium
	// ByzantiumBlock *big.Int `json:"byzantiumBlock,omitempty"` // Byzantium switch block (nil = no fork, 0 = already on byzantium)

	// Difficulty adjustment to target mean block time including uncles
	// https://github.com/ethereum/EIPs/issues/100
	EIP100FBlock *big.Int `json:"eip100FBlock,omitempty"`
	// Opcode REVERT
	// https://eips.ethereum.org/EIPS/eip-140
	EIP140FBlock *big.Int `json:"eip140FBlock,omitempty"`
	// Precompiled contract for bigint_modexp
	// https://github.com/ethereum/EIPs/issues/198
	EIP198FBlock *big.Int `json:"eip198FBlock,omitempty"`
	// Opcodes RETURNDATACOPY, RETURNDATASIZE
	// https://github.com/ethereum/EIPs/issues/211
	EIP211FBlock *big.Int `json:"eip211FBlock,omitempty"`
	// Precompiled contract for pairing check
	// https://github.com/ethereum/EIPs/issues/212
	EIP212FBlock *big.Int `json:"eip212FBlock,omitempty"`
	// Precompiled contracts for addition and scalar multiplication on the elliptic curve alt_bn128
	// https://github.com/ethereum/EIPs/issues/213
	EIP213FBlock *big.Int `json:"eip213FBlock,omitempty"`
	// Opcode STATICCALL
	// https://github.com/ethereum/EIPs/issues/214
	EIP214FBlock *big.Int `json:"eip214FBlock,omitempty"`
	// Metropolis diff bomb delay and reducing block reward
	// https://github.com/ethereum/EIPs/issues/649
	// note that this is closely related to EIP100.
	// In fact, EIP100 is bundled in
	eip649FInferred bool
	EIP649FBlock    *big.Int `json:"-"`
	// Transaction receipt status
	// https://github.com/ethereum/EIPs/issues/658
	EIP658FBlock *big.Int `json:"eip658FBlock,omitempty"`
	// NOT CONFIGURABLE: prevent overwriting contracts
	// https://github.com/ethereum/EIPs/issues/684
	// EIP684FBlock *big.Int `json:"eip684BFlock,omitempty"`

	// HF: Constantinople
	// ConstantinopleBlock *big.Int `json:"constantinopleBlock,omitempty"` // Constantinople switch block (nil = no fork, 0 = already activated)
	//
	// Opcodes SHR, SHL, SAR
	// https://eips.ethereum.org/EIPS/eip-145
	EIP145FBlock *big.Int `json:"eip145FBlock,omitempty"`
	// Opcode CREATE2
	// https://eips.ethereum.org/EIPS/eip-1014
	EIP1014FBlock *big.Int `json:"eip1014FBlock,omitempty"`
	// Opcode EXTCODEHASH
	// https://eips.ethereum.org/EIPS/eip-1052
	EIP1052FBlock *big.Int `json:"eip1052FBlock,omitempty"`
	// Constantinople difficulty bomb delay and block reward adjustment
	// https://eips.ethereum.org/EIPS/eip-1234
	eip1234FInferred bool
	EIP1234FBlock    *big.Int `json:"-"`
	// Net gas metering
	// https://eips.ethereum.org/EIPS/eip-1283
	EIP1283FBlock *big.Int `json:"eip1283FBlock,omitempty"`

	PetersburgBlock *big.Int `json:"petersburgBlock,omitempty"` // Petersburg switch block (nil = same as Constantinople)

	// HF: Istanbul
	// IstanbulBlock *big.Int `json:"istanbulBlock,omitempty"` // Istanbul switch block (nil = no fork, 0 = already on istanbul)
	//
	// EIP-152: Add Blake2 compression function F precompile
	EIP152FBlock *big.Int `json:"eip152FBlock,omitempty"`
	// EIP-1108: Reduce alt_bn128 precompile gas costs
	EIP1108FBlock *big.Int `json:"eip1108FBlock,omitempty"`
	// EIP-1344: Add ChainID opcode
	EIP1344FBlock *big.Int `json:"eip1344FBlock,omitempty"`
	// EIP-1884: Repricing for trie-size-dependent opcodes
	EIP1884FBlock *big.Int `json:"eip1884FBlock,omitempty"`
	// EIP-2028: Calldata gas cost reduction
	EIP2028FBlock *big.Int `json:"eip2028FBlock,omitempty"`
	// EIP-2200: Rebalance net-metered SSTORE gas cost with consideration of SLOAD gas cost change
	// It's a combined version of EIP-1283 + EIP-1706, with a structured definition so as to make it
	// interoperable with other gas changes such as EIP-1884.
	EIP2200FBlock        *big.Int `json:"eip2200FBlock,omitempty"`
	EIP2200DisableFBlock *big.Int `json:"eip2200DisableFBlock,omitempty"`

	// EIP-2384: Difficulty Bomb Delay (Muir Glacier)
	eip2384Inferred bool
	EIP2384FBlock   *big.Int `json:"eip2384FBlock,omitempty"`

	// EIP-3554: Difficulty Bomb Delay to December 2021
	// https://eips.ethereum.org/EIPS/eip-3554
	eip3554Inferred bool
	EIP3554FBlock   *big.Int `json:"eip3554FBlock,omitempty"`

	// EIP-4345: Difficulty Bomb Delay to June 2022
	// https://eips.ethereum.org/EIPS/eip-4345
	eip4345Inferred bool
	EIP4345FBlock   *big.Int `json:"eip4345FBlock,omitempty"`

	// EIP-1706: Resolves reentrancy attack vector enabled with EIP1283.
	// https://eips.ethereum.org/EIPS/eip-1706
	EIP1706FBlock *big.Int `json:"eip1706FBlock,omitempty"`

	// https://github.com/ethereum/EIPs/pull/2537: BLS12-381 curve operations
	EIP2537FBlock *big.Int `json:"eip2537FBlock,omitempty"`

	// EWASMBlock *big.Int `json:"ewasmBlock,omitempty"` // EWASM switch block (nil = no fork, 0 = already activated)

	ECIP1010PauseBlock *big.Int `json:"ecip1010PauseBlock,omitempty"` // ECIP1010 pause HF block
	ECIP1010Length     *big.Int `json:"ecip1010Length,omitempty"`     // ECIP1010 length
	ECIP1017FBlock     *big.Int `json:"ecip1017FBlock,omitempty"`
	ECIP1017EraRounds  *big.Int `json:"ecip1017EraRounds,omitempty"` // ECIP1017 era rounds
	ECIP1080FBlock     *big.Int `json:"ecip1080FBlock,omitempty"`

	ECIP1099FBlock *big.Int `json:"ecip1099FBlock,omitempty"` // ECIP1099 etchash HF block
	ECBP1100FBlock *big.Int `json:"ecbp1100FBlock,omitempty"` // ECBP1100:MESS artificial finality

	// EIP-2315: Simple Subroutines
	// https://eips.ethereum.org/EIPS/eip-2315
	EIP2315FBlock *big.Int `json:"eip2315FBlock,omitempty"`

	// TODO: Document me.
	EIP2565FBlock *big.Int `json:"eip2565FBlock,omitempty"`

	// EIP2718FBlock is typed tx envelopes
	EIP2718FBlock *big.Int `json:"eip2718FBlock,omitempty"`

	// EIP-2929: Gas cost increases for state access opcodes
	// https://eips.ethereum.org/EIPS/eip-2929
	EIP2929FBlock *big.Int `json:"eip2929FBlock,omitempty"`

	// EIP-3198: BASEFEE opcode
	// https://eips.ethereum.org/EIPS/eip-3198
	EIP3198FBlock *big.Int `json:"eip3198FBlock,omitempty"`

	// EIP-4399: RANDOM opcode (supplanting DIFFICULTY)
	EIP4399FBlock *big.Int `json:"eip4399FBlock,omitempty"`

	// EIP-2930: Access lists.
	EIP2930FBlock *big.Int `json:"eip2930FBlock,omitempty"`

	EIP1559FBlock *big.Int `json:"eip1559FBlock,omitempty"`
	EIP3541FBlock *big.Int `json:"eip3541FBlock,omitempty"`
	EIP3529FBlock *big.Int `json:"eip3529FBlock,omitempty"`

	EIP5133FBlock   *big.Int `json:"eip5133FBlock,omitempty"`
	eip5133Inferred bool

	// Shanghai
	EIP3651FTime *uint64 `json:"eip3651FTime,omitempty"` // EIP-3651: Warm COINBASE
	EIP3855FTime *uint64 `json:"eip3855FTime,omitempty"` // EIP-3855: PUSH0 instruction
	EIP3860FTime *uint64 `json:"eip3860FTime,omitempty"` // EIP-3860: Limit and meter initcode
	EIP4895FTime *uint64 `json:"eip4895FTime,omitempty"` // EIP-4895: Beacon chain push withdrawals as operations
	EIP6049FTime *uint64 `json:"eip6049FTime,omitempty"` // EIP-6049: Deprecate SELFDESTRUCT. Note: EIP-6049 does not change the behavior of SELFDESTRUCT in and of itself, but formally announces client developers' intention of changing it in future upgrades. It is recommended that software which exposes the SELFDESTRUCT opcode to users warn them about an upcoming change in semantics.

	// Shanghai with block activations
	EIP3651FBlock *big.Int `json:"eip3651FBlock,omitempty"` // EIP-3651: Warm COINBASE
	EIP3855FBlock *big.Int `json:"eip3855FBlock,omitempty"` // EIP-3855: PUSH0 instruction
	EIP3860FBlock *big.Int `json:"eip3860FBlock,omitempty"` // EIP-3860: Limit and meter initcode
	EIP4895FBlock *big.Int `json:"eip4895FBlock,omitempty"` // EIP-4895: Beacon chain push withdrawals as operations
	EIP6049FBlock *big.Int `json:"eip6049FBlock,omitempty"` // EIP-6049: Deprecate SELFDESTRUCT. Note: EIP-6049 does not change the behavior of SELFDESTRUCT in and of itself, but formally announces client developers' intention of changing it in future upgrades. It is recommended that software which exposes the SELFDESTRUCT opcode to users warn them about an upcoming change in semantics.

	// Cancun
	EIP4844FTime *uint64 `json:"eip4844FTime,omitempty"` // EIP-4844: Shard Blob Transactions https://eips.ethereum.org/EIPS/eip-4844
	EIP1153FTime *uint64 `json:"eip1153FTime,omitempty"` // EIP-1153: Transient Storage opcodes https://eips.ethereum.org/EIPS/eip-1153
	EIP5656FTime *uint64 `json:"eip5656FTime,omitempty"` // EIP-5656: MCOPY - Memory copying instruction https://eips.ethereum.org/EIPS/eip-5656
	EIP6780FTime *uint64 `json:"eip6780FTime,omitempty"` // EIP-6780: SELFDESTRUCT only in same transaction https://eips.ethereum.org/EIPS/eip-6780

	MergeNetsplitVBlock *big.Int `json:"mergeNetsplitVBlock,omitempty"` // Virtual fork after The Merge to use as a network splitter

	DisposalBlock *big.Int `json:"disposalBlock,omitempty"` // Bomb disposal HF block

	// Various consensus engines
	// Ethash    *ctypes.EthashConfig `json:"ethash,omitempty"`
	IsDevMode bool                 `json:"isDev,omitempty"`

	TrustedCheckpoint       TrustedCheckpoint      `json:"trustedCheckpoint,omitempty"`
	TrustedCheckpointOracle *CheckpointOracleConfig `json:"trustedCheckpointOracle,omitempty"`

	DifficultyBombDelaySchedule Uint64BigMapEncodesHex `json:"difficultyBombDelays,omitempty"` // JSON tag matches Parity's
	BlockRewardSchedule         Uint64BigMapEncodesHex `json:"blockReward,omitempty"`          // JSON tag matches Parity's

	RequireBlockHashes map[uint64]core.Hash `json:"requireBlockHashes"`

}

// TrustedCheckpoint represents a set of post-processed trie roots (CHT and
// BloomTrie) associated with the appropriate section index and head hash. It is
// used to start light syncing from this checkpoint and avoid downloading the
// entire header chain while still being able to securely access old headers/logs.
type TrustedCheckpoint struct {
	SectionIndex uint64      `json:"sectionIndex"`
	SectionHead  core.Hash `json:"sectionHead"`
	CHTRoot      core.Hash `json:"chtRoot"`
	BloomRoot    core.Hash `json:"bloomRoot"`
}

// CheckpointOracleConfig represents a set of checkpoint contract(which acts as an oracle)
// config which used for light client checkpoint syncing.
type CheckpointOracleConfig struct {
	Address   core.Address   `json:"address"`
	Signers   []core.Address `json:"signers"`
	Threshold uint64           `json:"threshold"`
}

// Uint64BigMapEncodesHex is a map that encodes and decodes w/ JSON hex format.
type Uint64BigMapEncodesHex map[uint64]*big.Int

// MapMeetsSpecification returns the block number at which a difficulty/+reward map meet specifications, eg. EIP649 and/or EIP1234, or EIP2384.
// This is a reverse lookup to extract EIP-spec'd parameters from difficulty and reward maps implementations.
func MapMeetsSpecification(difficulties Uint64BigMapEncodesHex, rewards Uint64BigMapEncodesHex, difficultySum, wantedReward *big.Int) *uint64 {
	var diffN *uint64
	var sl = []uint64{}

	// difficulty
	for k := range difficulties {
		sl = append(sl, k)
	}
	sort.Slice(sl, func(i, j int) bool {
		return sl[i] < sl[j]
	})

	var total = new(big.Int)
	for _, s := range sl {
		d := difficulties[s]
		if d == nil {
			panic(fmt.Sprintf("dnil difficulties: %v, sl: %v", difficulties, sl))
		}
		total.Add(total, d)
		if total.Cmp(difficultySum) >= 0 {
			diffN = &s //nolint:gosec,exportloopref
			break
		}
	}
	if diffN == nil {
		// difficulty bomb delay not configured,
		// then does not meet eip649/eip1234 spec
		return nil
	}

	if wantedReward == nil || rewards == nil {
		return diffN
	}

	reward, ok := rewards[*diffN]
	if !ok {
		return nil
	}
	if reward.Cmp(wantedReward) != 0 {
		return nil
	}

	return diffN
}

type UnsupportedConfigErr error

var (
	ErrUnsupportedConfigNoop  UnsupportedConfigErr = errors.New("unsupported config value (noop)")
	ErrUnsupportedConfigFatal UnsupportedConfigErr = errors.New("unsupported config value (fatal)")
)

type ErrUnsupportedConfig struct {
	Err    error
	Method string
	Value  interface{}
}

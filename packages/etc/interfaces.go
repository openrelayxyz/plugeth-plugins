package main

import (
	"math/big"
)

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
package main 

import (
	"math/big"
)

type PluginConfigurator struct {
}

func (p *PluginConfigurator) IsEnabled(fn func() *uint64, n *big.Int) bool {
	f := fn()
	if f == nil || n == nil {
		return false
	}
	return big.NewInt(int64(*f)).Cmp(n) <= 0
}

func (p *PluginConfigurator) IsEnabledByTime(fn func() *uint64, n *uint64) bool {
	f := fn()
	if f == nil || n == nil {
		return false
	}
	return *f <= *n
}

func (p *PluginConfigurator) EIP161FBlock() *uint64 {
	val := big.NewInt(8772000).Uint64()
	return &val
}

func (p *PluginConfigurator) GetEthashECIP1010PauseTransition() *uint64 {
	return nil
}

// ProtocolSpecifier defines protocol interfaces that are agnostic of consensus engine.
func (p *PluginConfigurator) GetAccountStartNonce() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetMaximumExtraDataSize() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetMinGasLimit() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetGasLimitBoundDivisor() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetNetworkID() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetChainID() *big.Int {
	return nil
}
func (p *PluginConfigurator) GetSupportedProtocolVersions() []uint {
	return nil
}
func (p *PluginConfigurator) GetMaxCodeSize() *uint64 {
	return nil
}

func (p *PluginConfigurator) GetElasticityMultiplier() uint64 {
	return 2
}
func (p *PluginConfigurator) GetBaseFeeChangeDenominator() uint64 {
	return 8
}

// Be careful with EIP2.
// It is a messy EIP, specifying diverse changes, like difficulty, intrinsic gas costs for contract creation,
// txpool management, and contract OoG handling.
// It is both Ethash-specific and _not_.
func (p *PluginConfigurator) GetEIP2Transition() *uint64 {
	return nil
}

func (p *PluginConfigurator) GetEIP7Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP150Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP152Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP160Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP161abcTransition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP161dTransition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP170Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP155Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP140Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP198Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP211Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP212Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP213Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP214Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP658Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP145Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP1014Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP1052Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP1283Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP1283DisableTransition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP1108Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP2200Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP2200DisableTransition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP1344Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP1884Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP2028Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetECIP1080Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP1706Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP2537Transition() *uint64 {
	return nil
}

func (p *PluginConfigurator) GetECBP1100Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP2315Transition() *uint64 {
	return nil
}

// ModExp gas cost
func (p *PluginConfigurator) GetEIP2565Transition() *uint64 {
	return nil
}

// Gas cost increases for state access opcodes
func (p *PluginConfigurator) GetEIP2929Transition() *uint64 {
	return nil
}

// Optional access lists
func (p *PluginConfigurator) GetEIP2930Transition() *uint64 {
	return nil
}

// Typed transaction envelope
func (p *PluginConfigurator) GetEIP2718Transition() *uint64 {
	return nil
}

func (p *PluginConfigurator) GetEIP1559Transition() *uint64 {
	return nil
}

func (p *PluginConfigurator) GetEIP3541Transition() *uint64 {
	return nil
}

func (p *PluginConfigurator) GetEIP3529Transition() *uint64 {
	return nil
}

func (p *PluginConfigurator) GetEIP3198Transition() *uint64 {
	return nil
}

// EIP4399 is the RANDOM opcode.
// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-4399.md
func (p *PluginConfigurator) GetEIP4399Transition() *uint64 {
	return nil
}

// Shanghai:
//
// EIP3651: Warm COINBASE
func (p *PluginConfigurator) GetEIP3651TransitionTime() *uint64 {
	return nil
}
// EIP3855: PUSH0 instruction
func (p *PluginConfigurator) GetEIP3855TransitionTime() *uint64 {
	return nil
}
// EIP3860: Limit and meter initcode
func (p *PluginConfigurator) GetEIP3860TransitionTime() *uint64 {
	return nil
}
// EIP4895: Beacon chain push withdrawals as operations
func (p *PluginConfigurator) GetEIP4895TransitionTime() *uint64 {
	return nil
}
// EIP6049: Deprecate SELFDESTRUCT
func (p *PluginConfigurator) GetEIP6049TransitionTime() *uint64 {
	return nil
}

// Shanghai expressed as block activation numbers:
func (p *PluginConfigurator) GetEIP3651Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP3855Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP3860Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP4895Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEIP6049Transition() *uint64 {
	return nil
}

// func (p *PluginConfigurator) GetMergeVirtualTransition is a Virtual fork after The Merge to use as a network splitter
func (p *PluginConfigurator) GetMergeVirtualTransition() *uint64 {
	return nil
}

// Cancun:
// EIP4844 - Shard Blob Transactions - https://eips.ethereum.org/EIPS/eip-4844
func (p *PluginConfigurator) GetEIP4844TransitionTime() *uint64 {
	return nil
}

// EIP1153 - Transient Storage opcodes - https://eips.ethereum.org/EIPS/eip-1153
func (p *PluginConfigurator) GetEIP1153TransitionTime() *uint64 {
	return nil
}

// EIP5656 - MCOPY - Memory copying instruction - https://eips.ethereum.org/EIPS/eip-5656
func (p *PluginConfigurator) GetEIP5656TransitionTime() *uint64 {
	return nil
}

// EIP6780 - SELFDESTRUCT only in same transaction - https://eips.ethereum.org/EIPS/eip-6780
func (p *PluginConfigurator) GetEIP6780TransitionTime() *uint64 {
	return nil
}

// type EthashConfigurator interface {
func (p *PluginConfigurator) GetEthashMinimumDifficulty() *big.Int {
	return nil
}
func (p *PluginConfigurator) GetEthashDifficultyBoundDivisor() *big.Int {
	return nil
}
func (p *PluginConfigurator) GetEthashDurationLimit() *big.Int {
	return nil
}
func (p *PluginConfigurator) GetEthashHomesteadTransition() *uint64 {
	return nil
}

// func (p *PluginConfigurator) GetEthashEIP779Transition should return the block if the node wants the fork.
// Otherwise, nil should be returned.
func (p *PluginConfigurator) GetEthashEIP779Transition() *uint64 {
	return nil
} // DAO

func (p *PluginConfigurator) GetEthashEIP649Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEthashEIP1234Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEthashEIP2384Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEthashEIP3554Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEthashEIP4345Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEthashECIP1010ContinueTransition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEthashECIP1017Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEthashECIP1017EraRounds() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEthashEIP100BTransition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEthashECIP1041Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEthashECIP1099Transition() *uint64 {
	return nil
}
func (p *PluginConfigurator) GetEthashEIP5133Transition() *uint64 {
	return nil
} // Gray Glacier difficulty bomb delay

func (p *PluginConfigurator) GetEthashTerminalTotalDifficulty() *big.Int {
	return nil
}

func (p *PluginConfigurator) GetEthashTerminalTotalDifficultyPassed() bool {
	return false
}

// IsTerminalPoWBlock(parentTotalDiff *big.Int, totalDiff *big.Int) bool

func (p *PluginConfigurator) GetEthashDifficultyBombDelaySchedule() Uint64BigMapEncodesHex {
	return nil
}
func (p *PluginConfigurator) GetEthashBlockRewardSchedule() Uint64BigMapEncodesHex {
	return nil
}

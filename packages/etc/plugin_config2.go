package main 

import (
	"math/big"

	"github.com/openrelayxyz/plugeth-utils/core"
)

func newU64(u uint64) *uint64 {
	return &u
}

func bigNewU64(i *big.Int) *uint64 {
	if i == nil {
		return nil
	}
	return newU64(i.Uint64())
}

// nolint: staticcheck
func setBig(i *big.Int, u *uint64) *big.Int {
	if u == nil {
		return nil
	}
	i = big.NewInt(int64(*u))
	return i
}

func (c *PluginConfigurator) GetNetworkID() *uint64 {
	return newU64(c.NetworkID)
}

func (c *PluginConfigurator) GetChainID() *big.Int {
	return c.ChainID
}

func (c *PluginConfigurator) GetMaxCodeSize() *uint64 {
	return GlobalConfigurator().GetMaxCodeSize()
}

func (c *PluginConfigurator) GetElasticityMultiplier() uint64 {
	return GlobalConfigurator().GetElasticityMultiplier()
}


func (c *PluginConfigurator) GetBaseFeeChangeDenominator() uint64 {
	return GlobalConfigurator().GetBaseFeeChangeDenominator()
}


func (c *PluginConfigurator) GetEIP7Transition() *uint64 {
	return bigNewU64(c.EIP7FBlock)
}


func (c *PluginConfigurator) GetEIP150Transition() *uint64 {
	return bigNewU64(c.EIP150Block)
}


func (c *PluginConfigurator) GetEIP152Transition() *uint64 {
	return bigNewU64(c.EIP152FBlock)
}

func (c *PluginConfigurator) GetEIP160Transition() *uint64 {
	return bigNewU64(c.EIP160FBlock)
}

func (c *PluginConfigurator) GetEIP161dTransition() *uint64 {
	return bigNewU64(c.EIP161FBlock)
}

func (c *PluginConfigurator) GetEIP161abcTransition() *uint64 {
	return bigNewU64(c.EIP161FBlock)
}

func (c *PluginConfigurator) GetEIP170Transition() *uint64 {
	return bigNewU64(c.EIP170FBlock)
}

func (c *PluginConfigurator) GetEIP155Transition() *uint64 {
	return bigNewU64(c.EIP155Block)
}

func (c *PluginConfigurator) GetEIP140Transition() *uint64 {
	return bigNewU64(c.EIP140FBlock)
}

func (c *PluginConfigurator) GetEIP198Transition() *uint64 {
	return bigNewU64(c.EIP198FBlock)
}

func (c *PluginConfigurator) GetEIP211Transition() *uint64 {
	return bigNewU64(c.EIP211FBlock)
}

func (c *PluginConfigurator) GetEIP212Transition() *uint64 {
	return bigNewU64(c.EIP212FBlock)
}

func (c *PluginConfigurator) GetEIP213Transition() *uint64 {
	return bigNewU64(c.EIP213FBlock)
}

func (c *PluginConfigurator) GetEIP214Transition() *uint64 {
	return bigNewU64(c.EIP214FBlock)
}

func (c *PluginConfigurator) GetEIP658Transition() *uint64 {
	return bigNewU64(c.EIP658FBlock)
}

func (c *PluginConfigurator) GetEIP145Transition() *uint64 {
	return bigNewU64(c.EIP145FBlock)
}

func (c *PluginConfigurator) GetEIP1014Transition() *uint64 {
	return bigNewU64(c.EIP1014FBlock)
}

func (c *PluginConfigurator) GetEIP1052Transition() *uint64 {
	return bigNewU64(c.EIP1052FBlock)
}

func (c *PluginConfigurator) GetEIP1283Transition() *uint64 {
	return bigNewU64(c.EIP1283FBlock)
}

func (c *PluginConfigurator) GetEIP1283DisableTransition() *uint64 {
	return bigNewU64(c.PetersburgBlock)
}

func (c *PluginConfigurator) GetEIP1108Transition() *uint64 {
	return bigNewU64(c.EIP1108FBlock)
}

func (c *PluginConfigurator) GetEIP2200Transition() *uint64 {
	return bigNewU64(c.EIP2200FBlock)
}

func (c *PluginConfigurator) GetEIP2200DisableTransition() *uint64 {
	return bigNewU64(c.EIP2200DisableFBlock)
}

func (c *PluginConfigurator) GetEIP1344Transition() *uint64 {
	return bigNewU64(c.EIP1344FBlock)
}

func (c *PluginConfigurator) GetEIP1884Transition() *uint64 {
	return bigNewU64(c.EIP1884FBlock)
}

func (c *PluginConfigurator) GetEIP2028Transition() *uint64 {
	return bigNewU64(c.EIP2028FBlock)
}

func (c *PluginConfigurator) GetECIP1080Transition() *uint64 {
	return bigNewU64(c.ECIP1080FBlock)
}

func (c *PluginConfigurator) GetEIP1706Transition() *uint64 {
	return bigNewU64(c.EIP1706FBlock)
}

func (c *PluginConfigurator) GetEIP2537Transition() *uint64 {
	return bigNewU64(c.EIP2537FBlock)
}

func (c *PluginConfigurator) GetECBP1100Transition() *uint64 {
	return bigNewU64(c.ECBP1100FBlock)
}

func (c *PluginConfigurator) GetEIP2315Transition() *uint64 {
	return bigNewU64(c.EIP2315FBlock)
}

func (c *PluginConfigurator) GetEIP2929Transition() *uint64 {
	return bigNewU64(c.EIP2929FBlock)
}

func (c *PluginConfigurator) GetEIP2930Transition() *uint64 {
	return bigNewU64(c.EIP2930FBlock)
}

func (c *PluginConfigurator) GetEIP2565Transition() *uint64 {
	return bigNewU64(c.EIP2565FBlock)
}

func (c *PluginConfigurator) GetEIP2718Transition() *uint64 {
	return bigNewU64(c.EIP2718FBlock)
}

func (c *PluginConfigurator) GetEIP1559Transition() *uint64 {
	return bigNewU64(c.EIP1559FBlock)
}

func (c *PluginConfigurator) GetEIP3541Transition() *uint64 {
	return bigNewU64(c.EIP3541FBlock)
}

func (c *PluginConfigurator) GetEIP3529Transition() *uint64 {
	return bigNewU64(c.EIP3529FBlock)
}

func (c *PluginConfigurator) GetEIP3198Transition() *uint64 {
	return bigNewU64(c.EIP3198FBlock)
}

func (c *PluginConfigurator) GetEIP4399Transition() *uint64 {
	return bigNewU64(c.EIP4399FBlock)
}

// EIP3651: Warm COINBASE
func (c *PluginConfigurator) GetEIP3651TransitionTime() *uint64 {
	return c.EIP3651FTime
}

// GetEIP3855TransitionTime EIP3855: PUSH0 instruction
func (c *PluginConfigurator) GetEIP3855TransitionTime() *uint64 {
	return c.EIP3855FTime
}

// GetEIP3860TransitionTime EIP3860: Limit and meter initcode
func (c *PluginConfigurator) GetEIP3860TransitionTime() *uint64 {
	return c.EIP3860FTime
}

// GetEIP4895TransitionTime EIP4895: Beacon chain push withdrawals as operations
func (c *PluginConfigurator) GetEIP4895TransitionTime() *uint64 {
	return c.EIP4895FTime
}

// GetEIP6049TransitionTime EIP6049: Deprecate SELFDESTRUCT
func (c *PluginConfigurator) GetEIP6049TransitionTime() *uint64 {
	return c.EIP6049FTime
}

// Shanghai by block
// EIP3651: Warm COINBASE
func (c *PluginConfigurator) GetEIP3651Transition() *uint64 {
	return bigNewU64(c.EIP3651FBlock)
}

// GetEIP3855Transition EIP3855: PUSH0 instruction
func (c *PluginConfigurator) GetEIP3855Transition() *uint64 {
	return bigNewU64(c.EIP3855FBlock)
}

// GetEIP3860Transition EIP3860: Limit and meter initcode
func (c *PluginConfigurator) GetEIP3860Transition() *uint64 {
	return bigNewU64(c.EIP3860FBlock)
}

// GetEIP4895Transition EIP4895: Beacon chain push withdrawals as operations
func (c *PluginConfigurator) GetEIP4895Transition() *uint64 {
	return bigNewU64(c.EIP4895FBlock)
}

// GetEIP6049Transition EIP6049: Deprecate SELFDESTRUCT
func (c *PluginConfigurator) GetEIP6049Transition() *uint64 {
	return bigNewU64(c.EIP6049FBlock)
}

// GetEIP4844TransitionTime EIP4844: Shard Blob Transactions
func (c *PluginConfigurator) GetEIP4844TransitionTime() *uint64 {
	return c.EIP4844FTime
}

// GetEIP1153TransitionTime EIP1153: Transient Storage opcodes
func (c *PluginConfigurator) GetEIP1153TransitionTime() *uint64 {
	return c.EIP1153FTime
}

// GetEIP5656TransitionTime EIP5656: MCOPY - Memory copying instruction
func (c *PluginConfigurator) GetEIP5656TransitionTime() *uint64 {
	return c.EIP5656FTime
}

// GetEIP6780TransitionTime EIP6780: SELFDESTRUCT only in same transaction
func (c *PluginConfigurator) GetEIP6780TransitionTime() *uint64 {
	return c.EIP6780FTime
}

func (c *PluginConfigurator) GetMergeVirtualTransition() *uint64 {
	return bigNewU64(c.MergeNetsplitVBlock)
}

func (c *PluginConfigurator) IsEnabled(fn func() *uint64, n *big.Int) bool {
	f := fn()
	if f == nil || n == nil {
		return false
	}
	return big.NewInt(int64(*f)).Cmp(n) <= 0
}

func (c *PluginConfigurator) IsEnabledByTime(fn func() *uint64, n *uint64) bool {
	f := fn()
	if f == nil || n == nil {
		return false
	}
	return *f <= *n
}

func (c *PluginConfigurator) GetForkCanonHash(n uint64) core.Hash {
	if c.RequireBlockHashes == nil {
		return core.Hash{}
	}
	for k, v := range c.RequireBlockHashes {
		if k == n {
			return v
		}
	}
	return core.Hash{}
}

func (c *PluginConfigurator) GetForkCanonHashes() map[uint64]core.Hash {
	return c.RequireBlockHashes
}

// func (c *PluginConfigurator) GetConsensusEngineType() ConsensusEngineT {
// 	if c.Ethash != nil {
// 		return ConsensusEngineT_Ethash
// 	}
// 	return ConsensusEngineT_Unknown
// }

func (c *PluginConfigurator) GetIsDevMode() bool {
	return c.IsDevMode
}

func (c *PluginConfigurator) GetEthashMinimumDifficulty() *big.Int {
	
	return GlobalConfigurator().GetEthashMinimumDifficulty()
}

func (c *PluginConfigurator) GetEthashDifficultyBoundDivisor() *big.Int {
	
	return GlobalConfigurator().GetEthashDifficultyBoundDivisor()
}

func (c *PluginConfigurator) GetEthashDurationLimit() *big.Int {
	
	return GlobalConfigurator().GetEthashDurationLimit()
}

func (c *PluginConfigurator) GetEthashHomesteadTransition() *uint64 {
	
	if c.EIP2FBlock == nil || c.EIP7FBlock == nil {
		return nil
	}
	return bigNewU64(BigMax(c.EIP2FBlock, c.EIP7FBlock))
}

func (c *PluginConfigurator) GetEIP2Transition() *uint64 {
	return bigNewU64(c.EIP2FBlock)
}

func (c *PluginConfigurator) GetEthashEIP779Transition() *uint64 {
	
	return bigNewU64(c.DAOForkBlock)
}

func (c *PluginConfigurator) GetEthashEIP649Transition() *uint64 {
	
	if c.eip649FInferred {
		return bigNewU64(c.EIP649FBlock)
	}

	var diffN *uint64
	defer func() {
		c.EIP649FBlock = setBig(c.EIP649FBlock, diffN)
		c.eip649FInferred = true
	}()

	// Get block number (key) from maps where EIP649 criteria is met.
	diffN = MapMeetsSpecification(
		c.DifficultyBombDelaySchedule,
		c.BlockRewardSchedule,
		EIP649DifficultyBombDelay,
		EIP649FBlockReward,
	)
	if diffN == nil {
		diffN = c.GetEthashEIP1234Transition()
	}
	return diffN
}

func (c *PluginConfigurator) GetEthashEIP1234Transition() *uint64 {
	
	if c.eip1234FInferred {
		return bigNewU64(c.EIP1234FBlock)
	}

	var diffN *uint64
	defer func() {
		c.EIP1234FBlock = setBig(c.EIP1234FBlock, diffN)
		c.eip1234FInferred = true
	}()

	// Get block number (key) from maps where EIP1234 criteria is met.
	diffN = MapMeetsSpecification(
		c.DifficultyBombDelaySchedule,
		c.BlockRewardSchedule,
		EIP1234DifficultyBombDelay,
		EIP1234FBlockReward,
	)
	return diffN
}

func (c *PluginConfigurator) GetEthashEIP2384Transition() *uint64 {
	
	if c.eip2384Inferred {
		return bigNewU64(c.EIP2384FBlock)
	}

	var diffN *uint64
	defer func() {
		c.EIP2384FBlock = setBig(c.EIP2384FBlock, diffN)
		c.eip2384Inferred = true
	}()

	// Get block number (key) from map where EIP2384 criteria is met.
	diffN = MapMeetsSpecification(c.DifficultyBombDelaySchedule, nil, EIP2384DifficultyBombDelay, nil)
	return diffN
}

func (c *PluginConfigurator) GetEthashEIP3554Transition() *uint64 {
	
	if c.eip3554Inferred {
		return bigNewU64(c.EIP3554FBlock)
	}

	var diffN *uint64
	defer func() {
		c.EIP3554FBlock = setBig(c.EIP3554FBlock, diffN)
		c.eip3554Inferred = true
	}()

	// Get block number (key) from map where EIP3554 criteria is met.
	diffN = MapMeetsSpecification(c.DifficultyBombDelaySchedule, nil, EIP3554DifficultyBombDelay, nil)
	return diffN
}

func (c *PluginConfigurator) GetEthashEIP4345Transition() *uint64 {
	
	if c.eip4345Inferred {
		return bigNewU64(c.EIP4345FBlock)
	}

	var diffN *uint64
	defer func() {
		c.EIP4345FBlock = setBig(c.EIP4345FBlock, diffN)
		c.eip4345Inferred = true
	}()

	// Get block number (key) from map where EIP4345 criteria is met.
	diffN = MapMeetsSpecification(c.DifficultyBombDelaySchedule, nil, EIP4345DifficultyBombDelay, nil)
	return diffN
}

func (c *PluginConfigurator) GetEthashECIP1010PauseTransition() *uint64 {
	
	return bigNewU64(c.ECIP1010PauseBlock)
}

func (c *PluginConfigurator) GetEthashECIP1010ContinueTransition() *uint64 {
	
	if c.ECIP1010PauseBlock == nil {
		return nil
	}
	if c.ECIP1010Length == nil {
		return nil
	}
	// transition = pause + length
	return bigNewU64(new(big.Int).Add(c.ECIP1010PauseBlock, c.ECIP1010Length))
}

func (c *PluginConfigurator) GetEthashECIP1017Transition() *uint64 {
	
	return bigNewU64(c.ECIP1017FBlock)
}

func (c *PluginConfigurator) GetEthashECIP1017EraRounds() *uint64 {
	
	return bigNewU64(c.ECIP1017EraRounds)
}

func (c *PluginConfigurator) GetEthashEIP100BTransition() *uint64 {
	
	return bigNewU64(c.EIP100FBlock)
}

func (c *PluginConfigurator) GetEthashECIP1041Transition() *uint64 {
	
	return bigNewU64(c.DisposalBlock)
}

func (c *PluginConfigurator) GetEthashECIP1099Transition() *uint64 {
	
	return bigNewU64(c.ECIP1099FBlock)
}

func (c *PluginConfigurator) GetEthashEIP5133Transition() *uint64 {
	
	if c.eip5133Inferred {
		return bigNewU64(c.EIP5133FBlock)
	}

	var diffN *uint64
	defer func() {
		c.EIP5133FBlock = setBig(c.EIP5133FBlock, diffN)
		c.eip5133Inferred = true
	}()

	// Get block number (key) from map where EIP5133 criteria is met.
	diffN = MapMeetsSpecification(c.DifficultyBombDelaySchedule, nil, EIP5133DifficultyBombDelay, nil)
	return diffN
}

func (c *PluginConfigurator) GetEthashDifficultyBombDelaySchedule() Uint64BigMapEncodesHex {
	
	return c.DifficultyBombDelaySchedule
}

func (c *PluginConfigurator) GetEthashBlockRewardSchedule() Uint64BigMapEncodesHex {
	
	return c.BlockRewardSchedule
}

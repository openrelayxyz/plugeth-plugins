package main 

import (
	"math/big"
)

type GlobalVarsConfigurator struct {
}

var gc = &GlobalVarsConfigurator{}

func GlobalConfigurator() *GlobalVarsConfigurator {
	return gc
}

func (_ GlobalVarsConfigurator) GetAccountStartNonce() *uint64 {
	return newU64(0)
}

func (_ GlobalVarsConfigurator) SetAccountStartNonce(n *uint64) error {
	if n == nil {
		return nil
	}
	if *n != 0 {
		return ErrUnsupportedConfigFatal
	}
	return nil
}

func (_ GlobalVarsConfigurator) GetMaximumExtraDataSize() *uint64 {
	return newU64(MaximumExtraDataSize)
}

func (_ GlobalVarsConfigurator) SetMaximumExtraDataSize(n *uint64) error {
	MaximumExtraDataSize = *n
	return nil
}

func (_ GlobalVarsConfigurator) GetMinGasLimit() *uint64 {
	return newU64(MinGasLimit)
}

func (_ GlobalVarsConfigurator) SetMinGasLimit(n *uint64) error {
	MinGasLimit = *n
	return nil
}

func (_ GlobalVarsConfigurator) GetGasLimitBoundDivisor() *uint64 {
	return newU64(GasLimitBoundDivisor)
}

func (_ GlobalVarsConfigurator) SetGasLimitBoundDivisor(n *uint64) error {
	GasLimitBoundDivisor = *n
	return nil
}

func (_ GlobalVarsConfigurator) GetMaxCodeSize() *uint64 {
	return newU64(MaxCodeSize)
}

func (_ GlobalVarsConfigurator) SetMaxCodeSize(n *uint64) error {
	if n == nil {
		return nil
	}
	MaxCodeSize = *n
	return nil
}

func (_ GlobalVarsConfigurator) GetElasticityMultiplier() uint64 {
	return DefaultElasticityMultiplier
}

func (_ GlobalVarsConfigurator) SetElasticityMultiplier(n uint64) error {
	// Noop.
	return nil
}

func (_ GlobalVarsConfigurator) GetBaseFeeChangeDenominator() uint64 {
	return DefaultBaseFeeChangeDenominator
}

func (_ GlobalVarsConfigurator) SetBaseFeeChangeDenominator(n uint64) error {
	// Noop.
	return nil
}

func (_ GlobalVarsConfigurator) GetEthashMinimumDifficulty() *big.Int {
	return MinimumDifficulty
}
func (_ GlobalVarsConfigurator) SetEthashMinimumDifficulty(i *big.Int) error {
	if i == nil {
		return ErrUnsupportedConfigFatal
	}
	MinimumDifficulty = i
	return nil
}

func (_ GlobalVarsConfigurator) GetEthashDifficultyBoundDivisor() *big.Int {
	return DifficultyBoundDivisor
}

func (_ GlobalVarsConfigurator) SetEthashDifficultyBoundDivisor(i *big.Int) error {
	if i == nil {
		return ErrUnsupportedConfigFatal
	}
	DifficultyBoundDivisor = i
	return nil
}

func (_ GlobalVarsConfigurator) GetEthashDurationLimit() *big.Int {
	return DurationLimit
}

func (_ GlobalVarsConfigurator) SetEthashDurationLimit(i *big.Int) error {
	if i == nil {
		return ErrUnsupportedConfigFatal
	}
	DurationLimit = i
	return nil
}
package main 

import (
	"fmt"
	"errors"
	"math/big"
	"bytes"
	"encoding/hex"

	"github.com/openrelayxyz/plugeth-utils/restricted/types"
)

// BigMax returns the larger of x or y.
func BigMax(x, y *big.Int) *big.Int {
	if x.Cmp(y) < 0 {
		return y
	}
	return x
}

// parent_diff_over_dbd is a  convenience fn for CalcDifficulty
func parent_diff_over_dbd(p *types.Header) *big.Int {
	return new(big.Int).Div(p.Difficulty, DifficultyBoundDivisor)
}

// parent_time_delta is a convenience fn for CalcDifficulty
func parent_time_delta(t uint64, p *types.Header) *big.Int {
	return new(big.Int).Sub(new(big.Int).SetUint64(t), new(big.Int).SetUint64(p.Time))
}

// VerifyDAOHeaderExtraData validates the extra-data field of a block header to
// ensure it conforms to DAO hard-fork rules.
//
// DAO hard-fork extension to the header validity:
//
//   - if the node is no-fork, do not accept blocks in the [fork, fork+10) range
//     with the fork specific extra-data set.
//   - if the node is pro-fork, require blocks in the specific range to have the
//     unique extra-data set.
func VerifyDAOHeaderExtraData(config ChainConfigurator, header *types.Header) error {
	// If the config wants the DAO fork, it should validate the extra data.
	// Otherwise, like any other block or any other config, it should not care.
	daoForkBlock := config.GetEthashEIP779Transition()
	if daoForkBlock == nil {
		return nil
	}
	daoForkBlockB := new(big.Int).SetUint64(*daoForkBlock)
	// Make sure the block is within the fork's modified extra-data range
	limit := new(big.Int).Add(daoForkBlockB, DAOForkExtraRange)
	if header.Number.Cmp(daoForkBlockB) < 0 || header.Number.Cmp(limit) >= 0 {
		return nil
	}
	if !bytes.Equal(header.Extra, DAOForkBlockExtra) {
		return ErrBadProDAOExtra
	}
	return nil

	// Leaving the "old" code in as dead commented code for reference.
	//
	// // Short circuit validation if the node doesn't care about the DAO fork
	// daoForkBlock := config.GetEthashEIP779Transition()
	// // Second clause catches test configs with nil fork blocks (maybe set dynamically or
	// // testing agnostic of chain config).
	// if daoForkBlock == nil && !generic.AsGenericCC(config).DAOSupport() {
	//	return nil
	// }
	//
	// if daoForkBlock == nil {
	//
	// }
	//
	// daoForkBlockB := new(big.Int).SetUint64(*daoForkBlock)
	//
	// // Make sure the block is within the fork's modified extra-data range
	// limit := new(big.Int).Add(daoForkBlockB, vars.DAOForkExtraRange)
	// if header.Number.Cmp(daoForkBlockB) < 0 || header.Number.Cmp(limit) >= 0 {
	//	return nil
	// }
	// // Depending on whether we support or oppose the fork, validate the extra-data contents
	// if generic.AsGenericCC(config).DAOSupport() {
	//	if !bytes.Equal(header.Extra, vars.DAOForkBlockExtra) {
	//		return ErrBadProDAOExtra
	//	}
	// } else {
	//	if bytes.Equal(header.Extra, vars.DAOForkBlockExtra) {
	//		return ErrBadNoDAOExtra
	//	}
	// }
	// // All ok, header has the same extra-data we expect
	// return nil
}

// FromHex returns the bytes represented by the hexadecimal string s.
// s may be prefixed with "0x".
func FromHex(s string) []byte {
	if has0xPrefix(s) {
		s = s[2:]
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return Hex2Bytes(s)
}

// Hex2Bytes returns the bytes represented by the hexadecimal string str.
func Hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}

// has0xPrefix validates str begins with '0x' or '0X'.
func has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns
// the difficulty that a new block should have when created at time
// given the parent block's time and difficulty.
func CalcDifficulty(config ChainConfigurator, time uint64, parent *types.Header) *big.Int {
	next := new(big.Int).Add(parent.Number, big1)
	out := new(big.Int)

	// TODO (meowbits): do we need this?
	// if config.IsEnabled(config.GetEthashTerminalTotalDifficulty, next) {
	// 	return big.NewInt(1)
	// }

	// ADJUSTMENT algorithms
	if config.IsEnabled(config.GetEthashEIP100BTransition, next) {
		// https://github.com/ethereum/EIPs/issues/100
		// algorithm:
		// diff = (parent_diff +
		//         (parent_diff / 2048 * max((2 if len(parent.uncles) else 1) - ((timestamp - parent.timestamp) // 9), -99))
		//        ) + 2^(periodCount - 2)
		out.Div(parent_time_delta(time, parent), EIP100FDifficultyIncrementDivisor)

		if parent.UncleHash == types.EmptyUncleHash {
			out.Sub(big1, out)
		} else {
			out.Sub(big2, out)
		}
		out.Set(BigMax(out, bigMinus99))
		out.Mul(parent_diff_over_dbd(parent), out)
		out.Add(out, parent.Difficulty)
	} else if config.IsEnabled(config.GetEIP2Transition, next) {
		// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-2.md
		// algorithm:
		// diff = (parent_diff +
		//         (parent_diff / 2048 * max(1 - (block_timestamp - parent_timestamp) // 10, -99))
		//        )
		out.Div(parent_time_delta(time, parent), EIP2DifficultyIncrementDivisor)
		out.Sub(big1, out)
		out.Set(BigMax(out, bigMinus99))
		out.Mul(parent_diff_over_dbd(parent), out)
		out.Add(out, parent.Difficulty)
	} else {
		// FRONTIER
		// algorithm:
		// diff =
		//   if parent_block_time_delta < params.DurationLimit
		//      parent_diff + (parent_diff // 2048)
		//   else
		//      parent_diff - (parent_diff // 2048)
		out.Set(parent.Difficulty)
		if parent_time_delta(time, parent).Cmp(DurationLimit) < 0 {
			out.Add(out, parent_diff_over_dbd(parent))
		} else {
			out.Sub(out, parent_diff_over_dbd(parent))
		}
	}

	// after adjustment and before bomb
	out.Set(BigMax(out, MinimumDifficulty))

	if config.IsEnabled(config.GetEthashECIP1041Transition, next) {
		return out
	}

	// EXPLOSION delays

	// exPeriodRef the explosion clause's reference point
	exPeriodRef := new(big.Int).Add(parent.Number, big1)

	if config.IsEnabled(config.GetEthashECIP1010PauseTransition, next) {
		ecip1010Explosion(config, next, exPeriodRef)
	} else if len(config.GetEthashDifficultyBombDelaySchedule()) > 0 {
		// This logic varies from the original fork-based logic (below) in that
		// configured delay values are treated as compounding values (-2000000 + -3000000 = -5000000@constantinople)
		// as opposed to hardcoded pre-compounded values (-5000000@constantinople).
		// Thus the Sub-ing.
		fakeBlockNumber := new(big.Int).Set(exPeriodRef)
		for activated, dur := range config.GetEthashDifficultyBombDelaySchedule() {
			if exPeriodRef.Cmp(big.NewInt(int64(activated))) < 0 {
				continue
			}
			fakeBlockNumber.Sub(fakeBlockNumber, dur)
		}
		exPeriodRef.Set(fakeBlockNumber)
	} else if config.IsEnabled(config.GetEthashEIP5133Transition, next) {
		// calcDifficultyEip4345 is the difficulty adjustment algorithm as specified by EIP 4345.
		// It offsets the bomb a total of 10.7M blocks.
		fakeBlockNumber := new(big.Int)
		delayWithOffset := new(big.Int).Sub(EIP5133DifficultyBombDelay, big1)
		if parent.Number.Cmp(delayWithOffset) >= 0 {
			fakeBlockNumber = fakeBlockNumber.Sub(parent.Number, delayWithOffset)
		}
		exPeriodRef.Set(fakeBlockNumber)
	} else if config.IsEnabled(config.GetEthashEIP4345Transition, next) {
		// calcDifficultyEip4345 is the difficulty adjustment algorithm as specified by EIP 4345.
		// It offsets the bomb a total of 10.7M blocks.
		fakeBlockNumber := new(big.Int)
		delayWithOffset := new(big.Int).Sub(EIP4345DifficultyBombDelay, big1)
		if parent.Number.Cmp(delayWithOffset) >= 0 {
			fakeBlockNumber = fakeBlockNumber.Sub(parent.Number, delayWithOffset)
		}
		exPeriodRef.Set(fakeBlockNumber)
	} else if config.IsEnabled(config.GetEthashEIP3554Transition, next) {
		// calcDifficultyEIP3554 is the difficulty adjustment algorithm for London (December 2021).
		// The calculation uses the Byzantium rules, but with bomb offset 9.7M.
		fakeBlockNumber := new(big.Int)
		delayWithOffset := new(big.Int).Sub(EIP3554DifficultyBombDelay, big1)
		if parent.Number.Cmp(delayWithOffset) >= 0 {
			fakeBlockNumber = fakeBlockNumber.Sub(parent.Number, delayWithOffset)
		}
		exPeriodRef.Set(fakeBlockNumber)
	} else if config.IsEnabled(config.GetEthashEIP2384Transition, next) {
		// calcDifficultyEIP2384 is the difficulty adjustment algorithm for Muir Glacier.
		// The calculation uses the Byzantium rules, but with bomb offset 9M.
		fakeBlockNumber := new(big.Int)
		delayWithOffset := new(big.Int).Sub(EIP2384DifficultyBombDelay, big1)
		if parent.Number.Cmp(delayWithOffset) >= 0 {
			fakeBlockNumber = fakeBlockNumber.Sub(parent.Number, delayWithOffset)
		}
		exPeriodRef.Set(fakeBlockNumber)
	} else if config.IsEnabled(config.GetEthashEIP1234Transition, next) {
		// calcDifficultyEIP1234 is the difficulty adjustment algorithm for Constantinople.
		// The calculation uses the Byzantium rules, but with bomb offset 5M.
		// Specification EIP-1234: https://eips.ethereum.org/EIPS/eip-1234
		// Note, the calculations below looks at the parent number, which is 1 below
		// the block number. Thus we remove one from the delay given

		// calculate a fake block number for the ice-age delay
		// Specification: https://eips.ethereum.org/EIPS/eip-1234
		fakeBlockNumber := new(big.Int)
		delayWithOffset := new(big.Int).Sub(EIP1234DifficultyBombDelay, big1)
		if parent.Number.Cmp(delayWithOffset) >= 0 {
			fakeBlockNumber = fakeBlockNumber.Sub(parent.Number, delayWithOffset)
		}
		exPeriodRef.Set(fakeBlockNumber)
	} else if config.IsEnabled(config.GetEthashEIP649Transition, next) {
		// The calculation uses the Byzantium rules, with bomb offset of 3M.
		// Specification EIP-649: https://eips.ethereum.org/EIPS/eip-649
		// Related meta-ish EIP-669: https://github.com/ethereum/EIPs/pull/669
		// Note, the calculations below looks at the parent number, which is 1 below
		// the block number. Thus we remove one from the delay given

		fakeBlockNumber := new(big.Int)
		delayWithOffset := new(big.Int).Sub(EIP649DifficultyBombDelay, big1)
		if parent.Number.Cmp(delayWithOffset) >= 0 {
			fakeBlockNumber = fakeBlockNumber.Sub(parent.Number, delayWithOffset)
		}
		exPeriodRef.Set(fakeBlockNumber)
	}

	// EXPLOSION

	// the 'periodRef' (from above) represents the many ways of hackishly modifying the reference number
	// (ie the 'currentBlock') in order to lie to the function about what time it really is
	//
	//   2^(( periodRef // EDP) - 2)
	//
	x := new(big.Int)
	x.Div(exPeriodRef, ExpDiffPeriod) // (periodRef // EDP)
	if x.Cmp(big1) > 0 {                     // if result large enough (not in algo explicitly)
		x.Sub(x, big2)      // - 2
		x.Exp(big2, x, nil) // 2^
	} else {
		x.SetUint64(0)
	}
	out.Add(out, x)
	return out
}

// VerifyGaslimit verifies the header gas limit according increase/decrease
// in relation to the parent gas limit.
func VerifyGaslimit(parentGasLimit, headerGasLimit uint64) error {
	// Verify that the gas limit remains within allowed bounds
	diff := int64(parentGasLimit) - int64(headerGasLimit)
	if diff < 0 {
		diff *= -1
	}
	limit := parentGasLimit / GasLimitBoundDivisor
	if uint64(diff) >= limit {
		return fmt.Errorf("invalid gas limit: have %d, want %d +-= %d", headerGasLimit, parentGasLimit, limit-1)
	}
	if headerGasLimit < MinGasLimit {
		return errors.New("invalid gas limit below 5000")
	}
	return nil
}

func ecip1010Explosion(config ChainConfigurator, next *big.Int, exPeriodRef *big.Int) {
	// https://github.com/ethereumproject/ECIPs/blob/master/ECIPs/ECIP-1010.md

	if next.Uint64() < *config.GetEthashECIP1010ContinueTransition() {
		exPeriodRef.SetUint64(*config.GetEthashECIP1010PauseTransition())
	} else {
		length := new(big.Int).SetUint64(*config.GetEthashECIP1010ContinueTransition() - *config.GetEthashECIP1010PauseTransition())
		exPeriodRef.Sub(exPeriodRef, length)
	}
}
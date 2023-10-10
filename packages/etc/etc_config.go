package main

import (
	"math/big"

	"github.com/openrelayxyz/plugeth-utils/core"
)

var etc_config = &PluginConfigurator {

	NetworkID:                 1,
		// Ethash:                    new(ctypes.EthashConfig),
		ChainID:                   big.NewInt(61),
		// SupportedProtocolVersions: vars.DefaultProtocolVersions,

		EIP2FBlock: big.NewInt(1150000),
		EIP7FBlock: big.NewInt(1150000),

		// DAOForkBlock:        big.NewInt(1920000),

		EIP150Block: big.NewInt(2500000),

		EIP155Block:        big.NewInt(3000000),
		EIP160FBlock:       big.NewInt(3000000),
		ECIP1010PauseBlock: big.NewInt(3000000),
		ECIP1010Length:     big.NewInt(2000000),

		ECIP1017FBlock:    big.NewInt(5000000),
		ECIP1017EraRounds: big.NewInt(5000000),

		DisposalBlock: big.NewInt(5900000),

		// EIP158~
		EIP161FBlock: big.NewInt(8772000),
		EIP170FBlock: big.NewInt(8772000),

		// Byzantium eq
		EIP100FBlock: big.NewInt(8772000),
		EIP140FBlock: big.NewInt(8772000),
		EIP198FBlock: big.NewInt(8772000),
		EIP211FBlock: big.NewInt(8772000),
		EIP212FBlock: big.NewInt(8772000),
		EIP213FBlock: big.NewInt(8772000),
		EIP214FBlock: big.NewInt(8772000),
		EIP658FBlock: big.NewInt(8772000),

		// Constantinople eq, aka Agharta
		EIP145FBlock:  big.NewInt(9573000),
		EIP1014FBlock: big.NewInt(9573000),
		EIP1052FBlock: big.NewInt(9573000),
		// EIP1283FBlock:   big.NewInt(9573000),
		// PetersburgBlock: big.NewInt(9573000),

		// Istanbul eq, aka Phoenix
		// ECIP-1088
		EIP152FBlock:  big.NewInt(10_500_839),
		EIP1108FBlock: big.NewInt(10_500_839),
		EIP1344FBlock: big.NewInt(10_500_839),
		EIP1884FBlock: big.NewInt(10_500_839),
		EIP2028FBlock: big.NewInt(10_500_839),
		EIP2200FBlock: big.NewInt(10_500_839), // RePetersburg (=~ re-1283)

		ECBP1100FBlock: big.NewInt(11_380_000), // ETA 09 Oct 2020
		ECIP1099FBlock: big.NewInt(11_700_000), // Etchash (DAG size limit)

		// Berlin eq, aka Magneto
		EIP2565FBlock: big.NewInt(13_189_133),
		EIP2718FBlock: big.NewInt(13_189_133),
		EIP2929FBlock: big.NewInt(13_189_133),
		EIP2930FBlock: big.NewInt(13_189_133),

		// London (partially), aka Mystique
		EIP3529FBlock: big.NewInt(14_525_000),
		EIP3541FBlock: big.NewInt(14_525_000),

		// Spiral, aka Shanghai (partially)
		// EIP4399FBlock: nil, // Supplant DIFFICULTY with PREVRANDAO. ETC does not spec 4399 because it's still PoW, and 4399 is only applicable for the PoS system.
		EIP3651FBlock: nil, // Warm COINBASE (gas reprice)
		EIP3855FBlock: nil, // PUSH0 instruction
		EIP3860FBlock: nil, // Limit and meter initcode
		// EIP4895FBlock: nil, // Beacon chain push withdrawals as operations
		EIP6049FBlock: nil, // Deprecate SELFDESTRUCT (noop)

		RequireBlockHashes: map[uint64]core.Hash{
			1920000: core.HexToHash("0x94365e3a8c0b35089c1d1195081fe7489b528a84b22199c916180db8b28ade7f"),
			2500000: core.HexToHash("0xca12c63534f565899681965528d536c52cb05b7c48e269c2a6cb77ad864d878a"),
		},

}
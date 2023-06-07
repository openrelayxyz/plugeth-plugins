module github.com/openrelayxyz/plugeth-plugins

go 1.19

require (
	github.com/hashicorp/golang-lru v0.5.5-0.20210104140557-80c98217689d
	github.com/holiman/uint256 v1.2.0
	github.com/inconshreveable/log15 v2.16.0+incompatible
	github.com/openrelayxyz/cardinal-rpc v1.1.0
	github.com/openrelayxyz/plugeth-utils v0.0.24
)

require (
	github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/mattn/go-colorable v0.1.0 // indirect
	github.com/mattn/go-isatty v0.0.5-0.20180830101745-3fb116b82035 // indirect
	github.com/openrelayxyz/cardinal-types v1.0.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	golang.org/x/crypto v0.1.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/term v0.1.0 // indirect
)

replace github.com/openrelayxyz/plugeth-utils => /home/philip/src/rivet/plugeth_superspace/plugeth-utils

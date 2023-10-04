module github.com/openrelayxyz/plugeth-plugins

go 1.19

require (
	github.com/deckarep/golang-set/v2 v2.3.1
	github.com/edsrzf/mmap-go v1.1.0
	github.com/hashicorp/golang-lru v0.5.5-0.20210104140557-80c98217689d
	github.com/holiman/uint256 v1.2.3
	github.com/openrelayxyz/plugeth-utils v1.3.0
	golang.org/x/crypto v0.12.0
	gonum.org/v1/gonum v0.14.0
)

require (
	github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	golang.org/x/exp v0.0.0-20230321023759-10a507213a29 // indirect
	golang.org/x/sys v0.11.0 // indirect
)

replace github.com/openrelayxyz/plugeth-utils => /home/philip/src/rivet/plugeth_superspace/plugeth-utils

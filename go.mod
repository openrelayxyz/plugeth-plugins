module github.com/openrelayxyz/plugeth-plugins

go 1.18

require (
	github.com/hashicorp/golang-lru v0.5.5-0.20210104140557-80c98217689d
	github.com/holiman/uint256 v1.2.0
	github.com/openrelayxyz/plugeth-utils v0.0.20
)

require (
	github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
)

replace github.com/openrelayxyz/plugeth-utils v0.0.20 => /home/philip/src/rivet/plugeth-utils/

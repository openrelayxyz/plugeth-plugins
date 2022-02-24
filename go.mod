module github.com/openrelayxyz/plugeth-plugins

go 1.16

require (
	github.com/hashicorp/golang-lru v0.5.5-0.20210104140557-80c98217689d
	github.com/holiman/uint256 v1.2.0
	github.com/openrelayxyz/plugeth-utils v0.0.15
	gopkg.in/urfave/cli.v1 v1.20.0
)

replace github.com/openrelayxyz/plugeth-utils v0.0.14 => ../plugeth-utils

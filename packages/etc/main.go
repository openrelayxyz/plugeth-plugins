package main

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/openrelayxyz/plugeth-utils/core"
)

var (
	ClassicBootnodes = []string{

		"enode://942bf2f0754972391467765be1d98206926fc8ad0be8a49cd65e1730420c37fa63355bddb0ae5faa1d3505a2edcf8fad1cf00f3c179e244f047ec3a3ba5dacd7@176.9.51.216:30355", // @q9f ceibo
		"enode://0b0e09d6756b672ac6a8b70895da4fb25090b939578935d4a897497ffaa205e019e068e1ae24ac10d52fa9b8ddb82840d5d990534201a4ad859ee12cb5c91e82@176.9.51.216:30365", // @q9f ceibo

		"enode://b9e893ea9cb4537f4fed154233005ae61b441cd0ecd980136138c304fefac194c25a16b73dac05fc66a4198d0c15dd0f33af99b411882c68a019dfa6bb703b9d@18.130.93.66:30303",
	}

	dnsPrefixETC string = "enrtree://AJE62Q4DUX4QMMXEHCSSCSC65TDHZYSMONSD64P3WULVLSF6MRQ3K@"

	ClassicDNSNetwork1 string = dnsPrefixETC + "all.classic.blockd.info"

	snapDiscoveryURLs []string
)

type ClassicService struct {
	backend core.Backend
	stack   core.Node
}

var (
	pl      core.PluginLoader
	backend restricted.Backend
	log     core.Logger
	events  core.Feed
)

var httpApiFlagName = "http.api"

func Initialize(ctx core.Context, loader core.PluginLoader, logger core.Logger) { 
	pl = loader
	events = pl.GetFeed()
	log = logger
	v := ctx.String(httpApiFlagName)
	if v != "" {
		ctx.Set(httpApiFlagName, v+",plugeth")
	} else {
		ctx.Set(httpApiFlagName, "eth,net,web3,plugeth")
		log.Info("Loaded consensus engine plugin")
	}
}

func GetAPIs(stack core.Node, backend core.Backend) []core.API {
	return []core.API{
		{
			Namespace: "plugeth",
			Version:   "1.0",
			Service:   &ClassicService{backend, stack},
			Public:    true,
		},
	}
}

func DefaultDataDir(path string) string {
	return filepath.Join(path, "classic")
}

func SetNetworkId() *uint64 {
	var networkId *uint64
	classicNetworkId := uint64(1)
	networkId = &classicNetworkId
	return networkId 
}

func SetBootstrapNodes() []string {
	result := ClassicBootnodes
	return result
}

func SetETHDiscoveryURLs(lightSync bool) []string {

	url := ClassicDNSNetwork1
	if lightSync == true {
		url = strings.ReplaceAll(url, "all", "les")
	}
	result := []string{url}
	snapDiscoveryURLs = result

	return result
}

func SetSnapDiscoveryURLs() []string {
	return snapDiscoveryURLs
}


func (service *ClassicService) Test(ctx context.Context) string {
	return "total classic"
}

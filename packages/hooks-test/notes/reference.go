package main

import (
	"context"
	// "reflect"
	"github.com/openrelayxyz/plugeth-utils/core"
	// "github.com/openrelayxyz/plugeth-utils/restricted"
)

type HookTestService struct {
	backend core.Backend
	stack   core.Node
}

var log core.Logger

var httpApiFlagName = "http.api"

func Initialize(ctx core.Context, loader core.PluginLoader, logger core.Logger) {
	log = logger
	v := ctx.String(httpApiFlagName)
	if v != "" {
		ctx.Set(httpApiFlagName, v+",plugeth")
	} else {
		ctx.Set(httpApiFlagName, "eth,net,web3,plugeth")
		log.Info("Loaded plugeth test plugin")
	}
}

func GetAPIs(stack core.Node, backend core.Backend) []core.API {
	return []core.API{
		{
			Namespace: "plugeth",
			Version:   "1.0",
			Service:   &IsSyncedService{backend, stack},
			Public:    true,
		},
	}
}




type hookCall struct {
	method string
	params []interface{}
}
  
var (
callch = make(chan hookCalls, 100)
)

var errs []errors

func main() {
// select {
// case callRecord := <- callch:
// 	//whatever
// case <-time.NewTimer(10 * time.second).C:
// 	// err
	
// }

err := client.Call("my_method", &whatever, params...)
// Check whatever is as expected and err == nil

if callRecord.method != "getRPC" {
	errs = append(errs, fmt.Errorf("Get RPC not called"))
}

if len(callRecord.params) != 2 {
	fmt.Println("call record params too smol")
}
// if _, ok := callRecord.params[1].(restricted.Params); !ok {
// 	errs = append(errs, fmt.Errorf("params not of right type"))
// }

// Use engine_newPlayloadV1 to insert a block
callRecord = <-callch
if callRecord.method != "XYZ" {
	errs = append()
}

if len(errs) > 0 {
	for _, err := range errs {
	log.Error("Error", "err", err)
	}
	os.Exit(1)
}
os.Exit(0)
}

// func GetRPCMethods(x, y, z string) {
// callch <- hookCall{"getRPC", []interface{x, y, z}}
// }

// func ProcessBlock() {
// callch <- hookCall{"processBlock"}
// }
  
// ===================================
  

// type TestObj struct {
// }

// type HookTestService struct {
// 	backend core.Backend
// 	stack   core.Node
// }

// var log core.Logger

// func Initialize(ctx core.Context, loader core.PluginLoader, logger core.Logger) {
// 	log = logger
// 	log.Info("loaded Get Rpc Calls plugin")
// }

// // var httpApiFlagName = "http.api"

// // func Initialize(ctx core.Context, loader core.PluginLoader, logger core.Logger) {
// // 	log = logger
// // 	v := ctx.String(httpApiFlagName)
// // 	if v != "" {
// // 		ctx.Set(httpApiFlagName, v+",plugeth")
// // 	} else {
// // 		ctx.Set(httpApiFlagName, "eth,net,web3,plugeth")
// // 		log.Info("loaded PluGeth hooks test plugin")
// // 	}
// // }

// func (t *TestObj) GetAPIs(stack core.Node, backend core.Backend) []core.API {
// 	return []core.API{
// 		{
// 			Namespace: "test",
// 			Version:   "1.0",
// 			Service:   &HookTestService{backend, stack},
// 			Public:    true,
// 		},
// 	}
// }

// // var test string

// func (service *HookTestService) TestAPI() []core.Api {
// 	result := GetAPIs(service.stack, service.backend)
// 	return result
// }

// func HookTester() {
// 	t := TestObj{}
// 	GetRPCCalls("2", "anyOldMethod", "goodbye horses")
// 	x := t.TestAPI()
// 	log.Error("")
// }

// // func InitializeNode(core.Node, core.Backend) 

// func GetRPCCalls(id, method, params string) {
// 	log.Info("Received RPC Call", "id", id, "method", method, "params", params)
// }
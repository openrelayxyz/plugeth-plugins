package main

import (
	// "fmt"
	"os"
	"context"
	"reflect"
	"github.com/openrelayxyz/plugeth-utils/core"
	// "github.com/openrelayxyz/plugeth-utils/restricted"
)

type HookTestService struct {
	backend core.Backend
	stack   core.Node
}

type hookCall struct {
	method string
	params []interface{}
}
  
var (
	callch = make(chan hookCall, 100)
)

var errs []error

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

var apis []core.API


func GetAPIs(stack core.Node, backend core.Backend) []core.API {
	apis =  []core.API{
		{
			Namespace: "plugeth",
			Version:   "1.0",
			Service:   &HookTestService{backend, stack},
			Public:    true,
		},
	}

	return apis
}


	
func HookTester() {
	
	cl := apis[0].Service.(*HookTestService).stack
	tp := reflect.TypeOf(cl)
	client, err := cl.Attach()
	if err != nil {
		errs = append(errs, err)
		log.Error("Error connecting with client")
	}

	var x interface{}
	err = client.Call(&x, "plugeth_isSynced")
	if err != nil {
		errs = append(errs, err)
		log.Error("failed to call method", "err", err)
	}
	log.Error("this is the return value", "type", tp, "obj", cl, "test", x, "len", len(errs))

	if len(errs) > 0 {
		for _, err := range errs {
		log.Error("Error", "err", err)
		}
	os.Exit(1)
	}

	os.Exit(0)
	
}

func (service *HookTestService) Test(ctx context.Context) string {
	return "you got me"
}




// type hookCall struct {
// 	method string
// 	params []interface{}
// }
  
// var (
// 	callch = make(chan hookCalls, 100)
// )

// var errs []errors

// func main() {
// // select {
// // case callRecord := <- callch:
// // 	//whatever
// // case <-time.NewTimer(10 * time.second).C:
// // 	// err
	
// // }

// err := client.Call("my_method", &whatever, params...)
// // Check whatever is as expected and err == nil

// if callRecord.method != "getRPC" {
// 	errs = append(errs, fmt.Errorf("Get RPC not called"))
// }

// if len(callRecord.params) != 2 {
// 	fmt.Println("call record params too smol")
// }
// // if _, ok := callRecord.params[1].(restricted.Params); !ok {
// // 	errs = append(errs, fmt.Errorf("params not of right type"))
// // }

// // Use engine_newPlayloadV1 to insert a block
// callRecord = <-callch
// if callRecord.method != "XYZ" {
// 	errs = append()
// }

// if len(errs) > 0 {
// 	for _, err := range errs {
// 	log.Error("Error", "err", err)
// 	}
// 	os.Exit(1)
// }
// os.Exit(0)
// }
package main

import (
	// "fmt"
	// "os"
	"context"
	// "reflect"
	"encoding/json"
	// "io"
	"io/ioutil"
	"github.com/openrelayxyz/plugeth-utils/core"
	// "github.com/openrelayxyz/plugeth-utils/restricted/types"
)

type HookTestService struct {
	backend core.Backend
	stack   core.Node
}

type hookCall struct {
	method string
	params []interface{}
}

// ForkchoiceUpdatedV2(update engine.ForkchoiceStateV1, payloadAttributes *engine.PayloadAttributes)
type ForkchoiceStateV1 struct {
	HeadBlockHash      core.Hash `json:"headBlockHash"`
	SafeBlockHash      core.Hash `json:"safeBlockHash"`
	FinalizedBlockHash core.Hash `json:"finalizedBlockHash"`
}
           
// type PayloadAttributes struct {
// 	Timestamp             uint64              `json:"timestamp"             gencodec:"required"`
// 	Random                common.Hash         `json:"prevRandao"            gencodec:"required"`
// 	SuggestedFeeRecipient common.Address      `json:"suggestedFeeRecipient" gencodec:"required"`
// 	Withdrawals           []*types.Withdrawal `json:"withdrawals"`
// }
  
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

func testDataDecompress() (map[string]json.RawMessage, error) {
	raw, err := ioutil.ReadFile("notes/test_data.json")
	if err != nil {
		return nil, err
	}

	// raw, err := ioutil.ReadAll(file) 
	// if err != nil {
	// 	return nil, err
	// }
	// if err == io.EOF || err == io.ErrUnexpectedEOF {
	// 	return nil, err
	// }
	var testData map[string]json.RawMessage
	json.Unmarshal(raw, &testData)
	return testData, nil
}

// func IsSynced(ctx context.Context) {
// 	callch <- hookCall{"getRPC", []interface{x, y, z}}
//   }


	
func HookTester() {
	
	cl := apis[0].Service.(*HookTestService).stack
	// tp := reflect.TypeOf(cl)
	client, err := cl.Attach()
	if err != nil {
		errs = append(errs, err)
		log.Error("Error connecting with client")
	}

	// var x interface{}
	// err = client.Call(&x, "plugeth_isSynced")
	// if err != nil {
	// 	errs = append(errs, err)
	// 	log.Error("failed to call method plugeth_isSynced", "err", err)
	// }
	// log.Error("this is the return value for isSynced", "test", x, "len", len(errs))

	var x interface{}
	err = client.Call(&x, "mynamespace_hello")
	if err != nil {
		errs = append(errs, err)
		log.Error("failed to call method mynamespace_hello", "err", err)
	}
	log.Error("this is the return value for hello", "test", x, "len", len(errs))


	block, err := testDataDecompress()
	if err != nil {
		log.Error("there was an error retrieving testdata", "err", err)
	}

	var y interface{}
	err = client.Call(&y, "engine_newPayloadV1", block)
	if err != nil {
		errs = append(errs, err)
		log.Error("failed to call method", "err", err)
	}
	log.Error("this is the return value for the engine call", "test", y)

	// var z interface{}
	// hash := core.HexToHash("0x5cd31a0a2b37532875307299b0dee57fbc03c4205c7b4db4db09f8fa32dca26c")
	// fd, err := testDataDecompress(")
	// if err != nil {
	// 	log.Error("there was an error retrieving finalized data", "err", err)
	// }
	// parms := []map[string]json.RawMessage{fd, nil}
	// err = client.Call(&z, "engine_forkchoiceUpdatedV2", parms)
	// if err != nil {
	// 	errs = append(errs, err)
	// 	log.Error("failed to call method", "err", err)
	// }
	// log.Error("this is the return value for the engine call", "test", y)

	if len(errs) > 0 {
		for _, err := range errs {
		log.Error("Error", "err", err)
		}
	// os.Exit(1)
	}

	// os.Exit(0)
	
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
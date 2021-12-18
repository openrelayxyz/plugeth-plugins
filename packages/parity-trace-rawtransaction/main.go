package main

import (
	"context"
	"reflect"
	"strings"

	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted/hexutil"
	"github.com/openrelayxyz/plugeth-utils/restricted/types"
	"gopkg.in/urfave/cli.v1"
)

type OuterResult struct {
	Output    hexutil.Bytes   `json:"output"`
	StateDiff *string         `json:"stateDiff"`
	Trace     []*ParityResult `json:"trace"`
	VMTrace   *string         `json:"vmTrace"`
}

// type InnerResult struct {
// 	GasUsed string        `json:"gasUsed"`
// 	Output  hexutil.Bytes `json:"output"`
// }

type Action struct {
	CallType string         `json:"callType"`
	From     string         `json:"from"`
	Gas      string         `json:"gas"`
	Input    string         `json:"input"`
	To       string         `json:"to"`
	Value    hexutil.Uint64 `json:"value"`
}

type ParityResult struct {
	Action Action `json:"action"`
	// Result        InnerResult `json:"result"`
	Error         string `json:"error"`
	SubTraces     int    `json:"subtraces"`
	TracerAddress []int  `json:"traceAddress"`
	Type          string `json:"type"`
}

type APIs struct {
	backend core.Backend
	stack   core.Node
}

var log core.Logger
var httpApiFlagName = "http.api"

func Initialize(ctx *cli.Context, loader core.PluginLoader, logger core.Logger) {
	log = logger
	v := ctx.GlobalString(httpApiFlagName)
	if v != "" {
		ctx.GlobalSet(httpApiFlagName, v+",trace")
	} else {
		ctx.GlobalSet(httpApiFlagName, "eth,net,web3,trace")
		log.Info("Loaded Open Ethereum tracer plugin")
	}
}

func GetAPIs(stack core.Node, backend core.Backend) []core.API {
	defer log.Info("APIs Initialized")
	return []core.API{
		{
			Namespace: "trace",
			Version:   "1.0",
			Service: &APIs{backend: backend,
				stack: stack,
			},
			Public: true,
		},
	}
}

type GethResponse struct {
	Type    string         `json:"type,omitempty"`
	From    string         `json:"from,omitempty"`
	To      string         `json:"to,omitempty"`
	Value   hexutil.Uint64 `json:"value,omitempty"`
	Gas     string         `json:"gas,omitempty"`
	GasUsed string         `json:"gasUsed,omitempty"`
	Input   string         `json:"input,omitempty"`
	Output  hexutil.Bytes  `json:"output,omitempty"`
	Error   string         `json:"error,omitempty"`
	Calls   []GethResponse `json:"calls,omitempty"`
}

func FilterPrecompileCalls(calls []GethResponse) []GethResponse {
	result := []GethResponse{}
	for _, call := range calls {
		//develop test case to see what parity does if is legit call to procompiled contract
		if !strings.HasPrefix(call.To, "0x000000000000000000000000000000000000") || call.Value != 0 {
			result = append(result, call)
		}
	}
	return result
}

func GethParity(gr GethResponse, address []int, t string) []*ParityResult {
	result := []*ParityResult{}
	calls := FilterPrecompileCalls(gr.Calls)
	// if gr.Output == "" {
	// 	gr.Output = "0x0"
	// }
	er := gr.Error
	if er == "execution reverted" {
		er = "Reverted"
	}
	addr := make([]int, len(address))
	copy(addr[:], address)
	result = append(result, &ParityResult{
		Action: Action{CallType: strings.ToLower(gr.Type),
			From:  gr.From,
			Gas:   gr.Gas,
			Input: gr.Input,
			To:    gr.To,
			Value: gr.Value},
		// Result: InnerResult{GasUsed: gr.GasUsed,
		// 	Output: gr.Output},
		Error:         er,
		SubTraces:     len(calls),
		TracerAddress: addr,
		Type:          t})
	for i, call := range calls {
		result = append(result, GethParity(call, append(address, i), t)...)
	}
	return result
}

func (ap *APIs) RawTransaction(ctx context.Context, data hexutil.Bytes, tracer []string) (interface{}, error) {
	// client, err := ap.stack.Attach()
	// if err != nil {
	// 	return nil, err
	// }
	tx := types.Transaction{}
	err := tx.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	txObject := make(map[string]interface{})
	// txObject["to"] = tx.To()
	// gs := hexutil.EncodeUint64(tx.Gas())
	// txObject["gas"] = gs
	// dt := hexutil.Encode(tx.Data())
	// txObject["data"] = dt
	// gp := hexutil.EncodeBig(tx.GasPrice())
	// txObject["gasPrice"] = gp
	// vl := hexutil.EncodeBig(tx.Value())
	// txObject["value"] = vl

	txObject["to"] = reflect.TypeOf(tx.To())
	txObject["gas"] = reflect.TypeOf(tx.Gas())
	txObject["data"] = reflect.TypeOf(tx.Data())
	txObject["gasPrice"] = reflect.TypeOf(tx.GasPrice())
	txObject["value"] = reflect.TypeOf(tx.Value())

	// gr := GethResponse{}
	// client.Call(&gr, "debug_traceCall", txObject, "latest", map[string]string{"tracer": "callTracer"})
	// tAddress := make([]int, 0)
	// gp := GethParity(gr, tAddress, strings.ToLower(gr.Type))
	// result := &OuterResult{
	// 	Output:    gr.Output,
	// 	StateDiff: nil,
	// 	Trace:     gp,
	// 	VMTrace:   nil,
	// }
	return txObject, err
}

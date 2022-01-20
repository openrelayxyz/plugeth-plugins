package main

import (
	"context"
	"strings"

	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted/hexutil"
	"gopkg.in/urfave/cli.v1"
)

type OuterResult struct {
	Output    string          `json:"output"`
	StateDiff *string         `json:"stateDiff"`
	Trace     []*ParityResult `json:"trace"`
	VMTrace   *string         `json:"vmTrace"`
}

type InnerResult struct {
	Address string `json:"address,omitempty"`
	Code    string `json:"code,omitempty"`
	GasUsed string `json:"gasUsed,omitempty"`
	Output  string `json:"output,omitempty"`
}

type Action struct {
	CallType string         `json:"callType,omitempty"`
	From     string         `json:"from,omitempty"`
	Gas      string         `json:"gas,omitempty"`
	Init     string         `json:"init,omitempty"`
	Input    string         `json:"input,omitempty"`
	To       string         `json:"to,omitempty"`
	Value    hexutil.Uint64 `json:"value"`
}

// type CreateAction struct {
// 	From  string         `json:"from"`
// 	Gas   string         `json:"gas"`
// 	Init  string         `json:"init"`
// 	Value hexutil.Uint64 `json:"value"`
// }
//
// type CreateInnerResult struct {
// 	Address string `json:"address"`
// 	Code    string `json:"code"`
// 	GasUsed string `json:"gasUsed"`
// }

type ParityResult struct {
	Action        *Action      `json:"action"`
	Error         string       `json:"error,omitempty"`
	Result        *InnerResult `json:"result,omitempty"`
	SubTraces     int          `json:"subtraces"`
	TracerAddress []int        `json:"traceAddress"`
	Type          string       `json:"type"`
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
	Output  string         `json:"output,omitempty"`
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
	addr := make([]int, len(address))
	copy(addr[:], address)
	if gr.Error == "execution reverted" {
		result = append(result, &ParityResult{
			Action: &Action{CallType: strings.ToLower(gr.Type),
				From:  gr.From,
				Gas:   gr.Gas,
				Input: gr.Input,
				To:    gr.To,
				Value: gr.Value},
			Error:         "Reverted",
			SubTraces:     len(calls),
			TracerAddress: addr,
			Type:          t})
	}
	if gr.Type == "CREATE" {
		result = append(result, &ParityResult{
			Action: &Action{
				From:  gr.From,
				Gas:   gr.Gas,
				Init:  gr.Input,
				Value: gr.Value},
			Result: &InnerResult{
				Address: gr.To,
				Code:    gr.Output,
				GasUsed: gr.GasUsed,
			},
			SubTraces:     len(calls),
			TracerAddress: addr,
			Type:          t})
	} else {
		result = append(result, &ParityResult{
			Action: &Action{CallType: strings.ToLower(gr.Type),
				From:  gr.From,
				Gas:   gr.Gas,
				Input: gr.Input,
				To:    gr.To,
				Value: gr.Value},
			Result: &InnerResult{GasUsed: gr.GasUsed,
				Output: gr.Output},
			SubTraces:     len(calls),
			TracerAddress: addr,
			Type:          t})
	}
	for i, call := range calls {
		if call.Type == "DELEGATECALL" {
			call.Value = gr.Value
		}
		result = append(result, GethParity(call, append(address, i), t)...)
	}
	return result
}

func (ap *APIs) ReplayTransaction(ctx context.Context, txHash core.Hash, types []string) (interface{}, error) {
	client, err := ap.stack.Attach()
	if err != nil {
		return nil, err
	}
	gr := GethResponse{}
	client.Call(&gr, "debug_traceTransaction", txHash, map[string]string{"tracer": "callTracer"})
	tAddress := make([]int, 0)
	gp := GethParity(gr, tAddress, strings.ToLower(gr.Type))
	result := &OuterResult{
		Output:    gr.Output,
		StateDiff: nil,
		Trace:     gp,
		VMTrace:   nil,
	}
	return result, err
}

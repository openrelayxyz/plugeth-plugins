package main

import (
	"context"
	"strings"

	"github.com/openrelayxyz/plugeth-utils/core"
	"gopkg.in/urfave/cli.v1"
)

type OuterResult struct {
	Output    string          `json:"output"`
	StateDiff *string         `json:"stateDiff"`
	Trace     []*ParityResult `json:"trace"`
	VMTrace   *string         `json:"vmTrace"`
}

type InnerResult struct {
	GasUsed string `json:"gasUsed"`
	Output  string `json:"output"`
}

type Action struct {
	CallType string `json:"callType"`
	From     string `json:"from"`
	Gas      string `json:"gas"`
	Input    string `json:"input"`
	To       string `json:"to"`
	Value    string `json:"value"`
}

type ParityResult struct {
	Action        Action `json:"action"`
	Error         string `json:"error,omitempty"`
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
	Value   string         `json:"value,omitempty"`
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
		if !strings.HasPrefix(call.To, "0x000000000000000000000000000000000000") || call.Value != "" {
			result = append(result, call)
		}
	}
	return result
}

func GethParity(gr GethResponse, address []int, t string) []*ParityResult {
	result := []*ParityResult{}
	calls := FilterPrecompileCalls(gr.Calls)
	if gr.Value == "" {
		gr.Value = "0x0"
	}
	if gr.Error == "execution reverted" {
		gr.Error = "Reverted"
	}
	result = append(result, &ParityResult{
		Action: Action{CallType: strings.ToLower(gr.Type),
			From:  gr.From,
			Gas:   gr.Gas,
			Input: gr.Input,
			To:    gr.To,
			Value: gr.Value},
		Error:         gr.Error,
		SubTraces:     len(calls),
		TracerAddress: address,
		Type:          t})
	for i, call := range calls {
		result = append(result, GethParity(call, append(address, i), t)...)
	}
	return result
}

func (ap *APIs) Call(ctx context.Context, txObject map[string]string, types []string) (interface{}, error) {
	client, err := ap.stack.Attach()
	if err != nil {
		return nil, err
	}
	// gethObject := make(map[string]interface{})
	// gethObject["from"] = txObject["from"]
	// gethObject["to"] = txObject["to"]
	// gethObject["gas"] = txObject["gas"]
	// gethObject["data"] = txObject["data"]
	// gethObject["gasPrice"] = txObject["gasPrice"]
	// gethObject["value"] = txObject["value"]
	gr := GethResponse{}
	client.Call(&gr, "debug_traceCall", txObject, "latest", map[string]string{"tracer": "callTracer"})
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

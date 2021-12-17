package main

import (
	"context"
	"strings"

	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted/hexutil"
	"github.com/openrelayxyz/plugeth-utils/restricted/rlp"
	"github.com/openrelayxyz/plugeth-utils/restricted/types"

	"gopkg.in/urfave/cli.v1"
)

type OuterResult struct {
	Output          string          `json:"output"`
	StateDiff       *string         `json:"stateDiff"`
	Trace           []*ParityResult `json:"trace"`
	TransactionHash core.Hash       `json:"transactionHash"`
	VMTrace         *string         `json:"vmTrace"`
}

type InnerResult struct {
	GasUsed string `json:"gasUsed"`
	Output  string `json:"output"`
}

type Action struct {
	CallType string         `json:"callType"`
	From     string         `json:"from"`
	Gas      string         `json:"gas"`
	Input    string         `json:"input"`
	To       string         `json:"to"`
	Value    hexutil.Uint64 `json:"value"`
}

type ParityResult struct {
	Action        Action      `json:"action"`
	Result        InnerResult `json:"result"`
	SubTraces     int         `json:"subtraces"`
	TracerAddress []int       `json:"traceAddress"`
	Type          string      `json:"type"`
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

type OuterGethResponse struct {
	Result GethResponse `json:"result"`
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
	if gr.Output == "" {
		gr.Output = "0x0"
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
		Result: InnerResult{GasUsed: gr.GasUsed,
			Output: gr.Output},
		SubTraces:     len(calls),
		TracerAddress: addr,
		Type:          t})
	for i, call := range calls {
		result = append(result, GethParity(call, append(address, i), t)...)
	}
	return result
}

func (ap *APIs) ReplayBlockTransactions(ctx context.Context, bkNumber string, tracer []string) (interface{}, error) {
	client, err := ap.stack.Attach()
	if err != nil {
		return nil, err
	}
	gr := []OuterGethResponse{}
	err = client.Call(&gr, "debug_traceBlockByNumber", bkNumber, map[string]string{"tracer": "callTracer"})
	if err != nil {
		return nil, err
	}
	blockNM, err := hexutil.DecodeUint64(bkNumber)
	if err != nil {
		return nil, err
	}
	block := &types.Block{}
	bkNB := int64(blockNM)
	rlpBlock, err := ap.backend.BlockByNumber(ctx, bkNB)
	if err != nil {
		return nil, err
	}
	if err := rlp.DecodeBytes(rlpBlock, block); err != nil {
		return nil, err
	}

	transactions := block.Transactions()

	pr := [][]*ParityResult{}
	for _, item := range gr {
		tAddress := make([]int, 0)
		pr = append(pr, GethParity(item.Result, tAddress, strings.ToLower(item.Result.Type)))
	}
	result := make([]OuterResult, len(pr))
	for i, _ := range result {
		result[i] = OuterResult{
			Output:          gr[i].Result.Output,
			StateDiff:       nil,
			Trace:           pr[i],
			TransactionHash: transactions[i].Hash(),
			VMTrace:         nil,
		}
	}
	return result, nil
}

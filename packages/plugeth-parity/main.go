package main



import (
	"context"

	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted"
	"github.com/openrelayxyz/plugeth-utils/restricted/hexutil"
	"github.com/openrelayxyz/plugeth-utils/restricted/rlp"
	"github.com/openrelayxyz/plugeth-utils/restricted/types"
)

type FinalResult struct {
	Output    string               `json:"output"`
	StateDiff map[string]*LayerTwo `json:"stateDiff"`
	Trace     []*ParityResult      `json:"trace"`
	TransactionHash *core.Hash      `json:"transactionHash,omitempty"`
	VMTrace   interface{}          `json:"vmTrace"`
}

type RawData struct {
        TraceVar [][]*ParityResult
        SDVar []struct{Result SDTracerService}
        VMVar []struct{Result VMTracerService}
        Outputs [][]string
}

type ParityTrace struct {
	backend restricted.Backend
	stack   core.Node
}

var Tracers = map[string]func(core.StateDB,  core.BlockContext) core.TracerResult{
	"plugethVMTracer": func(sdb core.StateDB, bctx core.BlockContext) core.TracerResult {
		return &VMTracerService{StateDB: sdb, log:log}
	},
	"plugethStateDiffTracer": func(sdb core.StateDB, bctx core.BlockContext) core.TracerResult {
		return &SDTracerService{stateDB: sdb, blockContext: bctx, log:log}
	},
}

func GetAPIs(stack core.Node, backend restricted.Backend) []core.API {
	return []core.API{
		{
			Namespace: "trace",
			Version:   "1.0",
			Service:   &ParityTrace{backend, stack},
			Public:    true,
		},
	}
}

var log core.Logger
var httpApiFlagName = "http.api"

func Initialize(ctx core.Context, loader core.PluginLoader, logger core.Logger) {
	log = logger
	v := ctx.String(httpApiFlagName)
	if v != "" {
		ctx.Set(httpApiFlagName, v+",trace")
	} else {
		ctx.Set(httpApiFlagName, "eth,net,web3,trace")
		log.Info("Loaded plugeth-parity plugin")
	}
}

func (pt *ParityTrace) blockStringProcessing (ctx context.Context, bkNum string) (string, error) {
	var err error
	var currentBlockString string
	var currentBlockUint64 uint64
	var result string
	client, err := pt.stack.Attach()
	if err != nil {
		return "", err
	}
	client.Call(&currentBlockString, "eth_blockNumber")
	currentBlockUint64, err = hexutil.DecodeUint64(currentBlockString)
	if err != nil {
		return "", err
	}
	switch bkNum {
	case "latest":
		result = hexutil.EncodeUint64(currentBlockUint64 -1)
	case "pending":
		result = hexutil.EncodeUint64(currentBlockUint64 -2)
	case "earliest":
		result = hexutil.EncodeUint64(0)
	default:
		result = bkNum
	}
	return result, nil
}

func (pt *ParityTrace) Call(ctx context.Context, txObject map[string]interface{}, tracerType []string, bkNum *string) (interface{}, error) {
	result := &FinalResult{}
	var output string
	var err error
	var bn string
	if bkNum == nil {
		b := "latest"
		bkNum = &b
	}
	bn, err = pt.blockStringProcessing(ctx, *bkNum)
	if err != nil {
		return nil, err
	}


	for _, typ := range tracerType {
		if typ == "trace" {
				result.Trace, output, err = pt.TraceVariantCall(ctx, txObject, bn)
				if err != nil {return nil, err}
			}
		if typ == "vmTrace" {
			  result.VMTrace, output, err = pt.VMTraceVariantCall(ctx, txObject, bn)
					if err != nil {return nil, err}
			}
		if typ == "stateDiff" {
			result.StateDiff, output, err = pt.StateDiffVariantCall(ctx, txObject, bn)
				if err != nil {return nil, err}
		}
	}
	result.Output = output
	return result, nil
}


func (pt *ParityTrace) RawTransaction(ctx context.Context, data hexutil.Bytes, tracerType []string) (interface{}, error) {
	result := &FinalResult{}
	var output string
	var err error
	tx := types.Transaction{}
	err = tx.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	config := pt.backend.ChainConfig()
	hs := types.LatestSigner(config)
	sender, err := hs.Sender(&tx)
	if err != nil {
		return nil, err
	}
	txObject := make(map[string]interface{})
	txObject["from"] = sender
	txObject["to"] = tx.To()
	gas := hexutil.EncodeUint64(tx.Gas())
	txObject["gas"] = gas
	dt := hexutil.Encode(tx.Data())
	txObject["data"] = dt
	price := hexutil.EncodeBig(tx.GasPrice())
	txObject["gasPrice"] = price
	vl := hexutil.EncodeBig(tx.Value())
	txObject["value"] = vl

	bn, err := pt.blockStringProcessing(ctx, "pending")
	if err != nil {
		return nil, err
	}

	for _, typ := range tracerType {
		if typ == "trace" {
				result.Trace, output, err = pt.TraceVariantCall(ctx, txObject, bn)
				if err != nil {return nil, err}
				}
		if typ == "vmTrace" {
			  result.VMTrace, output, err = pt.VMTraceVariantCall(ctx, txObject, bn)
					if err != nil {return nil, err}
		    }
		if typ == "stateDiff" {
			result.StateDiff, output, err = pt.StateDiffVariantCall(ctx, txObject, bn)
				if err != nil {return nil, err}
				    }
		}
	result.Output = output
	return result, nil
}

func (pt *ParityTrace) ReplayTransaction(ctx context.Context, txHash core.Hash, tracerType []string) (interface{}, error) {
	result := &FinalResult{}
	var output string
	var err error

	for _, typ := range tracerType {
		if typ == "trace" {
				result.Trace, output, err = pt.TraceVariantTransaction(ctx, txHash)
				if err != nil {return nil, err}
				}
		if typ == "vmTrace" {
			  result.VMTrace, output, err = pt.VMTraceVariantTransaction(ctx, txHash)
					if err != nil {return nil, err}
		    }
		if typ == "stateDiff" {
			result.StateDiff, output, err = pt.StateDiffVariantTransaction(ctx, txHash)
				if err != nil {return nil, err}
				    }
		}

	result.Output = output
	return result, nil
}

func (pt *ParityTrace) ReplayBlockTransactions(ctx context.Context, bkNum string, tracerType []string) (interface{}, error) {
	raw := RawData{}
	outputs := [][]string{}
	var traceOutputs []string
	var err error
	bn, err := pt.blockStringProcessing(ctx, bkNum)
	if err != nil {
		return nil, err
	}

	for _, typ := range tracerType {
		if typ == "trace" {
				raw.TraceVar, traceOutputs, err = pt.TraceVariantBlock(ctx, bn)
					if err != nil {return nil, err}
				for _, item := range traceOutputs {
					traceOutputs = append(traceOutputs, item)
					}
					outputs = append(outputs, traceOutputs)
				}
		if typ == "vmTrace" {
				raw.VMVar, err = pt.VMTraceVariantBlock(ctx, bn)
					if err != nil {return nil, err}
				traceOutputs := []string{}
				for _, item := range raw.VMVar {
					traceOutputs = append(traceOutputs, hexutil.Encode(item.Result.Output))
					}
					outputs = append(outputs, traceOutputs)
				}
		if typ == "stateDiff" {
			raw.SDVar, err = pt.StateDiffVariantBlock(ctx, bn)
				if err != nil {return nil, err}
			sdOutputs := []string{}
			for _, item := range raw.SDVar {
				sdOutputs = append(sdOutputs, hexutil.Encode(item.Result.Output))
			}
			outputs = append(outputs, sdOutputs)
						}
		}

	blockNM, err := hexutil.DecodeUint64(bn)
	if err != nil {
		return nil, err
	}
	block := &types.Block{}
	bkNB := int64(blockNM)
	rlpBlock, err := pt.backend.BlockByNumber(ctx, bkNB)
	if err != nil {
		return nil, err
	}
	if err := rlp.DecodeBytes(rlpBlock, block); err != nil {
		return nil, err
	}

	transactions := block.Transactions()
	results := make([]FinalResult, len(transactions))
	for i := range results {
		if outputs[0][i] == "" {
			outputs[0][i] = "0x"
		}
		txHash := transactions[i].Hash()
		results[i] = FinalResult{
			Output: outputs[0][i],
			TransactionHash: &txHash,
		}
		if len(raw.TraceVar) > 0 {
				results[i].Trace = raw.TraceVar[i]
			}
		if len(raw.VMVar) > 0 {
				results[i].VMTrace = raw.VMVar[i].Result.CurrentTrace
			}
		if len(raw.SDVar) > 0 {
			results[i].StateDiff = raw.SDVar[i].Result.ReturnObj
		}
	}
	return results, nil
}

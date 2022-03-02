package main



import (
	"context"

	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted"
	"github.com/openrelayxyz/plugeth-utils/restricted/hexutil"
	"github.com/openrelayxyz/plugeth-utils/restricted/rlp"
	"github.com/openrelayxyz/plugeth-utils/restricted/types"
	"gopkg.in/urfave/cli.v1"
)

type FinalResult struct {
	Output    string          `json:"output"`
	StateDiff map[string]*LayerTwo        `json:"stateDiff"`
	Trace     []*ParityResult `json:"trace"`
	TransactionHash core.Hash       `json:"transactionHash,omit empty"`
	VMTrace   interface{}        `json:"vmTrace"`
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
		return &VMTracerService{StateDB: sdb}
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

func Initialize(ctx *cli.Context, loader core.PluginLoader, logger core.Logger) {
	log = logger
	v := ctx.GlobalString(httpApiFlagName)
	if v != "" {
		ctx.GlobalSet(httpApiFlagName, v+",trace")
	} else {
		ctx.GlobalSet(httpApiFlagName, "eth,net,web3,trace")
		log.Info("Loaded tester plugin")
	}
}

func (pt *ParityTrace) Call(ctx context.Context, txObject map[string]interface{}, tracerType []string) (interface{}, error) {
	result := &FinalResult{}
	var output string
	var err error

	for _, typ := range tracerType {
		if typ == "trace" {
				result.Trace, output, err = pt.TraceVariantZero(ctx, txObject)
				if err != nil {return nil, err}
			}
		if typ == "vmTrace" {
			  result.VMTrace, output, err = pt.VMTraceVariantZero(ctx, txObject)
					if err != nil {return nil, err}
			}
		if typ == "stateDiff" {
			result.StateDiff, output, err = pt.StateDiffVariantZero(ctx, txObject)
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
	//config := *params.GoerliChainConfig
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


	for _, typ := range tracerType {
		if typ == "trace" {
				result.Trace, output, err = pt.TraceVariantZero(ctx, txObject)
				if err != nil {return nil, err}
				}
		if typ == "vmTrace" {
			  result.VMTrace, output, err = pt.VMTraceVariantZero(ctx, txObject)
					if err != nil {return nil, err}
		    }
		if typ == "stateDiff" {
			result.StateDiff, output, err = pt.StateDiffVariantZero(ctx, txObject)
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
				result.Trace, output, err = pt.TraceVariantOne(ctx, txHash)
				if err != nil {return nil, err}
				}
		if typ == "vmTrace" {
			  result.VMTrace, output, err = pt.VMTraceVariantOne(ctx, txHash)
					if err != nil {return nil, err}
		    }
		if typ == "stateDiff" {
			result.StateDiff, output, err = pt.StateDiffVariantOne(ctx, txHash)
				if err != nil {return nil, err}
				    }
		}
	result.Output = output
	return result, nil
}

func (pt *ParityTrace) ReplayBlockTransactions(ctx context.Context, bkNum string, tracerType []string) (interface{}, error) {
	raw := RawData{}
	outputs := [][]string{}
	var err error

	for _, typ := range tracerType {
		if typ == "trace" {
				raw.TraceVar, err = pt.TraceVariantTwo(ctx, bkNum)
					if err != nil {return nil, err}
				traceOutputs := []string{}
				for _, item := range raw.TraceVar {
					traceOutputs = append(traceOutputs, item[0].Result.Output)
					}
					outputs = append(outputs, traceOutputs)
				}
		if typ == "vmTrace" {
				raw.VMVar, err = pt.VMTraceVariantTwo(ctx, bkNum)
					if err != nil {return nil, err}
				traceOutputs := []string{}
				for _, item := range raw.VMVar {
					traceOutputs = append(traceOutputs, string(item.Result.Output))
					}
					outputs = append(outputs, traceOutputs)
				}
		if typ == "stateDiff" {
			raw.SDVar, err = pt.StateDiffVariantTwo(ctx, bkNum)
				if err != nil {return nil, err}
			sdOutputs := []string{}
			for _, item := range raw.SDVar {
				sdOutputs = append(sdOutputs, string(item.Result.Output))
			}
			outputs = append(outputs, sdOutputs)
						}
		}

	blockNM, err := hexutil.DecodeUint64(bkNum)
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
		results[i] = FinalResult{
			Output: outputs[0][i],
			TransactionHash: transactions[i].Hash(),
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

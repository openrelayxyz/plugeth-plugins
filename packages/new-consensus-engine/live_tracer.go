package main

import (
	"context"
	"math/big"
	// "strings"
	"time"

	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted"
	"github.com/openrelayxyz/plugeth-utils/restricted/hexutil"
)

// type LiveTracerResult struct {
// 	// backend core.Backend
// 	// stack core.Node
// 	CallStack []CallStack
// 	Results   []CallStack
// }

// type CallStack struct {
// 	Type    string         `json:"type"`
// 	From    core.Address   `json:"from"`
// 	To      core.Address   `json:"to"`
// 	Value   *big.Int       `json:"value,omitempty"`
// 	Gas     hexutil.Uint64 `json:"gas"`
// 	GasUsed hexutil.Uint64 `json:"gasUsed"`
// 	Input   hexutil.Bytes  `json:"input"`
// 	Output  hexutil.Bytes  `json:"output"`
// 	Time    string         `json:"time,omitempty"`
// 	Calls   []CallStack    `json:"calls,omitempty"`
// 	Results []CallStack    `json:"results,omitempty"`
// 	Error   string         `json:"error,omitempty"`
// }

// func (t *LiveTracerResult) TestLiveTracer(ctx context.Context) string {
// 	return "testLiveTracer"
// }

func (t *LiveTracerResult) TestLiveTracer(ctx context.Context) (<-chan []CallStack, error) {
	subch := make(chan []CallStack, 1000)
	rtrnch := make(chan []CallStack, 1000)
	go func() {
		log.Info("Subscription Block Tracer setup")
		sub := events.Subscribe(subch)
		for {
			select {
			case <-ctx.Done():
				sub.Unsubscribe()
				close(subch)
				close(rtrnch)
				return
			case t := <-subch:
				rtrnch <- t
			case <-sub.Err():
				sub.Unsubscribe()
				close(subch)
				close(rtrnch)
				return
			}
		}
	}()
	return rtrnch, nil
}

func GetLiveTracer(core.Hash, core.StateDB) core.BlockTracer {
	return &LiveTracerResult{}
}

func (r *LiveTracerResult) PreProcessBlock(hash core.Hash, number uint64, encoded []byte) {
	m := map[string]struct{}{
		"LivePreProcessBlock":struct{}{},
	}
	hookChan <- m
	r.Results = []CallStack{}
}

func (r *LiveTracerResult) PreProcessTransaction(tx core.Hash, block core.Hash, i int) {
	m := map[string]struct{}{
		"LivePreProcessTransaction":struct{}{},
	}
	hookChan <- m
}

func (r *LiveTracerResult) BlockProcessingError(tx core.Hash, block core.Hash, err error) {
	m := map[string]struct{}{
		"LiveBlockProcessingError":struct{}{},
	}
	hookChan <- m
}

func (r *LiveTracerResult) PostProcessTransaction(tx core.Hash, block core.Hash, i int, receipt []byte) {
	m := map[string]struct{}{
		"LivePostProcessTransaction":struct{}{},
	}
	hookChan <- m
}

func (r *LiveTracerResult) PostProcessBlock(block core.Hash) {
	m := map[string]struct{}{
		"LivePostProcessBlock":struct{}{},
	}
	hookChan <- m
	if len(r.Results) > 0 {
		events.Send(r.Results)
	}
}

func (r *LiveTracerResult) CaptureStart(from core.Address, to core.Address, create bool, input []byte, gas uint64, value *big.Int) {
	r.CallStack = []CallStack{}
	m := map[string]struct{}{
		"LiveCaptureStart":struct{}{},
	}
	hookChan <- m
}
func (r *LiveTracerResult) CaptureState(pc uint64, op core.OpCode, gas, cost uint64, scope core.ScopeContext, rData []byte, depth int, err error) {
	m := map[string]struct{}{
		"LiveCaptureState":struct{}{},
	}
	hookChan <- m
}

func (r *LiveTracerResult) CaptureFault(pc uint64, op core.OpCode, gas, cost uint64, scope core.ScopeContext, depth int, err error) {
	m := map[string]struct{}{
		"LiveCaptureFault":struct{}{},
	}
	hookChan <- m
}

func (r *LiveTracerResult) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) {
	m := map[string]struct{}{
		"LiveCaptureEnd":struct{}{},
	}
	hookChan <- m
	if len(r.CallStack) > 0 {
		r.Results = append(r.CallStack)
	}
}

func (r *LiveTracerResult) CaptureEnter(typ core.OpCode, from core.Address, to core.Address, input []byte, gas uint64, value *big.Int) {
	m := map[string]struct{}{
		"LiveCaptureEnter":struct{}{},
	}
	hookChan <- m
	r.CallStack = append(r.CallStack, CallStack{
		Type:  restricted.OpCode(typ).String(),
		From:  from,
		To:    to,
		Input: hexutil.Bytes(input),
		Gas:   hexutil.Uint64(gas),
		Calls: []CallStack{},
	})
}

func (r *LiveTracerResult) CaptureExit(output []byte, gasUsed uint64, err error) {
	m := map[string]struct{}{
		"LiveCaptureExit":struct{}{},
	}
	hookChan <- m
	if len(r.CallStack) > 1 {
		returnCall := r.CallStack[len(r.CallStack)-1]
		returnCall.GasUsed = hexutil.Uint64(gasUsed)
		returnCall.Output = output
		r.CallStack[len(r.CallStack)-2].Calls = append(r.CallStack[len(r.CallStack)-2].Calls, returnCall)
		r.CallStack = r.CallStack[:len(r.CallStack)-1]
	}
}

func (r *LiveTracerResult) Result() (interface{}, error) {
	m := map[string]struct{}{
		"LiveTracerResult":struct{}{},
	}
	hookChan <- m
	return "", nil
}


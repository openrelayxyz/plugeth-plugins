package main 

import (
	"math/big"
	"time"

	"github.com/openrelayxyz/plugeth-utils/core"
)

 type TracerService struct {}

var Tracers = map[string]func(core.StateDB) core.TracerResult{
    "testTracer": func(core.StateDB) core.TracerResult {
        return &TracerService{}
    },
}

func (b *TracerService) CaptureStart(from core.Address, to core.Address, create bool, input []byte, gas uint64, value *big.Int) {
	log.Error("capture start")
	m := map[string]struct{}{
		"CaptureStart": struct{}{},
	}
	hookChan <- m
}
func (b *TracerService) CaptureState(pc uint64, op core.OpCode, gas, cost uint64, scope core.ScopeContext, rData []byte, depth int, err error) {
	log.Error("capture state")
	m := map[string]struct{}{
		"CaptureState": struct{}{},
	}
	hookChan <- m
}
func (b *TracerService) CaptureFault(pc uint64, op core.OpCode, gas, cost uint64, scope core.ScopeContext, depth int, err error) {
	log.Error("capture fault")
	m := map[string]struct{}{
		"CaptureFault": struct{}{},
	}
	hookChan <- m
}
func (b *TracerService) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) {
	log.Error("capture end")
	m := map[string]struct{}{
		"CaptureEnd": struct{}{},
	}
	hookChan <- m
}
func (b *TracerService) CaptureEnter(typ core.OpCode, from core.Address, to core.Address, input []byte, gas uint64, value *big.Int) {
	log.Error("capture enter")
	m := map[string]struct{}{
		"CaptureEnter": struct{}{},
	}
	hookChan <- m
}
func (b *TracerService) CaptureExit(output []byte, gasUsed uint64, err error) {
	log.Error("capture exit")
	m := map[string]struct{}{
		"CaptureExit": struct{}{},
	}
	hookChan <- m
}
func (b *TracerService) Result() (interface{}, error) { 
	log.Error("result")
	m := map[string]struct{}{
		"Result": struct{}{},
	}
	hookChan <- m
	return "test complete", nil }
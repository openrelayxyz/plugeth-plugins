package main

import (
	"testing"
	"github.com/openrelayxyz/plugeth-utils/restricted/rlp"
	"github.com/openrelayxyz/plugeth-utils/restricted/hexutil"
	"github.com/openrelayxyz/plugeth-utils/core"
)

// type stateUpdate struct {
// 	Destructs map[core.Hash]struct{}
// 	Accounts map[core.Hash][]byte
// 	Storage map[core.Hash]map[core.Hash][]byte
// 	Code map[core.Hash][]byte
// }

func TestStateUpdateRLP(t *testing.T) {
	orig := &stateUpdate{
		Destructs: map[core.Hash]struct{}{
			core.HexToHash("0x01"): struct{}{},
			core.HexToHash("0xff"): struct{}{},
		},
		Accounts: map[core.Hash][]byte{
			core.HexToHash("0x02"): []byte("17"),
			core.HexToHash("0xfe"): []byte("f7"),
		},
		Storage: map[core.Hash]map[core.Hash][]byte{
			core.HexToHash("0x03"): map[core.Hash][]byte{
				core.HexToHash("0x04"): []byte("18"),
			},
			core.HexToHash("0xfd"): map[core.Hash][]byte{
				core.HexToHash("0xfc"): []byte("f8"),
			},
		},
		Code: map[core.Hash][]byte{
			core.HexToHash("0x05"): []byte("19"),
			core.HexToHash("0xfb"): []byte("f9"),
		},
	}
	data, err := rlp.EncodeToBytes(orig)
	if err != nil {
		t.Errorf("Error encoding: %v", err.Error())
	}
	loaded := new(stateUpdate)
	if err := rlp.DecodeBytes(data, loaded); err != nil {
		t.Errorf("Error decoding: %v", err.Error())
	}
	if _, ok := loaded.Destructs[core.HexToHash("0x01")]; !ok {
		t.Errorf("Destruct missing")
	}
	acct, ok := loaded.Accounts[core.HexToHash("0x02")]
	if !ok { t.Errorf("Account missing") }
	if string(acct) != "17" {
		t.Errorf("Unexpected acount value")
	}
	s, ok := loaded.Storage[core.HexToHash("0x03")]
	if !ok {
		t.Errorf("Account storage missing")
	}
	sv, ok := s[core.HexToHash("0x04")]
	if !ok {
		t.Errorf("Storage value missing")
	}
	if string(sv) != "18" {
		t.Errorf("Unexpected storage value")
	}
	code, ok := loaded.Code[core.HexToHash("0x05")]
	if !ok {
		t.Errorf("Code value missing")
	}
	if string(code) != "19" {
		t.Errorf("Unexpected code value")
	}

}

func TestSample(t *testing.T) {
	loaded := new(stateUpdate)
	data, _ := hexutil.Decode("0xf86bc0f866f2a0fe1981610ff568919eaebec883afe43d94f9558d9c666ca54b3d3bfef2311d3790cf108b01924b3e48efc42e452de58080f2a0c650ee8b209e7ed68bcb4c0603a07c2c47b4256a0cb0fd0c78db167863894d5790cf088b52e5e65630399d9b4c87288080c0c0")
	if err := rlp.DecodeBytes(data, loaded) ; err != nil {
		t.Errorf("Error decoding: %v", err.Error())
	}
}
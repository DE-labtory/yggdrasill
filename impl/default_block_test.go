package impl

import (
	"bytes"
	"testing"
	"time"
)

func TestNewEmptyBlock(t *testing.T) {
	validator := &DefaultValidator{}
	testingTime, _ := time.Parse("Jan 2, 2006 at 3:04pm (MST)", "Feb 3, 2013 at 7:54pm (PST)")
	blockCreator := []byte("testUser")
	genesisSeal := []byte("genesis")
	txList := getTxList(testingTime)

	block := NewEmptyBlock(genesisSeal, 0, blockCreator)
	block.SetTimestamp(testingTime)
	for _, tx := range txList {
		block.PutTx(tx)
	}

	txSeal, _ := validator.BuildTxSeal(convertType(txList))
	block.SetTxSeal(txSeal)

	seal, _ := validator.BuildSeal(block)
	block.SetSeal(seal)

	expected := []byte{99, 189, 233, 43, 92, 7, 47, 216, 209, 10, 54, 137, 155, 0, 36, 186, 188, 94, 159, 213, 238, 240, 23, 65, 122, 21, 140, 72, 117, 63, 24, 191}
	t.Run("Creating new block", func(t *testing.T) {
		if bytes.Compare(block.Seal(), expected) != 0 {
			t.Errorf("Seal = %v, want %v", block.Seal(), expected)
		}
	})
}

func getTxList(testingTime time.Time) []*DefaultTransaction {
	return []*DefaultTransaction{
		&DefaultTransaction{
			PeerID:    "p01",
			ID:        "tx01",
			Status:    0,
			Timestamp: testingTime,
			TxData: &TxData{
				Jsonrpc: "jsonRPC01",
				Method:  "invoke",
				Params: Params{
					Type:     0,
					Function: "function01",
					Args:     []string{"arg1", "arg2"},
				},
				ID: "txdata01",
			},
		},
		&DefaultTransaction{
			PeerID:    "p02",
			ID:        "tx02",
			Status:    0,
			Timestamp: testingTime,
			TxData: &TxData{
				Jsonrpc: "jsonRPC02",
				Method:  "invoke",
				Params: Params{
					Type:     0,
					Function: "function02",
					Args:     []string{"arg1", "arg2"},
				},
				ID: "txdata02",
			},
		},
		&DefaultTransaction{
			PeerID:    "p03",
			ID:        "tx03",
			Status:    0,
			Timestamp: testingTime,
			TxData: &TxData{
				Jsonrpc: "jsonRPC03",
				Method:  "invoke",
				Params: Params{
					Type:     0,
					Function: "function03",
					Args:     []string{"arg1", "arg2"},
				},
				ID: "txdata03",
			},
		},
		&DefaultTransaction{
			PeerID:    "p04",
			ID:        "tx04",
			Status:    0,
			Timestamp: testingTime,
			TxData: &TxData{
				Jsonrpc: "jsonRPC04",
				Method:  "invoke",
				Params: Params{
					Type:     0,
					Function: "function04",
					Args:     []string{"arg1", "arg2"},
				},
				ID: "txdata04",
			},
		},
	}
}

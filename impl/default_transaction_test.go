package impl

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewDefaultTransaction(t *testing.T) {
	testingTime := getTestingTime()
	params := NewParams(1, "functionName", make([]string, 0))
	txData := NewTxData("jsonrpc", Invoke, params, "contractID01")
	tx := NewDefaultTransaction("peerID01", "transactionID01", testingTime, txData)
	expectedByte := []byte{123, 34, 73, 68, 34, 58, 34, 116, 114, 97, 110, 115, 97, 99, 116, 105, 111, 110, 73, 68, 48, 49, 34, 44, 34, 83, 116, 97, 116, 117, 115, 34, 58, 48, 44, 34, 80, 101, 101, 114, 73, 68, 34, 58, 34, 112, 101, 101, 114, 73, 68, 48, 49, 34, 44, 34, 84, 105, 109, 101, 115, 116, 97, 109, 112, 34, 58, 34, 50, 48, 49, 51, 45, 48, 50, 45, 48, 51, 84, 49, 57, 58, 53, 52, 58, 48, 48, 90, 34, 44, 34, 84, 120, 68, 97, 116, 97, 34, 58, 123, 34, 74, 115, 111, 110, 114, 112, 99, 34, 58, 34, 106, 115, 111, 110, 114, 112, 99, 34, 44, 34, 77, 101, 116, 104, 111, 100, 34, 58, 34, 105, 110, 118, 111, 107, 101, 34, 44, 34, 80, 97, 114, 97, 109, 115, 34, 58, 123, 34, 84, 121, 112, 101, 34, 58, 49, 44, 34, 70, 117, 110, 99, 116, 105, 111, 110, 34, 58, 34, 102, 117, 110, 99, 116, 105, 111, 110, 78, 97, 109, 101, 34, 44, 34, 65, 114, 103, 115, 34, 58, 91, 93, 125, 44, 34, 73, 68, 34, 58, 34, 99, 111, 110, 116, 114, 97, 99, 116, 73, 68, 48, 49, 34, 125, 44, 34, 83, 105, 103, 110, 97, 116, 117, 114, 101, 34, 58, 110, 117, 108, 108, 125}

	serializedTx, err := tx.Serialize()
	assert.NoError(t, err)
	assert.Equal(t, serializedTx, expectedByte)
}

func TestDefaultTransaction_CalculateSeal(t *testing.T) {
	tx := getTestingTxList(0)[0]
	expectedSeal := []byte{51, 20, 199, 105, 161, 39, 78, 32, 195, 43, 149, 185, 235, 46, 217, 121, 58, 200, 87, 56, 177, 201, 240, 215, 195, 242, 13, 20, 169, 202, 38, 33}
	seal, err := tx.CalculateSeal()
	assert.NoError(t, err)

	assert.Equal(t, expectedSeal, seal)
}

func TestDefaultTransaction_Serialize(t *testing.T) {
	tx := getTestingTxList(0)[0]
	txBytes, err := tx.Serialize()

	assert.NoError(t, err)

	deserializedTx := &DefaultTransaction{}
	deserializedTx.Deserialize(txBytes)

	assert.Equal(t, deserializedTx, tx)
}

func getTestingTime() time.Time {
	const longForm = "Jan 2, 2006 at 3:04pm (MST)"
	testingTime, _ := time.Parse(longForm, "Feb 3, 2013 at 7:54pm (UTC)")

	return testingTime
}

func getTestingTxList(index int) []*DefaultTransaction {
	testingTime := getTestingTime()

	return [][]*DefaultTransaction{
		{
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
				Signature: nil,
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
				Signature: nil,
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
				Signature: nil,
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
				Signature: nil,
			},
		},
	}[index]
}

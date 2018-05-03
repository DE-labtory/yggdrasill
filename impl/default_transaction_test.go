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
	tx := getTestData()
	expectedSeal := []byte{20, 237, 92, 203, 196, 14, 152, 126, 192, 100, 16, 121, 150, 12, 129, 43, 250, 210, 188, 174, 149, 220, 17, 67, 195, 166, 82, 49, 230, 9, 160, 174}
	seal, err := tx.CalculateSeal()
	assert.NoError(t, err)

	assert.Equal(t, expectedSeal, seal)
}

func TestDefaultTransaction_Serialize(t *testing.T) {
	tx := getTestData()
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

func getTestData() *DefaultTransaction {
	testingTime := getTestingTime()

	testData := &DefaultTransaction{
		ID:        "transactionID01",
		PeerID:    "peerID01",
		Timestamp: testingTime,
		Status:    StatusTransactionInvalid,
		TxData: &TxData{
			Jsonrpc: "jsonrpc",
			Method:  Invoke,
			ID:      "contractID01",
			Params: Params{
				Type:     1,
				Function: "functionName",
				Args:     make([]string, 0),
			},
		},
		Signature: nil,
	}

	return testData
}

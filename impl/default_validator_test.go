package impl

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultValidator_BuildTxSeal(t *testing.T) {
	testData := getTestingData(0)
	expectedTxSealRoot := []byte{195, 17, 112, 227, 157, 68, 134, 162, 202, 81, 64, 22, 8, 206, 223, 48, 121, 236, 94, 40, 230, 158, 34, 224, 226, 75, 34, 57, 69, 239, 181, 239}

	validator := &DefaultValidator{}
	gotTxSeal, err := validator.BuildTxSeal(convertType(testData))

	assert.NoError(t, err)
	assert.Equal(t, expectedTxSealRoot, gotTxSeal[0])
}

func TestDefaultValidator_ValidateTxProof(t *testing.T) {
	testData := getTestingData(0)
	validator := &DefaultValidator{}
	convTestData := convertType(testData)

	txSeal, err := validator.BuildTxSeal(convTestData)
	assert.NoError(t, err)

	validationResult, err := validator.ValidateTxSeal(txSeal, convTestData)
	assert.Equal(t, true, validationResult)
}

func TestDefaultValidator_ValidateTransaction(t *testing.T) {
	notIncludedTxTime, _ := time.Parse("Jan 2, 2006 at 3:04pm (MST)", "Feb 3, 2013 at 7:54pm (PST)")
	notIncludedTx := &DefaultTransaction{
		PeerID:    "p05",
		ID:        "tx05",
		Status:    0,
		Timestamp: notIncludedTxTime,
		TxData: &TxData{
			Jsonrpc: "jsonRPC05",
			Method:  "invoke",
			Params: Params{
				Type:     0,
				Function: "function05",
				Args:     []string{"arg1", "arg2"},
			},
			ID: "txdata05",
		},
		Signature: nil,
	}
	testData := getTestingData(0)
	validator := &DefaultValidator{}
	txSeal, err := validator.BuildTxSeal(convertType(testData))
	assert.NoError(t, err)

	correctResult, err := validator.ValidateTransaction(txSeal, testData[1])
	assert.NoError(t, err)
	assert.Equal(t, true, correctResult)

	wrongResult, err := validator.ValidateTransaction(txSeal, notIncludedTx)
	assert.NoError(t, err)
	assert.Equal(t, false, wrongResult)
}

func getTestingData(index int) []*DefaultTransaction {
	const longForm = "Jan 2, 2006 at 3:04pm (MST)"
	testingTime, _ := time.Parse(longForm, "Feb 3, 2013 at 7:54pm (PST)")

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

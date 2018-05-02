package impl

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultValidator_BuildTxSeal(t *testing.T) {
	testData := getTestingData(0)
	expectedTxSealRoot := []byte{119, 178, 207, 195, 123, 230, 211, 193, 142, 68, 255, 99, 226, 172, 207, 211, 75, 251, 211, 128, 175, 230, 141, 51, 3, 186, 19, 179, 197, 104, 230, 29}

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
		},
	}[index]
}

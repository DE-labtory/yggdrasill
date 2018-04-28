package impl

import (
	"bytes"
	"testing"
	"time"

	"github.com/it-chain/yggdrasill/common"
)

func TestDefaultValidator_BuildTxSeal(t *testing.T) {
	testData := getTestingData(0)

	tests := []struct {
		name      string
		txList    []*DefaultTransaction
		wantProof []byte
		wantErr   bool
	}{
		{
			name:      "Create new merkle tree",
			txList:    testData,
			wantProof: []byte{119, 178, 207, 195, 123, 230, 211, 193, 142, 68, 255, 99, 226, 172, 207, 211, 75, 251, 211, 128, 175, 230, 141, 51, 3, 186, 19, 179, 197, 104, 230, 29},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := &DefaultValidator{}
			gotTxSeal, err := validator.BuildTxSeal(convertType(tt.txList))
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDefaultValidator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if bytes.Compare(gotTxSeal[0], tt.wantProof) != 0 {
				t.Errorf("NewDefaultValidator() = %v, want %v", gotTxSeal, tt.wantProof)
			}
		})
	}
}

func TestDefaultValidator_ValidateTxProof(t *testing.T) {
	testData := getTestingData(0)
	validator := &DefaultValidator{}
	txSeal, _ := validator.BuildTxSeal(convertType(testData))
	convTestData := convertType(testData)

	tests := []struct {
		name string
		want bool
	}{
		{
			name: "Test correct validation",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := validator.ValidateTxSeal(txSeal, convTestData); got != tt.want {
				t.Errorf("DefaultValidator.ValidateTxProof() = %v, want %v", got, tt.want)
			}
		})
	}
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
	txSeal, _ := validator.BuildTxSeal(convertType(testData))

	tests := []struct {
		name string
		t    *DefaultTransaction
		want bool
	}{
		{
			name: "Test true ValidationTransaction",
			t:    testData[1],
			want: true,
		},
		{
			name: "Test false ValidationTransaction",
			t:    notIncludedTx,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := validator.ValidateTransaction(txSeal, tt.t); got != tt.want {
				t.Errorf("DefaultValidator.ValidateTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getTestingData(index int) []*DefaultTransaction {
	const longForm = "Jan 2, 2006 at 3:04pm (MST)"
	testingTime, _ := time.Parse(longForm, "Feb 3, 2013 at 7:54pm (PST)")

	return [][]*DefaultTransaction{
		[]*DefaultTransaction{
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

func convertType(txList []*DefaultTransaction) []common.Transaction {
	convTxList := make([]common.Transaction, 0)
	for _, tx := range txList {
		convTxList = append(convTxList, tx)
	}

	return convTxList
}

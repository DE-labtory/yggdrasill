package validator

import (
	"bytes"
	"testing"
	"time"

	tx "github.com/it-chain/yggdrasill/transaction"
)

func TestMerkleTree_BuildProofAndTxProof(t *testing.T) {
	testData := getTestingData(0)

	tests := []struct {
		name         string
		txList       []*tx.DefaultTransaction
		wantRootHash []byte
		wantErr      bool
	}{
		{
			name:         "Create new merkle tree",
			txList:       testData,
			wantRootHash: []byte{119, 178, 207, 195, 123, 230, 211, 193, 142, 68, 255, 99, 226, 172, 207, 211, 75, 251, 211, 128, 175, 230, 141, 51, 3, 186, 19, 179, 197, 104, 230, 29},
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			merkleTree := &MerkleTree{}
			gotRoot, _, err := merkleTree.BuildProofAndTxProof(convertType(tt.txList))
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMerkleTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if bytes.Compare(gotRoot, tt.wantRootHash) != 0 {
				t.Errorf("NewMerkleTree() = %v, want %v", gotRoot, tt.wantRootHash)
			}
		})
	}
}

func TestMerkleTree_Validate(t *testing.T) {
	testData := getTestingData(0)
	merkleTree := &MerkleTree{}
	_, tree, _ := merkleTree.BuildProofAndTxProof(convertType(testData))
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
			if got, _ := merkleTree.Validate(tree, convTestData); got != tt.want {
				t.Errorf("MerkleTree.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMerkleTree_ValidateTransaction(t *testing.T) {
	notIncludedTxTime, _ := time.Parse("Jan 2, 2006 at 3:04pm (MST)", "Feb 3, 2013 at 7:54pm (PST)")
	notIncludedTx := &tx.DefaultTransaction{
		PeerID:    "p05",
		ID:        "tx05",
		Status:    0,
		Timestamp: notIncludedTxTime,
		TxData: &tx.TxData{
			Jsonrpc: "jsonRPC05",
			Method:  "invoke",
			Params: tx.Params{
				Type:     0,
				Function: "function05",
				Args:     []string{"arg1", "arg2"},
			},
			ID: "txdata05",
		},
	}
	testData := getTestingData(0)
	merkleTree := &MerkleTree{}
	_, tree, _ := merkleTree.BuildProofAndTxProof(convertType(testData))

	tests := []struct {
		name string
		t    *tx.DefaultTransaction
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
			if got, _ := merkleTree.ValidateTransaction(tree, tt.t); got != tt.want {
				t.Errorf("MerkleTree.ValidateTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getTestingData(index int) []*tx.DefaultTransaction {
	const longForm = "Jan 2, 2006 at 3:04pm (MST)"
	testingTime, _ := time.Parse(longForm, "Feb 3, 2013 at 7:54pm (PST)")

	return [][]*tx.DefaultTransaction{
		[]*tx.DefaultTransaction{
			&tx.DefaultTransaction{
				PeerID:    "p01",
				ID:        "tx01",
				Status:    0,
				Timestamp: testingTime,
				TxData: &tx.TxData{
					Jsonrpc: "jsonRPC01",
					Method:  "invoke",
					Params: tx.Params{
						Type:     0,
						Function: "function01",
						Args:     []string{"arg1", "arg2"},
					},
					ID: "txdata01",
				},
			},
			&tx.DefaultTransaction{
				PeerID:    "p02",
				ID:        "tx02",
				Status:    0,
				Timestamp: testingTime,
				TxData: &tx.TxData{
					Jsonrpc: "jsonRPC02",
					Method:  "invoke",
					Params: tx.Params{
						Type:     0,
						Function: "function02",
						Args:     []string{"arg1", "arg2"},
					},
					ID: "txdata02",
				},
			},
			&tx.DefaultTransaction{
				PeerID:    "p03",
				ID:        "tx03",
				Status:    0,
				Timestamp: testingTime,
				TxData: &tx.TxData{
					Jsonrpc: "jsonRPC03",
					Method:  "invoke",
					Params: tx.Params{
						Type:     0,
						Function: "function03",
						Args:     []string{"arg1", "arg2"},
					},
					ID: "txdata03",
				},
			},
			&tx.DefaultTransaction{
				PeerID:    "p04",
				ID:        "tx04",
				Status:    0,
				Timestamp: testingTime,
				TxData: &tx.TxData{
					Jsonrpc: "jsonRPC04",
					Method:  "invoke",
					Params: tx.Params{
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

func convertType(txList []*tx.DefaultTransaction) []tx.Transaction {
	convTxList := make([]tx.Transaction, 0)
	for _, tx := range txList {
		convTxList = append(convTxList, tx)
	}

	return convTxList
}

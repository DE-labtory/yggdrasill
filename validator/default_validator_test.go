package validator

import (
	"bytes"
	"testing"
	"time"

	tx "github.com/it-chain/yggdrasill/transaction"
)

func TestMerkleTree_BuildProof(t *testing.T) {
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
			wantRootHash: []byte{101, 200, 55, 65, 195, 166, 219, 48, 181, 132, 201, 148, 122, 187, 113, 151, 196, 136, 178, 241, 183, 21, 166, 213, 54, 196, 27, 116, 45, 204, 153, 17},
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			merkleTree := &MerkleTree{}
			got, err := merkleTree.BuildProof(convertType(tt.txList))
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMerkleTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if bytes.Compare(got[0], tt.wantRootHash) != 0 {
				t.Errorf("NewMerkleTree() = %v, want %v", got, tt.wantRootHash)
			}
		})
	}
}

func TestMerkleTree_Validate(t *testing.T) {
	testData := getTestingData(0)
	merkleTree := &MerkleTree{}
	proof, _ := merkleTree.BuildProof(convertType(testData))
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
			if got, _ := merkleTree.Validate(proof, convTestData); got != tt.want {
				t.Errorf("MerkleTree.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMerkleTree_ValidateTransaction(t *testing.T) {
	notIncludedTxTime, _ := time.Parse("Jan 2, 2006 at 3:04pm (MST)", "Feb 3, 2013 at 7:54pm (PST)")
	notIncludedTx := &tx.DefaultTransaction{
		InvokePeerID:      "p05",
		TransactionID:     "tx05",
		TransactionStatus: 0,
		TransactionType:   0,
		TransactionHash:   "hashValue",
		TimeStamp:         notIncludedTxTime,
		TxData: &tx.TxData{
			Jsonrpc: "jsonRPC05",
			Method:  "invoke",
			Params: tx.Params{
				ParamsType: 0,
				Function:   "function05",
				Args:       []string{"arg1", "arg2"},
			},
			ID: "txdata05",
		},
	}
	testData := getTestingData(0)
	merkleTree := &MerkleTree{}
	proof, _ := merkleTree.BuildProof(convertType(testData))

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
			if got, _ := merkleTree.ValidateTransaction(proof, tt.t); got != tt.want {
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
				InvokePeerID:      "p01",
				TransactionID:     "tx01",
				TransactionStatus: 0,
				TransactionType:   0,
				TransactionHash:   "hashValue",
				TimeStamp:         testingTime,
				TxData: &tx.TxData{
					Jsonrpc: "jsonRPC01",
					Method:  "invoke",
					Params: tx.Params{
						ParamsType: 0,
						Function:   "function01",
						Args:       []string{"arg1", "arg2"},
					},
					ID: "txdata01",
				},
			},
			&tx.DefaultTransaction{
				InvokePeerID:      "p02",
				TransactionID:     "tx02",
				TransactionStatus: 0,
				TransactionType:   0,
				TransactionHash:   "hashValue",
				TimeStamp:         testingTime,
				TxData: &tx.TxData{
					Jsonrpc: "jsonRPC02",
					Method:  "invoke",
					Params: tx.Params{
						ParamsType: 0,
						Function:   "function02",
						Args:       []string{"arg1", "arg2"},
					},
					ID: "txdata02",
				},
			},
			&tx.DefaultTransaction{
				InvokePeerID:      "p03",
				TransactionID:     "tx03",
				TransactionStatus: 0,
				TransactionType:   0,
				TransactionHash:   "hashValue",
				TimeStamp:         testingTime,
				TxData: &tx.TxData{
					Jsonrpc: "jsonRPC03",
					Method:  "invoke",
					Params: tx.Params{
						ParamsType: 0,
						Function:   "function03",
						Args:       []string{"arg1", "arg2"},
					},
					ID: "txdata03",
				},
			},
			&tx.DefaultTransaction{
				InvokePeerID:      "p04",
				TransactionID:     "tx04",
				TransactionStatus: 0,
				TransactionType:   0,
				TransactionHash:   "hashValue",
				TimeStamp:         testingTime,
				TxData: &tx.TxData{
					Jsonrpc: "jsonRPC04",
					Method:  "invoke",
					Params: tx.Params{
						ParamsType: 0,
						Function:   "function04",
						Args:       []string{"arg1", "arg2"},
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

package validator

import (
	"bytes"
	"testing"
	"time"

	tx "github.com/it-chain/yggdrasill/transaction"
)

func TestNewMerkleTree(t *testing.T) {
	const longForm = "Jan 2, 2006 at 3:04pm (MST)"
	testingTime, _ := time.Parse(longForm, "Feb 3, 2013 at 7:54pm (PST)")

	type args struct {
		txList []*tx.DefaultTransaction
	}
	tests := []struct {
		name         string
		args         args
		wantRootHash []byte
		wantErr      bool
	}{
		{
			name: "Create new merkle tree",
			args: args{
				txList: []*tx.DefaultTransaction{
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
				}},
			wantRootHash: []byte{101, 200, 55, 65, 195, 166, 219, 48, 181, 132, 201, 148, 122, 187, 113, 151, 196, 136, 178, 241, 183, 21, 166, 213, 54, 196, 27, 116, 45, 204, 153, 17},
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMerkleTree(tt.args.txList)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMerkleTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if bytes.Compare(got.data[0], tt.wantRootHash) != 0 {
				t.Errorf("NewMerkleTree() = %v, want %v", got, tt.wantRootHash)
			}
		})
	}
}

func TestMerkleTree_Serialize(t *testing.T) {
	merkleTree, _ := createTestingMerkleTree(0)
	tests := []struct {
		name string
		t    *MerkleTree
		want string
	}{
		{
			name: "Serialize test",
			t:    merkleTree,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.t.Serialize()
			newMerkleTree := &MerkleTree{data: nil}

			newMerkleTree.Deserialize(s)
			for i, bArr := range newMerkleTree.data {
				if bytes.Compare(bArr, tt.t.data[i]) != 0 {
					t.Errorf("Expected = %v, but Actual = %v", tt.t.data[i], bArr)
				}
			}
		})
	}
}

func createTestingMerkleTree(index int) (*MerkleTree, error) {
	const longForm = "Jan 2, 2006 at 3:04pm (MST)"
	testingTime, _ := time.Parse(longForm, "Feb 3, 2013 at 7:54pm (PST)")

	testData := [][]*tx.DefaultTransaction{
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
	}

	return NewMerkleTree(testData[index])
}

package validator

import tx "github.com/it-chain/yggdrasill/transaction"

type merkleNode struct {
	Hash []byte
	Tx   *tx.DefaultTransaction
}

type MerkleTree struct {
	merkleRoot []byte
	root       *merkleNode
	leafs      []*merkleNode
}

func (t *MerkleTree) Validate() bool {
	return false
}

func (t *MerkleTree) ValidateTransaction(proof []byte, tx *tx.Transaction) bool {
	return false
}

func (t *MerkleTree) GetProof() []byte {
	return t.merkleRoot
}

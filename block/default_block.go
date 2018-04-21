package block

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sort"
	"strings"
	"time"

	tx "github.com/it-chain/yggdrasill/transaction"
	"github.com/it-chain/yggdrasill/util"
)

type DefaultBlock struct {
	Header       *BlockHeader
	Proof        [][]byte
	Transactions []*tx.DefaultTransaction
}

type BlockHeader struct {
	Height             uint64
	PreviousHash       string
	Version            string
	MerkleTreeRootHash string
	TimeStamp          time.Time
	CreatorID          string
	Signature          []byte
	BlockHash          string
	MerkleTreeHeight   int
	TransactionCount   int
}

func (block *DefaultBlock) PutTransaction(transaction tx.Transaction) error {

	switch transaction.(type) {
	case *tx.DefaultTransaction:
		block.Transactions = append(block.Transactions, transaction.(*tx.DefaultTransaction))
		block.Header.TransactionCount++
	default:
		return InvalidTransactionTypeError
	}

	return nil
}

func (block *DefaultBlock) Serialize() ([]byte, error) {
	return util.Serialize(block)
}

func (block *DefaultBlock) GenerateHash() error {

	if block.Header.MerkleTreeRootHash == "" {
		return errors.New("no merkle tree root hash")
	}

	str := []string{block.Header.MerkleTreeRootHash, block.Header.TimeStamp.String(), block.Header.PreviousHash}
	block.Header.BlockHash = computeSHA256(str)

	return nil
}

func (block *DefaultBlock) GetHash() string {
	return block.Header.BlockHash
}

func (block *DefaultBlock) GetTransactions() []tx.Transaction {

	txs := make([]tx.Transaction, 0)

	for _, tx := range block.Transactions {
		txs = append(txs, tx)
	}

	return txs
}

func (block *DefaultBlock) GetHeight() uint64 {
	return block.Header.Height
}

func (block *DefaultBlock) IsPrev(serializedBlock []byte) bool {
	lastBlock := &DefaultBlock{}
	err := util.Deserialize(serializedBlock, lastBlock)

	if err != nil {
		return false
	}

	if (block.GetHeight() == lastBlock.GetHeight()+1) && (lastBlock.GetHash() == block.Header.PreviousHash) {
		return true
	}
	return false
}

func computeSHA256(data []string) string {

	sort.Strings(data)
	arg := strings.Join(data, ",")
	hash := sha256.New()
	hash.Write([]byte(arg))
	return hex.EncodeToString(hash.Sum(nil))
}

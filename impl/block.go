package impl

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/it-chain/yggdrasill/common"
	tx "github.com/it-chain/yggdrasill/transaction"
)

type DefaultBlock struct {
	seal       []byte
	prevSeal   []byte
	height     uint64
	txList     []*tx.DefaultTransaction
	txListSeal [][]byte
	timestamp  []byte
	CreatorID  string
}

func (block *DefaultBlock) GenerateSeal() ([]byte, error) {
	if block.prevSeal == nil || block.txListSeal == nil || block.timestamp == nil {
		return nil, common.ErrInsufficientFields
	}

	rootHash := block.txListSeal[0]
	combined := append(block.prevSeal, rootHash...)
	combined = append(combined, block.timestamp...)

	seal := calculateHash(combined)
	return seal, nil
}

func (block *DefaultBlock) SetPrevSeal(prevSeal []byte) {
	block.prevSeal = prevSeal
}

func (block *DefaultBlock) SetHeight(height uint64) {
	block.height = height
}

func (block *DefaultBlock) PutTx(transaction tx.Transaction) error {
	convTx, ok := transaction.(*tx.DefaultTransaction)
	if ok {
		if block.txList == nil {
			block.txList = make([]*tx.DefaultTransaction, 0)
		}

		block.txList = append(block.txList, convTx)
		fmt.Println(len(block.txList))

		return nil
	}

	return common.ErrTransactionType
}

func (block *DefaultBlock) SetTxListSeal(txListSeal [][]byte) {
	block.txListSeal = txListSeal
}

func (block *DefaultBlock) Seal() []byte {
	return block.seal
}

func (block *DefaultBlock) PrevSeal() []byte {
	return block.prevSeal
}

func (block *DefaultBlock) Height() uint64 {
	return block.height
}

func (block *DefaultBlock) TxList() []tx.Transaction {
	txList := make([]tx.Transaction, 0)
	for _, tx := range block.txList {
		txList = append(txList, tx)
	}
	return txList
}

func (block *DefaultBlock) TxListSeal() [][]byte {
	return block.txListSeal
}

func (block *DefaultBlock) Serialize() ([]byte, error) {
	data, err := json.Marshal(block)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (block *DefaultBlock) Deserialize(serializedBlock []byte) error {
	if len(serializedBlock) == 0 {
		return common.ErrDecodingEmptyBlock
	}

	err := json.Unmarshal(serializedBlock, block)
	if err != nil {
		return err
	}

	return nil
}

func (block *DefaultBlock) SetTimestamp(currentTime time.Time) error {
	timestamp, error := currentTime.MarshalBinary()
	if error != nil {
		return error
	}

	block.timestamp = timestamp
	return nil
}

func (block *DefaultBlock) Timestamp() (time.Time, error) {
	var time time.Time
	error := time.UnmarshalBinary(block.timestamp)
	if error != nil {
		return time, error
	}

	return time, nil
}

func calculateHash(b []byte) []byte {
	hashValue := sha256.New()
	hashValue.Write(b)
	return hashValue.Sum(nil)
}

package impl

import (
	"encoding/json"
	"time"

	"bytes"

	"github.com/it-chain/yggdrasill/common"
)

type DefaultBlock struct {
	Seal      []byte
	PrevSeal  []byte
	Height    uint64
	TxList    []*DefaultTransaction
	TxSeal    [][]byte
	Timestamp time.Time
	Creator   []byte
}

func (block *DefaultBlock) SetSeal(seal []byte) {
	block.Seal = seal
}

func (block *DefaultBlock) SetPrevSeal(prevSeal []byte) {
	block.PrevSeal = prevSeal
}

func (block *DefaultBlock) SetHeight(height uint64) {
	block.Height = height
}

func (block *DefaultBlock) PutTx(transaction common.Transaction) error {
	convTx, ok := transaction.(*DefaultTransaction)
	if ok {
		if block.TxList == nil {
			block.TxList = make([]*DefaultTransaction, 0)
		}

		block.TxList = append(block.TxList, convTx)

		return nil
	}

	return common.ErrTransactionType
}

func (block *DefaultBlock) SetTxSeal(txSeal [][]byte) {
	block.TxSeal = txSeal
}

func (block *DefaultBlock) SetCreator(creator []byte) {
	block.Creator = creator
}

func (block *DefaultBlock) SetTimestamp(currentTime time.Time) {
	block.Timestamp = currentTime
}

func (block *DefaultBlock) GetSeal() []byte {
	return block.Seal
}

func (block *DefaultBlock) GetPrevSeal() []byte {
	return block.PrevSeal
}

func (block *DefaultBlock) GetHeight() uint64 {
	return block.Height
}

func (block *DefaultBlock) GetTxList() []common.Transaction {
	txList := make([]common.Transaction, 0)
	for _, tx := range block.TxList {
		txList = append(txList, tx)
	}
	return txList
}

func (block *DefaultBlock) GetTxSeal() [][]byte {
	return block.TxSeal
}

func (block *DefaultBlock) GetCreator() []byte {
	return block.Creator
}

func (block *DefaultBlock) GetTimestamp() time.Time {
	return block.Timestamp
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

func (block *DefaultBlock) IsReadyToPublish() bool {
	return block.Seal != nil
}

func (block *DefaultBlock) IsPrev(serializedPrevBlock []byte) bool {
	prevBlock := &DefaultBlock{}
	prevBlock.Deserialize(serializedPrevBlock)

	return bytes.Compare(prevBlock.GetSeal(), block.GetPrevSeal()) == 0
}

func NewEmptyBlock(prevSeal []byte, height uint64, creator []byte) *DefaultBlock {
	block := &DefaultBlock{}

	block.SetPrevSeal(prevSeal)
	block.SetHeight(height)
	block.SetCreator(creator)

	return block
}

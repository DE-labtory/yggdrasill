package impl

import (
	"encoding/json"
	"time"

	"github.com/it-chain/yggdrasill/common"
)

type DefaultBlock struct {
	seal      []byte
	prevSeal  []byte
	height    uint64
	txList    []*DefaultTransaction
	txSeal    [][]byte
	timestamp []byte
	creator   []byte
}

func (block *DefaultBlock) SetSeal(seal []byte) {
	block.seal = seal
}

func (block *DefaultBlock) SetPrevSeal(prevSeal []byte) {
	block.prevSeal = prevSeal
}

func (block *DefaultBlock) SetHeight(height uint64) {
	block.height = height
}

func (block *DefaultBlock) PutTx(transaction common.Transaction) error {
	convTx, ok := transaction.(*DefaultTransaction)
	if ok {
		if block.txList == nil {
			block.txList = make([]*DefaultTransaction, 0)
		}

		block.txList = append(block.txList, convTx)

		return nil
	}

	return common.ErrTransactionType
}

func (block *DefaultBlock) SetTxSeal(txSeal [][]byte) {
	block.txSeal = txSeal
}

func (block *DefaultBlock) SetCreator(creator []byte) {
	block.creator = creator
}

func (block *DefaultBlock) SetTimestamp(currentTime time.Time) error {
	timestamp, error := currentTime.MarshalBinary()
	if error != nil {
		return error
	}

	block.timestamp = timestamp
	return nil
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

func (block *DefaultBlock) TxList() []common.Transaction {
	txList := make([]common.Transaction, 0)
	for _, tx := range block.txList {
		txList = append(txList, tx)
	}
	return txList
}

func (block *DefaultBlock) TxSeal() [][]byte {
	return block.txSeal
}

func (block *DefaultBlock) Creator() []byte {
	return block.creator
}

func (block *DefaultBlock) Timestamp() (time.Time, error) {
	var time time.Time
	error := time.UnmarshalBinary(block.timestamp)
	if error != nil {
		return time, error
	}

	return time, nil
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
	return block.Seal() != nil
}

func NewEmptyBlock(prevSeal []byte, height uint64, creator []byte) *DefaultBlock {
	block := &DefaultBlock{}

	block.SetPrevSeal(prevSeal)
	block.SetHeight(height)
	block.SetCreator(creator)

	return block
}

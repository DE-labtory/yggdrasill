package common

import (
	"errors"
)

var ErrTransactionType = errors.New("Wrong transaction type")
var ErrInsufficientFields = errors.New("Previous seal or transaction list seal is not set")
var ErrDecodingEmptyBlock = errors.New("Empty Block decoding failed")

// Block 인터페이스는 Block이 기본적으로 가져야 하는 기능들을 정의한다.
type Block interface {
	// Block의 required field setters
	GenerateSeal() ([]byte, error)
	SetPrevSeal(prevSeal []byte)
	SetHeight(height uint64)
	PutTx(tx Transaction) error
	SetTxListSeal(txListSeal [][]byte)

	// Block의 required field getters
	Seal() []byte
	PrevSeal() []byte
	Height() uint64
	TxList() []Transaction
	TxListSeal() [][]byte

	// Block을 저장을 위한 []byte로 변환 및 재변환
	Serialize() ([]byte, error)
	Deserialize(serializedBlock []byte) error
}

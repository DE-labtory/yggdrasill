package block

import (
	"errors"

	tx "github.com/it-chain/yggdrasill/transaction"
)

var InvalidTransactionTypeError = errors.New("Invalid Transaction Pointer Type Error")

// Block 인터페이스는 Block이 기본적으로 가져야 하는 기능들을 정의한다.
type Block interface {
	// 상태 변경 기능들
	PutTransaction(transaction tx.Transaction) error
	GenerateID() error

	// 상태 확인 기능들
	GetID() string
	GetTransactions() []tx.Transaction
	GetHeight() uint64
	IsPrev(serializedBlock []byte) bool

	// 기타 기능들
	Serialize() ([]byte, error)
}

// BaseBlock 구조체는 Block이 기본적으로 가져야 하는 맴버변수들을 정의한다.
type BaseBlock struct {
	ID     []byte
	PrevID []byte
	Proof  [][]byte
	TxList []tx.Transaction
	Height uint64
}

package block

import (
	"errors"

	tx "github.com/it-chain/yggdrasill/transaction"
)

var InvalidTransactionTypeError = errors.New("Invalid Transaction Pointer Type Error")

//interface에 맞춰 설계
//interface를 implement하는 모든 custom block을 사용 가능하게 구현.
type Block interface {
	PutTransaction(transaction tx.Transaction) error
	Serialize() ([]byte, error)
	GenerateHash() error
	GetHash() string
	GetTransactions() []tx.Transaction
	GetHeight() uint64
	IsPrev(serializedBlock []byte) bool
}

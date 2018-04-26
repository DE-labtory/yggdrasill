package block

import (
	"errors"

	tx "github.com/it-chain/yggdrasill/transaction"
	"github.com/it-chain/yggdrasill/validator"
)

// ErrInvalidTransactionType 에러는 Block에 들어온 Transaction의 Type이 의도한 것이 아닐 때 반환된다.
var ErrInvalidTransactionType = errors.New("Invalid Transaction Type Error")

// Block 구조체는 Block의 구조를 정의한다. Block의 구조를 확장하고 싶으면 BlockInf를 구현하여야 한다.
type Block struct {
	Proof [][]byte
	Manager
}

// Manager 인터페이스는 Block을 관리하기 위한 기능을 제공하는 함수들을 모아둔 인터페이스이다.
type Manager interface {
	PutTransaction(transaction tx.Transaction) error
	Serialize() ([]byte, error)
	GenerateHash() error
	GetHash() string
	GetTransactions() []tx.Transaction
	GetHeight() uint64
	IsPrev(serializedBlock []byte) bool
}

// SetProof 함수는 Block의 Proof 값을 주어진 Validator로 계산하여 할당한다.
func (b Block) SetProof(txList []tx.Transaction, validator validator.Validator) error {
	proof, error := validator.BuildProof(txList)
	if error != nil {
		return error
	}

	b.Proof = proof

	return nil
}

// NewBlock 함수는 새로운 Block을 생성한다.
func NewBlock(txList []tx.Transaction, customBlockManager Manager, validator validator.Validator) Block {
	newBlock := Block{}

	newBlock.SetProof(txList, validator)
	newBlock.Manager = customBlockManager

	return newBlock
}

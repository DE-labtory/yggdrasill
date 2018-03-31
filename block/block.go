package block

import tx "yggdrasill/transaction"

//interface에 맞춰 설계
//interface를 implement하는 모든 custom block을 사용 가능하게 구현.
type Block interface{
	PutTransaction(transaction tx.Transaction)
	FindTransactionIndexByHash(txHash string)
	Serialize() ([]byte, error)
	GenerateHash() error
	GetHash() string
	GetTransactions() []tx.Transaction
	GetHeight() uint64
	IsPrev(serializedBlock []byte) bool
}


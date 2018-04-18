package validator

import tx "github.com/it-chain/yggdrasill/transaction"

// Validator 는 Block의 Transaction들이 위변조 되지 않고 처음 작성된대로 저장되었음을 검증해준다.
// Default 구현체는 Merkle Tree이다.
type Validator interface {
	Validate() bool
	ValidateTransaction(proof []byte, tx *tx.Transaction) bool
	GetProof() []byte
	Serialize() string
	Deserialize(content string)
}

package common

// Validator 는 Block의 Transaction들이 위변조 되지 않고 처음 작성된대로 저장되었음을 검증해준다.
// Default 구현체는 Merkle Tree이다.
type Validator interface {
	ValidateProof(proof []byte, comparisonProof []byte) bool
	ValidateTxProof(txProof [][]byte, txList []Transaction) (bool, error)
	ValidateTransaction(txProof [][]byte, transaction Transaction) (bool, error)
	BuildProofAndTxProof(txList []Transaction) ([]byte, [][]byte, error)
}

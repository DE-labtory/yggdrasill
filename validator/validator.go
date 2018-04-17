package validator

// Validator 는 Block의 Transaction들이 위변조 되지 않고 처음 작성된대로 저장되었음을 검증해준다.
// Default 구현체는 Merkle Tree이다.
type Validator interface {
	Validate() bool
	ValidateContent(proof []byte, content *Content) bool
	GetProof() []byte
}

// Content 는 Validator가 검증하는 대상을 의미한다.
// it-chain은 블록체인이기 때문에 이는 Transaction, 또는 Transaction을 Wrapping한 객체가 된다.
// Default 구현체는 Merkle Tree의 Node로 DefaultTransaction을 Wrapping 한다.
type Content interface {
	CalculateHash() []byte
	Equals(content *Content) bool
}

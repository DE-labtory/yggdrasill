package common

// Transaction interface는 Transaction이 공통적으로 가져야 하는 함수들을 정의함.
type Transaction interface {
	// Transaction의 required field getters
	GetID() string
	GetContent() []byte
	GetSignature() []byte

	// Transaction의 required field setters
	CalculateSeal() ([]byte, error)
	SetSignature([]byte)

	// Transaction의 저장을 위한 []byte로 변환 및 재변환
	Serialize() ([]byte, error)
	Deserialize(serializedBytes []byte) error
}

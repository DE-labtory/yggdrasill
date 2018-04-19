package transaction

// Transaction interface는 it-chain에서 사용되는 모든 Transaction이 공통적으로 가져야 하는 함수들을 정의함.
type Transaction interface {
	Serialize() ([]byte, error)
	GetID() string
	CalculateHash() ([]byte, error)
}

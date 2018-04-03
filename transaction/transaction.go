package transaction

//interface에 맞춰 설계
//interface를 implement하는 모든 custom Transaction 사용 가능하게 구현.
type Transaction interface{
	Serialize() ([]byte, error)
	GetID() string
	CalculateHash() []byte
}


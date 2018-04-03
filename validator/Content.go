package validator

type Content interface{
	CalculateHash() []byte
	Equals(other Content) bool
}

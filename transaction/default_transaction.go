package transaction

import (
	"crypto/sha256"
	"time"

	"github.com/it-chain/yggdrasill/util"
)

// Status 변수는 Transaction의 상태를 Unconfirmed, Confirmed, Unknown 중 하나로 표현함.
type Status int

// TxDataType 변수는 Transaction의 함수가 invoke인지 query인지 표현한다.
type TxDataType string

// Type 은 현재 General 타입만 존재한다.
type Type int

// FunctionType 은 ...
type FunctionType string

// Transaction의 Status를 정의하는 상수들
// TODO: 필요한 것인지 논의가 필요함.
const (
	StatusTransactionUnconfirmed Status = 0
	StatusTransactionConfirmed   Status = 1
	StatusTransactionUnknown     Status = 2
)

// TxData의 Type을 정의하는 상수들
const (
	Invoke TxDataType = "invoke"
	Query  TxDataType = "query"
)

// Transaction의 Type을 정의하는 상수
// TODO: 필요한 것인지 논의가 필요함.
const (
	General Type = 0 + iota
)

// Params 구조체는 Jsonrpc에서 invoke하는 함수의 패러미터를 정의한다.
type Params struct {
	ParamsType int
	Function   string
	Args       []string
}

// TxData 구조체는 Jsonrpc에서 invoke하는 함수를 정의한다.
type TxData struct {
	Jsonrpc string
	Method  TxDataType
	Params  Params
	ID      string
}

// DefaultTransaction 구조체는 Transaction 인터페이스의 기본 구현체이다.
type DefaultTransaction struct {
	ID        string
	Status    Status
	Type      Type
	PeerID    string
	Timestamp time.Time
	TxData    *TxData
}

// Serialize 함수는 Transaction을 []byte 형태로 변환한다.
func (t DefaultTransaction) Serialize() ([]byte, error) {
	return util.Serialize(t)
}

// GetID 함수는 Transaction의 ID 값을 반환한다.
func (t DefaultTransaction) GetID() string {
	return t.ID
}

// CalculateHash 함수는 Transaction 고유의 Hash 값을 계산하여 반환한다.
func (t DefaultTransaction) CalculateHash() ([]byte, error) {
	serializedTx, error := util.Serialize(t)
	if error != nil {
		return nil, error
	}

	return calculateHash(serializedTx), nil
}

func calculateHash(b []byte) []byte {
	hashValue := sha256.New()
	hashValue.Write(b)
	return hashValue.Sum(nil)
}

//func CreateNewTransaction(peer_id string, tx_id string, tx_type Type, t time.Time, data *TxData) *Transaction{
//
//	transaction := &Transaction{PeerID:peer_id, ID:tx_id, Status:Status_TRANSACTION_UNKNOWN, Type:tx_type, Timestamp:t, TxData:data}
//
//	return transaction
//}
//
//func SetTxMethodParameters(params_type int, function string, args []string) Params{
//	return Params{params_type, function, args}
//}
//
//func SetTxData(jsonrpc string, method TxDataType, params Params, contract_id string) *TxData{
//	return &TxData{jsonrpc, method, params, contract_id}
//}

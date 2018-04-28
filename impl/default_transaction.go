package impl

import (
	"encoding/json"
	"time"
)

// Status 변수는 Transaction의 상태를 Unconfirmed, Confirmed, Unknown 중 하나로 표현함.
type Status int

// TxDataType 변수는 Transaction의 함수가 invoke인지 query인지 표현한다.
type TxDataType string

// FunctionType 은 ...
// TODO: FunctionType 타입에 대한 상수값들이 없음
type FunctionType string

// Transaction의 Status를 정의하는 상수들
// TODO: 필요한 것인지 논의가 필요함.
const (
	StatusTransactionInvalid Status = 0
	StatusTransactionValid   Status = 1
)

// TxData의 Type을 정의하는 상수들
const (
	Invoke TxDataType = "invoke"
	Query  TxDataType = "query"
)

// Params 구조체는 Jsonrpc에서 invoke하는 함수의 패러미터를 정의한다.
type Params struct {
	Type     int
	Function string
	Args     []string
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
	PeerID    string
	Timestamp time.Time
	TxData    *TxData
}

// GetID 함수는 Transaction의 ID 값을 반환한다.
func (t *DefaultTransaction) GetID() string {
	return t.ID
}

// CalculateSeal 함수는 Transaction 고유의 Hash 값을 계산하여 반환한다.
func (t *DefaultTransaction) CalculateSeal() ([]byte, error) {
	serializedTx, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	return calculateHash(serializedTx), nil
}

// Serialize 함수는 Transaction을 []byte 형태로 변환한다.
func (t *DefaultTransaction) Serialize() ([]byte, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (t *DefaultTransaction) Deserialize(serializedBytes []byte) error {
	if len(serializedBytes) == 0 {
		return nil
	}

	err := json.Unmarshal(serializedBytes, t)

	if err != nil {
		return err
	}

	return nil
}

// NewDefaultTransaction 함수는 새로운 DefaultTransaction를 반환한다.
func NewDefaultTransaction(peerID string, txID string, timestamp time.Time, txData *TxData) *DefaultTransaction {
	return &DefaultTransaction{
		ID:        txID,
		PeerID:    peerID,
		Timestamp: timestamp,
		TxData:    txData,
		Status:    StatusTransactionInvalid,
	}
}

// NewTxData 함수는 새로운 TxData 객체를 반환한다.
func NewTxData(jsonrpc string, method TxDataType, params Params, contractID string) *TxData {
	return &TxData{
		Jsonrpc: jsonrpc,
		Method:  method,
		Params:  params,
		ID:      contractID,
	}
}

// NewParams 함수는 새로운 Params 객체를 반환한다. (포인터가 아니라 객체 자체를 반환한다.)
func NewParams(paramsType int, function string, args []string) Params {
	return Params{
		Type:     paramsType,
		Function: function,
		Args:     args,
	}
}

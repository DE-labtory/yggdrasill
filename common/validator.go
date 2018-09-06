package common

import "time"

// Validator 는 Block의 Transaction들이 위변조 되지 않고 처음 작성된대로 저장되었음을 검증해준다.
// 위변조가 되지 않음을 증명하는 []byte 값을 Seal(인장)이라고 부르며, 흔히 Hash 값이 사용된다.
// Default 구현체는 Merkle Tree를 기반으로 Seal을 만들고, 검증한다.
type Validator interface {
	// Seal들을 작성해주는 함수들
	BuildSeal(timeStamp time.Time, prevSeal []byte, txSeal [][]byte, creator string) ([]byte, error)
	BuildTxSeal(txList []Transaction) ([][]byte, error)

	// Seal들을 검증해주는 함수들
	ValidateSeal(seal []byte, comparisonBlock Block) (bool, error)
	ValidateTxSeal(txSeal [][]byte, txList []Transaction) (bool, error)
	ValidateTransaction(txSeal [][]byte, transaction Transaction) (bool, error)
}

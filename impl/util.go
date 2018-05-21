package impl

import (
	"crypto/sha256"

	"github.com/it-chain/yggdrasill/common"
)

func calculateHash(b []byte) []byte {
	hashValue := sha256.New()
	hashValue.Write(b)
	return hashValue.Sum(nil)
}

func CalculateHash(b []byte) []byte {
	hashValue := sha256.New()
	hashValue.Write(b)
	return hashValue.Sum(nil)
}

func convertType(txList []*DefaultTransaction) []common.Transaction {
	convTxList := make([]common.Transaction, 0)
	for _, tx := range txList {
		convTxList = append(convTxList, tx)
	}

	return convTxList
}

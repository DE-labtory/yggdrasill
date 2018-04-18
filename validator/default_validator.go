package validator

import (
	"crypto/sha256"

	tx "github.com/it-chain/yggdrasill/transaction"
	"github.com/it-chain/yggdrasill/util"
)

type MerkleTree [][]byte

func (t MerkleTree) Validate() bool {
	return false
}

func (t MerkleTree) ValidateTransaction(proof []byte, tx *tx.Transaction) bool {
	return false
}

func (t MerkleTree) GetProof() []byte {
	if len(t) == 0 {
		return nil
	}

	return t[0]
}

func (t MerkleTree) Serialize() string {
	return ""
}

func (t MerkleTree) Deserialize(content string) {

}

// NewMerkleTree 는 DefaultTransaction 배열을 받아서 MerkleTree 객체를 생성하여 반환한다.
func NewMerkleTree(txList []*tx.DefaultTransaction) (MerkleTree, error) {
	leafNodeList := make([][]byte, 0)

	for _, tx := range txList {
		leafNode, error := calculateLeafNodeHash(tx)
		if error != nil {
			return nil, error
		}

		leafNodeList = append(leafNodeList, leafNode)
	}

	// leafNodeList의 개수는 짝수개로 맞춤. (홀수 일 경우 마지막 Tx를 중복 저장.)
	// TODO: 이래도 되는지 논의 필요.
	if len(leafNodeList)%2 != 0 {
		leafNodeList = append(leafNodeList, leafNodeList[len(leafNodeList)-1])
	}

	return buildTree(leafNodeList, leafNodeList)
}

func buildTree(nodeList [][]byte, fullNodeList [][]byte) ([][]byte, error) {
	intermediateNodeList := make([][]byte, 0)

	for i := 0; i < len(nodeList); i += 2 {
		leftIndex, rightIndex := i, i+1
		leftNode, rightNode := nodeList[leftIndex], nodeList[rightIndex]

		intermediateNode := calculateIntermediateNodeHash(leftNode, rightNode)

		intermediateNodeList = append(intermediateNodeList, intermediateNode)

		if len(nodeList) == 2 {
			return append(intermediateNodeList, fullNodeList...), nil
		}
	}

	newFullNodeList := append(intermediateNodeList, fullNodeList...)

	return buildTree(intermediateNodeList, newFullNodeList)
}

func calculateLeafNodeHash(tx *tx.DefaultTransaction) ([]byte, error) {
	serializedTx, error := util.Serialize(tx)
	if error != nil {
		return nil, error
	}

	return calculateHash(serializedTx), nil
}

func calculateIntermediateNodeHash(leftHash []byte, rightHash []byte) []byte {
	combinedHash := append(leftHash, rightHash...)

	return calculateHash(combinedHash)
}

func calculateHash(b []byte) []byte {
	hashValue := sha256.New()
	hashValue.Write(b)
	return hashValue.Sum(nil)
}

package validator

import (
	"bytes"
	"crypto/sha256"
	"errors"

	tx "github.com/it-chain/yggdrasill/transaction"
)

// SerializationJoinStr 상수는 Serialize()에서 배열을 구성하는 각 해시값들을 Join 할 때 쓰는 구분값
const SerializationJoinStr = " "

// ErrHashCalculationFailed 변수는 Hash 계산 중 발생한 에러를 정의한다.
var ErrHashCalculationFailed = errors.New("Hash Calculation Failed Error")

// MerkleTree 객체는 Validator interface를 구현한 객체.
type MerkleTree struct{}

// Validate 함수는 주어진 Transaction 리스트에 따라 MerkleTree 전체를 검증함.
func (t *MerkleTree) Validate(proof [][]byte, txList []tx.Transaction) (bool, error) {
	leafNodeIndex := 0
	for i, n := range proof {
		leftIndex, rightIndex := (i+1)*2-1, (i+1)*2
		if rightIndex >= len(proof) {
			// Check Leaf Node
			calculatedHash, error := txList[leafNodeIndex].CalculateHash()
			if error != nil {
				return false, ErrHashCalculationFailed
			}

			if bytes.Compare(n, calculatedHash) != 0 {
				return false, nil
			}
			leafNodeIndex++
		} else {
			// Check Intermediate Node
			leftNode, rightNode := proof[leftIndex], proof[rightIndex]
			calculatedHash := calculateIntermediateNodeHash(leftNode, rightNode)
			if bytes.Compare(n, calculatedHash) != 0 {
				return false, nil
			}
		}
	}

	return true, nil
}

// ValidateTransaction 함수는 주어진 Transaction이 이 merkletree에 올바로 있는지를 확인한다.
func (t *MerkleTree) ValidateTransaction(proof [][]byte, tx tx.Transaction) (bool, error) {
	hash, error := tx.CalculateHash()
	if error != nil {
		return false, error
	}

	index := -1
	for i, h := range proof {
		if bytes.Compare(h, hash) == 0 {
			index = i
		}
	}

	if index == -1 {
		return false, nil
	}

	var siblingIndex, parentIndex int
	for index > 0 {
		var isLeft bool
		if index%2 == 0 {
			siblingIndex = index - 1
			parentIndex = (index - 1) / 2
			isLeft = false
		} else {
			siblingIndex = index + 1
			parentIndex = index / 2
			isLeft = true
		}

		var parentHash []byte
		if isLeft {
			parentHash = calculateIntermediateNodeHash(proof[index], proof[siblingIndex])
		} else {
			parentHash = calculateIntermediateNodeHash(proof[siblingIndex], proof[index])
		}

		if bytes.Compare(parentHash, proof[parentIndex]) != 0 {
			return false, nil
		}

		index = parentIndex
	}

	return true, nil
}

// BuildProof 는 DefaultTransaction 배열을 받아서 MerkleTree 객체를 생성하여 반환한다.
func (t *MerkleTree) BuildProof(txList []tx.Transaction) ([][]byte, error) {
	leafNodeList := make([][]byte, 0)

	for _, tx := range txList {
		leafNode, error := tx.CalculateHash()
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

	tree, error := buildTree(leafNodeList, leafNodeList)
	if error != nil {
		return nil, error
	}

	return tree, nil
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

func calculateIntermediateNodeHash(leftHash []byte, rightHash []byte) []byte {
	combinedHash := append(leftHash, rightHash...)

	return calculateHash(combinedHash)
}

func calculateHash(b []byte) []byte {
	hashValue := sha256.New()
	hashValue.Write(b)
	return hashValue.Sum(nil)
}

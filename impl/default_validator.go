package impl

import (
	"bytes"
	"errors"

	"github.com/it-chain/yggdrasill/common"
)

// SerializationJoinStr 상수는 Serialize()에서 배열을 구성하는 각 해시값들을 Join 할 때 쓰는 구분값
const SerializationJoinStr = " "

// ErrHashCalculationFailed 변수는 Hash 계산 중 발생한 에러를 정의한다.
var ErrHashCalculationFailed = errors.New("Hash Calculation Failed Error")

// DefaultValidator 객체는 Validator interface를 구현한 객체.
type DefaultValidator struct{}

// ValidateProof 함수는 원래 Proof 값과 주어진 Proof 값(comparisonProof)을 비교하여, 올바른지 검증한다.
func (t *DefaultValidator) ValidateProof(proof []byte, comparisonProof []byte) bool {
	return bytes.Compare(proof, comparisonProof) == 0
}

// ValidateTxProof 함수는 주어진 Transaction 리스트에 따라 주어진 DefaultValidator 전체(proof)를 검증함.
func (t *DefaultValidator) ValidateTxProof(txProof [][]byte, txList []common.Transaction) (bool, error) {
	leafNodeIndex := 0
	for i, n := range txProof {
		leftIndex, rightIndex := (i+1)*2-1, (i+1)*2
		if rightIndex >= len(txProof) {
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
			leftNode, rightNode := txProof[leftIndex], txProof[rightIndex]
			calculatedHash := calculateIntermediateNodeHash(leftNode, rightNode)
			if bytes.Compare(n, calculatedHash) != 0 {
				return false, nil
			}
		}
	}

	return true, nil
}

// ValidateTransaction 함수는 주어진 Transaction이 이 merkletree(txProof)에 올바로 있는지를 확인한다.
func (t *DefaultValidator) ValidateTransaction(txProof [][]byte, tx common.Transaction) (bool, error) {
	hash, error := tx.CalculateHash()
	if error != nil {
		return false, error
	}

	index := -1
	for i, h := range txProof {
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
			parentHash = calculateIntermediateNodeHash(txProof[index], txProof[siblingIndex])
		} else {
			parentHash = calculateIntermediateNodeHash(txProof[siblingIndex], txProof[index])
		}

		if bytes.Compare(parentHash, txProof[parentIndex]) != 0 {
			return false, nil
		}

		index = parentIndex
	}

	return true, nil
}

// BuildProofAndTxProof 는 DefaultTransaction 배열을 받아서 DefaultValidator 객체와 Proof를 생성하여 반환한다.
// Proof는 주어진 txList의 위변조가 없다는 것을 증명할 []byte 값으로 DefaultValidator의 경우 루트 노드 값을 사용한다.
// TxProof는 개별 transaction들 각각에 대한 Proof 리스트를 의미한다.
func (t *DefaultValidator) BuildProofAndTxProof(txList []common.Transaction) ([]byte, [][]byte, error) {
	leafNodeList := make([][]byte, 0)

	for _, tx := range txList {
		leafNode, error := tx.CalculateHash()
		if error != nil {
			return nil, nil, error
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
		return nil, nil, error
	}

	// DefaultValidator 는 Merkle Tree의 루트노드(tree[0])를 Proof로 간주함
	return tree[0], tree, nil
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

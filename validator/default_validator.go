package validator

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"

	tx "github.com/it-chain/yggdrasill/transaction"
)

// SerializationJoinStr 상수는 Serialize()에서 배열을 구성하는 각 해시값들을 Join 할 때 쓰는 구분값
const SerializationJoinStr = " "

// ErrHashCalculationFailed 변수는 Hash 계산 중 발생한 에러를 정의한다.
var ErrHashCalculationFailed = errors.New("Hash Calculation Failed Error")

// MerkleTree 객체는 Validator interface를 구현한 객체.
// data 프로퍼티는 해시값([]byte)의 배열이다.
type MerkleTree struct {
	data [][]byte
}

// Validate 함수는 주어진 Transaction 리스트에 따라 MerkleTree 전체를 검증함.
func (t *MerkleTree) Validate(txList []tx.Transaction) (bool, error) {
	leafNodeIndex := 0
	for i, n := range t.data {
		leftIndex, rightIndex := (i+1)*2-1, (i+1)*2
		if rightIndex >= len(t.data) {
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
			leftNode, rightNode := t.data[leftIndex], t.data[rightIndex]
			calculatedHash := calculateIntermediateNodeHash(leftNode, rightNode)
			if bytes.Compare(n, calculatedHash) != 0 {
				return false, nil
			}
		}
	}

	return true, nil
}

// ValidateTransaction 함수는 주어진 Transaction이 이 merkletree에 올바로 있는지를 확인한다.
func (t *MerkleTree) ValidateTransaction(tx tx.Transaction) (bool, error) {
	hash, error := tx.CalculateHash()
	if error != nil {
		return false, error
	}

	index := -1
	for i, h := range t.data {
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
			parentHash = calculateIntermediateNodeHash(t.data[index], t.data[siblingIndex])
		} else {
			parentHash = calculateIntermediateNodeHash(t.data[siblingIndex], t.data[index])
		}

		if bytes.Compare(parentHash, t.data[parentIndex]) != 0 {
			return false, nil
		}

		index = parentIndex
	}

	return true, nil
}

// GetProof 함수는 MerkleTree의 루트 값을 반환함.
func (t *MerkleTree) GetProof() []byte {
	if len(t.data) == 0 || t.data[0] == nil {
		return nil
	}

	return t.data[0]
}

// Serialize 함수는 MerkleTree를 구성하는 각 노드(해시값)를 1차원 string 배열로 만든 후,
// 이를 다시 하나의 string 값으로 join 한 후, []byte로 변환하여 반환함.
// 이 반환값이 Block에 저장됨.
func (t *MerkleTree) Serialize() []byte {
	convStrArr := make([]string, 0)
	for _, byteArr := range t.data {
		s := hex.EncodeToString(byteArr)
		convStrArr = append(convStrArr, s)
	}

	return []byte(strings.Join(convStrArr, SerializationJoinStr))
}

// Deserialize 함수는 Serialize() 함수를 통해 변환되었던 string 값을 다시 [][]byte 값으로 변환한 후,
// 이를 MerkleTree 객체의 data 프로퍼티에 할당함.
func (t *MerkleTree) Deserialize(serialized []byte) error {
	if t.data == nil {
		t.data = make([][]byte, 0)
	}

	strArr := strings.Fields(string(serialized[:]))

	for _, str := range strArr {
		bArr, error := hex.DecodeString(str)
		if error != nil {
			return error
		}

		t.data = append(t.data, bArr)
	}

	return nil
}

// NewMerkleTree 는 DefaultTransaction 배열을 받아서 MerkleTree 객체를 생성하여 반환한다.
func NewMerkleTree(txList []tx.Transaction) (*MerkleTree, error) {
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

	return &MerkleTree{data: tree}, nil
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

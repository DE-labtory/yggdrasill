package validator

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	tx "github.com/it-chain/yggdrasill/transaction"
	"github.com/it-chain/yggdrasill/util"
)

// SerializationJoinStr 상수는 Serialize()에서 배열을 구성하는 각 해시값들을 Join 할 때 쓰는 구분값
const SerializationJoinStr = " "

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
			tx, ok := txList[leafNodeIndex].(*tx.DefaultTransaction)
			if ok {
				calculatedHash, error := calculateLeafNodeHash(tx)
				if error != nil {
					return false, errors.New("Hash Calculation Failed Error")
				}

				if bytes.Compare(n, calculatedHash) != 0 {
					fmt.Println(tx.TransactionID)
					return false, nil
				}
			} else {
				return false, errors.New("Type Conversion Failed Error")
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

func (t *MerkleTree) ValidateTransaction(proof []byte, tx tx.Transaction) bool {
	return false
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
func NewMerkleTree(txList []*tx.DefaultTransaction) (*MerkleTree, error) {
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

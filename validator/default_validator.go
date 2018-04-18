package validator

import (
	"crypto/sha256"

	tx "github.com/it-chain/yggdrasill/transaction"
	"github.com/it-chain/yggdrasill/util"
)

type merkleNode struct {
	Hash   []byte
	Parent *merkleNode
	Left   *merkleNode
	Right  *merkleNode
	IsLeaf bool
	Tx     *tx.DefaultTransaction
}

type MerkleTree struct {
	merkleRoot   []byte
	root         *merkleNode
	leafNodeList []*merkleNode
}

func (t *MerkleTree) Validate() bool {
	return false
}

func (t *MerkleTree) ValidateTransaction(proof []byte, tx *tx.Transaction) bool {
	return false
}

func (t *MerkleTree) GetProof() []byte {
	return t.merkleRoot
}

func (t *MerkleTree) Serialize() string {
	return ""
}

func (t *MerkleTree) Deserialize(content string) {

}

// NewMerkleTree 는 DefaultTransaction 배열을 받아서 MerkleTree 객체를 생성하여 반환한다.
func NewMerkleTree(txList []*tx.DefaultTransaction) (*MerkleTree, error) {
	leafNodeList := make([]*merkleNode, 0)

	for _, tx := range txList {
		hashValue, error := calculateLeafNodeHash(tx)
		if error != nil {
			return nil, error
		}

		leafNode := &merkleNode{
			Hash:   hashValue,
			Tx:     tx,
			IsLeaf: true,
		}

		leafNodeList = append(leafNodeList, leafNode)
	}

	// leafNodeList의 개수는 짝수개로 맞춤. (홀수 일 경우 마지막 Tx를 중복 저장.)
	// TODO: 이래도 되는지 논의 필요.
	if len(leafNodeList)%2 != 0 {
		leafNodeList = append(leafNodeList, leafNodeList[len(leafNodeList)-1])
	}

	rootNode, error := buildTree(leafNodeList)
	if error != nil {
		return nil, error
	}

	merkleTree := &MerkleTree{
		merkleRoot:   rootNode.Hash,
		root:         rootNode,
		leafNodeList: leafNodeList,
	}

	return merkleTree, nil
}

func buildTree(leafNodeList []*merkleNode) (*merkleNode, error) {
	intermediateNodeList := make([]*merkleNode, 0)

	for i := 0; i < len(leafNodeList); i += 2 {
		leftIndex, rightIndex := i, i+1
		leftNode, rightNode := leafNodeList[leftIndex], leafNodeList[rightIndex]

		hashValue := calculateIntermediateNodeHash(leftNode.Hash, rightNode.Hash)

		intermediateNode := &merkleNode{
			Hash:   hashValue,
			Left:   leftNode,
			Right:  rightNode,
			IsLeaf: false,
		}

		leftNode.Parent = intermediateNode
		rightNode.Parent = intermediateNode
		intermediateNodeList = append(intermediateNodeList, intermediateNode)

		if len(leafNodeList) == 2 {
			return intermediateNode, nil
		}
	}

	return buildTree(intermediateNodeList)
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

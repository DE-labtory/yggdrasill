package validator

import (
	"github.com/yggdrasill/transaction"
	"crypto/sha256"
	"errors"
	"github.com/yggdrasill/block"
	"bytes"
)

type MerkleTree struct {
	Root     *Node
	RootHash []byte
	Leafs    []*Node
}

type Node struct {
	Parent *Node
	Left   *Node
	Right  *Node
	leaf   bool
	Hash   []byte
	dup    bool
	tx     transaction.Transaction
}

func (n *Node) MakeNodeHash() []byte {
	if n.leaf {
		return n.tx.CalculateHash()
	}
	h := sha256.New()
	if n.Left.Hash != nil && n.Right.Hash != nil {
		h.Write(append(n.Left.Hash, n.Right.Hash...))
	}
	h.Write(append(n.Left.MakeNodeHash(), n.Right.MakeNodeHash()...))
	return h.Sum(nil)
}

func (m *MerkleTree) BuildTree(block block.Block) error {
	txList := block.GetTransactions()
	root, leafs, err := buildWithContent(txList)
	if err != nil {
		return err
	}
	m = &MerkleTree{
		Root:     root,
		RootHash: root.Hash,
		Leafs:    leafs,
	}
	return nil
}

func buildWithContent(txList []transaction.Transaction) (*Node, []*Node, error) {
	if len(txList) == 0 {
		return nil, nil, errors.New("Error: No Content")
	}
	var leafs []*Node
	for _, tx := range txList {
		leafs = append(leafs, &Node{
			Hash: tx.CalculateHash(),
			tx:   tx,
			leaf: true,
		})
	}
	if len(leafs)%2 == 1 {
		duplicate := &Node{
			Hash: leafs[len(leafs)-1].Hash,
			tx:   leafs[len(leafs)-1].tx,
			leaf: true,
			dup:  true,
		}
		leafs = append(leafs, duplicate)
	}
	root := buildIntermediate(leafs)
	return root, leafs, nil
}

func buildIntermediate(nodeLine []*Node) *Node {
	var nodes []*Node
	for i := 0; i < len(nodeLine); i += 2 {
		h := sha256.New()
		var left, right = i, i + 1
		hash := append(nodeLine[left].Hash, nodeLine[right].Hash...)
		h.Write(hash)
		n := &Node{
			Left:  nodeLine[left],
			Right: nodeLine[right],
			Hash:  h.Sum(nil),
		}
		nodes = append(nodes, n)
		nodeLine[left].Parent = n
		nodeLine[right].Parent = n
		if len(nodeLine) == 2 {
			return n
		}
	}
	return buildIntermediate(nodes)
}

func (m *MerkleTree) ReBuildTree() error {
	var txList []transaction.Transaction
	for _, n := range m.Leafs {
		txList = append(txList, n.tx)
	}
	root, leafs, err := buildWithContent(txList)
	if err != nil {
		return err
	}
	m.Root = root
	m.RootHash = root.Hash
	m.Leafs = leafs
	return nil
}

func (m *MerkleTree) VerifyTree() bool {
	calculatedMerkleRootHash := m.Root.MakeNodeHash()
	if bytes.Compare(m.RootHash, calculatedMerkleRootHash) == 0 {
		return true
	}
	return false
}

func (m *MerkleTree) VerifyTx(tx transaction.Transaction) (bool, error) {
	for _, n := range m.Leafs {
		if n.tx == tx {
			currentParent := n.Parent
			for currentParent != nil {
				h := sha256.New()
				h.Write(append(currentParent.Left.Hash, currentParent.Right.Hash...))
				if bytes.Compare(h.Sum(nil), currentParent.Hash) != 0 {
					return false, nil
				}
				currentParent = currentParent.Parent
			}
			return true, nil
		}
	}
	return false, errors.New("Error: Tx is not exist")
}




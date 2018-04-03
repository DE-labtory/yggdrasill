package validator

import "github.com/yggdrasill/transaction"

type MerkleTree struct {
	Root  *Node
	Leafs []*Node
}

type Node struct {
	Parent  *Node
	Left    *Node
	Right   *Node
	leaf    bool
	Hash    []byte
	Content transaction.DefaultTransaction
}


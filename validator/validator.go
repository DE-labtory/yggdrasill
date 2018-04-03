package validator

import (
	"github.com/yggdrasill/block"
	"github.com/yggdrasill/transaction"
)

type Validator interface{
	BuildTree(block block.Block) error
	ReBuildTree() error
	VerifyTx(tx transaction.Transaction) (bool, error)
	VerifyTree() bool
}


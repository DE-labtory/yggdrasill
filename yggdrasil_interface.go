package blockchaindb

import (
	"github.com/it-chain/yggdrasill/block"
	"github.com/it-chain/yggdrasill/transaction"
)

type YggdrasilInterface interface {
	Close()
	AddBlock(block block.Block) error
	GetBlockByNumber(block block.Block, blockNumber uint64) error
	GetBlockByHash(block block.Block, hash string) error
	GetLastBlock(block block.Block) error
	GetTransactionByTxID(transaction transaction.Transaction, txid string) error
	GetBlockByTxID(block block.Block, txid string) error
}
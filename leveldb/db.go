package leveldb

import (
	"errors"
	"yggdrasill/block"
	"github.com/it-chain/leveldb-wrapper"
	"yggdrasill/transaction"
	"fmt"
	"yggdrasill/util"
)

const (
	BLOCK_HASH_DB   = "block_hash"
	BLOCK_NUMBER_DB = "block_number"
	TRANSACTION_DB = "transaction"
	UTIL_DB        = "util"
	LAST_BLOCK_KEY = "last_block"
)

type BlockchainLevelDB struct {
	DBProvider *leveldbwrapper.DBProvider
}

func NewLevelDB(opts map[string]interface{}) *BlockchainLevelDB {
	dbPath := "./default_db"
	if val, ok := opts["db_path"]; ok {
		dbPath = val.(string)
	}
	return &BlockchainLevelDB{DBProvider: leveldbwrapper.CreateNewDBProvider(dbPath)}
}

func (b *BlockchainLevelDB) Close() {
	b.DBProvider.Close()
}

func (b *BlockchainLevelDB) AddBlock(block block.Block) error {
	utilDB := b.DBProvider.GetDBHandle(UTIL_DB)
	lastBlock, err := utilDB.Get([]byte(LAST_BLOCK_KEY))

	if err != nil {
		return err
	}
	if lastBlock != nil && !block.IsPrev(lastBlock) {
		return errors.New("height or prevHash is not matched")
	}

	serializedBlock, err := block.Serialize()
	if err != nil {
		return err
	}

	blockHashDB := b.DBProvider.GetDBHandle(BLOCK_HASH_DB)
	blockNumberDB := b.DBProvider.GetDBHandle(BLOCK_NUMBER_DB)
	transactionDB := b.DBProvider.GetDBHandle(TRANSACTION_DB)

	err = blockHashDB.Put([]byte(block.GetHash()), serializedBlock, true)
	if err != nil {
		return err
	}

	err = blockNumberDB.Put([]byte(fmt.Sprint(block.GetHeight())), []byte(block.GetHash()), true)
	if err != nil {
		return err
	}

	err = utilDB.Put([]byte(LAST_BLOCK_KEY), serializedBlock, true)
	if err != nil {
		return err
	}

	for _, tx := range block.GetTransactions() {
		serializedTX, err := tx.Serialize()
		if err != nil {
			return err
		}

		err = transactionDB.Put([]byte(tx.GetID()), serializedTX, true)
		if err != nil {
			return err
		}

		err = utilDB.Put([]byte(tx.GetID()), []byte(block.GetHash()), true)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *BlockchainLevelDB) GetBlockByNumber(block block.Block, blockNumber uint64) error {
	blockNumberDB := b.DBProvider.GetDBHandle(BLOCK_NUMBER_DB)

	blockHash, err := blockNumberDB.Get([]byte(fmt.Sprint(blockNumber)))
	if err != nil {
		return err
	}

	return b.GetBlockByHash(block, string(blockHash))
}

func (b *BlockchainLevelDB) GetBlockByHash(block block.Block, hash string) error {
	blockHashDB := b.DBProvider.GetDBHandle(BLOCK_HASH_DB)

	serializedBlock, err := blockHashDB.Get([]byte(hash))
	if err != nil {
		return err
	}

	err = util.Deserialize(serializedBlock, block)

	return err
}

func (b *BlockchainLevelDB) GetLastBlock(block block.Block) error {
	utilDB := b.DBProvider.GetDBHandle(UTIL_DB)

	serializedBlock, err := utilDB.Get([]byte(LAST_BLOCK_KEY))
	if serializedBlock == nil || err != nil {
		return err
	}

	err = util.Deserialize(serializedBlock, block)

	return err
}

func (b *BlockchainLevelDB) GetTransactionByTxID(transaction transaction.Transaction, txid string) error {
	transactionDB := b.DBProvider.GetDBHandle(TRANSACTION_DB)

	serializedTX, err := transactionDB.Get([]byte(txid))
	if err != nil {
		return err
	}

	err = util.Deserialize(serializedTX, transaction)

	return err
}

func (b *BlockchainLevelDB) GetBlockByTxID(block block.Block, txid string) error {
	utilDB := b.DBProvider.GetDBHandle(UTIL_DB)

	blockHash, err := utilDB.Get([]byte(txid))
	if err != nil {
		return err
	}

	return b.GetBlockByHash(block, string(blockHash))
}
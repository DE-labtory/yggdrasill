package blockchaindb

import (
	"fmt"

	"github.com/it-chain/leveldb-wrapper/key_value_db"
	"github.com/it-chain/yggdrasill/common"
	"github.com/it-chain/yggdrasill/transaction"
	"github.com/it-chain/yggdrasill/util"
	"github.com/it-chain/yggdrasill/validator"
)

const (
	BLOCK_HASH_DB   = "block_hash"
	BLOCK_NUMBER_DB = "block_number"
	TRANSACTION_DB  = "transaction"
	UTIL_DB         = "util"
	LAST_BLOCK_KEY  = "last_block"
)

type Yggdrasill struct {
	DBProvider *DBProvider
	validator  validator.Validator
}

func NewYggdrasill(keyValueDB key_value_db.KeyValueDB, validator validator.Validator, opts map[string]interface{}) *Yggdrasill {

	dbProvider := CreateNewDBProvider(keyValueDB)

	return &Yggdrasill{DBProvider: dbProvider, validator: validator}
}

func (y *Yggdrasill) Close() {
	y.DBProvider.Close()
}

func (y *Yggdrasill) createGenesisBlock() {

}

func (y *Yggdrasill) AddBlock(block common.Block) error {
	utilDB := y.DBProvider.GetDBHandle(UTIL_DB)

	// TODO: Check the last block
	// lastBlock, err := utilDB.Get([]byte(LAST_BLOCK_KEY))
	// if err != nil {
	// 	return err
	// }

	// if lastBlock != nil && !block.IsPrev(lastBlock) {
	// 	return NewBlockError("height or prevHash is not matched")
	// }

	serializedBlock, err := block.Serialize()
	if err != nil {
		return err
	}

	blockHashDB := y.DBProvider.GetDBHandle(BLOCK_HASH_DB)
	blockNumberDB := y.DBProvider.GetDBHandle(BLOCK_NUMBER_DB)
	transactionDB := y.DBProvider.GetDBHandle(TRANSACTION_DB)

	err = blockHashDB.Put(block.Seal(), serializedBlock, true)
	if err != nil {
		return err
	}

	err = blockNumberDB.Put([]byte(fmt.Sprint(block.Height())), block.Seal(), true)
	if err != nil {
		return err
	}

	err = utilDB.Put([]byte(LAST_BLOCK_KEY), serializedBlock, true)
	if err != nil {
		return err
	}

	for _, tx := range block.TxList() {
		serializedTX, err := tx.Serialize()
		if err != nil {
			return err
		}

		err = transactionDB.Put([]byte(tx.GetID()), serializedTX, true)
		if err != nil {
			return err
		}

		err = utilDB.Put([]byte(tx.GetID()), block.Seal(), true)
		if err != nil {
			return err
		}
	}

	return nil
}

func (y *Yggdrasill) GetBlockByNumber(block common.Block, height uint64) error {
	blockNumberDB := y.DBProvider.GetDBHandle(BLOCK_NUMBER_DB)

	blockHash, err := blockNumberDB.Get([]byte(fmt.Sprint(height)))
	if err != nil {
		return err
	}

	return y.GetBlockBySeal(block, blockHash)
}

func (y *Yggdrasill) GetBlockBySeal(block common.Block, seal []byte) error {
	blockHashDB := y.DBProvider.GetDBHandle(BLOCK_HASH_DB)

	serializedBlock, err := blockHashDB.Get(seal)
	if err != nil {
		return err
	}

	err = util.Deserialize(serializedBlock, block)

	return err
}

func (y *Yggdrasill) GetLastBlock(block common.Block) error {
	utilDB := y.DBProvider.GetDBHandle(UTIL_DB)

	serializedBlock, err := utilDB.Get([]byte(LAST_BLOCK_KEY))
	if serializedBlock == nil || err != nil {
		return err
	}

	err = util.Deserialize(serializedBlock, block)

	return err
}

func (y *Yggdrasill) GetTransactionByTxID(transaction transaction.Transaction, txid string) error {
	transactionDB := y.DBProvider.GetDBHandle(TRANSACTION_DB)

	serializedTX, err := transactionDB.Get([]byte(txid))
	if err != nil {
		return err
	}

	err = util.Deserialize(serializedTX, transaction)

	return err
}

func (y *Yggdrasill) GetBlockByTxID(block common.Block, txid string) error {
	utilDB := y.DBProvider.GetDBHandle(UTIL_DB)

	blockHash, err := utilDB.Get([]byte(txid))

	if err != nil {
		return err
	}

	return y.GetBlockBySeal(block, blockHash)
}

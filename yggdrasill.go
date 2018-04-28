package blockchaindb

import (
	"fmt"

	"github.com/it-chain/leveldb-wrapper/key_value_db"
	"github.com/it-chain/yggdrasill/common"
	"github.com/it-chain/yggdrasill/util"
)

const (
	BLOCK_SEAL_DB   = "block_seal"
	BLOCK_HEIGHT_DB = "block_height"
	TRANSACTION_DB  = "transaction"
	UTIL_DB         = "util"
	LAST_BLOCK_KEY  = "last_block"
)

type Yggdrasill struct {
	DBProvider *DBProvider
	validator  common.Validator
}

func NewYggdrasill(keyValueDB key_value_db.KeyValueDB, validator common.Validator, opts map[string]interface{}) *Yggdrasill {

	dbProvider := CreateNewDBProvider(keyValueDB)

	return &Yggdrasill{DBProvider: dbProvider, validator: validator}
}

func (y *Yggdrasill) Close() {
	y.DBProvider.Close()
}

func (y *Yggdrasill) AddBlock(block common.Block) error {
	utilDB := y.DBProvider.GetDBHandle(UTIL_DB)

	lastBlockByte, err := utilDB.Get([]byte(LAST_BLOCK_KEY))
	if err != nil {
		return err
	}

	if lastBlockByte != nil && !block.IsPrev(lastBlockByte) {
		return NewBlockError("height or prevSeal is not matched")
	}

	serializedBlock, err := block.Serialize()
	if err != nil {
		return err
	}

	blockSealDB := y.DBProvider.GetDBHandle(BLOCK_SEAL_DB)
	blockHeightDB := y.DBProvider.GetDBHandle(BLOCK_HEIGHT_DB)
	transactionDB := y.DBProvider.GetDBHandle(TRANSACTION_DB)

	err = blockSealDB.Put(block.GetSeal(), serializedBlock, true)
	if err != nil {
		return err
	}

	err = blockHeightDB.Put([]byte(fmt.Sprint(block.GetHeight())), block.GetSeal(), true)
	if err != nil {
		return err
	}

	err = utilDB.Put([]byte(LAST_BLOCK_KEY), serializedBlock, true)
	if err != nil {
		return err
	}

	for _, tx := range block.GetTxList() {
		serializedTX, err := tx.Serialize()
		if err != nil {
			return err
		}

		err = transactionDB.Put([]byte(tx.GetID()), serializedTX, true)
		if err != nil {
			return err
		}

		err = utilDB.Put([]byte(tx.GetID()), block.GetSeal(), true)
		if err != nil {
			return err
		}
	}

	return nil
}

func (y *Yggdrasill) GetBlockByHeight(block common.Block, height uint64) error {
	blockHeightDB := y.DBProvider.GetDBHandle(BLOCK_HEIGHT_DB)

	blockSeal, err := blockHeightDB.Get([]byte(fmt.Sprint(height)))
	if err != nil {
		return err
	}

	return y.GetBlockBySeal(block, blockSeal)
}

func (y *Yggdrasill) GetBlockBySeal(block common.Block, seal []byte) error {
	blockSealDB := y.DBProvider.GetDBHandle(BLOCK_SEAL_DB)

	serializedBlock, err := blockSealDB.Get(seal)
	if err != nil {
		return err
	}

	err = block.Deserialize(serializedBlock)

	return err
}

func (y *Yggdrasill) GetLastBlock(block common.Block) error {
	utilDB := y.DBProvider.GetDBHandle(UTIL_DB)

	serializedBlock, err := utilDB.Get([]byte(LAST_BLOCK_KEY))
	if serializedBlock == nil || err != nil {
		return err
	}

	err = block.Deserialize(serializedBlock)

	return err
}

func (y *Yggdrasill) GetTransactionByTxID(transaction common.Transaction, txID string) error {
	transactionDB := y.DBProvider.GetDBHandle(TRANSACTION_DB)

	serializedTX, err := transactionDB.Get([]byte(txID))
	if err != nil {
		return err
	}

	// TODO: add deserialize function to transaction
	err = util.Deserialize(serializedTX, transaction)

	return err
}

func (y *Yggdrasill) GetBlockByTxID(block common.Block, txID string) error {
	utilDB := y.DBProvider.GetDBHandle(UTIL_DB)

	blockSeal, err := utilDB.Get([]byte(txID))

	if err != nil {
		return err
	}

	return y.GetBlockBySeal(block, blockSeal)
}

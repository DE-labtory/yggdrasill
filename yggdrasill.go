package yggdrasill

import (
	"errors"
	"fmt"

	"github.com/it-chain/leveldb-wrapper/key_value_db"
	"github.com/it-chain/yggdrasill/common"
)

const (
	blockSealDB   = "block_seal"
	blockHeightDB = "block_height"
	transactionDB = "transaction"
	utilDB        = "util"
	lastBlockKey  = "last_block"
)

var ErrPrevSealMismatch = errors.New("PrevSeal value mismatch")
var ErrSealValidation = errors.New("seal validation failed")
var ErrTxSealValidation = errors.New("txSeal validation failed")
var ErrNoRequiredParameters = errors.New("required parameters not passed")
var ErrNoValidator = errors.New("validator not defined")

type Yggdrasill struct {
	DBProvider *DBProvider
	validator  common.Validator
}

func NewYggdrasill(keyValueDB key_value_db.KeyValueDB, validator common.Validator, opts map[string]interface{}) (*Yggdrasill, error) {
	if keyValueDB == nil || validator == nil {
		return nil, ErrNoRequiredParameters
	}

	dbProvider := CreateNewDBProvider(keyValueDB)

	return &Yggdrasill{DBProvider: dbProvider, validator: validator}, nil
}

func (y *Yggdrasill) Close() {
	y.DBProvider.Close()
}

func (y *Yggdrasill) AddBlock(block common.Block) error {
	serializedBlock, err := block.Serialize()
	if err != nil {
		return err
	}

	err = y.validateBlock(block)
	if err != nil {
		return err
	}

	utilDB := y.DBProvider.GetDBHandle(utilDB)
	blockSealDB := y.DBProvider.GetDBHandle(blockSealDB)
	blockHeightDB := y.DBProvider.GetDBHandle(blockHeightDB)
	transactionDB := y.DBProvider.GetDBHandle(transactionDB)

	err = blockSealDB.Put(block.GetSeal(), serializedBlock, true)
	if err != nil {
		return err
	}

	err = blockHeightDB.Put([]byte(fmt.Sprint(block.GetHeight())), block.GetSeal(), true)
	if err != nil {
		return err
	}

	err = utilDB.Put([]byte(lastBlockKey), serializedBlock, true)
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

func (y *Yggdrasill) validateBlock(block common.Block) error {
	if y.validator == nil {
		return ErrNoValidator
	}

	utilDB := y.DBProvider.GetDBHandle(UTIL_DB)

	lastBlockByte, err := utilDB.Get([]byte(LAST_BLOCK_KEY))
	if err != nil {
		return err
	}

	// Check if the new block has a correct pointer to the last block
	if lastBlockByte != nil && !block.IsPrev(lastBlockByte) {
		return ErrPrevSealMismatch
	}

	// Validate the Seal of the new block using the validator
	result, err := y.validator.ValidateSeal(block.GetSeal(), block)

	if err != nil {
		return err
	}

	if !result {
		return ErrSealValidation
	}

	// Validate the TxSeal of the new block using the validator
	result, err = y.validator.ValidateTxSeal(block.GetTxSeal(), block.GetTxList())
	if err != nil {
		return err
	}

	if !result {
		return ErrTxSealValidation
	}

	return nil
}

func (y *Yggdrasill) GetBlockByHeight(block common.Block, height uint64) error {
	blockHeightDB := y.DBProvider.GetDBHandle(blockHeightDB)

	blockSeal, err := blockHeightDB.Get([]byte(fmt.Sprint(height)))
	if err != nil {
		return err
	}

	return y.GetBlockBySeal(block, blockSeal)
}

func (y *Yggdrasill) GetBlockBySeal(block common.Block, seal []byte) error {
	blockSealDB := y.DBProvider.GetDBHandle(blockSealDB)

	serializedBlock, err := blockSealDB.Get(seal)
	if err != nil {
		return err
	}

	err = block.Deserialize(serializedBlock)

	return err
}

func (y *Yggdrasill) GetLastBlock(block common.Block) error {
	utilDB := y.DBProvider.GetDBHandle(utilDB)

	serializedBlock, err := utilDB.Get([]byte(lastBlockKey))
	if serializedBlock == nil || err != nil {
		return err
	}

	err = block.Deserialize(serializedBlock)

	return err
}

func (y *Yggdrasill) GetTransactionByTxID(transaction common.Transaction, txID string) error {
	transactionDB := y.DBProvider.GetDBHandle(transactionDB)

	serializedTX, err := transactionDB.Get([]byte(txID))
	if err != nil {
		return err
	}

	err = transaction.Deserialize(serializedTX)

	return err
}

func (y *Yggdrasill) GetBlockByTxID(block common.Block, txID string) error {
	utilDB := y.DBProvider.GetDBHandle(utilDB)

	blockSeal, err := utilDB.Get([]byte(txID))

	if err != nil {
		return err
	}

	return y.GetBlockBySeal(block, blockSeal)
}

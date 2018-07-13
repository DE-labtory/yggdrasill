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

type BlockStorageManager interface {
	Close()
	GetValidator() common.Validator
	AddBlock(block common.Block) error
	GetBlockByHeight(block common.Block, height uint64) error
	GetBlockBySeal(block common.Block, seal []byte) error
	GetBlockByTxID(block common.Block, txid string) error
	GetLastBlock(block common.Block) error
	GetTransactionByTxID(transaction common.Transaction, txid string) error
}

type BlockStorage struct {
	DBProvider *DBProvider
	validator  common.Validator
}

// NewBlockStorage 함수는 새로운 BlockStorage 객체를 생성한다. keyValueDB와 validator는 필수이며, opts는 현재 지원되지 않는다.
func NewBlockStorage(keyValueDB key_value_db.KeyValueDB, validator common.Validator, opts map[string]interface{}) (*BlockStorage, error) {
	if keyValueDB == nil || validator == nil {
		return nil, ErrNoRequiredParameters
	}

	dbProvider := CreateNewDBProvider(keyValueDB)

	return &BlockStorage{DBProvider: dbProvider, validator: validator}, nil
}

// Close 함수는 BlockStorage 객체의 DB를 닫는다.
func (y *BlockStorage) Close() {
	y.DBProvider.Close()
}

// AddBlock 함수는 새로운 Block을 Yggdrasill의 DB에 저장한다. 저장하기 전에 validator로 Block을 검증한다.
func (y *BlockStorage) AddBlock(block common.Block) error {
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

// GetBlockByHeight 함수는 BlockStorage 객체에 저장된 Block을 height 값으로 찾아 반환한다.
func (y *BlockStorage) GetBlockByHeight(block common.Block, height uint64) error {
	blockHeightDB := y.DBProvider.GetDBHandle(blockHeightDB)

	blockSeal, err := blockHeightDB.Get([]byte(fmt.Sprint(height)))
	if err != nil {
		return err
	}

	return y.GetBlockBySeal(block, blockSeal)
}

// GetBlockBySeal 함수는 BlockStorage 객체에 저장된 Block을 seal 값으로 찾아 반환한다.
func (y *BlockStorage) GetBlockBySeal(block common.Block, seal []byte) error {
	blockSealDB := y.DBProvider.GetDBHandle(blockSealDB)

	serializedBlock, err := blockSealDB.Get(seal)
	if err != nil {
		return err
	}

	err = block.Deserialize(serializedBlock)

	return err
}

// GetBlockByTxID 함수는 BlockStorage 객체에 저장된 Block을 Transaction의 ID 값으로 찾아 반환한다.
func (y *BlockStorage) GetBlockByTxID(block common.Block, txID string) error {
	utilDB := y.DBProvider.GetDBHandle(utilDB)

	blockSeal, err := utilDB.Get([]byte(txID))

	if err != nil {
		return err
	}

	return y.GetBlockBySeal(block, blockSeal)
}

// GetLastBlock 함수는 BlockStorage 객체에 저장된 마지막 block을 반환한다.
func (y *BlockStorage) GetLastBlock(block common.Block) error {
	utilDB := y.DBProvider.GetDBHandle(utilDB)

	serializedBlock, err := utilDB.Get([]byte(lastBlockKey))
	if serializedBlock == nil || err != nil {
		return err
	}

	err = block.Deserialize(serializedBlock)

	return err
}

// GetTransactionByTxID 함수는 BlockStorage 객체에 저장된 Block 안에 저장된 Transaction을 ID 값으로 찾아 반환한다.
func (y *BlockStorage) GetTransactionByTxID(transaction common.Transaction, txID string) error {
	transactionDB := y.DBProvider.GetDBHandle(transactionDB)

	serializedTX, err := transactionDB.Get([]byte(txID))
	if err != nil {
		return err
	}

	err = transaction.Deserialize(serializedTX)

	return err
}

func (y *BlockStorage) GetValidator() common.Validator {
	return y.validator
}

func (y *BlockStorage) validateBlock(block common.Block) error {
	if y.validator == nil {
		return ErrNoValidator
	}

	utilDB := y.DBProvider.GetDBHandle(utilDB)

	lastBlockByte, err := utilDB.Get([]byte(lastBlockKey))
	if err != nil {
		return err
	}
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

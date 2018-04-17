package blockchaindb

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/it-chain/leveldb-wrapper"
	"github.com/it-chain/yggdrasill/block"
	"github.com/it-chain/yggdrasill/transaction"
	"github.com/stretchr/testify/assert"
)

func TestYggDrasill_AddBlock(t *testing.T) {

	dbPath := "./.db"
	opts := map[string]interface{}{
		"db_path": dbPath,
	}

	db := leveldbwrapper.CreateNewDB(dbPath)
	y := NewYggdrasil(db, nil, opts)
	defer func() {
		y.Close()
		os.RemoveAll(dbPath)
	}()

	firstBlock := &block.DefaultBlock{Header: &block.BlockHeader{Height: 0, CreatorID: "test"}}
	tx := &transaction.DefaultTransaction{TransactionID: "123"}
	err := firstBlock.PutTransaction(tx)
	assert.NoError(t, err)

	err = y.AddBlock(firstBlock)
	assert.NoError(t, err)

	lastBlock := &block.DefaultBlock{}
	err = y.GetLastBlock(lastBlock)
	assert.NoError(t, err)

	assert.Equal(t, firstBlock.GetHeight(), lastBlock.GetHeight())
	assert.Equal(t, uint64(0), firstBlock.GetHeight())
	assert.Equal(t, uint64(0), lastBlock.GetHeight())
	assert.Equal(t, "test", lastBlock.Header.CreatorID)
	assert.Equal(t, "123", lastBlock.Transactions[0].TransactionID)

	//fmt.Print(lastBlock)
}

//when height did not matched
func TestYggDrasill_AddBlock2(t *testing.T) {

	dbPath := "./.db"
	opts := map[string]interface{}{
		"db_path": dbPath,
	}

	db := leveldbwrapper.CreateNewDB(dbPath)
	y := NewYggdrasil(db, nil, opts)

	defer func() {
		y.Close()
		os.RemoveAll(dbPath)
	}()

	block1 := &block.DefaultBlock{Header: &block.BlockHeader{Height: 0, CreatorID: "test"}}
	block2 := &block.DefaultBlock{Header: &block.BlockHeader{Height: 2, CreatorID: "test"}}

	err := y.AddBlock(block1)
	assert.NoError(t, err)

	err = y.AddBlock(block2)
	assert.Error(t, err)
}

func TestYggdrasil_GetBlockByNumber(t *testing.T) {

	dbPath := "./.db"
	opts := map[string]interface{}{
		"db_path": dbPath,
	}

	db := leveldbwrapper.CreateNewDB(dbPath)
	y := NewYggdrasil(db, nil, opts)
	defer func() {
		y.Close()
		os.RemoveAll(dbPath)
	}()

	for i := 0; i < 100; i++ {
		tmpBlock := &block.DefaultBlock{Header: &block.BlockHeader{Height: uint64(i), CreatorID: fmt.Sprintf("test_%d", i), BlockHash: fmt.Sprintf("hash_%d", i)}}
		if i > 0 {
			tmpBlock.Header.PreviousHash = fmt.Sprintf("hash_%d", i-1)
		}

		err := y.AddBlock(tmpBlock)
		assert.NoError(t, err)
	}

	randomNumber := uint64(rand.Intn(100))

	retrievedBlock := &block.DefaultBlock{}
	err := y.GetBlockByNumber(retrievedBlock, randomNumber)

	assert.NoError(t, err)
	assert.Equal(t, randomNumber, retrievedBlock.GetHeight())
	assert.Equal(t, fmt.Sprintf("test_%d", randomNumber), retrievedBlock.Header.CreatorID)
}

func TestYggdrasil_GetBlockByHash(t *testing.T) {

	dbPath := "./.db"
	opts := map[string]interface{}{
		"db_path": dbPath,
	}

	db := leveldbwrapper.CreateNewDB(dbPath)
	y := NewYggdrasil(db, nil, opts)
	defer func() {
		y.Close()
		os.RemoveAll(dbPath)
	}()

	for i := 0; i < 100; i++ {
		tmpBlock := &block.DefaultBlock{Header: &block.BlockHeader{Height: uint64(i), CreatorID: fmt.Sprintf("test_%d", i), BlockHash: fmt.Sprintf("hash_%d", i)}}
		if i > 0 {
			tmpBlock.Header.PreviousHash = fmt.Sprintf("hash_%d", i-1)
		}

		err := y.AddBlock(tmpBlock)
		assert.NoError(t, err)
	}

	randomNumber := uint64(rand.Intn(100))

	retrievedBlock := &block.DefaultBlock{}
	err := y.GetBlockByHash(retrievedBlock, fmt.Sprintf("hash_%d", randomNumber))

	assert.NoError(t, err)
	assert.Equal(t, randomNumber, retrievedBlock.GetHeight())
	assert.Equal(t, fmt.Sprintf("test_%d", randomNumber), retrievedBlock.Header.CreatorID)
}

func TestYggdrasil_GetLastBlock(t *testing.T) {

	dbPath := "./.db"
	opts := map[string]interface{}{
		"db_path": dbPath,
	}

	db := leveldbwrapper.CreateNewDB(dbPath)
	y := NewYggdrasil(db, nil, opts)
	defer func() {
		y.Close()
		os.RemoveAll(dbPath)
	}()

	for i := 0; i < 100; i++ {
		tmpBlock := &block.DefaultBlock{Header: &block.BlockHeader{Height: uint64(i), CreatorID: fmt.Sprintf("test_%d", i), BlockHash: fmt.Sprintf("hash_%d", i)}}
		if i > 0 {
			tmpBlock.Header.PreviousHash = fmt.Sprintf("hash_%d", i-1)
		}

		err := y.AddBlock(tmpBlock)
		assert.NoError(t, err)
	}

	retrievedBlock := &block.DefaultBlock{}
	err := y.GetLastBlock(retrievedBlock)

	assert.NoError(t, err)
	assert.Equal(t, uint64(99), retrievedBlock.GetHeight())
	assert.Equal(t, "test_99", retrievedBlock.Header.CreatorID)
}

func TestYggdrasil_GetTransactionByTxID(t *testing.T) {

	dbPath := "./.db"
	opts := map[string]interface{}{
		"db_path": dbPath,
	}

	db := leveldbwrapper.CreateNewDB(dbPath)
	y := NewYggdrasil(db, nil, opts)
	defer func() {
		y.Close()
		os.RemoveAll(dbPath)
	}()
}

func TestYggdrasil_GetBlockByTxID(t *testing.T) {

	dbPath := "./.db"
	opts := map[string]interface{}{
		"db_path": dbPath,
	}

	db := leveldbwrapper.CreateNewDB(dbPath)
	y := NewYggdrasil(db, nil, opts)
	defer func() {
		y.Close()
		os.RemoveAll(dbPath)
	}()
}

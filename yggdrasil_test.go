package blockchainleveldb

import (
	"fmt"
	"os"
	"testing"

	"github.com/it-chain/yggdrasill/block"
	"github.com/stretchr/testify/assert"
)

func TestYggDrasill_AddBlock(t *testing.T) {
	dbPath := "./.db"

	defer os.RemoveAll(dbPath)
	y := NewYggdrasil(dbPath, nil)
	block1 := block.DefaultBlock{Header: block.BlockHeader{Height: 0, CreatorID: "jun"}}
	err := y.AddBlock(block1)

	if err != nil {
		fmt.Print(err.Error())
	}

	block2 := &block.DefaultBlock{}

	if err != nil {
		fmt.Print(err.Error())
	}

	y.GetLastBlock(block2)

	assert.Equal(t, block1.GetHeight(), block2.GetHeight())
	assert.Equal(t, uint64(0), block1.GetHeight())
	assert.Equal(t, uint64(0), block2.GetHeight())
	assert.Equal(t, "jun", block2.Header.CreatorID)
}

//when height did not matched
func TestYggDrasill_AddBlock2(t *testing.T) {
	dbPath := "./.db"

	defer os.RemoveAll(dbPath)
	y := NewYggdrasil(dbPath, nil)
	block1 := block.DefaultBlock{Header: block.BlockHeader{Height: 0, CreatorID: "jun"}}
	block2 := block.DefaultBlock{Header: block.BlockHeader{Height: 2, CreatorID: "jun"}}
	err := y.AddBlock(block1)

	if err != nil {
		fmt.Print(err.Error())
	}

	err = y.AddBlock(block2)

	fmt.Print(err.Error())
	assert.Error(t, err)
}

//func TestBlockchainLevelDB_AddBlock(t *testing.T) {
//	path := "./test_db"
//	blockchainLevelDB := CreateNewBlockchainLevelDB(path)
//	defer func(){
//		blockchainLevelDB.Close()
//		os.RemoveAll(path)
//	}()
//
//	block := domain.CreateNewBlock(nil, "test")
//
//	err := blockchainLevelDB.AddBlock(block)
//	assert.NoError(t, err)
//}
//
//func TestBlockchainLevelDB_GetBlockByNumber(t *testing.T) {
//	path := "./test_db"
//	blockchainLevelDB := CreateNewBlockchainLevelDB(path)
//	defer func(){
//		blockchainLevelDB.Close()
//		os.RemoveAll(path)
//	}()
//
//	block := domain.CreateNewBlock(nil, "test")
//	blockNumber := block.Header.Number
//
//	err := blockchainLevelDB.AddBlock(block)
//	assert.NoError(t, err)
//
//	retrievedBlock, err := blockchainLevelDB.GetBlockByNumber(blockNumber)
//	assert.NoError(t, err)
//	assert.Equal(t, block, retrievedBlock)
//}
//
//func TestBlockchainLevelDB_GetBlockByHash(t *testing.T) {
//	path := "./test_db"
//	blockchainLevelDB := CreateNewBlockchainLevelDB(path)
//	defer func(){
//		blockchainLevelDB.Close()
//		os.RemoveAll(path)
//	}()
//
//	block := domain.CreateNewBlock(nil, "test")
//	blockHash := block.Header.BlockHash
//
//	err := blockchainLevelDB.AddBlock(block)
//	assert.NoError(t, err)
//
//	retrievedBlock, err := blockchainLevelDB.GetBlockByHash(blockHash)
//	assert.NoError(t, err)
//	assert.Equal(t, block, retrievedBlock)
//}
//
//func TestBlockchainLevelDB_GetLastBlock(t *testing.T) {
//	path := "./test_db"
//	blockchainLevelDB := CreateNewBlockchainLevelDB(path)
//	defer func(){
//		blockchainLevelDB.Close()
//		os.RemoveAll(path)
//	}()
//
//	block1 := domain.CreateNewBlock(nil, "test1")
//	block2 := domain.CreateNewBlock(nil, "test2")
//
//	err := blockchainLevelDB.AddBlock(block1)
//	assert.NoError(t, err)
//
//	lastBlock, err := blockchainLevelDB.GetLastBlock()
//	assert.NoError(t, err)
//	assert.Equal(t, block1, lastBlock)
//
//	err = blockchainLevelDB.AddBlock(block2)
//	assert.NoError(t, err)
//
//	lastBlock, err = blockchainLevelDB.GetLastBlock()
//	assert.NoError(t, err)
//	assert.Equal(t, block2, lastBlock)
//}
//
//func TestBlockchainLevelDB_GetTransactionByTxID(t *testing.T) {
//	path := "./test_db"
//	blockchainLevelDB := CreateNewBlockchainLevelDB(path)
//	defer func(){
//		blockchainLevelDB.Close()
//		os.RemoveAll(path)
//	}()
//
//	block := domain.CreateNewBlock(nil, "test")
//	tx := domain.CreateNewTransaction(
//		"test",
//		"test",
//		0,
//		time.Now().Round(0),
//		&domain.TxData{})
//	tx.GenerateHash()
//	err :=block.PutTranscation(tx)
//	assert.NoError(t, err)
//
//	err = blockchainLevelDB.AddBlock(block)
//	assert.NoError(t, err)
//
//	retrievedTx, err := blockchainLevelDB.GetTransactionByTxID("test")
//	assert.NoError(t, err)
//	assert.Equal(t, tx, retrievedTx)
//}
//
//func TestBlockchainLevelDB_GetBlockByTxID(t *testing.T) {
//	path := "./test_db"
//	blockchainLevelDB := CreateNewBlockchainLevelDB(path)
//	defer func(){
//		blockchainLevelDB.Close()
//		os.RemoveAll(path)
//	}()
//
//	block := domain.CreateNewBlock(nil, "test")
//	tx := domain.CreateNewTransaction(
//		"test",
//		"test",
//		0,
//		time.Now().Round(0),
//		&domain.TxData{})
//	tx.GenerateHash()
//	err := block.PutTranscation(tx)
//	assert.NoError(t, err)
//
//	err = blockchainLevelDB.AddBlock(block)
//	assert.NoError(t, err)
//
//	retrievedBlock, err := blockchainLevelDB.GetBlockByTxID("test")
//	assert.NoError(t, err)
//	assert.Equal(t, block, retrievedBlock)
//}

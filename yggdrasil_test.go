package blockchainleveldb

import (
	"testing"
	"github.com/it-chain/yggdrasill/block"
	"fmt"
)

func TestYggDrasill_GetLastBlock(t *testing.T) {
	//dbPath := "~/.db"
	//y := NewYggdrasil(dbPath,nil)
	//y.AddBlock()


}

func TestDeserialize(t *testing.T){

	var block1 block.Block

	block1 = block.DefaultBlock{Header:block.BlockHeader{Height:1,CreatorID:"jun"}}

	b, err := block1.Serialize()

	if err != nil{

	}

	block2 := &block.DefaultBlock{}

	err = deserialize(b,block2)

	if err != nil{

	}

	fmt.Print(block1)
	fmt.Print(block2)
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
package blockchaindb

import (
	"os"
	"testing"

	"github.com/it-chain/leveldb-wrapper"
	"github.com/it-chain/yggdrasill/impl"
	"github.com/it-chain/yggdrasill/transaction"
	"github.com/stretchr/testify/assert"
)

func TestYggDrasill_AddBlock(t *testing.T) {

	dbPath := "./.db"
	opts := map[string]interface{}{
		"db_path": dbPath,
	}

	db := leveldbwrapper.CreateNewDB(dbPath)
	y := NewYggdrasill(db, nil, opts)
	defer func() {
		y.Close()
		os.RemoveAll(dbPath)
	}()

	firstBlock := &impl.DefaultBlock{CreatorID: "test"}
	tx := &transaction.DefaultTransaction{ID: "123"}
	err := firstBlock.PutTx(tx)
	assert.NoError(t, err)

	err = y.AddBlock(firstBlock)
	assert.NoError(t, err)

	lastBlock := &impl.DefaultBlock{}
	err = y.GetLastBlock(lastBlock)
	assert.NoError(t, err)

	assert.Equal(t, firstBlock.Height(), lastBlock.Height())
	assert.Equal(t, uint64(0), firstBlock.Height())
	assert.Equal(t, uint64(0), lastBlock.Height())
	assert.Equal(t, "test", lastBlock.CreatorID)
	assert.Equal(t, "123", lastBlock.TxList()[0].GetID())

	//fmt.Print(lastBlock)
}

//when height did not matched
func TestYggDrasill_AddBlock2(t *testing.T) {

	dbPath := "./.db"
	opts := map[string]interface{}{
		"db_path": dbPath,
	}

	db := leveldbwrapper.CreateNewDB(dbPath)
	y := NewYggdrasill(db, nil, opts)

	defer func() {
		y.Close()
		os.RemoveAll(dbPath)
	}()

	block1 := &impl.DefaultBlock{CreatorID: "test"}
	block1.SetHeight(0)
	block2 := &impl.DefaultBlock{CreatorID: "test"}
	block2.SetHeight(2)

	err := y.AddBlock(block1)
	assert.NoError(t, err)

	err = y.AddBlock(block2)
	assert.Error(t, err)
}

// func TestYggdrasil_GetBlockByNumber(t *testing.T) {

// 	dbPath := "./.db"
// 	opts := map[string]interface{}{
// 		"db_path": dbPath,
// 	}

// 	db := leveldbwrapper.CreateNewDB(dbPath)
// 	y := NewYggdrasill(db, nil, opts)
// 	defer func() {
// 		y.Close()
// 		os.RemoveAll(dbPath)
// 	}()

// 	for i := 0; i < 100; i++ {
// 		tmpBlock := &impl.DefaultBlock{CreatorID: fmt.Sprintf("test_%d", i)}
// 		tmpBlock.SetHeight(uint64(i))
// 		if i > 0 {
// 			tmpBlock.SetPrevSeal([]byte(fmt.Sprintf("hash_%d", i-1)))
// 			tmpBlock.SetTxListSeal([][]byte{[]byte(fmt.Sprintf("txListSeal_%d", i-1))})
// 			tmpBlock.SetTimestamp(time.Now())
// 		}
// 		tmpBlock.GenerateSeal()

// 		err := y.AddBlock(tmpBlock)
// 		assert.NoError(t, err)
// 	}

// 	randomNumber := uint64(rand.Intn(100))

// 	retrievedBlock := &impl.DefaultBlock{}
// 	err := y.GetBlockByNumber(retrievedBlock, randomNumber)

// 	assert.NoError(t, err)
// 	assert.Equal(t, randomNumber, retrievedBlock.GetHeight())
// 	assert.Equal(t, fmt.Sprintf("test_%d", randomNumber), retrievedBlock.Header.CreatorID)
// }

// func TestYggdrasil_GetBlockBySeal(t *testing.T) {

// 	dbPath := "./.db"
// 	opts := map[string]interface{}{
// 		"db_path": dbPath,
// 	}

// 	db := leveldbwrapper.CreateNewDB(dbPath)
// 	y := NewYggdrasill(db, nil, opts)
// 	defer func() {
// 		y.Close()
// 		os.RemoveAll(dbPath)
// 	}()

// 	for i := 0; i < 100; i++ {
// 		tmpBlock := &impl.DefaultBlock{Header: &impl.BlockHeader{Height: uint64(i), CreatorID: fmt.Sprintf("test_%d", i), BlockHash: fmt.Sprintf("hash_%d", i)}}
// 		if i > 0 {
// 			tmpBlock.Header.PreviousHash = fmt.Sprintf("hash_%d", i-1)
// 		}

// 		err := y.AddBlock(tmpBlock)
// 		assert.NoError(t, err)
// 	}

// 	randomNumber := uint64(rand.Intn(100))

// 	retrievedBlock := &impl.DefaultBlock{}
// 	err := y.GetBlockBySeal(retrievedBlock, []byte(fmt.Sprintf("hash_%d", randomNumber)))

// 	assert.NoError(t, err)
// 	assert.Equal(t, randomNumber, retrievedBlock.GetHeight())
// 	assert.Equal(t, fmt.Sprintf("test_%d", randomNumber), retrievedBlock.Header.CreatorID)
// }

// func TestYggdrasil_GetLastBlock(t *testing.T) {

// 	dbPath := "./.db"
// 	opts := map[string]interface{}{
// 		"db_path": dbPath,
// 	}

// 	db := leveldbwrapper.CreateNewDB(dbPath)
// 	y := NewYggdrasill(db, nil, opts)
// 	defer func() {
// 		y.Close()
// 		os.RemoveAll(dbPath)
// 	}()

// 	for i := 0; i < 100; i++ {
// 		tmpBlock := &impl.DefaultBlock{Header: &impl.BlockHeader{Height: uint64(i), CreatorID: fmt.Sprintf("test_%d", i), BlockHash: fmt.Sprintf("hash_%d", i)}}
// 		if i > 0 {
// 			tmpBlock.Header.PreviousHash = fmt.Sprintf("hash_%d", i-1)
// 		}

// 		err := y.AddBlock(tmpBlock)
// 		assert.NoError(t, err)
// 	}

// 	retrievedBlock := &impl.DefaultBlock{}
// 	err := y.GetLastBlock(retrievedBlock)

// 	assert.NoError(t, err)
// 	assert.Equal(t, uint64(99), retrievedBlock.GetHeight())
// 	assert.Equal(t, "test_99", retrievedBlock.Header.CreatorID)
// }

// func TestYggdrasil_GetTransactionByTxID(t *testing.T) {

// 	//given
// 	dbPath := "./.db"
// 	opts := map[string]interface{}{
// 		"db_path": dbPath,
// 	}

// 	db := leveldbwrapper.CreateNewDB(dbPath)
// 	y := NewYggdrasill(db, nil, opts)
// 	defer func() {
// 		y.Close()
// 		os.RemoveAll(dbPath)
// 	}()

// 	firstBlock := &impl.DefaultBlock{Header: &impl.BlockHeader{Height: 0, CreatorID: "test"}}
// 	tx := &transaction.DefaultTransaction{TransactionID: "123"}
// 	err := firstBlock.PutTransaction(tx)
// 	assert.NoError(t, err)

// 	err = y.AddBlock(firstBlock)
// 	assert.NoError(t, err)

// 	//when
// 	retrievedTx := &transaction.DefaultTransaction{}
// 	err = y.GetTransactionByTxID(retrievedTx, tx.TransactionID)
// 	assert.NoError(t, err)

// 	//then
// 	assert.Equal(t, retrievedTx, tx)
// }

// func TestYggdrasil_GetBlockByTxID(t *testing.T) {

// 	//given
// 	dbPath := "./.db"
// 	opts := map[string]interface{}{
// 		"db_path": dbPath,
// 	}

// 	db := leveldbwrapper.CreateNewDB(dbPath)
// 	y := NewYggdrasill(db, nil, opts)
// 	defer func() {
// 		y.Close()
// 		os.RemoveAll(dbPath)
// 	}()

// 	firstBlock := &impl.DefaultBlock{Header: &impl.BlockHeader{Height: 0, CreatorID: "test"}}
// 	tx := &transaction.DefaultTransaction{TransactionID: "123"}
// 	err := firstBlock.PutTransaction(tx)
// 	assert.NoError(t, err)

// 	err = y.AddBlock(firstBlock)
// 	assert.NoError(t, err)

// 	//when
// 	retrievedBlock := &impl.DefaultBlock{}
// 	err = y.GetBlockByTxID(retrievedBlock, tx.TransactionID)
// 	assert.NoError(t, err)

// 	//then
// 	assert.Equal(t, firstBlock, retrievedBlock)
// }

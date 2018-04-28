package blockchaindb

import (
	"os"
	"testing"

	"time"

	"github.com/it-chain/leveldb-wrapper"
	"github.com/it-chain/yggdrasill/common"
	"github.com/it-chain/yggdrasill/impl"
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

	firstBlock := getNewBlock([]byte("genesis"), 0)
	err := y.AddBlock(firstBlock)
	assert.NoError(t, err)

	lastBlock := &impl.DefaultBlock{}
	err = y.GetLastBlock(lastBlock)
	assert.NoError(t, err)

	assert.Equal(t, firstBlock.GetHeight(), lastBlock.GetHeight())
	assert.Equal(t, uint64(0), firstBlock.GetHeight())
	assert.Equal(t, uint64(0), lastBlock.GetHeight())
	assert.Equal(t, []byte("testUser"), lastBlock.GetCreator())
	assert.Equal(t, "tx01", lastBlock.GetTxList()[0].GetID())

	//fmt.Print(lastBlock)
}

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

	block1 := getNewBlock([]byte("genesis"), 0)
	block2 := getNewBlock(block1.GetSeal(), 1)

	err := y.AddBlock(block1)
	assert.NoError(t, err)

	err = y.AddBlock(block2)
	assert.NoError(t, err)
}

// PrevSeal 값을 잘못 입력해서 에러를 출력.
func TestYggDrasill_AddBlock3(t *testing.T) {

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

	block1 := getNewBlock([]byte("genesis"), 0)
	block2 := getNewBlock([]byte("genesis"), 1)

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
// 	err := y.GetBlockByHeight(retrievedBlock, randomNumber)

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

func getNewBlock(prevSeal []byte, height uint64) *impl.DefaultBlock {
	validator := &impl.DefaultValidator{}
	testingTime, _ := time.Parse("Jan 2, 2006 at 3:04pm (MST)", "Feb 3, 2013 at 7:54pm (UTC)")
	blockCreator := []byte("testUser")
	txList := getTxList(testingTime)

	block := impl.NewEmptyBlock(prevSeal, height, blockCreator)
	block.SetTimestamp(testingTime)
	for _, tx := range txList {
		block.PutTx(tx)
	}

	txSeal, _ := validator.BuildTxSeal(convertTxListType(txList))
	block.SetTxSeal(txSeal)

	seal, _ := validator.BuildSeal(block)
	block.SetSeal(seal)

	return block
}

func getTxList(testingTime time.Time) []*impl.DefaultTransaction {
	return []*impl.DefaultTransaction{
		{
			PeerID:    "p01",
			ID:        "tx01",
			Status:    0,
			Timestamp: testingTime,
			TxData: &impl.TxData{
				Jsonrpc: "jsonRPC01",
				Method:  "invoke",
				Params: impl.Params{
					Type:     0,
					Function: "function01",
					Args:     []string{"arg1", "arg2"},
				},
				ID: "txdata01",
			},
		},
		{
			PeerID:    "p02",
			ID:        "tx02",
			Status:    0,
			Timestamp: testingTime,
			TxData: &impl.TxData{
				Jsonrpc: "jsonRPC02",
				Method:  "invoke",
				Params: impl.Params{
					Type:     0,
					Function: "function02",
					Args:     []string{"arg1", "arg2"},
				},
				ID: "txdata02",
			},
		},
		{
			PeerID:    "p03",
			ID:        "tx03",
			Status:    0,
			Timestamp: testingTime,
			TxData: &impl.TxData{
				Jsonrpc: "jsonRPC03",
				Method:  "invoke",
				Params: impl.Params{
					Type:     0,
					Function: "function03",
					Args:     []string{"arg1", "arg2"},
				},
				ID: "txdata03",
			},
		},
		{
			PeerID:    "p04",
			ID:        "tx04",
			Status:    0,
			Timestamp: testingTime,
			TxData: &impl.TxData{
				Jsonrpc: "jsonRPC04",
				Method:  "invoke",
				Params: impl.Params{
					Type:     0,
					Function: "function04",
					Args:     []string{"arg1", "arg2"},
				},
				ID: "txdata04",
			},
		},
	}
}

func convertTxListType(txList []*impl.DefaultTransaction) []common.Transaction {
	convTxList := make([]common.Transaction, 0)
	for _, tx := range txList {
		convTxList = append(convTxList, tx)
	}

	return convTxList
}

package yggdrasill

import (
	"math/rand"
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

func TestYggdrasill_GetBlockByHeight(t *testing.T) {

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

	prevSeal := []byte("genesis")
	for i := 0; i < 100; i++ {
		tmpBlock := getNewBlock(prevSeal, uint64(i))
		err := y.AddBlock(tmpBlock)
		assert.NoError(t, err)

		prevSeal = tmpBlock.GetSeal()
	}

	randomNumber := uint64(rand.Intn(100))

	retrievedBlock := &impl.DefaultBlock{}
	err := y.GetBlockByHeight(retrievedBlock, randomNumber)

	assert.NoError(t, err)
	assert.Equal(t, randomNumber, retrievedBlock.GetHeight())
}

func TestYggdrasil_GetBlockBySeal(t *testing.T) {

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

	prevSeal := []byte("genesis")
	randomNumber := uint64(rand.Intn(100))
	var testSeal []byte
	for i := 0; i < 100; i++ {
		tmpBlock := getNewBlock(prevSeal, uint64(i))
		err := y.AddBlock(tmpBlock)

		assert.NoError(t, err)

		prevSeal = tmpBlock.GetSeal()
		if uint64(i) == randomNumber {
			testSeal = tmpBlock.GetSeal()
		}
	}

	retrievedBlock := &impl.DefaultBlock{}
	err := y.GetBlockBySeal(retrievedBlock, testSeal)

	assert.NoError(t, err)
	assert.Equal(t, randomNumber, retrievedBlock.GetHeight())
	assert.Equal(t, testSeal, retrievedBlock.GetSeal())
}

func TestYggdrasil_GetLastBlock(t *testing.T) {

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

	prevSeal := []byte("genesis")
	var lastSeal []byte
	for i := 0; i < 100; i++ {
		tmpBlock := getNewBlock(prevSeal, uint64(i))
		err := y.AddBlock(tmpBlock)

		assert.NoError(t, err)

		prevSeal = tmpBlock.GetSeal()

		if i == 99 {
			lastSeal = tmpBlock.GetSeal()
		}
	}

	retrievedBlock := &impl.DefaultBlock{}
	err := y.GetLastBlock(retrievedBlock)

	assert.NoError(t, err)
	assert.Equal(t, uint64(99), retrievedBlock.GetHeight())
	assert.Equal(t, lastSeal, retrievedBlock.GetSeal())
}

func TestYggdrasil_GetTransactionByTxID(t *testing.T) {

	//given
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

	//when
	retrievedTx := &impl.DefaultTransaction{}
	err = y.GetTransactionByTxID(retrievedTx, "tx01")
	assert.NoError(t, err)

	//then
	assert.Equal(t, retrievedTx, getTxList(getTime())[0])
}

func TestYggdrasil_GetBlockByTxID(t *testing.T) {

	//given
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

	//when
	retrievedBlock := &impl.DefaultBlock{}
	err = y.GetBlockByTxID(retrievedBlock, "tx01")
	assert.NoError(t, err)

	//then
	assert.Equal(t, firstBlock, retrievedBlock)
}

func getNewBlock(prevSeal []byte, height uint64) *impl.DefaultBlock {
	validator := &impl.DefaultValidator{}
	testingTime := getTime()
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

func getTime() time.Time {
	testingTime, _ := time.Parse("Jan 2, 2006 at 3:04pm (MST)", "Feb 3, 2013 at 7:54pm (UTC)")
	return testingTime
}

func convertTxListType(txList []*impl.DefaultTransaction) []common.Transaction {
	convTxList := make([]common.Transaction, 0)
	for _, tx := range txList {
		convTxList = append(convTxList, tx)
	}

	return convTxList
}

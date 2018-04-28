package impl

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewEmptyBlock(t *testing.T) {
	block := getNewBlock()

	expected := []byte{170, 156, 92, 136, 64, 227, 248, 194, 78, 168, 107, 144, 205, 66, 234, 40, 204, 27, 117, 52, 199, 24, 32, 245, 115, 97, 146, 217, 14, 104, 227, 165}
	assert.Equal(t, block.GetSeal(), expected)
}

func TestSerializeAndDeserialize(t *testing.T) {
	block := getNewBlock()

	serializedBlock, err := block.Serialize()
	assert.NoError(t, err)

	deserializedBlock := &DefaultBlock{}
	err = deserializedBlock.Deserialize(serializedBlock)
	assert.NoError(t, err)

	assert.Equal(t, deserializedBlock, block)
}

func getNewBlock() *DefaultBlock {
	validator := &DefaultValidator{}
	testingTime, _ := time.Parse("Jan 2, 2006 at 3:04pm (MST)", "Feb 3, 2013 at 7:54pm (UTC)")
	blockCreator := []byte("testUser")
	genesisSeal := []byte("genesis")
	txList := getTxList(testingTime)

	block := NewEmptyBlock(genesisSeal, 0, blockCreator)
	block.SetTimestamp(testingTime)
	for _, tx := range txList {
		block.PutTx(tx)
	}

	txSeal, _ := validator.BuildTxSeal(convertType(txList))
	block.SetTxSeal(txSeal)

	seal, _ := validator.BuildSeal(block)
	block.SetSeal(seal)

	return block
}

func getTxList(testingTime time.Time) []*DefaultTransaction {
	return []*DefaultTransaction{
		{
			PeerID:    "p01",
			ID:        "tx01",
			Status:    0,
			Timestamp: testingTime,
			TxData: &TxData{
				Jsonrpc: "jsonRPC01",
				Method:  "invoke",
				Params: Params{
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
			TxData: &TxData{
				Jsonrpc: "jsonRPC02",
				Method:  "invoke",
				Params: Params{
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
			TxData: &TxData{
				Jsonrpc: "jsonRPC03",
				Method:  "invoke",
				Params: Params{
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
			TxData: &TxData{
				Jsonrpc: "jsonRPC04",
				Method:  "invoke",
				Params: Params{
					Type:     0,
					Function: "function04",
					Args:     []string{"arg1", "arg2"},
				},
				ID: "txdata04",
			},
		},
	}
}

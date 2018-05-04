package impl

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewEmptyBlock(t *testing.T) {
	block := getNewBlock()

	expected := []byte{26, 156, 70, 177, 186, 43, 248, 224, 3, 35, 95, 141, 188, 119, 78, 150, 234, 255, 250, 238, 211, 69, 72, 231, 88, 240, 25, 253, 75, 86, 74, 253}
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
	txList := getTestingTxList(0)

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

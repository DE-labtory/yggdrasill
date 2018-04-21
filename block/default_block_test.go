package block

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewBlock(t *testing.T) {
	_, err := CreateNewBlock(nil, "")
	assert.Error(t, err)

	block2, err := CreateNewBlock(nil, "Genesis")
	assert.Equal(t, "Genesis", block2.Header.CreatorID)

	block3, err := CreateNewBlock(block2, "JunkSound")
	assert.Equal(t, uint64(1), block3.Header.Height)
	assert.Equal(t, "JunkSound", block3.Header.CreatorID)
	assert.Equal(t, "", block3.Header.PreviousHash)
	assert.Equal(t, 0, block3.Header.MerkleTreeHeight)
	assert.Equal(t, 0, block3.Header.TransactionCount)

}

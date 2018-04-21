package block

import (
	"fmt"
	"testing"
	"time"

	tx "github.com/it-chain/yggdrasill/transaction"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewBlock(t *testing.T) {
	_, err := CreateNewBlock(nil, "")
	assert.Error(t, err)

	block1, err := CreateNewBlock(nil, "Genesis")
	assert.Equal(t, "Genesis", block1.Header.CreatorID)

	block2, err := CreateNewBlock(block1, "JunkSound")
	assert.Equal(t, uint64(1), block2.Header.Height)
	assert.Equal(t, "JunkSound", block2.Header.CreatorID)
	assert.Equal(t, "", block2.Header.PreviousHash)
	assert.Equal(t, 0, block2.Header.MerkleTreeHeight)
	assert.Equal(t, 0, block2.Header.TransactionCount)
}

func TestCreateGenesisBlock(t *testing.T) {
	GenesisBlock, err := CreateGenesisBlock()
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), GenesisBlock.Header.Height)
	assert.Equal(t, "", GenesisBlock.Header.PreviousHash)
	assert.Equal(t, "", GenesisBlock.Header.Version)
	assert.Equal(t, "", GenesisBlock.Header.MerkleTreeRootHash)
	assert.Equal(t, time.Now().String()[:19], GenesisBlock.Header.TimeStamp.String()[:19])
	assert.Equal(t, "Genesis", GenesisBlock.Header.CreatorID)
	fmt.Println()
	assert.Equal(t, make([]uint8, 0), GenesisBlock.Header.Signature)
	assert.Equal(t, "", GenesisBlock.Header.BlockHash)
	assert.Equal(t, 0, GenesisBlock.Header.MerkleTreeHeight)
	assert.Equal(t, 0, GenesisBlock.Header.TransactionCount)
	assert.Equal(t, make([][]string, 0), GenesisBlock.MerkleTree)
	assert.Equal(t, make([]*tx.Transaction, 0), GenesisBlock.Transactions)
}

func TestConfigFromJson(t *testing.T) {
	_, err := ConfigFromJson("WrongFileName")
	assert.Error(t, err)
}

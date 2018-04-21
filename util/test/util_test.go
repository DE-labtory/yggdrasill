package test

import (
	"testing"

	"github.com/it-chain/yggdrasill/block"
	"github.com/it-chain/yggdrasill/util"
	"github.com/stretchr/testify/assert"
)

func TestDeserialize(t *testing.T) {

	//var block1 block.Block

	block1 := block.DefaultBlock{Header: block.BlockHeader{Height: 1, CreatorID: "jun"}}

	b, err := block1.Serialize()

	if err != nil {

	}

	block2 := &block.DefaultBlock{}

	err = util.Deserialize(b, block2)

	if err != nil {

	}

	assert.Equal(t, block2.GetHeight(), block1.GetHeight())
	assert.Equal(t, block1.Header.CreatorID, block1.Header.CreatorID)
}

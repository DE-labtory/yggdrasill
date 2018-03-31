package blockchaindb

import(
	"github.com/it-chain/yggdrasill/validator"
	"github.com/it-chain/yggdrasill/leveldb"
	"github.com/it-chain/yggdrasill/block"
	"github.com/it-chain/yggdrasill/transaction"
	"fmt"
)

type Yggdrasil struct {
	dbType    string
	db        YggdrasilInterface
	validator validator.Validator
}

func NewYggdrasil(dbType string, validator validator.Validator, opts map[string]interface{}) *Yggdrasil {
	var db YggdrasilInterface

	switch dbType {
	case "leveldb":
		db = leveldb.NewLevelDB(opts)
		break
	default :
		panic(fmt.Sprint("Unsupported DB Type"))
	}

	return &Yggdrasil{dbType: dbType, db: db, validator: validator}
}

func (y *Yggdrasil) Close() {
	y.db.Close()
}

func (y *Yggdrasil) AddBlock(block block.Block) error {
	return y.db.AddBlock(block)
}

func (y *Yggdrasil) GetBlockByNumber(block block.Block, blockNumber uint64) error {
	return y.db.GetBlockByNumber(block, blockNumber)
}

func (y *Yggdrasil) GetBlockByHash(block block.Block, hash string) error {
	return y.db.GetBlockByHash(block, hash)
}

func (y *Yggdrasil) GetLastBlock(block block.Block) error {
	return y.db.GetLastBlock(block)
}

func (y *Yggdrasil) GetTransactionByTxID(transaction transaction.Transaction, txid string) error {
	return y.db.GetTransactionByTxID(transaction, txid)
}

func (y *Yggdrasil) GetBlockByTxID(block block.Block, txid string) error {
	return y.db.GetBlockByTxID(block, txid)
}
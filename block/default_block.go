package block

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"sort"
	"strings"
	"time"

	"encoding/json"
	"io/ioutil"

	tx "github.com/it-chain/yggdrasill/transaction"
	"github.com/it-chain/yggdrasill/util"
)

type DefaultBlock struct {
	Header       *BlockHeader
	MerkleTree   [][]string
	Transactions []*tx.Transaction
}

type BlockHeader struct {
	Height             uint64    `json:"Height"`
	PreviousHash       string    `json:"PreviousHash"`
	Version            string    `json:"Version"`
	MerkleTreeRootHash string    `json:"MerkleTreeRootHash"`
	TimeStamp          time.Time `json:"TimeStamp"`
	CreatorID          string    `json:"CreatorID"`
	Signature          []byte    `json:"Signature"`
	BlockHash          string    `json:"BlockHash"`
	MerkleTreeHeight   int       `json:"MerkleTreeHeight"`
	TransactionCount   int       `json:"TransactionCount"`
}

func (block *DefaultBlock) PutTransaction(transaction *tx.Transaction) {

	block.Transactions = append(block.Transactions, transaction)
	block.Header.TransactionCount++
}

func (block *DefaultBlock) FindTransactionIndexByHash(txHash string) {

}

func (block *DefaultBlock) Serialize() ([]byte, error) {
	return util.Serialize(block)
}

func (block *DefaultBlock) GenerateHash() error {

	if block.Header.MerkleTreeRootHash == "" {
		return errors.New("no merkle tree root hash")
	}

	str := []string{block.Header.MerkleTreeRootHash, block.Header.PreviousHash}
	block.Header.BlockHash = computeSHA256(str)

	return nil
}

func (block *DefaultBlock) GetHash() string {
	return block.Header.BlockHash
}

func (block *DefaultBlock) GetTransactions() []*tx.Transaction {
	return block.Transactions
}

func (block *DefaultBlock) GetHeight() uint64 {
	return block.Header.Height
}

func (block *DefaultBlock) IsPrev(serializedBlock []byte) bool {
	return true
}

func computeSHA256(data []string) string {

	sort.Strings(data)
	arg := strings.Join(data, ",")
	hash := sha256.New()
	hash.Write([]byte(arg))
	return hex.EncodeToString(hash.Sum(nil))
}

func CreateNewBlock(prevBlock *DefaultBlock, createPeerId string) (*DefaultBlock, error) {
	var header BlockHeader

	if createPeerId == "" {
		return &DefaultBlock{}, errors.New("You have to put createPeerId")
	}

	if prevBlock == nil {
		header.Height = 0
		header.PreviousHash = ""
		header.Version = ""

	} else {
		header.Height = prevBlock.Header.Height + 1
		header.PreviousHash = prevBlock.Header.BlockHash
		header.Version = prevBlock.Header.Version

	}
	header.CreatorID = createPeerId
	header.MerkleTreeHeight = 0
	header.TimeStamp = time.Now()
	header.TransactionCount = 0
	header.MerkleTreeRootHash = ""
	header.BlockHash = ""
	header.Signature = make([]uint8, 0)

	return &DefaultBlock{Header: &header, MerkleTree: make([][]string, 0), Transactions: make([]*tx.Transaction, 0)}, nil
}

func CreateGenesisBlock() (*DefaultBlock, error) {
	byteValue, err := ConfigFromJson("GenesisBlockConfig.json")
	if err != nil {
		return nil, err
	}

	var GenesisBlock *DefaultBlock
	json.Unmarshal(byteValue, &GenesisBlock)
	GenesisBlock.Header.TimeStamp = time.Now()
	return GenesisBlock, nil
}

func ConfigFromJson(filename string) ([]uint8, error) {
	folderpath := "../config/"
	jsonFile, err := os.Open(folderpath + filename)
	defer jsonFile.Close()
	if err != nil {
		return nil, err
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	return byteValue, nil
}

//
//func (s *Block) PutTranscation(tx *transaction.Transaction) error{
//
//	//todo 이부분은 아직 보류
//	//if transaction.Validate() == false{
//	//	return errors.New("invalid transaction")
//	//}
//	//if transaction.TransactionStatus == Status_TRANSACTION_UNKNOWN {
//	//	if true { // Docker에서 실행하고 return이 true면 Confirmed 나중에 수정할 것.
//	//		transaction.TransactionStatus = Status_TRANSACTION_CONFIRMED
//	//	} else {
//	//		transaction.TransactionStatus = Status_TRANSACTION_UNCONFIRMED
//	//	}
//	//}
//
//	//todo 다른경우가 있을 수 있기 때문에 무시하도록 해야함
//	for _, confirmedTx := range s.Transactions{
//		if confirmedTx.TransactionID == tx.TransactionID{
//			return nil
//		}
//	}
//
//	s.Transactions = append(s.Transactions, tx)
//	s.TransactionCount++
//
//	return nil
//}
//
//func (s Block) FindTransactionIndex(hash string) (idx int, err error){
//	for idx = 0; idx < s.TransactionCount; idx++{
//		if hash == s.Transactions[idx].TransactionHash{
//			return idx, nil
//		}
//	}
//	return -1, errors.New("txHash is not here")
//}
//
//func (s *Block) MakeMerkleTree(){
//
//	if s.TransactionCount == 0 {
//		s.Header.MerkleTreeRootHash = ""
//		return
//	}
//
//	var mtList []string
//
//	for _, h := range s.Transactions{
//		mtList = append(mtList, h.TransactionHash)
//	}
//	for {
//		treeLength := len(mtList)
//		s.MerkleTreeHeight++
//		if treeLength <= 1 {
//			s.MerkleTree = append(s.MerkleTree, mtList)
//			break
//		} else if treeLength % 2 == 1 {
//			mtList = append(mtList, mtList[treeLength - 1])
//			treeLength++
//		}
//		s.MerkleTree = append(s.MerkleTree, mtList)
//		var tmpMtList []string
//		for x := 0; x < treeLength/2; x++{
//			idx := x * 2
//			hashArg := []string{mtList[idx], mtList[idx+1]}
//			mkHash := common.ComputeSHA256(hashArg)
//			tmpMtList = append(tmpMtList, mkHash)
//		}
//		mtList = tmpMtList
//	}
//	if len(mtList) == 1 {
//		s.Header.MerkleTreeRootHash = mtList[0]
//	}
//}
//
//func (s Block) MakeMerklePath(idx int) (path []string){
//	for i := 0; i < s.MerkleTreeHeight-1; i++{
//		path = append(path, s.MerkleTree[i][(idx >> uint(i)) ^ 1])
//	}
//	return path
//}
//
//func (s *Block) GenerateBlockHash() error{
//
//	if s.Header.MerkleTreeRootHash == "" {
//		return errors.New("no merkle tree root hash")
//	}
//
//	str := []string{s.Header.MerkleTreeRootHash, s.Header.TimeStamp.String(), s.Header.PreviousHash}
//	s.Header.BlockHash = common.ComputeSHA256(str)
//	return nil
//}
//
//func (s Block) BlockSerialize() ([]byte, error){
//	return common.Serialize(s)
//}
//
//func BlockDeserialize(by []byte) (Block, error) {
//	block := Block{}
//	err := common.Deserialize(by, &block)
//	return block, err
//}
//
//// 해당 트랜잭션이 정당한지 머클패스로 검사함
//func (s Block) VerifyTx(tx transaction.Transaction) (bool, error) {
//
//	if s.Header.BlockHeight == 0 && s.TransactionCount == 0 {
//		return true, nil;
//	}
//
//	hash := tx.TransactionHash
//	idx, err := s.FindTransactionIndex(hash)
//
//	if err != nil {
//		return false, err
//	}
//
//	merklePath := s.MakeMerklePath(idx)
//
//	for _, sibling_hash := range merklePath{
//		str := []string{hash, sibling_hash}
//		hash = common.ComputeSHA256(str)
//	}
//
//	if hash == s.Header.MerkleTreeRootHash{
//		return true, nil
//	} else {
//		return false, errors.New("transaction is invalid")
//	}
//}
//
//// 블럭내의 모든 트랜잭션들이 정당한지 머클패스로 검사함
//func (s Block) VerifyBlock() (bool, error) {
//	for idx := 0; idx < s.TransactionCount; idx++{
//		txVarification, txErr := s.VerifyTx(*s.Transactions[idx])
//		if txVarification == false  {
//			err := errors.New("block is invalid --- " + strconv.Itoa(idx) + "'s " + txErr.Error())
//			return false, err
//		}
//	}
//	return true, nil
//}

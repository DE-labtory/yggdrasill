package block

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sort"
	"strings"
	"time"

	tx "github.com/it-chain/yggdrasill/transaction"
	"github.com/it-chain/yggdrasill/util"
)

type DefaultBlock struct {
	Header       BlockHeader
	MerkleTree   [][]string
	Transactions []tx.Transaction
}

type BlockHeader struct {
	Height             uint64
	PreviousHash       string
	Version            string
	MerkleTreeRootHash string
	TimeStamp          time.Time
	CreatorID          string
	Signature          []byte
	BlockHash          string
	MerkleTreeHeight   int
	TransactionCount   int
}

func (block DefaultBlock) PutTransaction(transaction tx.Transaction) {

	block.Transactions = append(block.Transactions, transaction)
	block.Header.TransactionCount++
}

func (block DefaultBlock) FindTransactionIndexByHash(txHash string) {

}

func (block DefaultBlock) Serialize() ([]byte, error) {
	return util.Serialize(block)
}

func (block DefaultBlock) GenerateHash() error {

	if block.Header.MerkleTreeRootHash == "" {
		return errors.New("no merkle tree root hash")
	}

	str := []string{block.Header.MerkleTreeRootHash, block.Header.TimeStamp.String(), block.Header.PreviousHash}
	block.Header.BlockHash = computeSHA256(str)

	return nil
}

func (block DefaultBlock) GetHash() string {
	return block.Header.BlockHash
}

func (block DefaultBlock) GetTransactions() []tx.Transaction {
	return block.Transactions
}

func (block DefaultBlock) GetHeight() uint64 {
	return block.Header.Height
}

func (block DefaultBlock) IsPrev(serializedBlock []byte) bool {
	return true
}

func computeSHA256(data []string) string {

	sort.Strings(data)
	arg := strings.Join(data, ",")
	hash := sha256.New()
	hash.Write([]byte(arg))
	return hex.EncodeToString(hash.Sum(nil))
}

//
//func CreateNewBlock(prevBlock *Block, createPeerId string) *Block{
//
//	var header BlockHeader
//	if prevBlock == nil{
//		header.Number = 0
//		header.PreviousHash = ""
//		header.Version = ""
//		header.BlockHeight = 0
//	} else {
//		header.Number = prevBlock.Header.Number + 1
//		header.PreviousHash = prevBlock.Header.BlockHash
//		header.Version = prevBlock.Header.Version
//		header.BlockHeight = prevBlock.Header.BlockHeight + 1
//	}
//	header.CreatedPeerID = createPeerId
//	header.TimeStamp = time.Now().Round(0)
//	header.BlockStatus = Status_BLOCK_UNCONFIRMED
//
//	return &Block{Header:&header, MerkleTree:make([][]string, 0), MerkleTreeHeight:0, TransactionCount:0, Transactions:make([]*Transaction, 0)}
//}
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

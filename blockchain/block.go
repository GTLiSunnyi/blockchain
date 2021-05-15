package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/GTLiSunnyi/blockchain/tx"
	"github.com/GTLiSunnyi/blockchain/types"
	"github.com/GTLiSunnyi/blockchain/utils"
)

type Block struct {
	Header *Header
	TXs    []tx.TX // 交易数据
}

type Header struct {
	PreBlockHash [32]byte // 前区块哈希
	TimeStamp    uint64   // 时间戳，从1970.1.1至今的秒数
	Interval     uint64   // 一般出块的间隔时间
	MerkleRoot   [32]byte // 梅克尔树
	Address      string   // 打包区块的人
	Hash         [32]byte // 当前区块哈希
	Height       int      // 区块高度
}

// 创建区块
func (bc *BC) CreateBlock(addr string, prikey *ecdsa.PrivateKey, txs []tx.TX, preBlockHash [32]byte, block *Block, c chan bool) {
	var header = &Header{
		PreBlockHash: preBlockHash,
		TimeStamp:    uint64(time.Now().Unix()),
		Interval:     5,
		Address:      addr,
	}

	block.Header = header
	block.TXs = txs

	block.setTxSignatureAndMerkle(prikey, bc)
	block.Header.Hash = block.GetBlockHash()
	block.Header.Height = types.Height

	types.Height++
	fmt.Printf("打包完成，区块高度：%d，区块哈希：%s。\n", block.Header.Height, block.Header.Hash)

	if c == nil {
		return
	} else {
		c <- true
	}
}

// 设置梅克尔根，取哈希后签名
func (block *Block) setTxSignatureAndMerkle(prikey *ecdsa.PrivateKey, bc *BC) {
	// 遍历交易
	for i, tx := range block.TXs {
		signature := block.Sign(&tx, prikey)
		block.TXs[i].Signature = signature
	}

	bc.TxPool = nil

	block.setMerkle()
}

// 对交易进行签名
func (block *Block) Sign(tx *tx.TX, prikey *ecdsa.PrivateKey) [][]byte {
	hashText := sha1.Sum(append(tx.DenomTX, tx.FileTX...))

	//数字签名
	r, s, _ := ecdsa.Sign(rand.Reader, prikey, hashText[:])

	rText, _ := r.MarshalText()
	sText, _ := s.MarshalText()

	return [][]byte{rText, sText}
}

// 校验交易签名是否正确
func (block *Block) IsValid(pubkey *ecdsa.PublicKey) bool {
	// 遍历交易
	for _, tx := range block.TXs {
		var r, s big.Int
		r.UnmarshalText(tx.Signature[0])
		s.UnmarshalText(tx.Signature[1])

		tx.Signature = nil
		hashText := sha1.Sum(append(tx.DenomTX, tx.FileTX...))

		//认证
		res := ecdsa.Verify(pubkey, hashText[:], &r, &s)
		if !res {
			return false
		}
	}

	return true
}

// 生成区块哈希
func (block *Block) GetBlockHash() [32]byte {
	tmp := [][]byte{
		block.Header.PreBlockHash[:],
		utils.UintToByte(block.Header.TimeStamp),
		utils.UintToByte(block.Header.Interval),
		block.Header.MerkleRoot[:],
		[]byte(block.Header.Address),
	}

	data := bytes.Join(tmp, []byte{})

	return sha256.Sum256(data)

}

// 设置梅克尔根
func (block *Block) setMerkle() {
	var bytes []byte
	// 遍历交易
	for _, tx := range block.TXs {
		txByte, _ := json.Marshal(tx)
		bytes = append(bytes, txByte...)
	}

	block.Header.MerkleRoot = sha256.Sum256(bytes)
}

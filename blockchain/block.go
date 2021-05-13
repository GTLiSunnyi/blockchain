package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha1"
	"math/big"
	"time"

	"mybc/tx"
	"mybc/utils"
)

const Interval uint64 = 5

type Block struct {
	Header
	TXs *[]tx.TX // 交易数据
}

type Header struct {
	PreBlockHash []byte   // 前区块哈希
	MerkleRoot   []byte   // Merkle根，交易的哈希值就是Merkle根
	TimeStamp    uint64   // 时间戳，从1970.1.1至今的秒数
	Interval     uint64   // 一般出块的间隔时间
	Address      string   // 打包区块的人
	Signature    [][]byte // 打包者的签名
	Hash         []byte   // 当前区块哈希
}

// 创建区块
func CreateBlock(addr string, prikey *ecdsa.PrivateKey, txs []tx.TX, preBlockHash []byte) *Block {
	var Header = Header{
		PreBlockHash: preBlockHash,
		TimeStamp:    uint64(time.Now().Unix()),
		Interval:     Interval,
		Address:      addr,
	}

	block := Block{
		Header: Header,
		TXs:    &txs,
	}

	block.Signature = block.setMerkleAndGetSignature(prikey)
	block.Hash = block.GetBlockHash()
	return &block
}

// 设置梅克尔根，取哈希后签名
func (block *Block) setMerkleAndGetSignature(prikey *ecdsa.PrivateKey) [][]byte {
	var hash []byte
	// 遍历交易
	for _, tx := range *block.TXs {
		hash = append(hash, tx.Signature...)
	}
	block.MerkleRoot = hash

	hashText := sha1.Sum([]byte(hash))
	//数字签名
	r, s, _ := ecdsa.Sign(rand.Reader, prikey, hashText[:])

	rText, _ := r.MarshalText()
	sText, _ := s.MarshalText()
	return [][]byte{rText, sText}
}

// 生成区块哈希
func (block *Block) GetBlockHash() []byte {
	tmp := [][]byte{
		block.Header.PreBlockHash,
		block.Header.MerkleRoot,
		[]byte(block.Header.Address),
		block.Signature[0],
		block.Signature[1],
		utils.UintToByte(block.Header.TimeStamp),
		utils.UintToByte(block.Header.Interval),
	}

	data := bytes.Join(tmp, []byte{})
	return data
}

// 校验挖矿是否有效
func (block *Block) IsValid(pubkey *ecdsa.PublicKey) bool {
	var hash []byte
	// 遍历交易
	for _, tx := range *block.TXs {
		hash = append(hash, tx.Id...)
	}

	if string(hash) != string(block.MerkleRoot) {
		return false
	}

	hashText := sha1.Sum(hash)

	var r, s big.Int
	r.UnmarshalText(block.Signature[0])
	s.UnmarshalText(block.Signature[1])

	//认证
	res := ecdsa.Verify(pubkey, hashText[:], &r, &s)
	return res
}

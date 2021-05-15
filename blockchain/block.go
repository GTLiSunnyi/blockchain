package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha1"
	"math/big"
	"time"

	"mybc/tx"
	"mybc/types"
	"mybc/utils"
)

type Block struct {
	Header
	TXs []tx.TX // 交易数据
}

type Header struct {
	PreBlockHash []byte // 前区块哈希
	TimeStamp    uint64 // 时间戳，从1970.1.1至今的秒数
	Interval     uint64 // 一般出块的间隔时间
	Address      string // 打包区块的人
	Hash         []byte // 当前区块哈希
}

// 创建区块
func CreateBlock(addr string, prikey *ecdsa.PrivateKey, txs []tx.TX, preBlockHash []byte, block *Block, c chan bool) {
	var header = Header{
		PreBlockHash: preBlockHash,
		TimeStamp:    uint64(time.Now().Unix()),
		Interval:     types.Interval,
		Address:      addr,
	}

	block.Header = header
	block.TXs = txs

	block.setMerkleAndTxSignature(prikey)
	block.Hash = block.GetBlockHash()

	if c == nil {
		return
	} else {
		c <- true
	}
}

// 设置梅克尔根，取哈希后签名
func (block *Block) setMerkleAndTxSignature(prikey *ecdsa.PrivateKey) {
	// 遍历交易
	for i, tx := range block.TXs {
		hashText := sha1.Sum(append(tx.DenomTX, tx.FileTX...))
		//数字签名
		r, s, _ := ecdsa.Sign(rand.Reader, prikey, hashText[:])

		rText, _ := r.MarshalText()
		sText, _ := s.MarshalText()

		block.TXs[i].Signature = [][]byte{rText, sText}
	}
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
func (block *Block) GetBlockHash() []byte {
	tmp := [][]byte{
		block.Header.PreBlockHash,
		[]byte(block.Header.Address),
		utils.UintToByte(block.Header.TimeStamp),
		utils.UintToByte(block.Header.Interval),
	}

	data := bytes.Join(tmp, []byte{})
	return data
}

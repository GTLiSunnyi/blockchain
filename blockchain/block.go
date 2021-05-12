package blockchain

import (
	"time"

	"mybc/consensus"
	"mybc/tx"
)

type Block struct {
	PreBlockHash []byte  // 前区块哈希
	MerkleRoot   []byte  // Merkle根，交易的哈希值就是Merkle根
	TimeStamp    uint64  // 时间戳，从1970.1.1至今的秒数
	Difficulty   uint64  // 挖矿的难度值，v2时使用
	Nonce        uint64  // 随机数，挖矿找的就是它
	TXs          []tx.TX // 交易数据
	Hash         []byte  // 当前区块哈希，区块中本不存在的字段，为了方便操作添加进来了
}

// 创建区块
func NewBlock(txs []tx.TX, preBlockHash []byte) (*Block, bool) {
	block := Block{
		PreBlockHash: preBlockHash,
		TimeStamp:    uint64(time.Now().Unix()),
		Difficulty:   16,
		TXs:          txs,
	}
	block.setMerkle()
	pow := consensus.NewPow(&block)
	block.Hash, block.Nonce = pow.Run()
	return &block, pow.IsValid()
}

// 设置梅克尔根
func (block *Block) setMerkle() {
	var hash []byte
	// 遍历交易
	for _, tx := range block.TXs {
		hash = append(hash, tx.Id...)
	}
	block.MerkleRoot = hash
}

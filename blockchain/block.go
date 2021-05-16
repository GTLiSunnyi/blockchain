package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/GTLiSunnyi/blockchain/tx"
	"github.com/GTLiSunnyi/blockchain/types"
	"github.com/GTLiSunnyi/blockchain/utils"
	"github.com/GTLiSunnyi/blockchain/wallet"
)

type Block struct {
	Header *Header
	TXs    []tx.TX // 交易数据
}

type Header struct {
	PreBlockHash [32]byte      // 前区块哈希
	TimeStamp    uint64        // 时间戳，从1970.1.1至今的秒数
	Interval     time.Duration // 一般出块的间隔时间
	MerkleRoot   [32]byte      // 梅克尔树
	Address      string        // 打包区块的人
	Hash         [32]byte      // 当前区块哈希
	Height       int           // 区块高度
}

// 创建区块
func (bc *BC) CreateBlock(address string, ws *wallet.Wallets, preBlockHash [32]byte, block *Block, c chan bool, packagers []string) {
	var header = &Header{
		PreBlockHash: preBlockHash,
		TimeStamp:    uint64(time.Now().Unix()),
		Interval:     types.Interval,
		Address:      address,
	}

	block.Header = header

	block.setTxSignatureAndMerkle(ws.Gather[address].PriKey, bc)
	block.Header.Hash = block.GetBlockHash()
	block.Header.Height = types.Height

	types.Height++

	isOK := block.Verify(packagers, ws)

	if isOK {
		fmt.Printf("打包完成，区块高度：%d，区块哈希：%X。\n", block.Header.Height, block.Header.Hash)
		fmt.Println("打包者地址：", address)
		if c == nil {
			return
		}
		c <- true
	} else {
		fmt.Println("区块验证失败，取消上链！")
		c <- false
	}
}

func (block *Block) Verify(packagers []string, ws *wallet.Wallets) bool {
	for _, tx := range block.TXs {
		var sum int
		for _, packager := range packagers {
			isTrue := tx.IsValid(ws.Gather[packager].PubKey)
			if isTrue {
				sum++
			}
		}
		if sum < len(packagers)*2/3 {
			return false
		}
	}

	return true
}

// 设置梅克尔根，取哈希后签名
func (block *Block) setTxSignatureAndMerkle(prikey *ecdsa.PrivateKey, bc *BC) {
	// 遍历交易
	for _, tx := range block.TXs {
		tx.Sign(prikey)
	}

	bc.TxPool = nil

	block.setMerkle()
}

// 生成区块哈希
func (block *Block) GetBlockHash() [32]byte {
	tmp := [][]byte{
		block.Header.PreBlockHash[:],
		utils.UintToByte(block.Header.TimeStamp),
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

package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/GTLiSunnyi/blockchain/account"
	"github.com/GTLiSunnyi/blockchain/tx"
	"github.com/GTLiSunnyi/blockchain/types"
	"github.com/GTLiSunnyi/blockchain/utils"
)

type Block struct {
	Header *Header
	Txs    []tx.Tx // 交易数据
}

type Header struct {
	PreBlockHash [32]byte      // 前区块哈希
	TimeStamp    uint64        // 时间戳，从1970.1.1至今的秒数
	Interval     time.Duration // 一般出块的间隔时间
	MerkleRoot   [32]byte      // 梅克尔树
	Address      string        // 打包区块的人
	Hash         [32]byte      // 当前区块哈希
	Height       int           // 区块高度
	Version      string
}

// 创建区块
func (bc *BC) CreateBlock(address string, accounts *account.Accounts, preBlockHash [32]byte, block *Block, c chan bool, packagers []string) {
	var header = &Header{
		PreBlockHash: preBlockHash,
		TimeStamp:    uint64(time.Now().Unix()),
		Interval:     types.Interval,
		Address:      address,
		Version:      types.Version,
	}

	block.Header = header

	block.setTxSignatureAndMerkle(accounts.Gather[address].PriKey, bc)
	block.Header.Height = types.Height
	block.Header.Hash = block.GetBlockHash()

	types.Height++

	isOK := block.Verify(packagers, accounts)

	if isOK {
		bc.LastBlockHash = block.Header.Hash

		fmt.Printf("打包完成，区块高度：%d，区块哈希：%X。\n", block.Header.Height, block.Header.Hash)
		fmt.Println("打包者地址：", address)
		fmt.Printf("交易信息：%+v\n", block.Txs)

		if c == nil {
			return
		}
		c <- true
	} else {
		fmt.Println("区块验证失败，取消上链！")
		c <- false
	}
}

func (block *Block) Verify(packagers []string, accounts *account.Accounts) bool {
	for _, tx := range block.Txs {
		var sum int
		for _, packager := range packagers {
			isTrue := tx.IsValid(accounts.Gather[packager].PubKey)
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
	for _, tx := range block.Txs {
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
		[]byte(block.Header.Version),
	}

	data := bytes.Join(tmp, []byte{})

	return sha256.Sum256(data)

}

// 设置梅克尔根
func (block *Block) setMerkle() {
	var bytes []byte
	// 遍历交易
	for _, tx := range block.Txs {
		txByte, _ := json.Marshal(tx)
		bytes = append(bytes, txByte...)
	}

	block.Header.MerkleRoot = sha256.Sum256(bytes)
}

package consensus

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

type Pow struct {
	Block      *Block
	Difficulty *big.Int
}

// 创建工作量证明
func NewPow(block *Block) *Pow {
	pow := Pow{
		Block: block,
	}
	// 创建大数
	// 0000000000000000000000000000000000000000000000000000000000000001
	bigIntTmp := big.NewInt(1)
	// 将1向坐移256 - Bits位
	// 0001000000000000000000000000000000000000000000000000000000000000(需要移4位)
	bigIntTmp.Lsh(bigIntTmp, 256-uint(block.Difficulty))
	pow.Difficulty = bigIntTmp
	return &pow
}

// 挖矿
func (pow *Pow) Run() ([]byte, uint64) {
	// 1.获取block数据
	// 2.获取哈希
	// 3.与难度值比较
	// a.哈希值大于难度值，nonce++
	// b.哈希值小于难度值，挖矿成功，退出
	var nonce uint64
	var hash [32]byte
	for {
		hash = sha256.Sum256(pow.blockHash(nonce))
		//将hash（数组类型）转成big.int, 然后与pow.target进行比较, 需要引入局部变量
		var bigIntTmp big.Int
		bigIntTmp.SetBytes(hash[:])
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		//   x              y
		if bigIntTmp.Cmp(pow.Difficulty) == -1 {
			//此时x < y ， 挖矿成功！
			fmt.Printf("挖矿成功！nonce: %d, 哈希值为: %x\n", nonce, hash)
			break
		} else {
			nonce++
		}
	}
	return hash[:], nonce
}

// 校验挖矿是否有效
func (pow *Pow) IsValid() bool {
	//在校验的时候，block的数据是完整的，我们要做的是校验一下，Hash，block数据，和Nonce是否满足难度值要求
	data := pow.blockHash(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	var tmp big.Int
	tmp.SetBytes(hash[:])
	return tmp.Cmp(pow.Difficulty) == -1
}

// 生成区块哈希
func (pow *Pow) blockHash(Nonce uint64) []byte {
	block := pow.Block
	tmp := [][]byte{
		block.PreBlockHash,
		block.MerkleRoot,
		UintToByte(block.TimeStamp),
		UintToByte(block.Difficulty),
		UintToByte(Nonce),
	}
	data := bytes.Join(tmp, []byte{})
	return data
}

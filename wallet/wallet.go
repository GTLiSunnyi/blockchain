package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"

	"github.com/GTLiSunnyi/blockchain/types"
)

type Wallet struct {
	types.AccountType
	PubKey *ecdsa.PublicKey
	PriKey *ecdsa.PrivateKey
}

// 获取钱包地址
func (wallet *Wallet) GetAddress() string {
	// 将公钥进行一系列处理后得到地址
	pubKeyHash /*20字节*/ := HashPubKey(wallet.PubKey)
	version /*1字节*/ := 0x00
	// 21字节数据
	payload := append([]byte{byte(version)}, pubKeyHash...)
	// 创建checksum(4字节校验码)
	checksum := Checksum(payload)
	// 25字节数据
	payload = append(payload, checksum...)
	// base58处理
	address := base58.Encode(payload)
	fmt.Println("钱包地址为：", address)

	return address
}

// 检验地址是否正确
func IsValidAddress(address string) bool {
	// 1.将地址解码
	decodeInfo := base58.Decode(address)
	// 2.如果没有25个字节，则给的地址错误
	if len(decodeInfo) != 25 {
		fmt.Println("1")
		return false
	}
	// 3.取出前21个字节，运行CheckSum函数得到checksum1
	payload := decodeInfo[:len(decodeInfo)-4]
	checksum1 := Checksum(payload)
	// 4.取出后4个字节得到checksum2
	checksum2 := decodeInfo[len(decodeInfo)-4:]
	// 5.比较checksum1和checksum2，相同则地址有效
	return bytes.Equal(checksum1, checksum2)
}

// 将公钥处理为公钥哈希
func HashPubKey(pubKey *ecdsa.PublicKey) []byte {
	hash := sha256.Sum256(append(pubKey.X.Bytes(), pubKey.Y.Bytes()...))
	// 生成publicKeyHash
	// 1.创建hash160对象
	rip160Hash := ripemd160.New()
	// 2.向hash160中write数据
	_, err := rip160Hash.Write(hash[:])
	if err != nil {
		log.Panic(err)
	}
	// 3.做哈希运算
	return rip160Hash.Sum(nil)
}

// 计算checksum(4字节校验码)
func Checksum(payload []byte) []byte {
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])
	return second[:4]
}

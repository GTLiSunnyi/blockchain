package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"mybc/types"
	"os"

	"github.com/boltdb/bolt"
)

var Ws *Wallets

// 所有钱包集合
// Wallets对外，WalletKeyPair对内，Wallets调用WalletKeyPair
type Wallets struct {
	Gather map[string]*Wallet // 地址=>钱包
}

// 钱包集合文件名字
const FileName = "./wallets.dat"

func (ws *Wallets) NewWallets() {
	Ws.Gather = make(map[string]*Wallet)
	Ws.LoadFile()
}

func (ws *Wallets) NewWallet() *Wallet {
	// 创建钱包，并保存到钱包集合中，最后保存到本地文件
	// 创建私钥
	priKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println("创建私钥失败，创建钱包失败")
		log.Panic(err)
	}
	// 由私钥创建公钥
	pubKey := priKey.PublicKey

	wallet := &Wallet{types.NodeTypes, &pubKey, priKey}

	address := wallet.GetAddress()
	Ws.Gather[address] = wallet
	Ws.SaveFile()

	return wallet
}

// 打印钱包集合
func (ws *Wallets) GetList() {
	for address, _ := range ws.Gather {
		fmt.Println(address)
	}
}

// 读取本地文件
func (ws *Wallets) LoadFile() {
	// 读取本地文件，返回解码后的信息
	_, err := os.Stat(FileName)
	if err == nil {
		data, err := ioutil.ReadFile(FileName)
		if err != nil {
			fmt.Println("读取本地文件失败")
			return
		}
		// 涉及到“序列化/反序列化类型是interface或者struct中某些字段是interface”
		// 所以解码方式比较特别
		gob.Register(elliptic.P256())
		decoder := gob.NewDecoder(bytes.NewReader(data))
		err = decoder.Decode(&ws)
		if err != nil {
			fmt.Println("解码读取文件数据失败")
		}
	}
}

// 保存到本地文件
func (ws *Wallets) SaveFile() {
	var buffer bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(ws)
	if err != nil {
		fmt.Println("序列化失败")
	}
	err = ioutil.WriteFile(FileName, buffer.Bytes(), 0600)
	if err != nil {
		fmt.Println("保存到本地失败")
	}
}

// 创建普通节点
func CreateNodeAccount() string {
	key := Ws.NewWallet()

	db, err := bolt.Open(types.DBName, 0600, nil)
	if err != nil {
		panic(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(types.NodeAccountBucketName))
		if b == nil {
			// 桶不存在则创建
			b, err = tx.CreateBucket([]byte(types.NodeAccountBucketName))
			if err != nil {
				panic(err)
			}
		}

		b.Put([]byte("NodeAccount"), []byte(key.GetAddress()))
		return nil
	})

	return key.GetAddress()
}

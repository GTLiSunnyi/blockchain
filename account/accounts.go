package account

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/boltdb/bolt"

	"github.com/GTLiSunnyi/blockchain/types"
)

// 所有钱包集合
// Wallets对外，WalletKeyPair对内，Wallets调用WalletKeyPair
type Accounts struct {
	Gather map[string]*Account // 地址=>钱包
	DB     *bolt.DB
}

func NewAccounts(db *bolt.DB) *Accounts {
	gather := make(map[string]*Account)
	accounts := &Accounts{Gather: gather}

	superWallet, superAddress := accounts.NewAccount()
	superWallet.AccountType = types.SuperTypes
	types.CurrentUsers = superAddress

	// 账户数据库
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(types.AccountBucketName))
		if b == nil {
			// 桶不存在则创建
			var err error
			b, err = tx.CreateBucket([]byte(types.AccountBucketName))
			if err != nil {
				panic(err)
			}
		}

		b.Put([]byte(types.CurrentUsers), []byte(types.SuperTypes))
		return nil
	})

	accounts.DB = db
	accounts.LoadFile()

	return accounts
}

func (accounts *Accounts) NewAccount() (*Account, string) {
	// 创建钱包，并保存到钱包集合中，最后保存到本地文件
	// 创建私钥
	priKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println("创建私钥失败，创建钱包失败")
		log.Panic(err)
	}

	// 由私钥创建公钥
	pubKey := priKey.PublicKey
	fmt.Println("公钥为：", pubKey)

	// pubKeyByte := append(pubKey.X.Bytes(), pubKey.Y.Bytes()...)

	account := &Account{types.NodeTypes, "", &pubKey, priKey}
	address := account.GetAddress()

	accounts.Gather[address] = account
	accounts.SaveFile()

	return account, address
}

// 查询所有的地址及用户权限
func (accounts *Accounts) QueryAccount() {
	for k, v := range accounts.Gather {
		fmt.Println(k, v.AccountType)
	}
}

// 打印钱包集合
func (accounts *Accounts) GetList() {
	for address, _ := range accounts.Gather {
		fmt.Println(address)
	}
}

// 读取本地文件
func (accounts *Accounts) LoadFile() {
	// 读取本地文件，返回解码后的信息
	_, err := os.Stat(types.FileName)
	if err == nil {
		data, err := ioutil.ReadFile(types.FileName)
		if err != nil {
			fmt.Println("读取本地文件失败")
			return
		}
		// 涉及到“序列化/反序列化类型是interface或者struct中某些字段是interface”
		// 所以解码方式比较特别
		gob.Register(elliptic.P256())
		decoder := gob.NewDecoder(bytes.NewReader(data))
		err = decoder.Decode(&accounts)
		if err != nil {
			fmt.Println("解码读取文件数据失败")
		}
	}
}

// 保存到本地文件
func (accounts *Accounts) SaveFile() {
	var buffer bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(accounts)
	if err != nil {
		fmt.Println("序列化失败")
	}
	err = ioutil.WriteFile(types.FileName, buffer.Bytes(), 0600)
	if err != nil {
		fmt.Println("保存到本地失败")
	}
}

// 创建普通节点
func (accounts *Accounts) CreateNodeAccount() string {
	_, address := accounts.NewAccount()

	accounts.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(types.AccountBucketName))
		b.Put([]byte(address), []byte(types.NodeTypes))
		return nil
	})
	fmt.Println("新的普通节点的地址为：", address)
	fmt.Println("节点成功加入网络！")

	return address
}

func (accounts *Accounts) AddPms(address string) {
	fmt.Println(address)
	if accounts.IsInAccounts(address) {
		accounts.Gather[address].AccountType = types.AdminTypes
		fmt.Println("成功将节点权限提升为管理员！")
	} else {
		fmt.Errorf("账户不存在！")
	}
}

func (accounts *Accounts) RmPms(address string) {
	if accounts.IsInAccounts(address) {
		accounts.Gather[address].AccountType = types.NodeTypes
		fmt.Println("成功将节点权降低为普通节点！")
	} else {
		fmt.Errorf("账户不存在！")
	}
}

func (accounts *Accounts) IsInAccounts(addrress string) bool {
	fmt.Printf("%+v", accounts.Gather[addrress])
	return accounts.Gather[addrress] != nil
}

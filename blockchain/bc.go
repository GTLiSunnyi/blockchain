package blockchain

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/boltdb/bolt"

	"github.com/GTLiSunnyi/blockchain/account"
	"github.com/GTLiSunnyi/blockchain/tx"
	"github.com/GTLiSunnyi/blockchain/types"
)

type BC struct {
	DB     *bolt.DB
	TxPool []tx.Tx
	Iterator
	LastBlockHash [32]byte
}

// 创建区块链
func NewBC() (*BC, *bolt.DB) {
	db, err := bolt.Open(types.DBName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(types.BlockChainBucketName))
		if b == nil {
			// 桶不存在则创建
			_, err = tx.CreateBucket([]byte(types.BlockChainBucketName))
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
	bc := &BC{DB: db, LastBlockHash: [32]byte{}, TxPool: []tx.Tx{}}

	return bc, db
}

func (bc *BC) RunBC(accounts *account.Accounts, c chan string) {
	it := bc.NewIterator(accounts)

	var block Block
	bc.CreateBlock(types.CurrentUsers, accounts, [32]byte{}, &block, nil, nil)
	bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(types.BlockChainBucketName))
		blockinfo, _ := json.Marshal(block)
		b.Put([]byte{byte(block.Header.Height)}, bc.LastBlockHash[:])
		b.Put(bc.LastBlockHash[:], blockinfo)

		fmt.Printf("打包完成，区块高度：%d\n", block.Header.Height)
		fmt.Printf("区块哈希：%X\n", block.Header.Hash)
		fmt.Println("打包者地址：", block.Header.Address)
		fmt.Printf("交易信息：%+v\n", block.Txs)

		return nil
	})

	// 定时执行共识任务
	go func() {
		for range types.Ticker.C {
			it.Run(bc, accounts)
		}
	}()
}

// 迭代器
type Iterator struct {
	Packagers          []string
	CurrentPackagerNum int
	DB                 *bolt.DB
	Chan               chan bool
}

func (bc *BC) NewIterator(accounts *account.Accounts) *Iterator {
	it := Iterator{
		CurrentPackagerNum: 0,
		DB:                 bc.DB,
		Chan:               make(chan bool),
	}

	it.UpdatePackagers(accounts)
	return &it
}

// 运行迭代器
func (it *Iterator) Run(bc *BC, accounts *account.Accounts) {
	fmt.Println("当前打包区块队列：", it.Packagers)
	currentPackager := it.Packagers[it.CurrentPackagerNum]
	it.CurrentPackagerNum++

	var block Block
	block.Txs = bc.TxPool
	go bc.CreateBlock(currentPackager, accounts, bc.LastBlockHash, &block, it.Chan, it.Packagers)

	go func() {
		select {
		case isOk, ok := <-it.Chan:
			if ok {
				if isOk {
					bc.LastBlockHash = block.Header.Hash
					bc.DB.Update(func(tx *bolt.Tx) error {
						b := tx.Bucket([]byte(types.BlockChainBucketName))
						blockinfo, _ := json.Marshal(block)
						b.Put([]byte{byte(block.Header.Height)}, bc.LastBlockHash[:])
						b.Put(bc.LastBlockHash[:], blockinfo)

						fmt.Printf("打包完成，区块高度：%d\n", block.Header.Height)
						fmt.Printf("区块哈希：%X\n", block.Header.Hash)
						fmt.Println("打包者地址：", block.Header.Address)
						fmt.Printf("交易信息：%+v\n\n\n", block.Txs)

						return nil
					})
				} else {
					// 取消这次打包，并撤销管理员职位
					it.DB.Update(func(tx *bolt.Tx) error {
						b := tx.Bucket([]byte(types.AccountBucketName))
						b.Put([]byte(currentPackager), []byte(types.NodeTypes))
						return nil
					})

					fmt.Println("区块验证失败, 你已经失去管理员资格！")
				}
			}
		case <-time.After(types.Interval):
			// 多少秒内没处理完毕
			// 取消这次打包，并撤销管理员职位
			bc.DB.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte(types.AccountBucketName))
				b.Put([]byte(currentPackager), []byte(types.NodeTypes))
				return nil
			})

			fmt.Println("未在指定时间内打包完成，你已经失去管理员资格！")
		}

		if len(it.Packagers) == it.CurrentPackagerNum {
			it.UpdatePackagers(accounts)
			it.CurrentPackagerNum = 0
		}
	}()
}

func (it *Iterator) UpdatePackagers(accounts *account.Accounts) {
	it.Packagers = nil
	for _, v := range accounts.Gather {
		if v.AccountType == types.AdminTypes || v.AccountType == types.SuperTypes {
			it.Packagers = append(it.Packagers, v.Address)
		}
	}
}

func (bc *BC) QueryBlock(order string) {
	height, err := strconv.Atoi(order)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("查询高度：%d", height)

	bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(types.BlockChainBucketName))

		blockHash := b.Get([]byte{byte(height)})

		if blockHash == nil {
			fmt.Println("该区块高度不存在！")
			return nil
		}

		blockInfo := b.Get(blockHash)
		fmt.Println("以下是区块信息：")
		fmt.Println(string(blockInfo))

		return nil
	})
}

func (bc *BC) SendTx(tx *tx.Tx) {
	bc.TxPool = append(bc.TxPool, *tx)
	fmt.Println("交易发送成功！")
}

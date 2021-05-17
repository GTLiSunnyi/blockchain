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
	it := bc.NewIterator()
	it.UpdatePackagers(bc)

	var block Block
	bc.CreateBlock(types.CurrentUsers, accounts, [32]byte{}, &block, nil, nil)

	// 定时执行共识任务
	go func() {
		for _ = range types.Ticker.C {
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

func (bc *BC) NewIterator() *Iterator {
	it := Iterator{
		CurrentPackagerNum: 0,
		DB:                 bc.DB,
		Chan:               make(chan bool),
	}

	it.Packagers = it.UpdatePackagers(bc)
	return &it
}

// 运行迭代器
func (it *Iterator) Run(bc *BC, accounts *account.Accounts) {
	currentPackager := it.Packagers[it.CurrentPackagerNum]
	defer func() {
		it.CurrentPackagerNum++
		if len(it.Packagers) == it.CurrentPackagerNum {
			it.UpdatePackagers(bc)
			it.CurrentPackagerNum = 0
		}
	}()

	var block Block
	block.Txs = bc.TxPool
	go bc.CreateBlock(currentPackager, accounts, bc.LastBlockHash, &block, it.Chan, it.Packagers)

	go func() {
		select {
		case _, ok := <-it.Chan:
			if ok {
				it.DB.View(func(tx *bolt.Tx) error {
					b := tx.Bucket([]byte(types.BlockChainBucketName))
					blockinfo, _ := json.Marshal(block)
					b.Put([]byte{byte(block.Header.Height)}, bc.LastBlockHash[:])
					b.Put(bc.LastBlockHash[:], blockinfo)

					return nil
				})
			}
		case <-time.After(types.Interval):
			// 多少秒内没处理完毕
			// 取消这次打包，并撤销管理员职位
			it.DB.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte(types.AccountBucketName))
				b.Put([]byte(currentPackager), []byte(types.NodeTypes))
				return nil
			})

			fmt.Errorf("未在指定时间内打包完成，你已经失去管理员资格！\n")
		}
	}()
}

func (it *Iterator) UpdatePackagers(bc *BC) []string {
	var packagers []string
	bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(types.AccountBucketName))
		_ = b.ForEach(func(k, v []byte) error {
			if string(v) == string(types.AdminTypes) || string(v) == string(types.SuperTypes) {
				packagers = append(packagers, string(k))
			}
			return nil
		})
		return nil
	})

	return packagers
}

func (bc *BC) QueryBlock(order string) {
	height, err := strconv.Atoi(order)
	if err != nil {
		fmt.Println(err)
	}

	bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(types.BlockChainBucketName))

		blockHash := b.Get([]byte{byte(height)})

		if blockHash == nil {
			fmt.Errorf("该区块高度不存在！")
			return nil
		}

		var block *Block
		blockInfo := b.Get(blockHash)
		json.Unmarshal(blockInfo, &block)
		fmt.Printf("区块信息：%+v\n", block)

		return nil
	})
}

func (bc *BC) SendTx(tx *tx.Tx) {
	bc.TxPool = append(bc.TxPool, *tx)
	fmt.Println("交易发送成功！")
}

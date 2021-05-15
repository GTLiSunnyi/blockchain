package blockchain

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"

	"mybc/tx"
	"mybc/types"
	"mybc/wallet"
)

type BC struct {
	DB     *bolt.DB
	TxPool []tx.TX
	Iterator
	LastBlockHash []byte
}

// 创建区块链
func NewBC(ws *wallet.Wallets) *BC {
	db, err := bolt.Open(types.DBName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(types.BlockChainBucketName))
		if b == nil {
			// 桶不存在则创建
			b, err = tx.CreateBucket([]byte(types.BlockChainBucketName))
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})

	var block *Block
	go CreateBlock(types.CurrentUsers, ws.Gather[types.CurrentUsers].PriKey, nil, []byte(""), block, nil)
	bc := &BC{DB: db, LastBlockHash: []byte("")}

	return bc
}

func (bc *BC) RunBC(ws *wallet.Wallets) {
	it := bc.NewIterator()

	// 定时执行共识任务
	ticker := time.NewTicker(time.Duration(types.Interval) * time.Second)
	go func() {
		for _ = range ticker.C {
			it.Run(bc, ws)
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
		CurrentPackagerNum: 1,
		DB:                 bc.DB,
		Chan:               make(chan bool),
	}

	it.Packagers = it.UpdatePackagers(bc)
	return &it
}

// 运行迭代器
func (it *Iterator) Run(bc *BC, ws *wallet.Wallets) {
	currentPackager := it.Packagers[it.CurrentPackagerNum]
	if len(it.Packagers) == it.CurrentPackagerNum {
		it.CurrentPackagerNum = 0
	} else {
		it.CurrentPackagerNum++
	}

	var block *Block
	go CreateBlock(currentPackager, ws.Gather[currentPackager].PriKey, bc.TxPool, bc.LastBlockHash, block, it.Chan)
	it.Chan <- true

	go func() {
		select {
		case _, ok := <-it.Chan:
			if ok {
				it.DB.View(func(tx *bolt.Tx) error {
					b := tx.Bucket([]byte(types.BlockChainBucketName))
					blockinfo, _ := json.Marshal(block)
					b.Put(bc.LastBlockHash, blockinfo)

					bc.LastBlockHash = block.Hash
					return nil
				})
			}
		case <-time.After(time.Second * time.Duration(types.Interval)):
			// 多少秒内没处理完毕
			// 取消这次打包，并撤销管理员职位
			it.DB.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte(types.AccountBucketName))
				b.Put([]byte(currentPackager), []byte(types.NodeTypes))
				return nil
			})

			fmt.Errorf("未在指定时间内打包完成，你已经失去管理员资格！")
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

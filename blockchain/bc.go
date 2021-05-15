package blockchain

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/boltdb/bolt"

	"github.com/GTLiSunnyi/blockchain/cmd"
	"github.com/GTLiSunnyi/blockchain/tx"
	"github.com/GTLiSunnyi/blockchain/types"
	"github.com/GTLiSunnyi/blockchain/wallet"
)

type BC struct {
	DB     *bolt.DB
	TxPool []tx.TX
	Iterator
	LastBlockHash [32]byte
}

// 创建区块链
func NewBC(address string, ws *wallet.Wallets, db *bolt.DB) *BC {
	var block Block
	bc := &BC{DB: db, LastBlockHash: [32]byte{}, TxPool: []tx.TX{}}

	bc.CreateBlock(address, ws.Gather[types.CurrentUsers].PriKey, nil, [32]byte{}, &block, nil)

	return bc
}

func (bc *BC) RunBC(ws *wallet.Wallets, cmd *cmd.Cmd) {
	it := bc.NewIterator()

	// 定时执行共识任务
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for _ = range ticker.C {
			select {
			case <-cmd.ChanList:

			default:
				it.UpdatePackagers(bc)
				it.Run(bc, ws)
			}

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
func (it *Iterator) Run(bc *BC, ws *wallet.Wallets) {
	currentPackager := it.Packagers[it.CurrentPackagerNum]
	defer func() {
		it.CurrentPackagerNum++
		if len(it.Packagers) == it.CurrentPackagerNum {
			it.CurrentPackagerNum = 0
		}
	}()

	var block Block
	go bc.CreateBlock(currentPackager, ws.Gather[currentPackager].PriKey, bc.TxPool, bc.LastBlockHash, &block, it.Chan)

	go func() {
		select {
		case _, ok := <-it.Chan:
			if ok {
				it.DB.View(func(tx *bolt.Tx) error {
					b := tx.Bucket([]byte(types.BlockChainBucketName))
					blockinfo, _ := json.Marshal(block)
					b.Put(bc.LastBlockHash[:], blockinfo)

					bc.LastBlockHash = block.Header.Hash
					return nil
				})
			}
		case <-time.After(time.Second * 5):
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

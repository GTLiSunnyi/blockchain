package main

import (
	"github.com/boltdb/bolt"

	"mybc/blockchain"
	"mybc/cmd"
	"mybc/types"
	"mybc/wallet"
)

func main() {
	ws := wallet.NewWallets()

	w := ws.NewWallet()
	types.CurrentUsers = w.GetAddress()

	// 账户数据库
	ws.DB.Update(func(tx *bolt.Tx) error {
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

	// 创建区块链
	bc := blockchain.NewBC(ws)

	defer bc.DB.Close()
	defer ws.DB.Close()

	// 运行区块链
	command := cmd.Cmd{}
	command.SuperRun(ws)

	select {}
}

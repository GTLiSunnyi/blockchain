package main

import (
	"log"

	"github.com/boltdb/bolt"

	"github.com/GTLiSunnyi/blockchain/blockchain"
	"github.com/GTLiSunnyi/blockchain/cmd"
	"github.com/GTLiSunnyi/blockchain/types"
	"github.com/GTLiSunnyi/blockchain/wallet"
)

func main() {
	db, err := bolt.Open(types.DBName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(types.AccountBucketName))
		if b == nil {
			// 桶不存在则创建
			b, err = tx.CreateBucket([]byte(types.AccountBucketName))
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
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

	ws := wallet.NewWallets(db)

	_, address := ws.NewWallet()
	types.CurrentUsers = address

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

	// 创建cmd
	command := cmd.NewCmd()

	// 创建区块链
	bc := blockchain.NewBC(address, ws, db)
	bc.RunBC(ws, command)

	command.SuperRun(ws)

	select {}
}

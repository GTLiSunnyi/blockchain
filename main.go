package main

import (
	"github.com/boltdb/bolt"

	"mybc/accounts"
	"mybc/blockchain"
	"mybc/cmd"
	"mybc/types"
	"mybc/wallet"
)

func main() {
	wallet.Ws.NewWallets()

	key := wallet.Ws.NewWallet()
	types.CurrentUsers = key.GetAddress()

	// 创建超级管理员数据库
	db1, err := bolt.Open(types.DBName, 0600, nil)
	if err != nil {
		panic(err)
	}
	db1.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(types.SuperAccountBucketName))
		if b == nil {
			// 桶不存在则创建
			b, err = tx.CreateBucket([]byte(types.SuperAccountBucketName))
			if err != nil {
				panic(err)
			}
		}
		b.Put([]byte("SuperAccont"), []byte(types.CurrentUsers))
		return nil
	})

	// 创建管理员数据库
	db2, err := bolt.Open(types.DBName, 0600, nil)
	if err != nil {
		panic(err)
	}
	db2.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(types.AdminAccountBucketName))
		if b == nil {
			// 桶不存在则创建
			b, err = tx.CreateBucket([]byte(types.AdminAccountBucketName))
			if err != nil {
				panic(err)
			}
		}

		return nil
	})

	// 创建普通节点数据库
	db3, err := bolt.Open(types.DBName, 0600, nil)
	if err != nil {
		panic(err)
	}
	db3.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(types.NodeAccountBucketName))
		if b == nil {
			// 桶不存在则创建
			b, err = tx.CreateBucket([]byte(types.NodeAccountBucketName))
			if err != nil {
				panic(err)
			}
		}

		return nil
	})

	// 创建区块链
	blockchain.BlockChain = blockchain.NewBC()

	accounts.SuperAdmin = &accounts.Super{Address: types.CurrentUsers,
		AdminDB: db2,
		NodeDB:  db3,
	}

	defer db1.Close()
	defer db2.Close()
	defer db3.Close()
	defer blockchain.BlockChain.DB.Close()

	// 运行区块链
	command := cmd.Cmd{}
	command.SuperRun()
}

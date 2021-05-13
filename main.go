package main

import (
	"net"

	"github.com/boltdb/bolt"

	"mybc/accounts"
	"mybc/blockchain"
	"mybc/cmd"
	"mybc/types"
	"mybc/wallet"
)

func main() {
	key := wallet.NewWallet()

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

		b.Put([]byte("SuperAccont"), []byte(key.GetAddress()))
		return nil
	})

	types.CurrentUsers = key.GetAddress()

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

	accounts.SuperAdmin = &accounts.Super{Key: key,
		SuperDB: db1,
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

func checkIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

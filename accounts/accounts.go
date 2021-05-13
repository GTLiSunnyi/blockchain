package accounts

import (
	"github.com/boltdb/bolt"

	"mybc/tx"
	"mybc/types"
	"mybc/wallet"
)

var SuperAdmin *Super

type Super struct {
	Key     *wallet.Wallet
	Tx      *[]tx.TX
	SuperDB *bolt.DB
	AdminDB *bolt.DB
	NodeDB  *bolt.DB
}

func IsInAccounts(addr string) string {
	addrbyte := []byte(addr)
	var accountType []byte

	SuperAdmin.NodeDB.View(
		func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(types.NodeAccountBucketName))
			accountType = b.Get(addrbyte)

			if accountType != nil {
				accountType = []byte("普通节点")
			}
			return nil
		},
	)

	SuperAdmin.AdminDB.View(
		func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(types.AdminAccountBucketName))
			accountType := b.Get(addrbyte)

			if accountType != nil {
				accountType = []byte("管理员")
			}
			return nil
		},
	)

	SuperAdmin.SuperDB.View(
		func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(types.SuperAccountBucketName))
			accountType := b.Get(addrbyte)

			if accountType != nil {
				accountType = []byte("超级管理员")
			}
			return nil
		},
	)

	return string(accountType)
}

// 创建普通节点
func CreateNodeAccount() *wallet.Wallet {
	key := wallet.NewWallet()

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

	return key
}

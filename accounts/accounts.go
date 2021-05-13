package accounts

import (
	"github.com/boltdb/bolt"

	"mybc/tx"
	"mybc/types"
)

var SuperAdmin *Super

type Super struct {
	Address string
	Tx      *[]tx.TX
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
			accountType = b.Get(addrbyte)

			if accountType != nil {
				accountType = []byte("管理员")
			}
			return nil
		},
	)

	if addr == SuperAdmin.Address {
		return "超级管理员"
	}

	return string(accountType)
}

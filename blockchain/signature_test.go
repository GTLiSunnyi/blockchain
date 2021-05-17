package blockchain

import (
	"log"
	"testing"

	"github.com/boltdb/bolt"

	"github.com/GTLiSunnyi/blockchain/account"
	"github.com/GTLiSunnyi/blockchain/tx"
	"github.com/GTLiSunnyi/blockchain/types"
)

func TestSignature(t *testing.T) {
	db, err := bolt.Open("test", 0600, nil)
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(types.AccountBucketName))
		if b == nil {
			// 桶不存在则创建
			_, err = tx.CreateBucket([]byte(types.AccountBucketName))
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
	defer db.Close()

	var accounts = account.NewAccounts(db)
	var w, address = accounts.NewAccount(types.NodeTypes)

	var block Block
	block.Txs = []tx.Tx{*tx.NewTx("asdfghjk", address)}

	block.Txs[0].Sign(w.PriKey)
	isValid := block.Txs[0].IsValid(w.PubKey)
	if isValid {
		t.Log("success")
	} else {
		t.Error("failed")
	}

	var newAccount, _ = accounts.NewAccount(types.NodeTypes)
	isValid = block.Txs[0].IsValid(newAccount.PubKey)
	if !isValid {
		t.Log("success")
	}
}

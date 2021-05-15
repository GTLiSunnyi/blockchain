package blockchain

import (
	"log"
	"testing"

	"github.com/boltdb/bolt"

	"github.com/GTLiSunnyi/blockchain/tx"
	"github.com/GTLiSunnyi/blockchain/types"
	"github.com/GTLiSunnyi/blockchain/wallet"
)

func TestSignature(t *testing.T) {
	var block Block
	block.TXs = []tx.TX{*tx.NewFileTx([]byte("asdfghjk"))}

	db, err := bolt.Open("test", 0600, nil)
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
	defer db.Close()

	var ws = wallet.NewWallets(db)
	var w, _ = ws.NewWallet()

	block.Sign(&block.TXs[0], w.PriKey)
	isValid := block.IsValid(w.PubKey)
	if isValid {
		t.Log("success")
	} else {
		t.Error("failed")
	}

	var newW, _ = ws.NewWallet()
	isValid = block.IsValid(newW.PubKey)
	if !isValid {
		t.Log("success")
	}
}

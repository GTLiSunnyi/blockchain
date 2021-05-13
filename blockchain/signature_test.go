package blockchain

import (
	"testing"

	"mybc/tx"
	"mybc/wallet"
)

func TestSignature(t *testing.T) {
	var block Block
	block.MerkleRoot = []byte("1234567890-=asdfghjklzxcvbnm,./")
	block.TXs = &[]tx.TX{tx.TX{Id: []byte("asdfghjkl")}}
	var wallet = wallet.NewWallet()
	_ = block.setMerkleAndGetSignature(wallet.PriKey)
	isValid := block.IsValid(wallet.PubKey)
	if isValid {
		t.Log("success")
	} else {
		t.Error("failed")
	}
}

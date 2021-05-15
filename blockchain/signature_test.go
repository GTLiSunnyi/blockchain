package blockchain

import (
	"testing"

	"mybc/tx"
	"mybc/wallet"
)

func TestSignature(t *testing.T) {
	var block Block
	block.TXs = []tx.TX{*tx.NewFileTx([]byte("asdfghjk"))}

	var ws = wallet.NewWallets()
	var w = ws.NewWallet()

	block.setMerkleAndTxSignature(w.PriKey)
	isValid := block.IsValid(w.PubKey)
	if isValid {
		t.Log("success")
	} else {
		t.Error("failed")
	}

	block.setMerkleAndTxSignature(w.PriKey)
	isValid = block.IsValid(ws.NewWallet().PubKey)
	if !isValid {
		t.Log("success")
	}
}

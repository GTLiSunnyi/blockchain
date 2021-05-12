package accounts

import (
	"mybc/blockchain"
	"mybc/wallet"
)

type SuperAccount struct {
	Key     *wallet.Wallet  // 自己的key
	Wallets *wallet.Wallets // 所有人的key
	BC      *blockchain.BC
}

type AdminAccounts struct {
	AdminAccounts []string
}

type node struct {
	Key *wallet.Wallet
}

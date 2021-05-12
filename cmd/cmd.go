package cmd

import (
	"fmt"
	"time"
)

type Cmd struct {
}

func (cmd *Cmd) Run(addr string) {
	// 提示信息
	const prompt = `
*******************
send     转账交易
balance	 获取余额
wallet   创建钱包
list     钱包集合
print    打印区块
quit     退出程序
*******************
`

	var order string // 接收命令
	fmt.Println(prompt)
	for {
		fmt.Scan(&order)
		switch order {
		case "send":
			bc.addBlock(true)
		case "balance":
			fmt.Println("请输入获取对象")
			fmt.Scan(&order)
			balance, isOrderExist := bc.GetBalance(order)
			if isOrderExist {
				fmt.Printf("%v的余额为%v\n", order, balance)
			}
		case "wallet":
			wallets := NewWallets()
			wallets.CreateWallets()
		case "list":
			wallets := NewWallets()
			wallets.GetList()
		case "print":
			it := bc.NewIterator()
			// 打印区块链
			for {
				block := it.Run()
				fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++")
				fmt.Printf("preBlockHash: %v\n", block.PreBlockHash)
				fmt.Printf("merkleRoot: %v\n", block.MerkleRoot)
				timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
				fmt.Printf("timeStamp : %s\n", timeFormat)
				fmt.Printf("difficulty: %v\n", block.Difficulty)
				fmt.Printf("nonce: %v\n", block.Nonce)
				fmt.Printf("hash: %v\n", block.Hash)
				if block.PreBlockHash == nil {
					break
				}
			}
		case "quit":
			return
		default:
			fmt.Println("请输入正确的命令")
		}
	}
}

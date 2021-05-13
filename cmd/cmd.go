package cmd

import (
	"fmt"

	"mybc/accounts"
	"mybc/query"
	"mybc/types"
)

type Cmd struct {
}

func (cmd *Cmd) SuperRun() {
	// 提示信息
	const prompt = `
*******************
use      切换用户
addPermissions 增加权限
createNode 创建节点
query    查询区块、节点
quit     退出程序
*******************
`

	go run()

	var order string // 接收命令
	fmt.Println(prompt)
	for {
		fmt.Scan(&order)
		switch order {
		case "use":
			// 切换用户
			fmt.Println("请输入切换用户的地址：")
			fmt.Scan(&order)
			cmd.SwitchUsers(true, order)
		case "addPermissions":

		case "createNode":

		case "query":
			query.Query()
		case "quit":
			fmt.Println("exit...")
			return
		default:
			fmt.Println("请输入正确的命令")
		}
	}
}

func (cmd *Cmd) NodeRun() {
	// 提示信息
	const prompt = `
*******************
send     发送文件
query    查询区块
quit     退出程序
*******************
`

	var order string // 接收命令
	fmt.Println(prompt)
	for {
		fmt.Scan(&order)
		switch order {
		case "send":
			// bc.addBlock(true)
		case "query":
			// it := bc.NewIterator()
			// // 打印区块链
			// for {
			// 	block := it.Run()
			// 	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++")
			// 	fmt.Printf("preBlockHash: %v\n", block.PreBlockHash)
			// 	fmt.Printf("merkleRoot: %v\n", block.MerkleRoot)
			// 	timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
			// 	fmt.Printf("timeStamp : %s\n", timeFormat)
			// 	fmt.Printf("difficulty: %v\n", block.Difficulty)
			// 	fmt.Printf("nonce: %v\n", block.Nonce)
			// 	fmt.Printf("hash: %v\n", block.Hash)
			// 	if block.PreBlockHash == nil {
			// 		break
			// 	}
			// }
		case "quit":
			return
		default:
			fmt.Println("请输入正确的命令")
		}
	}
}

func (cmd *Cmd) AdminRun() {
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
		// 	bc.addBlock(true)
		// case "balance":
		// 	fmt.Println("请输入获取对象")
		// 	fmt.Scan(&order)
		// 	balance, isOrderExist := bc.GetBalance(order)
		// 	if isOrderExist {
		// 		fmt.Printf("%v的余额为%v\n", order, balance)
		// 	}
		case "wallet":
		// 	wallets := NewWallets()
		// 	wallets.CreateWallets()
		// case "list":
		// 	wallets := NewWallets()
		// 	wallets.GetList()
		// case "print":
		// 	it := bc.NewIterator()
		// 	// 打印区块链
		// 	for {
		// 		block := it.Run()
		// 		fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++")
		// 		fmt.Printf("preBlockHash: %v\n", block.PreBlockHash)
		// 		fmt.Printf("merkleRoot: %v\n", block.MerkleRoot)
		// 		timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
		// 		fmt.Printf("timeStamp : %s\n", timeFormat)
		// 		fmt.Printf("difficulty: %v\n", block.Difficulty)
		// 		fmt.Printf("nonce: %v\n", block.Nonce)
		// 		fmt.Printf("hash: %v\n", block.Hash)
		// 		if block.PreBlockHash == nil {
		// 			break
		// 		}
		// 	}
		case "quit":
			return
		default:
			fmt.Println("请输入正确的命令")
		}
	}
}

func (cmd *Cmd) SwitchUsers(isSuper bool, addr string) {
	accountType := accounts.IsInAccounts(addr)
	if accountType == "" {
		if isSuper {
			nodeAddr := accounts.CreateNodeAccount()
			fmt.Println("切换到一个新的普通节点，addr：", nodeAddr)
			accountType = "普通节点"
			types.CurrentUsers = addr
		} else {
			fmt.Println("输入的账户不存在！请使用超级管理员创建账户。")
			return
		}
	} else {
		types.CurrentUsers = addr
		fmt.Sprintf("切换成功！\n当前用户：%s, %s", accountType, addr)
	}

	switch accountType {
	case "超级管理员":
		cmd.SuperRun()
	case "管理员":
		cmd.AdminRun()
	case "普通节点":
		cmd.NodeRun()
	}
}

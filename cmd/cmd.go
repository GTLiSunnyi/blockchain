package cmd

import (
	"fmt"
	"os"

	"github.com/GTLiSunnyi/blockchain/types"
	"github.com/GTLiSunnyi/blockchain/wallet"
)

type Cmd struct {
	ChanList chan string
}

func NewCmd() *Cmd {
	cmd := &Cmd{make(chan string, 100)}
	return cmd
}

func (cmd *Cmd) SuperRun(ws *wallet.Wallets) {
	// 提示信息
	const prompt = `
*******************************
add            添加用户
use            切换用户         
addPms         增加权限       
createNode     创建节点       
query          查询区块、节点 
quit           退出程序
*******************************
`

	var order string // 接收命令
	fmt.Println(prompt)
	for {
		fmt.Scan(&order)
		switch order {
		case "add":
		case "use":
			// 切换用户
			fmt.Println("请输入切换用户的地址：")
			cmd.ChanList <- "stop"
			fmt.Scan(&order)
			cmd.SwitchUsers(true, ws, order)
			cmd.ChanList <- "run"
		case "addPms":
			fmt.Println("请输入切换用户的地址：")
			fmt.Scan(&order)
			ws.AddPms(order)
		case "createNode":
			ws.CreateNodeAccount()
		case "query":
			fmt.Println("请输入需要查询的事物：")
			fmt.Println(`
	*******************************
	account      查询账户
	block        查询区块
	*******************************
	`)
			fmt.Scan(&order)
			switch order {
			case "account":
				ws.QueryAccount()
			case "block":
			default:
				fmt.Println("请输入正确的命令")
			}
		case "quit":
			fmt.Println("exit...")
			os.Exit(-1)
		default:
			fmt.Println("请输入正确的命令")
		}
	}
}

func (cmd *Cmd) AdminRun(ws *wallet.Wallets) {
	// 提示信息
	const prompt = `
*******************************
use       切换用户
send      发送普通交易
denom	  denom
query     查询区块、节点
quit      退出程序
*******************************
`

	var order string // 接收命令
	fmt.Println(prompt)
	for {
		fmt.Scan(&order)
		switch order {
		case "use":
			// 切换用户
			fmt.Println("请输入切换用户的地址：")
			fmt.Scan(&order)
			cmd.SwitchUsers(false, ws, order)
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
			fmt.Println("exit...")
			os.Exit(-1)
		default:
			fmt.Println("请输入正确的命令")
		}
	}
}

func (cmd *Cmd) NodeRun(ws *wallet.Wallets) {
	// 提示信息
	const prompt = `
*******************************
use      	切换用户
send        发送普通交易
query       查询区块、节点
quit        退出程序
*******************************
`

	var order string // 接收命令
	fmt.Println(prompt)
	for {
		fmt.Scan(&order)
		switch order {
		case "use":
			// 切换用户
			fmt.Println("请输入切换用户的地址：")
			fmt.Scan(&order)
			cmd.SwitchUsers(false, ws, order)
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
			fmt.Println("exit...")
			os.Exit(-1)
		default:
			fmt.Println("请输入正确的命令")
		}
	}
}

func (cmd *Cmd) SwitchUsers(isSuper bool, ws *wallet.Wallets, addr string) {
	var accountType string

	if !ws.IsInAccounts(addr) {
		if isSuper {
			nodeAddr := ws.CreateNodeAccount()
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
		cmd.SuperRun(ws)
	case "管理员":
		cmd.AdminRun(ws)
	case "普通节点":
		cmd.NodeRun(ws)
	}
}

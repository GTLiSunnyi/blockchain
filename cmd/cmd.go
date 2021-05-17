package cmd

import (
	"fmt"
	"os"

	"github.com/GTLiSunnyi/blockchain/account"
	"github.com/GTLiSunnyi/blockchain/blockchain"
	"github.com/GTLiSunnyi/blockchain/denom"
	"github.com/GTLiSunnyi/blockchain/tx"
	"github.com/GTLiSunnyi/blockchain/types"
)

type Cmd struct {
	*blockchain.BC
	*account.Accounts
	*denom.Denoms
	ChanList chan string
}

func NewCmd() *Cmd {
	bc, db := blockchain.NewBC()
	accounts := account.NewAccounts(db)
	denoms := denom.NewDenoms(db)
	return &Cmd{bc, accounts, denoms, make(chan string)}
}

func (cmd *Cmd) SuperRun() {
	// 提示信息
	const prompt = `
*******************************
add            添加用户
use            切换用户         
addPms         增加权限 
rmPms          撤销权限
denom          denom/nft           
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
			types.Ticker.Stop()
			_ = cmd.Accounts.CreateNodeAccount()
			fmt.Println("添加用户成功！")
			types.Ticker.Reset(types.Interval)
		case "use":
			// 切换用户
			types.Ticker.Stop()
			fmt.Println("请输入切换用户的地址：")
			fmt.Scan(&order)
			cmd.SwitchUsers(order)
			types.Ticker.Reset(types.Interval)
		case "addPms":
			types.Ticker.Stop()
			fmt.Println("请输入用户的地址：")
			fmt.Scan(&order)
			cmd.Accounts.AddPms(order)
			types.Ticker.Reset(types.Interval)
		case "rmPms":
			types.Ticker.Stop()
			fmt.Println("请输入用户地址：")
			fmt.Scan(&order)
			cmd.Accounts.RmPms(order)
			types.Ticker.Reset(types.Interval)
		case "query":
			types.Ticker.Stop()
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
				cmd.Accounts.QueryAccount()
			case "block":
				fmt.Println("请输入区块高度：")
				fmt.Scan(&order)
				cmd.BC.QueryBlock(order)
			default:
				fmt.Println("请输入正确的命令")
			}
			types.Ticker.Reset(types.Interval)
		case "quit":
			types.Ticker.Stop()
			fmt.Println("exit...")
			os.Exit(-1)
		default:
			fmt.Println("请输入正确的命令")
		}
	}
}

func (cmd *Cmd) AdminRun() {
	// 提示信息
	const prompt = `
*******************************
use       切换用户
send      发送普通交易
denom	  denom
query     查询区块、denom
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
			types.Ticker.Stop()
			fmt.Println("请输入切换用户的地址：")
			fmt.Scan(&order)
			cmd.SwitchUsers(order)
			types.Ticker.Reset(types.Interval)
		case "send":
			types.Ticker.Stop()
			fmt.Println("请输入交易数据：")
			fmt.Scan(&order)
			cmd.BC.SendTx(tx.NewTx(order, types.CurrentUsers))
			types.Ticker.Reset(types.Interval)
		case "query":
			types.Ticker.Stop()
			fmt.Println("请输入需要查询的事物：")
			fmt.Println(`
	*******************************
	block        查询区块
	denom        查询denom
	*******************************
	`)
			fmt.Scan(&order)
			switch order {
			case "block":
				fmt.Println("请输入区块高度：")
				fmt.Scan(&order)
				cmd.BC.QueryBlock(order)
			case "denom":
				cmd.Denoms.Query(types.CurrentUsers)
			default:
				fmt.Println("请输入正确的命令")
			}
			types.Ticker.Reset(types.Interval)
		case "quit":
			types.Ticker.Stop()
			fmt.Println("exit...")
			os.Exit(-1)
		default:
			fmt.Println("请输入正确的命令")
		}
	}
}

func (cmd *Cmd) NodeRun() {
	// 提示信息
	const prompt = `
*******************************
use      	切换用户
send        发送普通交易
query       查询区块、denom
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
			cmd.ChanList <- "stop"
			fmt.Scan(&order)
			cmd.SwitchUsers(order)
			cmd.ChanList <- "start"
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

func (cmd *Cmd) SwitchUsers(address string) {
	if !cmd.Accounts.IsInAccounts(address) {
		fmt.Println("输入的账户不存在！请使用超级管理员创建账户。")
		return
	} else {
		types.CurrentUsers = address
		fmt.Printf("切换成功！\n当前用户：%s, %s\n", cmd.Accounts.Gather[address].AccountType, address)
	}

	switch cmd.Accounts.Gather[address].AccountType {
	case types.SuperTypes:
		cmd.SuperRun()
	case types.AdminTypes:
		cmd.AdminRun()
	case types.NodeTypes:
		cmd.NodeRun()
	default:
		fmt.Println("切换节点失败")
	}
}

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
		case "send":
			types.Ticker.Stop()
			fmt.Println("请输入交易数据：")
			fmt.Scan(&order)
			cmd.BC.SendTx(tx.NewTx(order, types.CurrentUsers))
			types.Ticker.Reset(types.Interval)
		case "denom":
			types.Ticker.Stop()
			fmt.Println("请输入要执行的事物：")
			fmt.Println(`
	*******************************
	denom        创建denom
	nft          创建nft
	transfer     交易nft
	*******************************
	`)
			fmt.Scan(&order)
			switch order {
			case "denom":
				fmt.Println("请输入denom名称：")
				fmt.Scan(&order)
				cmd.Denoms.CreateDenom(order, cmd.Accounts.Gather[types.CurrentUsers])
			case "nft":
				var nftName string
				var denomName string
				fmt.Println("请输入nft名称：")
				fmt.Scan(&nftName)
				fmt.Println("请输入denom名称：")
				fmt.Scan(&denomName)
				fmt.Println("请输入uri：")
				fmt.Scan(&order)
				cmd.Denoms.MintNft(nftName, denomName, order)
			case "transfer":
				var nftName string
				fmt.Println("请输入nft名称：")
				fmt.Scan(&nftName)
				fmt.Println("请输入denom名称：")
				fmt.Scan(&order)
				cmd.Denoms.TransferNft(nftName, order, cmd.Accounts.Gather[types.CurrentUsers])
			default:
				fmt.Println("请输入正确的命令")
			}
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
			fmt.Println("请输入区块高度：")
			fmt.Scan(&order)
			cmd.BC.QueryBlock(order)
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
		types.Ticker.Reset(types.Interval)
		cmd.SuperRun()
	case types.AdminTypes:
		types.Ticker.Reset(types.Interval)
		cmd.AdminRun()
	case types.NodeTypes:
		types.Ticker.Reset(types.Interval)
		cmd.NodeRun()
	default:
		fmt.Println("切换节点失败")
	}
}

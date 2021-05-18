package main

import (
	"github.com/GTLiSunnyi/blockchain/cmd"
)

func main() {
	// 创建cmd
	command := cmd.NewCmd()
	// TODO denom相关的操作要发送交易

	defer command.BC.DB.Close()

	command.BC.RunBC(command.Accounts, command.ChanList)
	command.SuperRun()

	select {}
}

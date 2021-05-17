package main

import (
	"github.com/GTLiSunnyi/blockchain/cmd"
)

func main() {
	// 创建cmd
	command := cmd.NewCmd()
	// super is node
	// pack block tx

	defer command.BC.DB.Close()
	defer command.Accounts.DB.Close()

	command.BC.RunBC(command.Accounts, command.ChanList)
	command.SuperRun()

	select {}
}

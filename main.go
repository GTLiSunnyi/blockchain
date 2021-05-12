package main

import (
	"flag"
	"fmt"
	"net"

	"mybc/blockchain"
	"mybc/cmd"
	"mybc/wallet"
)

var IP string
var Port string
var IsSuper bool

func main() {
	defer BC.DB.Close()
	var addr string

	flag.StringVar(&IP, "ip", "", "-ip=:127.0.0.1")
	flag.StringVar(&Port, "port", "", "-port=:8080")
	flag.StringVar(&addr, "addr", "", "-addr=abc")
	flag.BoolVar(&IsSuper, "isSuper", false, "-isSuper")
	flag.Parse()

	if IP == "" && Port == "" && addr == "" {
		return
	}
	if addr == "" && Port == "" || IP == "" {
		return
	} else {
		ip := checkIP()
		if ip == "" {
			fmt.Errorf("Unable to get a avilable ip")
			return
		}
	}

	if !IsSuper {
		IsSuper = true
	}
	if addr == "" {
		key = wallet.NewWallet()
		Accounts = append(Accounts, key.GetAddress())
		addr = key.GetAddress()
	} else if !isInAccounts(addr) {
		fmt.Errorf("addr is not in the system!")
		return
	}

	comond := cmd.Cmd{}
	if addr == SuperAccount {
		// 创建链
		BC = blockchain.NewBC()
		// 创建创世块
		BC.AddBlock(false)
	}

	comond.Run(addr)
}

func checkIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func isInAccounts(addr string) bool {
	for _, v := range Accounts {
		if v == addr {
			return true
		}
	}
	for _, v := range AdminAccounts {
		if v == addr {
			return true
		}
	}
	return addr == SuperAccount
}

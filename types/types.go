package types

import "time"

var CurrentUsers string

var Height int = 0

const DBName = "blockchain.db"
const AccountBucketName = "Account"
const BlockChainBucketName = "block"
const DenomBucketName = "denom"

type AccountType string

const SuperTypes AccountType = "超级管理员"
const AdminTypes AccountType = "管理员"
const NodeTypes AccountType = "普通节点"

// 钱包集合文件名字
const FileName = "./accounts.dat"

const Interval = 6 * time.Second
const Version string = "v1.0.0"

var Ticker = time.NewTicker(Interval)

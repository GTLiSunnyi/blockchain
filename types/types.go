package types

import "time"

var CurrentUsers string

var Height int = 0

const DBName = "blockchain.db"
const AccountBucketName = "Account"
const BlockChainBucketName = "block"

type AccountType string

const SuperTypes AccountType = "超级管理员"
const AdminTypes AccountType = "管理员"
const NodeTypes AccountType = "普通节点"

// 钱包集合文件名字
const FileName = "./wallets.dat"

const Interval = 10 * time.Second

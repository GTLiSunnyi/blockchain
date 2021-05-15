package types

var CurrentUsers string

const DBName = "mybc.db"

const AccountBucketName = "Account"

const BlockChainBucketName = "block"

type AccountType string

const SuperTypes AccountType = "超级管理员"
const AdminTypes AccountType = "管理员"
const NodeTypes AccountType = "普通节点"

const Interval uint64 = 5

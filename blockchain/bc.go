package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"

	"github.com/boltdb/bolt"

	"mybc/tx"
	"mybc/utils"
	"mybc/wallet"
)

type BC struct {
	DB            *bolt.DB
	LastBlockHash []byte
}

const DBName = "blockChain.db"
const BlockBucketName = "block"

// 创建区块链
func NewBC() *BC {
	db, err := bolt.Open(DBName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucketName))
		if b == nil {
			// 桶不存在则创建
			b, err = tx.CreateBucket([]byte(BlockBucketName))
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
	return &BC{DB: db}
}

// 添加区块
func (bc *BC) AddBlock(isTransfer bool) {
	var txs []tx.TX // 交易
	if isTransfer {
		// 是交易，则添加交易和挖矿交易
		var from, to, miner string
		var amount float64
		fmt.Println("请分别输入转账者、收款人、金额和挖矿者")
		fmt.Scan(&from, &to, &amount, &miner)
		TX := *tx.NewTX(from, to, amount, bc)
		if TX.Id == nil {
			return
		}
		txs = append(txs, TX)
		TX = *tx.NewCoinbaseTX(miner)
		if TX.Id == nil {
			return
		}
		txs = append(txs, TX)
	} else {
		// 不是交易，则创建创世地址和创世块
		ws := wallet.NewWallets()
		address := ws.CreateWallets()
		TX := *tx.NewCoinbaseTX(address)
		if TX.Id == nil {
			return
		}
		txs = append(txs, TX)
	}
	for _, tx := range txs {
		if !bc.VerifyTX(&tx) {
			fmt.Println("发现无效交易")
			return
		}
	}
	// 最后一个区块的哈希就是当前区块的PreBlockHash
	block, isVaild := NewBlock(txs, bc.LastBlockHash)
	if isVaild {
		bc.DB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(BlockBucketName))
			b.Put(block.Hash, utils.Serialize(block))
			b.Put([]byte("lastBlockHash"), block.Hash)
			bc.LastBlockHash = block.Hash
			return nil
		})
	} else {
		fmt.Println("挖矿失败")
		return
	}
}

// 迭代器
type Iterator struct {
	DB          *bolt.DB
	CurrentHash []byte
}

func (bc *BC) NewIterator() *Iterator {
	it := Iterator{
		DB:          bc.DB,
		CurrentHash: bc.LastBlockHash,
	}
	return &it
}

// 运行迭代器
func (it *Iterator) Run() *Block {
	var block Block
	it.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucketName))
		blockInfo := b.Get(it.CurrentHash)
		err := json.Unmarshal(blockInfo, &block)
		if err != nil {
			fmt.Println("迭代器反序列化错误")
			log.Panic(err)
		}
		it.CurrentHash = block.PreBlockHash
		return nil
	})
	return &block
}

// 余额信息结构体(存储的都是当前块的信息)
type UTXOInfo struct {
	TXId   []byte    // 交易id
	Index  int64     // output的索引值
	Output tx.Output // output本身
}

func (bc *BC) FindMyUtoxs(pubKey []byte) []UTXOInfo {
	// 记录input=owner的交易，遍历除了记录过的交易外的output=owner的总金额
	var UTXOInfos []UTXOInfo //新的返回结构
	it := bc.NewIterator()
	// 这是标识已经消耗过的utxo的结构，key是交易id，value是这个id里面的output索引的数组
	spentUTXOs := make(map[string][]int64)
	// 1.遍历账本
	for {
		block := it.Run()
		// 2.遍历交易
		for _, tx := range block.TXs {
			// 遍历交易输入:inputs
			for _, input := range tx.Inputs {
				// 判断当前被使用input是否为目标地址所有
				if bytes.Equal(input.PrePubKey, pubKey) {
					key := string(input.PreId)
					spentUTXOs[key] = append(spentUTXOs[key], input.PreIndex)
				}
			}
			key := string(tx.Id)
			indexes /*[]int64{0,1}*/ := spentUTXOs[key]
		OUTPUT:
			// 3.遍历output
			for i, output := range tx.Outputs {
				if len(indexes) != 0 {
					for _, j /*0, 1*/ := range indexes {
						if int64(i) == j {
							continue OUTPUT
						}
					}
				}
				// 4.找到属于我的所有output
				if bytes.Equal(output.PubKeyHash, tx.HashPubKey(pubKey)) {
					utxoinfo := UTXOInfo{tx.Id, int64(i), output}
					UTXOInfos = append(UTXOInfos, utxoinfo)
				}
			}
		}
		if len(block.PreBlockHash) == 0 {
			break
		}
	}
	return UTXOInfos
}

// 找到付款人合适的UTXOs
func (bc *BC) FindNeedUTXOs(pubKey []byte, amount float64) (map[string][]int64, float64) {
	needUtxos := make(map[string][]int64) // 标识能用的utxo, //返回的结构
	var resValue float64                  // 统计的金额
	// 复用FindMyUtxo函数，这个函数已经包含了所有信息
	utxoinfos := bc.FindMyUtoxs(pubKey)
	for _, utxoinfo := range utxoinfos {
		key := string(utxoinfo.TXId)
		needUtxos[key] = append(needUtxos[key], int64(utxoinfo.Index))
		resValue += utxoinfo.Output.Amount
		// 2.找到足够的金额则返回，否则就会最终返回账户余额
		if resValue >= amount {
			break
		}
	}
	return needUtxos, resValue
}

// 获取余额
func (bc *BC) GetBalance(owner string) (float64, bool) {
	Wallet, ok := wallet.NewWallets().Gather[owner]
	if !wallet.IsValidAddress(owner) {
		fmt.Println("地址错误")
		return 0.0, false
	} else if !ok {
		fmt.Println("该地址不存在")
		return 0.0, false
	}
	pubKey := Wallet.PubKey
	utxos := bc.FindMyUtoxs(pubKey)
	var balance = 0.00
	for _, utxo := range utxos {
		balance += utxo.Output.Amount
	}
	return balance, true
}

// 交易签名
func (bc *BC) SignTX(TX *tx.TX, priKey *ecdsa.PrivateKey) {
	// 1.遍历账本找到所有引用交易
	preTXs := make(map[string]*tx.TX)
	// 遍历tx的inputs，通过id查找所引用的交易
	for _, input := range TX.Inputs {
		preTX := bc.FindTX(input.PreId)
		if preTX == nil {
			fmt.Println("没有找到交易")
		} else {
			// 保存找到的交易
			preTXs[string(input.PreId)] = preTX
		}
	}
	TX.Sign(priKey, preTXs)
}

// 矿工校验流程
// 1. 找到交易input所引用的所有的交易prevTXs
// 2. 对交易进行校验
func (bc *BC) VerifyTX(TX *tx.TX) bool {

	// 校验的时候，如果是挖矿交易，直接返回true
	if TX.Inputs[0].PreIndex == -1 {
		return true
	}

	prevTXs := make(map[string]*tx.TX)
	//遍历tx的inputs，通过id去查找所引用的交易
	for _, input := range TX.Inputs {
		prevTx := bc.FindTX(input.PreId)

		if prevTx == nil {
			fmt.Printf("没有找到交易: %x\n", input.PreId)
		} else {
			//把找到的引用交易保存起来
			//0x222
			//0x333
			prevTXs[string(input.PreId)] = prevTx
		}
	}

	return TX.Verify(prevTXs)
}

func (bc *BC) FindTX(id []byte) *tx.TX {
	// 遍历区块链交易，通过对比id来识别
	it := bc.NewIterator()
	for {
		block := it.Run()
		for _, tx := range block.TXs {
			// 找到相同id交易
			if bytes.Equal(tx.Id, id) {
				fmt.Println("找到了所引用的交易")
				return &tx
			}
		}
		if block.PreBlockHash == nil {
			break
		}
	}
	return nil
}

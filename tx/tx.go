package tx

import (
	"github.com/btcsuite/btcutil/base58"
)

type TX struct {
	Signature [][]byte
	DenomTX   []byte
	FileTX    []byte
}

// 将地址转换为公钥哈希
func GetPubKeyHash(address string) []byte {
	data := base58.Decode(address)
	pubKeyHash := data[1:21]
	return pubKeyHash
}

// 创建文件交易
func NewFileTx(file []byte) *TX {
	return &TX{FileTX: file}
}

// 创建denom交易
func NewDenomTX(denomRecord []byte) *TX {
	return &TX{FileTX: denomRecord}
}

// 对交易进行签名
// func (tx *TX) Sign(priKey *ecdsa.PrivateKey, preTXs map[string]*TX) {
// 	fmt.Printf("对交易进行签名...\n")
// 	// 1.拷贝一份交易txCopy
// 	//  做相应裁剪：把每一个input的sign和pubkey设置为nil
// 	//  output不改变
// 	txCopy := tx.TrimmedCopy()
// 	// 2.遍历txCopy.inputs
// 	//  把这个input所引用的output的公钥哈希拿过来，赋值给pubkey
// 	for i, input := range txCopy.Inputs {
// 		// 找到引用的交易
// 		preTX := preTXs[string(input.PreId)]
// 		output := preTX.Outputs[input.PreIndex]
// 		// for循环迭代出来的数据是一个副本，对这个input进行修改，不会影响到原始数据
// 		// 所以我们这里需要使用下标方式修改

// 		// input.PubKey = output.PubKeyHash
// 		txCopy.Inputs[i].PrePubKey = output.PubKeyHash
// 		// 签名主要对数据的hash进行签名
// 		// 我们的数据都在交易中，求交易的哈希
// 		// TX的SetID函数就是对交易的哈希
// 		// 所以我们可以使用交易id作为我们签名的内容
// 		// 3.生成要签名的数据（哈希）
// 		txCopy.Id = utils.Serialize(txCopy)
// 		Id := txCopy.Id
// 		// 清理，原理同上
// 		// input.PubKey = nil
// 		txCopy.Inputs[i].PrePubKey = nil
// 		fmt.Printf("要签名的数据，Id：%x\n", Id)
// 		// 4.对数据进行签名r，s
// 		r, s, err := ecdsa.Sign(rand.Reader, priKey, Id)
// 		if err != nil {
// 			fmt.Printf("交易签名失败，err：%v\n", err)
// 		}
// 		// 5.拼接r，s为字节流
// 		signature := append(r.Bytes(), s.Bytes()...)
// 		// 6.赋值给原始交易的Signature字段
// 		tx.Inputs[i].Signature = signature
// 	}
// }

// func (tx *TX) Verify(preTXs map[string]*TX) bool {
// 	fmt.Printf("对交易进行校验...\n")
// 	// 1.拷贝修剪副本
// 	txCopy := tx.TrimmedCopy()
// 	// 2.遍历原始交易（注意不是txCopy）
// 	for i, input := range tx.Inputs {
// 		// 3.遍历原始交易的input所引用的前交易preTX
// 		preTX := preTXs[string(input.PreId)]
// 		output := preTX.Outputs[input.PreIndex]
// 		// 4.找到output的公钥哈希，赋值给txCopy对应的input
// 		txCopy.Inputs[i].PrePubKey = output.PubKeyHash
// 		// 5.还原签名数据
// 		txCopy.Id = utils.Serialize(txCopy)
// 		// 清理动作，重要！！！
// 		txCopy.Inputs[i].PrePubKey = nil

// 		verifyData := txCopy.Id
// 		fmt.Printf("verifyData: %x\n", verifyData)
// 		// 6.校验
// 		// 还原签名为r，s
// 		signature := input.Signature
// 		// 公钥字节流
// 		pubKeyBytes := input.PrePubKey
// 		r := big.Int{}
// 		s := big.Int{}
// 		rData := signature[:len(signature)/2]
// 		sData := signature[len(signature)/2:]
// 		r.SetBytes(rData)
// 		s.SetBytes(sData)
// 		// type PublicKey struct {
// 		// 	elliptic.Curve
// 		// 	X, Y *big.Int
// 		// }

// 		// 还原公钥为curve，x，y
// 		x := big.Int{}
// 		y := big.Int{}
// 		xData := pubKeyBytes[:len(pubKeyBytes)/2]
// 		yData := pubKeyBytes[len(pubKeyBytes)/2:]
// 		x.SetBytes(xData)
// 		y.SetBytes(yData)
// 		curve := elliptic.P256()
// 		publicKey := ecdsa.PublicKey{curve, &x, &y}
// 		// 数据、签名、公钥准备完毕，开始校验
// 		// func Verify(pub &*PublicKey, hash []byte, r, s *big.Int) bool
// 		if !ecdsa.Verify(&publicKey, verifyData, &r, &s) {
// 			return false
// 		}
// 	}
// 	return true
// }

// trim：裁剪
// func (tx *TX) TrimmedCopy() *TX {
// 	var inputs []Input
// 	for _, input := range tx.Inputs {
// 		inputs = append(inputs, Input{input.PreId, input.PreIndex, nil, nil})
// 	}
// 	return &TX{tx.Id, inputs, tx.Outputs}
// }

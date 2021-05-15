package tx

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha1"
	"math/big"

	"github.com/btcsuite/btcutil/base58"
)

type TX struct {
	Signature  [][]byte
	DenomTX    []byte
	OrdinaryTX []byte
}

// 将地址转换为公钥哈希
func GetPubKeyHash(address string) []byte {
	data := base58.Decode(address)
	pubKeyHash := data[1:21]
	return pubKeyHash
}

// 创建文件交易
func NewFileTx(ordinaryTX []byte) *TX {
	return &TX{OrdinaryTX: ordinaryTX}
}

// 创建denom交易
func NewDenomTX(denomRecord []byte) *TX {
	return &TX{DenomTX: denomRecord}
}

// 对交易进行签名
func (tx *TX) Sign(prikey *ecdsa.PrivateKey) {
	hashText := sha1.Sum(append(tx.DenomTX, tx.OrdinaryTX...))

	//数字签名
	r, s, _ := ecdsa.Sign(rand.Reader, prikey, hashText[:])

	rText, _ := r.MarshalText()
	sText, _ := s.MarshalText()

	tx.Signature = [][]byte{rText, sText}
}

// 校验交易签名是否正确
func (tx *TX) IsValid(pubkey *ecdsa.PublicKey) bool {
	var r, s big.Int
	r.UnmarshalText(tx.Signature[0])
	s.UnmarshalText(tx.Signature[1])

	signatureCopy := tx.Signature
	tx.Signature = nil
	hashText := sha1.Sum(append(tx.DenomTX, tx.OrdinaryTX...))

	//认证
	res := ecdsa.Verify(pubkey, hashText[:], &r, &s)
	if !res {
		return false
	}

	tx.Signature = signatureCopy
	return true
}

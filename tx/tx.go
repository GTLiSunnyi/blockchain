package tx

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha1"
	"math/big"

	"github.com/btcsuite/btcutil/base58"
)

type Tx struct {
	Signature [][]byte
	Data      string
	Address   string
}

func NewTx(txData string, address string) *Tx {
	return &Tx{nil, txData, address}
}

// 将地址转换为公钥哈希
func GetPubKeyHash(address string) []byte {
	data := base58.Decode(address)
	pubKeyHash := data[1:21]
	return pubKeyHash
}

// 对交易进行签名
func (tx *Tx) Sign(prikey *ecdsa.PrivateKey) {
	hashText := sha1.Sum([]byte(tx.Data))

	//数字签名
	r, s, _ := ecdsa.Sign(rand.Reader, prikey, hashText[:])

	rText, _ := r.MarshalText()
	sText, _ := s.MarshalText()

	tx.Signature = [][]byte{rText, sText}
}

// 校验交易签名是否正确
func (tx *Tx) IsValid(pubkey *ecdsa.PublicKey) bool {
	var r, s big.Int
	r.UnmarshalText(tx.Signature[0])
	s.UnmarshalText(tx.Signature[1])

	signatureCopy := tx.Signature
	tx.Signature = nil
	hashText := sha1.Sum([]byte(tx.Data))

	//认证
	res := ecdsa.Verify(pubkey, hashText[:], &r, &s)
	if !res {
		return false
	}

	tx.Signature = signatureCopy
	return true
}

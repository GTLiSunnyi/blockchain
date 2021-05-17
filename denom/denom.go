package denom

import (
	"encoding/json"
	"fmt"

	"github.com/GTLiSunnyi/blockchain/account"
	"github.com/GTLiSunnyi/blockchain/types"
	"github.com/boltdb/bolt"
)

type Denoms struct {
	Denom map[string]*Denom
	DB    *bolt.DB
}

type Denom struct {
	Name         string
	OwnerAccount *account.Account
	Nfts         map[string]*Nft
}

type Nft struct {
	Name         string
	DenomName    string
	Uri          string
	OwnerAccount *account.Account
}

func NewDenoms(db *bolt.DB) *Denoms {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(types.DenomBucketName))
		if b == nil {
			// 桶不存在则创建
			var err error
			b, err = tx.CreateBucket([]byte(types.DenomBucketName))
			if err != nil {
				panic(err)
			}
		}

		b.Put([]byte(types.CurrentUsers), []byte(types.SuperTypes))
		return nil
	})

	return &Denoms{nil, db}
}

func (denoms *Denoms) CreateDenom(name string, account *account.Account) {
	if account.AccountType == types.AdminTypes {
		denoms := &Denoms{}
		denoms.Denom = make(map[string]*Denom)
		denoms.Denom[name] = &Denom{name, account, nil}

		var denomName []byte
		denoms.DB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(types.DenomBucketName))
			denomName = b.Get([]byte(denoms.Denom[name].Name))

			if denomName == nil {
				fmt.Println("该denom名称已经存在！")
				return nil
			} else {
				d, _ := json.Marshal(denoms)
				b.Put([]byte("1"), d)
			}

			return nil
		})

		if denomName == nil {
			return
		} else {
			fmt.Println("创建denom成功！")
		}
	} else {
		fmt.Println("只有管理员才能创建denom！")
		return
	}
}

func (denoms *Denoms) MintNft(nftName string, denomName string, uri string) *Nft {
	if ok := denoms.Denom[denomName].Nfts[nftName]; ok != nil {
		fmt.Println("该denom下的nft名称已经存在！")
	}

	nft := &Nft{
		Name:         nftName,
		DenomName:    denomName,
		Uri:          uri,
		OwnerAccount: denoms.Denom[denomName].OwnerAccount,
	}

	denoms.Denom[denomName].Nfts[nft.Name] = nft

	denoms.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(types.DenomBucketName))
		d, _ := json.Marshal(denoms)
		b.Put([]byte(denomName), d)
		return nil
	})

	return nft
}

func (denoms *Denoms) TransferNft(nftName string, denomName string, account *account.Account) {
	if nftName == "" {
		fmt.Println("nft名称不能是空！")
		return
	}

	// 转移对象是不是管理员
	if account.AccountType != types.AdminTypes {
		fmt.Println("不能和非管理员交易")
		return
	}

	if account.Address == denoms.Denom[denomName].Nfts[nftName].OwnerAccount.Address {
		fmt.Println("不能给自己！")
		return
	} else {
		denoms.Denom[denomName].Nfts[denomName].OwnerAccount = account

		denoms.DB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(types.DenomBucketName))
			d, _ := json.Marshal(denoms)
			b.Put([]byte(denomName), d)
			return nil
		})

		fmt.Println("nft交易成功！")
	}
}

func (denoms *Denoms) Query(address string) {
	for denomName, denom := range denoms.Denom {
		if denom.OwnerAccount.Address == address {
			fmt.Printf("与您相关的denom信息如下：%+v\n", denoms.Denom[denomName])
		}
	}
}

func (denoms *Denoms) RmPms(address string) {
	// 去除denom
	for denomName, denom := range denoms.Denom {
		if denom.OwnerAccount.Address == address {
			delete(denoms.Denom, denomName)
		}
	}
	fmt.Println("去除denom成功！")
}

// func (nft *Nft) EditNftUri(uri string, account *account.Account) {
// 	if account.Address == nft.Owner {
// 		nft.Uri = uri

// 		fmt.Println("编辑nft成功！")
// 	} else {
// 		fmt.Println("该地址没有此nft，", account.Address)
// 	}
// }

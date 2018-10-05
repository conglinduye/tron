package service

import (
	"encoding/hex"
	"fmt"

	"github.com/wlcy/tron/explorer/lib/log"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/wlcy/tron/explorer/core/utils"
	"github.com/wlcy/tron/explorer/web/ext/entity"
)

//CreateAccount 只生成地址，不存数据库
func CreateAccount() (*entity.CreateAccount, error) {
	privKey, err := ethcrypto.GenerateKey()
	if nil != err {
		log.Errorf("gen privKey err for createAccount:[%v]", err)
		return nil, err
	}
	privKeyByte := ethcrypto.FromECDSA(privKey)
	hexPrivKey := hex.EncodeToString(privKeyByte)
	hexPubKey, err := utils.GetTronPublickey(hexPrivKey)
	if nil != err {
		log.Errorf("gen publicKey err for createAccount:[%v]", err)
		return nil, err
	}
	hexAddr, err := utils.GetTronHexAddress(hexPubKey)
	if nil != err {
		log.Errorf("gen hexAddress err for createAccount:[%v]", err)
		return nil, err
	}
	base58Addr := utils.Base58EncodeAddr(utils.HexDecode(hexAddr))
	fmt.Printf("CreateAccount done: hexPrivKey:[%v],hexAddr:[%v],base58Addr:[%v]", hexPrivKey, hexAddr, base58Addr)
	account := &entity.CreateAccount{}
	account.Address = base58Addr
	account.Key = hexPrivKey
	return account, nil
}

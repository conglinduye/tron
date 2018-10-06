package service

import (
	"encoding/hex"
	"fmt"

	"github.com/wlcy/tron/explorer/lib/log"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/wlcy/tron/explorer/core/utils"
	"github.com/wlcy/tron/explorer/web/ext/entity"
	"github.com/wlcy/tron/explorer/web/ext/module"
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

//QueryAccountBalance 获取账户余额信息
func QueryAccountBalance(address string) (*entity.AccountBalance, error) {
	var filterSQL string
	strSQL := fmt.Sprintf(`
	select account_name,acc.address,acc.balance as totalBalance,frozen,create_time,latest_operation_time,votes ,
		wit.url,wit.is_job,acc.allowance,acc.latest_withdraw_time,acc.is_witness,acc.net_usage,
		acc.free_net_used,acc.free_net_limit,acc.net_used,acc.net_limit,acc.asset_net_used,acc.asset_net_limit,
        ass.asset_name as token_name,ass.creator_address,ass.balance,asset.owner_address
    from tron_account acc
	left join account_asset_balance ass on ass.address=acc.address
	left join asset_issue asset on asset.asset_name=ass.asset_name
    left join witness wit on wit.address=acc.address
			where 1=1 `)

	//按传入条件拼接sql
	if address != "" {
		filterSQL = fmt.Sprintf(" and (acc.address='%v' or acc.account_name='%v')", address, address)
	}
	return module.QueryAccountRealize(strSQL, filterSQL, address)
}

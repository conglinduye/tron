package module

import (
	"encoding/json"

	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/lib/util"
	"github.com/wlcy/tron/explorer/web/ext/entity"
)

//QueryAccountRealize 操作数据库
func QueryAccountRealize(strSQL, filterSQL, address string) (*entity.AccountBalance, error) {
	strFullSQL := strSQL + " " + filterSQL
	log.Sql(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	if err != nil {
		log.Errorf("QueryAccountRealize error :[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryAccountRealize dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	var account = &entity.AccountBalance{}
	var oldBalance = make([]*entity.BalanceInfoDB, 0)
	var apiBalance = make([]*entity.FrozenBalance, 0)
	var frozenInfo = &entity.FrozenInfo{Total: 0, Balances: apiBalance}
	var balances = make([]*entity.Balance, 0)
	var totalFrozen = int64(0)
	var addressIn = ""
	var accountBalance = int64(0)

	accountTokenMap := make(map[string][]*entity.Balance, 0) //保存每个账户的token信息
	if dataPtr.ResNum() == 0 {
		account = &entity.AccountBalance{
			Allowance: 0,
			Entropy:   0,
			Balances:  balances,
			Frozen:    frozenInfo,
		}
	} else {
		//填充数据
		for dataPtr.NextT() {
			if addressIn == "" {
				addressIn = dataPtr.GetField("address")
				accountBalance = mysql.ConvertDBValueToInt64(dataPtr.GetField("totalBalance"))
				account.Allowance = mysql.ConvertDBValueToInt64(dataPtr.GetField("allowance"))
				account.Entropy = mysql.ConvertDBValueToInt64(dataPtr.GetField("net_usage"))
				frozen := dataPtr.GetField("frozen")
				if frozen != "" {
					if err := json.Unmarshal([]byte(frozen), &oldBalance); err != nil {
						log.Errorf("Unmarshal data failed:[%v]-[%v]", err, frozen)
					}
				}

				for _, blanceFrozen := range oldBalance {
					apiFrozen := &entity.FrozenBalance{}
					apiFrozen.Amount = blanceFrozen.Amount
					apiFrozen.Expires = blanceFrozen.Expires
					apiBalance = append(apiBalance, apiFrozen)
					totalFrozen += blanceFrozen.Amount
				}
				frozenInfo.Total = totalFrozen
				frozenInfo.Balances = apiBalance
				account.Frozen = frozenInfo
			}
			balance := &entity.Balance{}
			balance.Name = dataPtr.GetField("token_name")
			balance.Balance = mysql.ConvertDBValueToFloat64(dataPtr.GetField("balance"))

			if addressIn != "" {
				if tokenInfo, ok := accountTokenMap[addressIn]; ok {
					if balance.Name != "" {
						tokenInfo = append(tokenInfo, balance)
						accountTokenMap[addressIn] = tokenInfo
					}
				} else {
					tokenArr := make([]*entity.Balance, 0)
					if accountBalance > 0 { //如果账户有trx余额，那也计入token信息
						ownbalance := &entity.Balance{}
						ownbalance.Name = "TRX"
						ownbalance.Balance = float64(accountBalance) / 1000000 //单位换算，页面按照TRX显示
						tokenArr = append(tokenArr, ownbalance)
					}
					if balance.Balance > 0 {
						tokenArr = append(tokenArr, balance)
					}

					accountTokenMap[addressIn] = tokenArr
				}
			}
		}

		//拼接tokeninfo列表
		if tokenInfo, ok := accountTokenMap[addressIn]; ok {
			account.Balances = tokenInfo
		}
	}
	return account, nil

}

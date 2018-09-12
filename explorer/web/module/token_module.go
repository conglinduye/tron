package module

import (
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/lib/util"
	"encoding/json"
)

//QueryTokens
func QueryTokensRealize(strSQL, filterSQL, sortSQL, pageSQL string) (*entity.TokenResp, error){
	strFullSQL := strSQL + " " + filterSQL + " " + sortSQL + " " + pageSQL
	log.Debug(strFullSQL)
	dataPtr, err := mysql.QueryTableData(strFullSQL)
	if err != nil {
		log.Errorf("QueryTokens error:[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	if dataPtr == nil {
		log.Errorf("QueryTokens dataPtr is nil ")
		return nil, util.NewErrorMsg(util.Error_common_internal_error)
	}
	tokenResp := &entity.TokenResp{}
	tokens := make([]*entity.TokenInfo, 0)

	for dataPtr.NextT() {
		token := &entity.TokenInfo{}
		token.OwnerAddress = dataPtr.GetField("owner_address")
		token.Name = dataPtr.GetField("asset_name")
		token.TotalSupply = mysql.ConvertDBValueToInt64(dataPtr.GetField("total_supply"))
		token.TrxNum = mysql.ConvertDBValueToInt64(dataPtr.GetField("trx_num"))
		token.Num = mysql.ConvertDBValueToInt64(dataPtr.GetField("num"))
		token.EndTime = mysql.ConvertDBValueToInt64(dataPtr.GetField("end_time"))
		token.StartTime = mysql.ConvertDBValueToInt64(dataPtr.GetField("start_time"))
		token.VoteScore = mysql.ConvertDBValueToInt64(dataPtr.GetField("vote_score"))
		token.Description = dataPtr.GetField("asset_desc")
		token.Url = dataPtr.GetField("utl")
		frozenJson := dataPtr.GetField("frozen_supply")
		var tokenFrozen []entity.TokenFrozen
		if err := json.Unmarshal([]byte(frozenJson), &tokenFrozen); err == nil {
			token.Frozen = tokenFrozen
		}
		token.Abbr = dataPtr.GetField("asset_abbr")

		tokens = append(tokens, token)
	}

	var total = int64(len(tokens))
	total, err = mysql.QuerySQLViewCount(strSQL)
	if err != nil {
		log.Errorf("query view count error:[%v], SQL:[%v]", err, strSQL)
	}
	tokenResp.Total = total
	tokenResp.Data = tokens

	return tokenResp, nil
}

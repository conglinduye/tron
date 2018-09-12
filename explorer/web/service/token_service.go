package service

import (
	"github.com/wlcy/tron/explorer/web/entity"
	"fmt"
	"github.com/wlcy/tron/explorer/web/module"
)

//QueryTokens
func QueryTokens(req *entity.Token) (*entity.TokenResp, error) {
	var filterSQL, sortSQL, pageSQL string

	strSQL := fmt.Sprintf(`
			select owner_address, asset_name, asset_abbr, total_supply, frozen_supply,
			trx_num, num, start_time, end_time, order_num, vote_score, asset_desc, url
			from asset_issue
			where 1=1 `)

	if req.Owner != "" {
		filterSQL = fmt.Sprintf(" and owner_address=%v", req.Owner)
	}
	if req.Name != "" {
		filterSQL = fmt.Sprintf(" and asset_name=%v", req.Name)

	}

	tokenResp, err := module.QueryTokensRealize(strSQL, filterSQL, sortSQL, pageSQL)
	if err != nil {
		return tokenResp, err
	}

	calculate(tokenResp)

	return tokenResp, nil
}

// calculate
func calculate(tokenResp *entity.TokenResp) {
	tokens := tokenResp.Data
	for index1 := range tokens {
		token := tokens[index1]
		var tokenFrozens []entity.TokenFrozen

		var frozenSupply int64 = 0
		for index2 := range tokenFrozens {
			frozenSupply = frozenSupply + tokenFrozens[index2].Amount
		}
		totalSupply := token.TotalSupply
		availableSupply := totalSupply - frozenSupply

		var availableTokens int64 = 0
		//to do

		var totalParticipateAmount int64 = 0
		// to do
		issuedTokens := availableSupply - availableTokens

		issuedPercentage := float64(issuedTokens) / float64(totalSupply) * 100
		remainingTokens := totalSupply - frozenSupply - issuedTokens
		percentage := float64(remainingTokens) / float64(totalSupply) * 100
		frozenSupplyPercentage := float64(frozenSupply) / float64(totalSupply) * 100

		price := token.TrxNum / token.Num

		token.Price = price
		token.Issued = issuedTokens
		token.Participated = totalParticipateAmount
		token.IssuedPercentage = issuedPercentage
		token.Available = availableTokens
		token.AvailableSupply = availableSupply
		token.Remaining = remainingTokens
		token.RemainingPercentage = percentage
		token.Percentage = percentage
		token.FrozenTotal = frozenSupply
		token.FrozenPercentage = frozenSupplyPercentage

	}

}

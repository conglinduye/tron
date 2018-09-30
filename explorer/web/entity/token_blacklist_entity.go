package entity


type AssetBlacklistReq struct {
	Start  					string 				`json:"start"`  				// start
	Limit  					string 				`json:"limit"`  				// limit
	OwnerAddress 			string 				`json:"ownerAddress"` 			// ownerAddress
	TokenName 				string 				`json:"tokenName"`				// tokenName
}


type AssetBlacklistResp struct {
	Total 					int64				`json:"total"`					// total
	Data					[]*AssetBlacklist	`json:"data"`					// data
}

type AssetBlacklist struct {
	Id 						int64				`json:"id"`						// id
	OwnerAddress     		string 				`json:"ownerAddress"`    		// ownerAddress
	TokenName				string 				`json:"assetName"`				// tokenName
	CreateTime				string 				`json:"createTime"`				// createTime
}

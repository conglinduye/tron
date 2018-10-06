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
	Id 						string				`json:"id"`						// id
	OwnerAddress     		string 				`json:"ownerAddress"`    		// ownerAddress
	TokenName				string 				`json:"assetName"`				// tokenName
	CreateTime				string 				`json:"createTime"`				// createTime
}


type AssetExtInfo struct {
	Id 						string				`json:"id"`
	Address					string				`json:"address"`
	TokenName				string 				`json:"tokenName"`
	TokenId					string				`json:"tokenId"`
	Brief					string 				`json:"brief"`
	Website					string 				`json:"website"`
	WhitePaper				string 				`json:"whitePaper"`
	Github					string 				`json:"github"`
	Country					string 				`json:"country"`
	Credit					string				`json:"credit"`
	Reddit					string 				`json:"reddit"`
	Twitter					string 				`json:"twitter"`
	Facebook				string				`json:"facebook"`
	Telegram				string 				`json:"telegram"`
	Steam					string 				`json:"steam"`
	Medium					string 				`json:"medium"`
	Webchat					string 				`json:"webchat"`
	Weibo					string 				`json:"weibo"`
	Review					string 				`json:"review"`
	Status					string 				`json:"status"`
	UpdateTime				string				`json:"updateTime"`
}

type AssetExtLogo struct {
	Id 						string				`json:"id"`
	Address					string 				`json:"address"`
	LogoUrl					string 				`json:"logoUrl"`
	UpdateTime				string				`json:"updateTime"`
}
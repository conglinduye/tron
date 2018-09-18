package entity

//WitnessInfo 查询超级代表的结果
type WitnessInfo struct {
	Address           	string 		`json:"address"`           //:"TDo2qwRLEkTZCXzcotiEQLt2wLeTttKTnZ",//地址
	Name             	string 		`json:"name"`              //:"TronSchool",//名称
	URL               	string 		`json:"url"`               //:"http://www.tron.school/",//url
	Producer          	bool   		`json:"producer"`          //:false,//是否出块
	LatestBlockNumber 	int64  		`json:"latestBlockNumber"` //:0,//最近快高
	LatestSlotNumber  	int64  		`json:"latestSlotNumber"`  //:0,//岁进slot数
	MissedTotal       	int64  		`json:"missedTotal"`       //:0,//丢块总数
	ProducedTotal     	int64  		`json:"producedTotal"`     //:0,//产出块总数
	ProducedTrx       	int64  		`json:"producedTrx"`       //:0,//获得trx
	Votes             	int64  		`json:"votes"`             //:362616//得票数
	ProducePercentage 	float64		`json:"producePercentage"` //出块效率
	VotesPercentage		float64		`json:"votesPercentage"`	//票数占比
}

//WitnessStatisticInfo ...
type WitnessStatisticInfo struct {
	Address       string  `json:"address"`       //:"TGzz8gjYiYRqpfmDwnLxfgPuLVNmpCswVp",//地址
	Name          string  `json:"name"`          //:"Sesameseed",//名称
	URL           string  `json:"url"`           //:"https://www.sesames‘eed.org",//url
	BlockProduced int64   `json:"blockProduced"` //:309,//当前一轮各超级代表的出块
	Total         int64   `json:"total"`         //:8130,//当前一轮各超级代表的出块总和
	Percentage    float64 `json:"percentage"`    //:0.03800738007380074// blockProduced/total
}

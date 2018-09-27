package service

import (
	"testing"

	"github.com/wlcy/tron/explorer/core/utils"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/web/entity"
)

func TestTransactions(t *testing.T) {
	Init()
	req := &entity.Transactions{}
	req.Sort = "-number"
	req.Limit = 5
	req.Start = 0

	//req.Number = "2287351"

	resp, err := QueryTransactions(req)
	if err != nil {
		log.Error(err)
	}
	ss, _ := mysql.JSONObjectToString(resp)
	log.Printf("total:%v", ss)
	/*for _, value := range resp.Data {
		log.Printf("data:%#v", value)
	}*/

}

func TestTransaction(t *testing.T) {
	Init()
	req := &entity.Transactions{}
	/*req.Sort = "-number"
	req.Limit = "5"
	req.Start = "0"
	*/
	req.Hash = "086cd2282f698c0f72b6eb4b3eb880c2eb4a2bd8249c6ae644dc82f52b82490a"

	resp, err := QueryTransaction(req)
	if err != nil {
		log.Error(err)
	}
	ss, _ := mysql.JSONObjectToString(resp)
	log.Printf("total:%v", ss)

}

func TestA(t *testing.T) {
	ss := "0a0a48756f6269546f6b656e121541b7a3dd3b45f5a30cb108b90cc12cee3a70ca4e861a1541c1b94b6cf7b946db06de3253ecabeb9b01e2b1f42001"
	ss = "0a1541e2c45340b37b97e313c91217e340765e4579bcc3121a0a1541beab998551416b02f6721129bb01b51fceceba08108d79"
	addr, contract := utils.GetContractInfoStr2(2, utils.HexDecode(ss))

	log.Printf("total:%v-%v", addr, contract)
}

func TestPostTrans(t *testing.T) {
	req := &entity.PostTransaction{
		//Transaction: "0A84010A025006220880DDBCA411E6159840E8F7E1BDDC2C5204484148415A67080112630A2D747970652E676F6F676C65617069732E636F6D2F70726F746F636F6C2E5472616E73666572436F6E747261637412320A1541E552F6487585C2B58BC2C9BB4492BC1F17132CD012154190919CBA90CE96F9B9A63AFDE5AC66453D3F690E18C0843D124183A239CD8B1A3998B56DF45667E230B6BED13C10889BED456F40716DC5558F58715DD1E31DDF57419C6FD5F4864C5BD8995E20306C44609FA3C4ED556DA6DCE600",
		Transaction: "0a91010a024a4f2208c5c74b0d009b98ba40f0a2c5d3e12c5a730802126f0a32747970652e676f6f676c65617069732e636f6d2f70726f746f636f6c2e5472616e736665724173736574436f6e747261637412390a074e4241436f696e121541096e8fd0a9f7848fda8ee6066e9df2035c4b13f81a1541e552f6487585c2b58bc2c9bb4492bc1f17132cd0200170bedfc2d3e12c1241dee4438d2a710cd4893f986d2f6e8c8417d13df24e534bdb09abb5d6ccf9d07a0a20010cc49bd2376264ed6f0bda02556b59d4ed9b1bb1f201867f5e54c186c100",
	}

	resp, _ := PostTransaction(req, "1")
	ss, _ := mysql.JSONObjectToString(resp)
	log.Printf("total:%v", ss)
	//fmt.Printf("result:[%v],err:[%v]", result, err)
}

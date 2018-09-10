package grpcclient

import (
	"fmt"
	"testing"

	"github.com/wlcy/tron/explorer/core/utils"
)

func TestWalletSolidity(*testing.T) {

	client := NewWalletSolidity(fmt.Sprintf("%s:50051", utils.GetRandSolidityNodeAddr()))

	err := client.Connect()
	if nil != err {
		fmt.Println(err)
		return
	}

	fmt.Println(client.GetState(), client.Target())

	// account, _ := client.GetAccount("TGo9Me13BSagSHXmKZDbZrLaFW9PXYYs3T")
	// utils.VerifyCall(client.GetAssetIssueList())

	// fmt.Printf("%s\n%s\n%s\n", account.AssetIssuedName, utils.Base64Encode(account.AssetIssuedName), utils.Base64Decode("VE9O"))
	// fmt.Println(utils.Base64Encode(utils.Base58DecodeAddr("TGo9Me13BSagSHXmKZDbZrLaFW9PXYYs3T")))

	block, _ := client.GetNowBlock()

	blockHash := utils.CalcBlockHash(block)
	blockNum := block.BlockHeader.RawData.Number
	blockNumByte := utils.BinaryBigEndianEncodeInt64(blockNum)

	fmt.Printf("%v\n%v\t%v\n%v\n", blockNum, blockNumByte, blockNumByte[6:8], utils.BinaryBigEndianDecodeUint64(blockNumByte))

	fmt.Printf("%v\nblock hash:[%v]\n%v\n", blockNum, utils.HexEncode(blockHash), blockHash[8:16])

	for _, tran := range block.Transactions {
		trxHash := utils.HexEncode(utils.CalcTransactionHash(tran))
		fmt.Printf("transaction hash:[%v](%v)[%v]-->%#v\n\n", trxHash, tran.RawData.RefBlockHash, tran.RawData.RefBlockBytes, tran.RawData)
	}

	return

	// addr := "TGzz8gjYiYRqpfmDwnLxfgPuLVNmpCswVp" // super witness
	// account, err := client.GetAccount(addr)
	// utils.VerifyCall(account, err)
	// fmt.Printf("account_name:[%s], account_id:[%s]\n", account.AccountName, account.AccountId)

	// list1, err := client.GetPaginatedAssetIssueList(1, 3)
	// utils.VerifyCall(list1, err)
	// list2, err := client.GetPaginatedAssetIssueList(2, 3)
	// utils.VerifyCall(list2, err)
	// list3, err := client.GetPaginatedAssetIssueList(1, 6)
	// utils.VerifyCall(list3, err)

	// block, err := client.GetNowBlock()
	// utils.VerifyCall(block, err)
	// blockID := utils.CalcBlockID(block)
	// fmt.Printf("blockID:%v\n[%v]\n", blockID, utils.HexEncode(blockID))

	// block, err = client.GetBlockByNum(2099805)
	// utils.VerifyCall(block, err)
	// blockID = utils.CalcBlockID(block)
	// fmt.Printf("blockID:%v\n[%v]\n", blockID, utils.HexEncode(blockID))

	// for _, tran := range block.Transactions {
	// 	utils.VerifyCall(tran, nil)
	// 	tranHash := utils.CalcTransactionHash(tran)
	// 	fmt.Printf("transaction hash:%v\n[%v]\n", tranHash, utils.HexEncode(tranHash))
	// }

	// fmt.Printf("\n\n\n")
	// trans, err := client.GetTransactionByID("2455e8a5263b17e9dbac4c08a3c969d5263367628cace44e5df5facb0f93f387")
	// utils.VerifyCall(trans, err)
	// tranHash := utils.CalcTransactionHash(trans)
	// fmt.Printf("transaction hash:%v\n[%v]\n", tranHash, utils.HexEncode(tranHash))

	// transInfo, err := client.GetTransactionInfoByID("2455e8a5263b17e9dbac4c08a3c969d5263367628cace44e5df5facb0f93f387")
	// utils.VerifyCall(transInfo, err)

	// utils.VerifyCall(client.GetTransactionCountByBlockNum(2099805))

	// utils.VerifyCall(client.GetBlockByNum2(2099805))
	// utils.VerifyCall(client.ListWitnesses())
	// utils.VerifyCall(client.GetTransactionByID("13935662a6c71bab0b86bec47c906a7a5aa9492cdf516fdb0ee6925c47fa483a"))

	// dd := `ChVBPvOMnl/UK3kfykJ6e3Ls56BGQsESFUGowCjVhWt/V1+94U2H9mffVCJ7/hoEREFDQyCAg6y4Eg=="}}],"timestamp":1532690261908},"signature":["DeLG7puiOuTnPNaHjBy8kNC8zJEl081Wuv1qTnFieGM78WqzRM/6mfFtSWGqWNjIOz2gljGct6gFKNxclDLFPQA`
	// pp := &core.ParticipateAssetIssueContract{}

	// rawData := utils.Base64Decode(dd)
	// err = proto.Unmarshal(rawData, pp)
	// fmt.Printf("%s\n%v\n", proto.MarshalTextString(pp), err)

}

func TestRW(*testing.T) {
	client := GetRandomSolidity()

	account, _ := client.GetAccount("TKoU7MkprWw8q142Sd199XU4B5fUaMVBNm")
	fmt.Printf("%#v\n", account)

	// "QWvZtdNgH9MCNjy7j27oD+YL4lt6"
}

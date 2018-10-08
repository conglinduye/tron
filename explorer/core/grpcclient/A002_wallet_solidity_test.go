package grpcclient

import (
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/tronprotocol/grpc-gateway/core"
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

	fmt.Printf("blockNum:%v\nblockNum bytes:%v\tblockNum byte[6~8]%v\nblockNum again:%v\n", blockNum, blockNumByte, blockNumByte[6:], utils.BinaryBigEndianDecodeUint64(blockNumByte))

	fmt.Printf("blockNum:%v\nblock hash:[%v]\nblock hash ref[8~16]%v\n", blockNum, utils.HexEncode(blockHash), blockHash[8:16])

	for _, tran := range block.Transactions {
		trxHash := utils.HexEncode(utils.CalcTransactionHash(tran))
		fmt.Printf("transaction hash:[%v](refBlockByte:%v)[refBlockHash:%v]-->%#v\n\n", trxHash, tran.RawData.RefBlockHash, tran.RawData.RefBlockBytes, tran.RawData)
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
	client1 := GetRandomWallet()

	addr := "TDGmmTC7xDgQGwH4FYRGuE7SFH2MePHYeH"

	account, _ := client.GetAccount(addr)
	// fmt.Printf("%#v\n", account)
	utils.VerifyCall(account, nil)

	// "QWvZtdNgH9MCNjy7j27oD+YL4lt6"
	utils.VerifyCall(client1.GetAccountNet(addr))
}

func TestBlock(*testing.T) {

	// client := GetRandomSolidity()

	client1 := GetRandomWallet()

	// utils.VerifyCall(client.GetBlockByNum(2224654))

	ttt, err := client1.GetNextMaintenanceTime()

	fmt.Printf("%v\n%v\n%v\n", err, ttt, utils.ConverTimestampStr(ttt))
	// wg := sync.WaitGroup{}

	// wg.Add(1)
	// go func() {
	// 	block, _ := client.GetNowBlock()
	// 	fmt.Printf("solidity now block:%v\n", block.BlockHeader.RawData.Number)
	// 	wg.Done()
	// }()
	// block1, _ := client1.GetNowBlock()
	// fmt.Printf("full now block:%v\n", block1.BlockHeader.RawData.Number)
	// utils.VerifyCall(client.GetBlockByNum(block1.BlockHeader.RawData.Number))
	// wg.Wait()
}

func TestGetBlock(*testing.T) {
	var num int64 = 2270833 // 2270833 2271567

	getBlockSolidity(num)
	getBlockFull(num)
}

func getBlockSolidity(num int64) {
	client := GetRandomSolidity()

	block, err := client.GetBlockByNum(num)

	if nil != err || nil == block {
		fmt.Println(err)
		return
	}

	data, _ := proto.Marshal(block)
	fmt.Printf("solidity block:[%v], hash:%v, size:%v\n", num, utils.HexEncode(utils.CalcBlockHash(block)), len(data))
	showBlockTrx(block)
	fmt.Printf("\n\n")
}

func getBlockFull(num int64) {
	client := GetRandomWallet()

	// block, err := client.GetBlockByNum(num)
	blocks, err := client.GetBlockByLimitNext(num, num+1)

	if nil != err || nil == blocks {
		fmt.Println(err)
		return
	}

	data, _ := proto.Marshal(blocks[0])
	fmt.Printf("fullnode block:[%v], hash:%v, size:%v\n", num, utils.HexEncode(utils.CalcBlockHash(blocks[0])), len(data))
	showBlockTrx(blocks[0])
	fmt.Printf("\n\n")
}

func showBlockTrx(block *core.Block) {
	if nil == block {
		return
	}
	for _, trx := range block.Transactions {
		ctxOwner, _ := utils.GetContractInfoStr(trx.RawData.Contract[0])
		fmt.Printf("trx_hash:%64v\ttype:%-30v\towner_address:%v\ttimestamp:%30v\texpire:%30v\n",
			utils.HexEncode(utils.CalcTransactionHash(trx)), trx.RawData.Contract[0].Type, ctxOwner,
			utils.ConverTimestampStr(trx.RawData.Timestamp), utils.ConverTimestampStr(trx.RawData.Expiration))
	}
}

/*
mysql> select * from blocks where block_id = 2271567\G
*************************** 1. row ***************************
       block_id: 2271567
     block_hash: 000000000022a94f5868c99d3c430559876adc5b028dda1c78c9b694a2d93ed6
    parent_hash: 000000000022a94efee06aeb9f3708740690f7fa8ee5465bf7209412b107fc19
witness_address: TSzoLaVCdSNDpNxgChcFt9rSRF5wWAZiR4
   tx_trie_hash: 058cf1f0c18d992538789c9e156c3cd21d01fbcb776c35230d55cb086a6899e5
     block_size: 1928
transaction_num: 9
      confirmed: 1
    create_time: 1536722433000
  modified_time: 2018-09-12 03:21:42.449452

*/

func TestGetAccount(*testing.T) {
	// utils.TestNet = true

	sw := GetRandomSolidity()
	addr := "TCaCaa6DhXkaRXu4T2BwCfu14WPyBwx2Ef" // smart contract addr
	addr = "TAHAyhFkRz37o6mYepBtikrrnCNETEFtW5"  // common addr
	addr = "TNTYHpa71hW9y2i55FovHYwt2MmPdD76nt"
	addr = "TYmns3F1Je57RLCC8rGSRHCuhM4XKEseEU"
	addr = "TA2XgvKJXHvoFC25TTfTVmfsYp5sriDrY3"
	addr = "TQw7CVALENBqCP3uZuo1KxiTnvCqGnoeod"
	addr = "TM5jRzWJqCSRDgSpdqS15Ejxrg2oA568US"
	account, err := sw.GetAccount(addr)
	fmt.Printf("%v\n%v\n", utils.ToJSONStr(account), err)
}
func TestGetAccountNet(*testing.T) {
	sw := GetRandomWallet()

	addr := "TM5jRzWJqCSRDgSpdqS15Ejxrg2oA568US"
	account, err := sw.GetAccountNet(addr)
	fmt.Printf("%v\n%v\n", utils.ToJSONStr(account), err)

	account1, err1 := sw.GetAccount(addr)
	fmt.Printf("%v\n%v\n", utils.ToJSONStr(account1), err1)
}

func TestGetAssetIssue(*testing.T) {
	sw := GetRandomWallet()

	ctx, err := sw.GetAssetIssueByName("ZZZZZB")

	fmt.Printf("%v\n%v\n", utils.ToJSONStr(ctx), err)
}

// 1536737988000
// 1538033451000

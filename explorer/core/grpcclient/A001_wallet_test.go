package grpcclient

import (
	"fmt"
	"testing"
	"time"

	"github.com/tronprotocol/grpc-gateway/core"
	"github.com/wlcy/tron/explorer/core/utils"
)

func TestWallet(*testing.T) {

	client := NewWallet(fmt.Sprintf("%s:50051", utils.GetRandFullNodeAddr()))

	err := client.Connect()
	if nil != err {
		fmt.Println(err)
		return
	}

	fmt.Println(client.GetState(), client.Target())

	// addr := "TGzz8gjYiYRqpfmDwnLxfgPuLVNmpCswVp" // super witness
	// addr = "TGo9Me13BSagSHXmKZDbZrLaFW9PXYYs3T"
	// utils.VerifyCall(client.TotalTransaction())

	// fmt.Println(utils.HexEncode(utils.Base64Decode("pjuwcFDis+Y=")))

	block, _ := client.GetNowBlock()

	blockHash := utils.CalcBlockHash(block)
	blockNum := block.BlockHeader.RawData.Number
	blockNumByte := utils.BinaryBigEndianEncodeInt64(blockNum)

	fmt.Printf("%v\n%v\t%v\n%v\n", blockNum, blockNumByte, blockNumByte[6:8], utils.BinaryBigEndianDecodeUint64(blockNumByte))

	fmt.Printf("%v\nblock hash:[%v]\n%v\n", blockNum, utils.HexEncode(blockHash), blockHash[8:16])

	for _, tran := range block.Transactions {
		trxHash := utils.HexEncode(utils.CalcTransactionHash(tran))
		fmt.Printf("transaction hash:[%v](%v)[%v]-->%#v\n\n", trxHash, tran.RawData.RefBlockHash, tran.RawData.RefBlockBytes, tran.RawData)
		utils.VerifyCall(tran, nil)
		fmt.Printf("time:%v--%v\n", time.Unix(0, tran.RawData.Timestamp*1000000), time.Unix(0, tran.RawData.Expiration*1000000))
	}

	return
	// utils.VerifyCall(client.GetAssetIssueByAccount(addr))

	// utils.VerifyCall(client.GetAssetIssueByName("ZYON"))
	// utils.VerifyCall(client.GetBlockByID("0000000000200f1f71ef0935b64d019298c83a06603f9cc0b603118bfe58ba51"))
	// utils.VerifyCall(client.GetBlockByLimitNext(0, 3))
	// utils.VerifyCall(client.DeployContract("TGzz8gjYiYRqpfmDwnLxfgPuLVNmpCswVp", nil))
	// utils.VerifyCall(client.TriggerContract("TGzz8gjYiYRqpfmDwnLxfgPuLVNmpCswVp", "TGzz8gjYiYRqpfmDwnLxfgPuLVNmpCswVp", 123, nil))
	// utils.VerifyCall(client.ListWitnesses())
	// utils.VerifyCall(client.ListProposals())
	// utils.VerifyCall(client.GetPaginatedAssetIssueList())
}

func TestJ(*testing.T) {
	var w *Wallet

	fmt.Printf("[%v]\n", utils.ToJSONStr(w))
	fmt.Printf("[%v]\n", utils.ToJSONStr(int64(0)))

}

func TestGenTrx(*testing.T) {
	trx := new(core.Transaction)

	trx.RawData = new(core.TransactionRaw)

	tran := new(core.TransferContract)
	tran.OwnerAddress = utils.Base58DecodeAddr("TMLQYUXX6R3tMyEG2CNkjWigLgLxNQ5Uj2")
	tran.ToAddress = utils.Base58DecodeAddr("TDPgbSpKrLnaBMF79QUg3aigsG1tsWoxLJ")
	tran.Amount = 1230 // sun
	trx.RawData.Timestamp = 888888888
	trx.RawData.Expiration = time.Now().Unix() * 1000
}

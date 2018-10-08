package grpcclient

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
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

	tran := new(core.TransferContract)
	tran.OwnerAddress = utils.Base58DecodeAddr("TMLQYUXX6R3tMyEG2CNkjWigLgLxNQ5Uj2")
	tran.ToAddress = utils.Base58DecodeAddr("TDPgbSpKrLnaBMF79QUg3aigsG1tsWoxLJ")
	tran.Amount = 1230 // sun

	trx, err := utils.BuildTransaction(core.Transaction_Contract_TransferContract, tran, []byte("test"))
	if nil != err {
		fmt.Printf("build trx failed:%v\n", err)
		return
	}
	client := GetRandomWallet()
	block, err := client.GetNowBlock()
	if nil != err {
		fmt.Printf("getBlock failed:%v\n", err)
		return
	}

	blockHash := utils.CalcBlockHash(block)
	trx.RawData.RefBlockHash = blockHash[8:16]
	trx.RawData.RefBlockBytes = utils.BinaryBigEndianEncodeInt64(block.BlockHeader.RawData.Number)[6:8]
	trx.RawData.Timestamp = time.Now().UTC().UnixNano() / 1000000
	trx.RawData.Expiration = time.Now().UTC().Add(5*time.Minute).UnixNano() / 1000000

	sign, err := utils.SignTransaction(trx, "3D7237EC6B1C7532453B76A5E95EBFD1BA6E739E7A1C9B3FF4F60492A90F70A3")
	if nil != err {
		fmt.Printf("sign failed:%v\n", err)
		return
	}

	trx.Signature = append(trx.Signature, sign)

	ret, err := client.BroadcastTransaction(trx)
	if nil != ret {
		fmt.Printf("%v\n%s\n%v\n", err, ret.Message, ret.Code)
	}

}

func TestExtTrx(*testing.T) {
	a := `0A89010A02361622081222F564E33ACC0840C8D492A1E52C520B626C62626C626C626C626C5A65080112610A2D747970652E676F6F676C65617069732E636F6D2F70726F746F636F6C2E5472616E73666572436F6E747261637412300A15417CABC0D927D00FC6DCA6587E3F645FBC86CF638D121541A0D1341AFB4406F0135724925630789614606D85186412418D1D185A85150B27DFF058DBEC91D7098359F65404AD0ECC3DAEEDE08B610AE8E8B0611E00D0D0DA496D65A6827EFB2C173CC2FAA5F14A2CDF64A6088E9CE5A601`

	trx := new(core.Transaction)
	err := proto.Unmarshal(utils.HexDecode(a), trx)

	if nil != err {
		fmt.Printf("err:%v\n", err)
		return
	}

	fmt.Printf("trx:%v\n", utils.ToJSONStr(trx))
	ctxType, ctx, err := utils.GetTransactionContract(trx)
	ctxObj, err := utils.GetContractInfoObj(trx.RawData.Contract[0])
	fmt.Printf("ctxType:%v\n%v\n%v\n", ctxType.String(), ctx, ctxObj)

}

package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang/protobuf/ptypes"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/tronprotocol/grpc-gateway/core"
	"github.com/wlcy/tron/explorer/core/grpcclient"
	"github.com/wlcy/tron/explorer/core/utils"
	"github.com/wlcy/tron/explorer/lib/log"
)

// TBTransfer transfer contract
func TBTransfer(ctx *gin.Context) {
	ctxReq, err := getContextBody(ctx)
	if nil != err || nil == ctxReq {
		ctx.JSON(http.StatusOK, newTBResponse(err))
		return
	}

	if ctxReq.Broadcast {

	}
}

// TBTransferAsset transfer asset contract
func TBTransferAsset(ctx *gin.Context) {
}

// TBAccountCreate account create contract
func TBAccountCreate(ctx *gin.Context) {

}

// TBAccountUpdate account update contract
func TBAccountUpdate(ctx *gin.Context) {

}

// TBWithdrawBalance withdraw balance contract
func TBWithdrawBalance(ctx *gin.Context) {

}

func getContextBody(ctx *gin.Context) (*TBRequestType, error) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if nil != err {
		return nil, err
	}

	tbReq := new(TBRequestType)

	err = json.Unmarshal(body, tbReq)
	if nil != err {
		return nil, err
	}

	err = tbReq.ConvertContract()

	if nil != err {
		return nil, err
	}

	if tbReq.Broadcast {
		err = signAndBroadcastContrat(tbReq.ContractType, tbReq.RealContract, tbReq.Data, tbReq.Key)
	}

	return tbReq, err
}

func signAndBroadcastContract(ctxType core.Transaction_Contract_ContractType, contract interface{}, data string, privatekey string) (err error) {
	if 0 == len(privatekey) {
		log.Errorf("broadcast need private key for signature:%v:%#v", ctxType, contract)
		return errInvalidRequest
	}

	trx := new(core.Transaction)
	trx.RawData = new(core.TransactionRaw)

	// fill contract data
	if 0 < len(data) {
		trx.RawData.Data = []byte(data)
	}

	// fill contract detail
	contractRaw := new(core.Transaction_Contract)
	contractRaw.Type = ctxType
	pbMsg, ok := contract.(proto.Message)
	if ok {
		contractRaw.Parameter, err = ptypes.MarshalAny(pbMsg)
	} else {
		return errInvalidRequest
	}
	trx.RawData.Contract = append(trx.RawData.Contract, contractRaw)

	// set transaction timestamp, in millisecon
	trx.RawData.Timestamp = time.Now().UTC().UnixNano() / int64(time.Millisecond)

	// sign transaction
	sign, err := utils.SignTransaction(trx, privatekey)
	if nil != err {
		return errInvalidRequest
	}
	trx.Signature = append(trx.Signature, sign)

	// broadcast
	client := grpcclient.GetRandomWallet()
	resp, err := client.BroadcastTransaction(trx)

	return
}

func newTBResponse(err error) *TBResponseType {
	resp := new(TBResponseType)

	resp.Success = false
	resp.Result = TBResultType{}
	resp.Result.Code = "Failed"
	resp.Result.Message = err.Error()

	return resp
}

// TBRequestType ...
type TBRequestType struct {
	Contract  map[string]interface{} `json:"contract"`
	Key       string                 `json:"key"`
	Broadcast bool                   `json:"broadcast"`
	Data      string                 `json:"data"`

	RealContract interface{}                            `json:"-"`
	ContractType core.Transaction_Contract_ContractType `json:"-"`
}

var (
	errInvalidRequest = fmt.Errorf("Invalid request body")
)

// ConvertContract ...
func (tbR *TBRequestType) ConvertContract() (err error) {
	err = errInvalidRequest
	if nil == tbR.Contract {
		return
	}

	ctxMap := tbR.Contract

	ownerAddress, ownerOK := ctxMap["ownerAddress"]
	toAddress, toOK := ctxMap["toAddress"]
	amount, amountOK := ctxMap["amount"]
	assetName, assetOK := ctxMap["assetName"]
	accountAddress, accountAddrOK := ctxMap["accountAddress"]
	accountName, accountNameOK := ctxMap["accountAddress"]

	defer func() { // interface convert failed exception
		if panErr := recover(); nil != panErr {
			log.Errorf("parse transfer builder parameter failed:%v", panErr)
		}
	}()

	if ownerOK && toOK && amountOK && assetOK { // transferAsset
		realCtx := new(core.TransferAssetContract)
		realCtx.OwnerAddress = utils.Base58DecodeAddr(ownerAddress.(string))
		realCtx.ToAddress = utils.Base58DecodeAddr(toAddress.(string))
		realCtx.AssetName = []byte(assetName.(string))
		realCtx.Amount = amount.(int64)

		tbR.RealContract = realCtx
		tbR.ContractType = core.Transaction_Contract_TransferAssetContract

	} else if ownerOK && toOK && amountOK { // transfer
		realCtx := new(core.TransferContract)
		realCtx.OwnerAddress = utils.Base58DecodeAddr(ownerAddress.(string))
		realCtx.ToAddress = utils.Base58DecodeAddr(toAddress.(string))
		realCtx.Amount = amount.(int64)

		tbR.RealContract = realCtx
		tbR.ContractType = core.Transaction_Contract_TransferContract

	} else if ownerOK && accountAddrOK { // accountCreate
		realCtx := new(core.AccountCreateContract)
		realCtx.OwnerAddress = utils.Base58DecodeAddr(ownerAddress.(string))
		realCtx.AccountAddress = utils.Base58DecodeAddr(accountAddress.(string))

		tbR.RealContract = realCtx
		tbR.ContractType = core.Transaction_Contract_AccountCreateContract

	} else if ownerOK && accountNameOK { // accountUpdate
		realCtx := new(core.AccountUpdateContract)
		realCtx.OwnerAddress = utils.Base58DecodeAddr(ownerAddress.(string))
		realCtx.AccountName = []byte(accountName.(string))

		tbR.RealContract = realCtx
		tbR.ContractType = core.Transaction_Contract_AccountUpdateContract

	} else if ownerOK { // balance
		realCtx := new(core.WithdrawBalanceContract)
		realCtx.OwnerAddress = utils.Base58DecodeAddr(ownerAddress.(string))

		tbR.RealContract = realCtx
		tbR.ContractType = core.Transaction_Contract_WithdrawBalanceContract

	} else {
		return errInvalidRequest
	}

	return nil
}

// TBResponseType ...
type TBResponseType struct {
	Success bool         `json:"success"`
	Result  TBResultType `json:"result"`
}

// TBResultType ...
type TBResultType struct {
	Code    string      `json:"code"`
	Message interface{} `json:"message"`
}

package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tronprotocol/grpc-gateway/core"
	"github.com/wlcy/tron/explorer/core/utils"
	"github.com/wlcy/tron/explorer/lib/log"
)

// TBTransfer transfer contract
func TBTransfer(ctx *gin.Context) {
	ctxReq, err := getContextBody(ctx)
	if nil != err || nil == ctxReq {
		ctx.JSON(http.StatusOK, newTBResponse(nil, err))
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

	tbReq.Trx, err = utils.BuildTransaction(tbReq.ContractType, tbReq.RealContract, []byte(tbReq.Data))
	if nil != err {
		return nil, err
	}

	if tbReq.Broadcast {
		resp, err := signAndBroadcastTrx(tbReq.Trx, tbReq.Key)
		if nil != err {
			return nil, err
		}
		tbReq.Resp = resp
	} else {
		tbReq.Resp = extractTrx(tbReq.Trx)
	}

	tbResp := newTBResponse(tbReq, err)

	ctx.JSON(http.StatusOK, tbResp)

	return tbReq, err
}

func extractTrx(trx *core.Transaction) interface{} {
	if nil == trx || nil == trx.RawData || 0 == len(trx.RawData.Contract) {
		return nil
	}

	ret := new(TrxOutputType)
	ret.Hash = utils.Base64Encode(utils.CalcTransactionHash(trx))
	ret.Timestamp = trx.RawData.Timestamp
	ret.Data = string(trx.RawData.Data)

	ret.Signature = extractTrxSign(trx)

	ret.Contracts = extractTrxCtx(trx)

	return ret
}

func extractTrxCtx(trx *core.Transaction) (ret []map[string]interface{}) {
	for _, ctx := range trx.RawData.Contract {
		_, tmp := utils.GetContractInfoStr3(int32(ctx.Type), ctx.Parameter.Value)
		tmpMap, ok := tmp.(map[string]interface{})
		if ok {
			tmpMap["contractType"] = ctx.Type.String()
			tmpMap["contractTypeId"] = int32(ctx.Type)
			ret = append(ret, tmpMap)
		}
	}
	return
}

func extractTrxSign(trx *core.Transaction) (ret []*SignOutputType) {
	for _, sign := range trx.Signature {
		signOut := new(SignOutputType)
		signOut.Bytes = utils.Base64Encode(sign)

		pubKey, err := utils.GetSignedPublicKey(trx)
		if nil != err {
			continue
		}
		signOut.Address, _ = utils.GetTronHexAddress(utils.Base64Encode(pubKey))

		ret = append(ret, signOut)
	}
	return
}

// TrxOutputType ...
type TrxOutputType struct {
	Hash      string                   `json:"hash"`
	Timestamp int64                    `json:"timestamp"`
	Contracts []map[string]interface{} `json:"contracts"`
	Data      string                   `json:"data"`
	Signature []*SignOutputType        `json:"signature"`
}

// SignOutputType ...
type SignOutputType struct {
	Bytes   string `json:"bytes"`
	Address string `json:"address"`
}

func signAndBroadcastTrx(trx *core.Transaction, privatekey string) (resp interface{}, err error) {
	if 0 == len(privatekey) {
		log.Errorf("broadcast need private key for signature:%#v", trx)
		return nil, errInvalidRequest
	}

	// sign transaction
	sign, err := utils.SignTransaction(trx, privatekey)
	if nil != err {
		return nil, errInvalidRequest
	}
	trx.Signature = append(trx.Signature, sign)

	resp, err = GetWalletClient().BroadcastTransaction(trx)

	return
}

func newTBResponse(req *TBRequestType, err error) *TBResponseType {
	resp := new(TBResponseType)

	if nil != err || nil == req {
		resp.Success = false
		resp.Result = TBResultType{}
		resp.Result.Code = "Failed"
		resp.Result.Message = err.Error()
		return resp
	}

	resp.Success = true
	resp.Result = TBResultType{}
	resp.Result.Message = req.Resp
	resp.Result.Code = "Success"

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
	Trx          *core.Transaction                      `json:"-"`
	Resp         interface{}                            `json:"-"`
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

package service

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/wlcy/tron/explorer/lib/mysql"

	"github.com/golang/protobuf/proto"
	"github.com/tronprotocol/grpc-gateway/core"
	"github.com/wlcy/tron/explorer/core/utils"
	"github.com/wlcy/tron/explorer/lib/config"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/util"
	"github.com/wlcy/tron/explorer/web/buffer"
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
)

//QuerySystemStatus ...
func QuerySystemStatus() (*entity.SystemStatusResp, error) {
	var solidityProcess, fullnodeProcess float64
	var systemStatusResp = &entity.SystemStatusResp{}
	netConfig := config.NetType
	if netConfig == "" {
		netConfig = "mainnet"
	}
	netType := &entity.Network{Type: netConfig} //从配置文件中获取网络配置

	blockBuffer := buffer.GetBlockBuffer()

	// db中最大确认块
	confirmBlockIDDb := blockBuffer.GetMaxConfirmedBlockID()
	//查询数据库获取最大块包含非确认--从buffer获取（由于数据库中没存未确认的块，所以从fullnode获取）
	latestdBlockIDDb := blockBuffer.GetMaxBlockID()
	//从buffer中获取fullnode块高
	fullnodeNowBlockID := blockBuffer.GetFullNodeMaxBlockID() //realMaxBlockID
	//从buffer中获取db的确认过的块高
	solidityNowBlockID := blockBuffer.GetSolidityNodeMaxBlockID()
	if solidityNowBlockID > 0 {
		//计算总进度
		solidityProcess = float64(confirmBlockIDDb) / float64(solidityNowBlockID) * 100
	}
	if fullnodeNowBlockID > 0 {
		fullnodeProcess = float64(latestdBlockIDDb) / float64(fullnodeNowBlockID) * 100
	}
	totalProcess := (solidityProcess + fullnodeProcess) / 2
	log.Debugf("QuerySystemStatus,confirmBlockIDDb:[%v],latestdBlockIDDb:[%v],solidityNowBlockID:[%v],fullnodeNowBlockID:[%v]\n",
		confirmBlockIDDb, latestdBlockIDDb, solidityNowBlockID, fullnodeNowBlockID)
	log.Debugf("QuerySystemStatus,solidityProcess:[%.f],fullnodeProcess:[%.f],totalProcess:[%.f]\n",
		solidityProcess, fullnodeProcess, totalProcess)
	//拼接返回数据
	systemStatusResp.Network = netType
	systemStatusResp.Sync = &entity.Sync{Progress: totalProcess}
	systemStatusResp.Full = &entity.BlockNode{Block: fullnodeNowBlockID}
	systemStatusResp.Solidity = &entity.BlockNode{Block: solidityNowBlockID}
	systemStatusResp.Database = &entity.DataBase{Block: latestdBlockIDDb, ConfirmedBlock: confirmBlockIDDb}

	return systemStatusResp, nil
}

//QueryMarketsBuffer ... 从buffer获取市场信息
func QueryMarketsBuffer() ([]*entity.MarketInfo, error) {

	marketBuffer := buffer.GetMarketBuffer()
	markets := marketBuffer.GetMarket()

	return markets, nil
}

//QueryAuth ...
func QueryAuth(req *entity.Auth) (*entity.AuthResp, error) {
	if req == nil || req.Transaction == "" {
		return nil, util.NewErrorMsg(util.Error_common_parameter_invalid)
	}
	jsonData := req.Transaction
	tranHexData := utils.HexDecode(jsonData)
	transaction := &core.Transaction{}
	if err := proto.Unmarshal(tranHexData, transaction); err != nil {
		log.Errorf("pb unmarshal err:[%v];hexData:[%v]", err, tranHexData)
		return nil, err
	}
	//获取rawHash
	rawHash := utils.CalcTransactionHash(transaction)
	log.Debugf("rawHash:[%v]", rawHash)
	//计算地址
	pubKey, _ := utils.GetSignedPublicKey(transaction)
	signatureAddress, err := utils.GetTronBase58Address(utils.HexEncode(pubKey))
	log.Debugf("signatureAddress:[%v]", signatureAddress)
	if nil != err {
		return nil, err
	}
	//解析transaction中的contract为witnessUpdateContract结构
	witnessUpdateContract := &core.WitnessUpdateContract{}
	contractData := transaction.RawData.Contract[0].Parameter.Value
	if err := proto.Unmarshal(contractData, witnessUpdateContract); err != nil || witnessUpdateContract == nil {
		log.Errorf("pb unmarshal to WitnessUpdateContract err:[%v];contractData:[%v]", err, contractData)
		return nil, err
	}
	witnessOwnerAddress := utils.Base58EncodeAddr(witnessUpdateContract.OwnerAddress)
	log.Debugf("witnessOwnerAddress:[%v],signatureAddress:[%v]", witnessOwnerAddress, signatureAddress)
	if witnessOwnerAddress == signatureAddress { //验证通过，计算token
		newToken, err := GenWebToken(signatureAddress)
		log.Debugf("gen web newToken:[%v],err:[%v]", newToken, err)
		return &entity.AuthResp{Token: newToken}, err
	}
	return nil, nil
}

//QueryTestRequestCoin ...
func QueryTestRequestCoin(req *entity.TestCoin, ip string) (*entity.TestCoinResp, error) {
	requestResult := &entity.TestCoinResp{}
	//1. 校验ip最近一小时是否申请过，如果申请过，则退出
	if module.FindByRecentIP(ip) {
		requestResult.Success = false
		requestResult.Code = "ALREADY_REQUESTED_IP"
		requestResult.Message = "Already requested TRX from IP recently"
	} else if module.FindByAddress(req.Address) { //1. 校验address是否申请过，如果申请过，则退出
		requestResult.Success = false
		requestResult.Code = "ALREADY_REQUESTED_IP"
		requestResult.Message = fmt.Sprintf("Already requested for address %v", req.Address)
	} else {
		if verifyCode(req.CaptchaCode) {
			fromAccount := config.TestPk
			amount := mysql.ConvertStringToInt64(config.TestAmount, 0)
			//计算地址

			hexPubKey, hexAddr, base58Addr, err := utils.GetTronPublicInfoByPrivateKey(fromAccount)
			if nil != err {
				log.Errorf(" GetTronPublicInfoByPrivateKey err :[%v]", err)
				return nil, err
			}

			log.Debugf("CreateAccount done: hexPubKey:[%v],hexAddr:[%v],base58Addr:[%v]", hexPubKey, hexAddr, base58Addr)
			transferCtx := &core.TransferContract{}
			transferCtx.OwnerAddress = utils.HexDecode(base58Addr)
			transferCtx.ToAddress = utils.Base58DecodeAddr(req.Address)
			transferCtx.Amount = amount
			ctx, _ := proto.Marshal(transferCtx)
			log.Debugf("ctx:%v", ctx)
			transaction := &core.Transaction{}
			//向主网发布广播
			result, err := GetWalletClient().BroadcastTransaction(transaction)
			if err != nil {
				log.Errorf("call broadcastTransaction err[%v],transaction:[%#v]", err, transaction)
				return requestResult, err
			}
			//解析主网接口返回
			if result.Result { //如果成功则写trxRequest
				log.Debugf("trx request result:%v", result.Result)
				err = module.InsertTrxRequest(req.Address, ip)
				log.Debugf("trx request insert db result:%v", err)
			}
			requestResult.Success = result.Result
			requestResult.Amount = amount
			requestResult.Code = result.Code.String()
			requestResult.Message = string(result.Message)
			/*
			   val fromAccount = config.get[String]("testnet.trx-distribution.pk")
			                   val amount = config.get[Long]("testnet.trx-distribution.amount")
			                   val fromAccountKey = ECKey.fromPrivate(ByteArray.fromHexString(fromAccount))

			                   val transfer = transactionBuilder.buildTrxTransfer(
			                     fromAccountKey.getAddress,
			                     to,
			                     amount)

			                   await(for {
			                     transactionWithRef <- transactionBuilder.setReference(transfer)
			                     signedTransaction = transactionBuilder.sign(transactionWithRef, ByteArray.fromHexString(fromAccount))
			                     result <- wallet.broadcastTransaction(signedTransaction)
			                   } yield {

			                     if (result.result) {
			                       trxRequestModelRepository.insertAsync(TrxRequestModel(
			                         address = to,
			                         ip = ip,
			                       ))
			                     }

			                     Ok(Json.obj(
			                       "success" -> result.result.asJson,
			                       "amount" -> amount.asJson,
			                       "code" -> result.code.toString.asJson,
			                       "message" -> new String(result.message.toByteArray).toString.asJson,
			                     ))
			*/
		} else {
			requestResult.Success = false
			requestResult.Code = "WRONG_CAPTCHA"
			requestResult.Message = "Wrong Captcha Code"

		}
	}
	return requestResult, nil
}

//verifyCode 校验code
func verifyCode(code string) bool {
	var result bool
	var verifyURL = "https://www.google.com/recaptcha/api/siteverify"
	siteCode := config.TestCaptchaSiteKey
	verifyCode := &entity.VerifyCode{Secret: siteCode, Response: code}
	requestBytes, err := json.Marshal(verifyCode)
	if err != nil {
		log.Error(err)
		return result
	}

	postData := bytes.NewBuffer(requestBytes)
	resp, err := util.SendRequest(verifyURL, "POST", "", postData)
	if err != nil {
		log.Error(err)
		return result
	}

	//fmt.Println(string(resp))
	var response entity.VerifyCodeResp

	err = json.Unmarshal(resp.Bytes(), &response)
	if err != nil {
		log.Error(err)
		return result
	}
	log.Debugf("verifyCode return :[%#v]", response)
	result = response.Success
	//始终返回true，暂定
	return true
}

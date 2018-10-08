package utils

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes"

	"github.com/golang/protobuf/proto"
	"github.com/tronprotocol/grpc-gateway/core"
)

var (
	errInvalidContract    = fmt.Errorf("Invalid contract type")
	errInvalidTransaction = fmt.Errorf("Invalid transaction")

	invalidContractType = int32(-1)
)

// BuildTransaction build transaction with giving contract type and contract pb object
//	it will check if the giving type match the contract object
// note:
//	the return transaction has no signature or timestamp
// in:
//	ctxType: contract type
//	contract: contract pb object, should be proto.Message type, should match ctxType, or it will return error
//	data: append data to the transaction
// out:
//	trx: new transaction contain contract, timestamp is empty,
//	err: error message
func BuildTransaction(ctxType core.Transaction_Contract_ContractType, contract interface{}, data []byte) (trx *core.Transaction, err error) {
	pbMsg, ok := contract.(proto.Message)
	if !ok || nil == pbMsg {
		return nil, errInvalidContract
	}

	any, err := ptypes.MarshalAny(pbMsg)
	if nil != err || nil == any {
		return nil, err
	}

	if !strings.HasSuffix(any.TypeUrl, ctxType.String()) {
		return nil, fmt.Errorf("contract type do not match: [%v]<-->[%v]", ctxType.String(), any.TypeUrl)
	}

	trx = new(core.Transaction)
	trx.RawData = new(core.TransactionRaw)
	contractRaw := new(core.Transaction_Contract)
	contractRaw.Type = ctxType
	contractRaw.Parameter = any
	trx.RawData.Contract = append(trx.RawData.Contract, contractRaw)

	if 0 != len(data) {
		trx.RawData.Data = data
	}
	return
}

// GetTransactionContract get first contract from transaction
//	checkt error and contract before check ctxType as the default ctxType is 0 which is a valid contract type
func GetTransactionContract(trx *core.Transaction) (ctxType core.Transaction_Contract_ContractType, contract interface{}, err error) {
	ctxTypes, contracts, _, err := ExtractTransactionContracts(trx)
	if nil != err || 0 == len(contracts) || 0 == len(ctxTypes) {
		return
	}
	return ctxTypes[0], contracts[0], nil
}

// ExtractTransactionContracts extract contract object from transaction
func ExtractTransactionContracts(trx *core.Transaction) (ctxType []core.Transaction_Contract_ContractType, contract []interface{}, data []byte, err error) {
	if nil == trx || nil == trx.RawData || 0 == len(trx.RawData.Contract) {
		return nil, nil, nil, errInvalidTransaction
	}

	if 0 < len(trx.RawData.Data) {
		data = make([]byte, len(trx.RawData.Data))
		copy(trx.RawData.Data, data)
	}

	var pbMsg proto.Message
	for _, ctxRaw := range trx.RawData.Contract {

		pbMsg, err = ptypes.Empty(ctxRaw.Parameter)
		if nil != err {
			return
		}

		err = ptypes.UnmarshalAny(ctxRaw.Parameter, pbMsg)
		if nil != err {
			return
		}

		ctxType = append(ctxType, ctxRaw.Type)
		contract = append(contract, pbMsg)
	}

	return
}

package utils

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
	"github.com/tronprotocol/grpc-gateway/core"
)

func TestSignIssue(t *testing.T) {
	trx1 := &core.Transaction{}
	trx1.RawData = &core.TransactionRaw{}
	raw := trx1.RawData
	trx1.RawData.Contract = make([]*core.Transaction_Contract, 0)

	raw.Expiration = 1535984307000
	raw.Timestamp = 1535984248734
	raw.RefBlockBytes = Base64Decode("6wU=")
	raw.RefBlockHash = Base64Decode("kHFzV3aTNCs=")

	ctx := &core.Transaction_Contract{}
	trx1.RawData.Contract = append(trx1.RawData.Contract, ctx)

	ctx.Type = 1
	ctx.Parameter = &any.Any{}
	ctx.Parameter.TypeUrl = "type.googleapis.com/protocol.TransferContract"
	// ctx.Parameter.Value = base64Decode("ChVBfKvA2SfQD8bcplh+P2RfvIbPY40SFUGI/UILGFB3//k91EPeCUmys3F1jRhj")

	transferCtx := &core.TransferContract{}

	privKey, pubKey, hexAddr, base58Addr, _ := newAccount()
	fmt.Printf("%v\n%v\n%v\n%v\n", privKey, pubKey, hexAddr, base58Addr)

	transferCtx.OwnerAddress = HexDecode(hexAddr)
	transferCtx.ToAddress = Base58DecodeAddr(base58Addr)
	transferCtx.Amount = 100999
	ctx.Parameter.Value, _ = proto.Marshal(transferCtx)

	jsonStr, _ := json.Marshal(trx1)
	fmt.Printf("%s\n\n", jsonStr)

	sign, err := SignTransaction(trx1, privKey)
	_ = err
	fmt.Printf("sign:%v\n%v\n%v\n", sign, hex.EncodeToString(sign), base64.StdEncoding.EncodeToString(sign))
	trx1.Signature = append(trx1.Signature, sign)

	sign1 := trx1.Signature[0]
	fmt.Printf("sign:%v\n%v\n%v\n", sign1, hex.EncodeToString(sign1), base64.StdEncoding.EncodeToString(sign1))

	fmt.Printf("verify result:%v\n", VerifySign(trx1, pubKey))

	recPubKey, err := GetSignedPublicKey(trx1)
	fmt.Printf("recPubKey:%v\n%v\n", recPubKey, HexEncode(recPubKey))

	fmt.Println(Base58DecodeAddr("7YxAaK71utTpYJ8u4Zna7muWxd1pQwimpGxy8"))

}

func TestX(*testing.T) {

	fmt.Println(ToJSONStr(nil))
	fmt.Println(GetContractInfoStr2(2, HexDecode("0a0449504653121541d13433f53fdf88820c2e530da7828ce15d6585cb1a154198f4b89409bb65edbcebb26d46d28cd00bb002ed2001")))
}

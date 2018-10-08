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
	// _, a := GetContractInfoStr2(2, HexDecode("0a0449504653121541d13433f53fdf88820c2e530da7828ce15d6585cb1a154198f4b89409bb65edbcebb26d46d28cd00bb002ed2001"))
	// _, a := GetContractInfoStr2(4, HexDecode("0a1541630405349beb81b61dc5df48f3b70aeb62684b7a121a0a1541beab998551416b02f6721129bb01b51fceceba0810e807121a0a15412fb5abdf8a1670f533c219e7251fe30b8984935910e807121a0a15418a445facc2aa94d72292ebbcb2a611e9fd8a6c6e10ed07121a0a1541b3eec71481e8864f0fc1f601b836b74c4054828710ed07121a0a15414a193c92cd631c1911b99ca964da8fd342f4cddd10ec07121a0a1541d1dbde8b8f71b48655bec4f6bb532a0142b88bc010ed07121a0a15417b88db9da8aacae0a7e967d24c0fc00129e815f610e807121a0a154167e39013be3cdd3814bed152d7439fb5b679140910e807121a0a1541c4bc4d7f64df4fd3670ce38e1a60080a50da85cf10e807121a0a1541243accc5241d97ce79272b06952ee88a34d8e1f910d80b121a0a1541c189fa6fc9ed7a3580c3fe291915d5c6a6259be710e807121a0a15414d1ef8673f916debb7e2515a8f3ecaf2611034aa10ed07121a0a1541d49bf5202b3dba65d46a5be73396b6b66d3555aa10e807"))
	_, a := GetContractInfoStr2(4, HexDecode("0a154118f13591da0077a0f57a6c9a4759960714a23fec121a0a15410694981b116304ed21e05896fb16a6bc2e91c92c10e807121a0a15411d7aba13ea199a63d1647e58e39c16a9bb9da68910e807121a0a15411661f25387370c9cd3a9a5d97e60ca90f4844e7e10e807121a0a15417ad0ee1300d0366e901fa613a929137dde1d222410e807121a0a15418a445facc2aa94d72292ebbcb2a611e9fd8a6c6e10e807121a0a1541746e6af4ac9db3473c0c955f1fca11d4013f32ed10e807121a0a15414a193c92cd631c1911b99ca964da8fd342f4cddd10a01f121a0a1541f70386347e689e6308e4172ed7319c49c0f66e0b10e807121a0a1541c189fa6fc9ed7a3580c3fe291915d5c6a6259be710e807121a0a154116440834509c59de4ee6ba4933678626f451befe10e807121a0a15412fb5abdf8a1670f533c219e7251fe30b8984935910e807121a0a1541d32b3fa8ca0b4896257fdf1821ac8d116da84c4510e807121a0a1541c4bc4d7f64df4fd3670ce38e1a60080a50da85cf10e807121a0a1541e40de6895c142ade8b86194063bcdbaa6c9360b610e807121a0a15414593d27b70d21454b39ab60bf13291dae8dc032610bb09121a0a1541f29f57614a6b201729473c837e1d2879e9f90b8e10e807121a0a1541d49bf5202b3dba65d46a5be73396b6b66d3555aa10e807121a0a154193a8bc2e7d6bb1bd75fb2d74107ffbda81af439d10e807121a0a154127a6419bbe59f4e64a064d710787e578a150d6a710e807121a0a1541d1dbde8b8f71b48655bec4f6bb532a0142b88bc010e807121a0a15416419765bacf1dc441f722cabc8b661140558bb5d10e807121a0a154167e39013be3cdd3814bed152d7439fb5b679140910e807121a0a15415863f6091b8e71766da808b1dd3159790f61de7d10e807121a0a1541243accc5241d97ce79272b06952ee88a34d8e1f910e807121a0a15415095d4f4d26ebc672ca12fc0e3a48d6ce3b169d210e807121a0a1541b3eec71481e8864f0fc1f601b836b74c4054828710e807121a0a154172fd5dfb8ab36eb28df8e4aee97966a60ebf9efe10e807121a0a154138e3e3a163163db1f6cfceca1d1c64594dd1f0ca10e807"))

	fmt.Println(formatContractJSONStr(a))
}

func TestJSONNil(*testing.T) {
	acc := new(core.Account)
	fmt.Printf("%v\n", ToJSONStr(acc.FrozenSupply))
	fmt.Printf("%v\n", ToJSONStr(acc.LatestAssetOperationTime))
}

func TestBuildTrx(*testing.T) {
	ctx := new(core.TransferContract)
	ctx.OwnerAddress = []byte("abcdefg")
	ctx.ToAddress = []byte("gfedcba")
	ctx.Amount = 1 * 1000000 // in sun

	trx, err := BuildTransaction(core.Transaction_Contract_TransferContract, ctx, nil)
	fmt.Printf("%v\n%#v\n", err, ToJSONStr(trx))

	ctxTypes, contracts, data, err := ExtractTransactionContracts(trx)
	fmt.Printf("%#v\n%#v\n%v\n%v\n", ctxTypes, contracts, data, err)

	ctxType, contract, err := GetTransactionContract(trx)
	fmt.Printf("%#v\n%#v\n%v\n", ctxType, contract, err)
}

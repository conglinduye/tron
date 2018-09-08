package utils

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/golang/protobuf/proto"
	"github.com/tronprotocol/grpc-gateway/core"
)

// GetTronPublicInfoByPrivateKey 根据私钥获取公共信息
// param:
//	in: hex encoding private key
// return:
//	publicKey: 公钥
//	hexAddr: hex address
//	base58Addr: base58 address
func GetTronPublicInfoByPrivateKey(in string) (publicKey, hexAddr, base58Addr string, err error) {
	publicKey, err = getPublicKey(in)
	if nil != err {
		return
	}
	hexAddr, err = GetTronHexAddress(publicKey)
	if nil != err {
		return
	}
	base58Addr = Base58EncodeAddr(HexDecode(hexAddr))

	return
}

// GetTronPublickey 根据 hex encoding private key 生成 hex encoding public key
//	in: hex encoding private key
//	out: hex encoding public key (uncompressed public key)
func GetTronPublickey(in string) (out string, err error) {
	priKey, err := getPrivateKey(in)
	if nil != err {
		return "", err
	}
	// TRON 使用 ECDSA 生成公钥，增加了前缀 04，将公钥的（X，Y）拼接
	// return strings.ToUpper("04" + hex.EncodeToString(priKey.PublicKey.X.Bytes()) + hex.EncodeToString(priKey.PublicKey.Y.Bytes())), nil
	out = HexEncode(ethcrypto.FromECDSAPub(&priKey.PublicKey))
	return
}

// GetTronHexAddress 根据 hex encoding public key生成 hex encoding address
//	in: hex encoding public key (uncompressed public key)
//	out: hex encoding address
func GetTronHexAddress(in string) (out string, err error) {
	pubBytes, err := hex.DecodeString(in)
	if nil != err {
		return "", err
	}
	rawPubKey := pubBytes[1:] // remove prefix byte

	sha3Hash := sha3.NewKeccak256() // use sha3 keccad256
	sha3Hash.Write(rawPubKey)
	hashRet := sha3Hash.Sum(nil)

	hashRetStr := HexEncode(hashRet)                              // covert to hex string
	out = fmt.Sprintf("%s%s", AddressPrefixMain, hashRetStr[24:]) // address prefix + hash remove first 24 length

	return
}

// GetTronBase58Address 根据 hex encoding public key生成 hex encoding address
//	in: hex encoding public key (uncompressed public key)
//	out: base58 encoding address
func GetTronBase58Address(in string) (out string, err error) {
	hexAddr, err := GetTronHexAddress(in)
	if nil != err {
		return "", err
	}

	out = Base58EncodeAddr(HexDecode(hexAddr))

	return
}

// SignTransaction 根据私钥对交易进行签名，签名不回填transaction对象
//	transaction: 交易
//	hexPrivKey: hex encoding 私钥
func SignTransaction(transaction *core.Transaction, hexPrivKey string) ([]byte, error) {
	rawData, _ := proto.Marshal(transaction.RawData)
	hash0 := sha256.Sum256(rawData)

	privKey, err := getPrivateKey(hexPrivKey)
	if nil != err || nil == privKey {
		return nil, err
	}

	signData, err := ethcrypto.Sign(hash0[:], privKey)
	// fmt.Printf("sign data, raw hash:%v\ngen sign:%v\n", hash0, signData)
	// transaction.Signature = append(transaction.Signature, signData)
	return signData, err
}

// VerifySign 验证签名
//	transaction: 交易对象
//	pubKey: hex encoding 公钥 (uncompressed key)
func VerifySign(transaction *core.Transaction, pubKey string) bool {
	sign := transaction.Signature[0]
	if len(sign) != 65 { // sign check
		return false
	}
	rawData, _ := proto.Marshal(transaction.RawData)
	hash0 := sha256.Sum256(rawData)
	// fmt.Printf("verifying sign, raw hash:%v\n", hash0)
	// fmt.Printf("origin sign:%v\n", sign)
	// _ = sign
	// _ = hash0
	return ethcrypto.VerifySignature(HexDecode(pubKey), hash0[:], sign[:64])
}

// GetSignedPublicKey 获取交易签名账户的公钥
func GetSignedPublicKey(transaction *core.Transaction) ([]byte, error) {
	sign := transaction.Signature[0]
	if len(sign) != 65 { // sign check
		return nil, ErrorInvalidSign
	}
	rawData, _ := proto.Marshal(transaction.RawData)
	hash0 := sha256.Sum256(rawData)

	return ethcrypto.Ecrecover(hash0[:], sign)
}

// getPrivateKey 根据hexEncoding私钥生成私钥对象
//	in: hex encoding private key
func getPrivateKey(in string) (*ecdsa.PrivateKey, error) {
	return ethcrypto.HexToECDSA(in)
}

// getPublicKey 根据hexEncoding私钥生成hexEncoding公钥
//	in: hex encoding private key
func getPublicKey(in string) (string, error) {
	return GetTronPublickey(in)
}

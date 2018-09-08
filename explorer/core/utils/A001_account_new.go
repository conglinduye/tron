package utils

import (
	"encoding/json"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/satori/go.uuid"
)

// newAccount 创建新的账户，返回私钥，公钥，地址信息
func newAccount() (hexPrivKey, hexPubKey, hexAddr, base58Addr string, err error) {
	privKey, err := ethcrypto.GenerateKey()
	if nil != err {
		return
	}
	privKeyByte := ethcrypto.FromECDSA(privKey)
	hexPrivKey = HexEncode(privKeyByte)
	hexPubKey, _ = getPublicKey(hexPrivKey)
	hexAddr, _ = GetTronHexAddress(hexPubKey)
	base58Addr = Base58EncodeAddr(HexDecode(hexAddr))

	return
}

// KeyStorage 保存加密的私钥
/*
{
    "address": "TQdXjU831NbBc2H9jSjZ952F6q5zmDvjYe",
    "key": "92fe6e64f35f15898d6c4356d1351cded76d66f3435f1b6f2c2aa9ba383895aa46165ecb2e2152c47fd78d01f1f125bb52c7233df4b8634a5b36e6f8729ff4f5",
    "salt": "4d079ef8-02b3-421e-afdf-9de7b5cd602a",
    "version": 1
}
*/
type KeyStorage struct {
	Version       interface{} `json:"version"` // version
	Address       string      `json:"address"` // account address in base58 encoding
	Key           string      `json:"key"`     // encrypted private key
	Salt          string      `json:"salt"`    // salt used to encrypt private key
	privateHexKey string      // 私钥
	publicKey     string      // 公钥
	hexAddr       string      // hex address
}

// genPrivateKeyStoreage 使用密码对私钥信息进行加密，生成私钥秘文
//	生成字符串写入文件
func genPrivateKeyStoreage(password string, privKey string) (string, error) {
	storage := new(KeyStorage)
	storage.Version = 1

	pubKey, err := getPublicKey(privKey)
	if nil != err {
		return "", err
	}

	storage.Address, err = GetTronBase58Address(pubKey)
	if nil != err {
		return "", err
	}

	uuidVal, err := uuid.NewV4()
	storage.Salt = uuidVal.String()

	storage.Key = HexEncode(AesCtrEncrypt(password, storage.Salt, []byte(privKey)))

	rawData, err := json.Marshal(storage)
	return HexEncode(rawData), err
}

// readPrivateKeyStorage 读取私钥加密数据并解密
// param
//	password: 加密密码
//	content: 私钥加密数据
// return
//	storage: 解密数据
//	result: 是否成功
//	err: 错误信息
func readPrivateKeyStorage(password string, content string) (storage *KeyStorage, result bool, err error) {
	storage = new(KeyStorage)

	err = json.Unmarshal(HexDecode(content), storage)
	if nil != err {
		return nil, false, err
	}

	privKey := string(AesCtrDecrypt(password, storage.Salt, HexDecode(storage.Key)))

	pubKey, hexAddr, base58Addr, err := GetTronPublicInfoByPrivateKey(privKey)
	if nil != err {
		return nil, false, err
	}
	if base58Addr != storage.Address {
		return storage, false, ErrorDecrypt
	}
	storage.publicKey = pubKey
	storage.privateHexKey = privKey
	storage.hexAddr = hexAddr

	return storage, true, nil
}

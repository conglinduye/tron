package utils

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/btcsuite/btcutil/base58"
)

// Base54DecodeAddr 将base58地址转换为原始字节
func Base58DecodeAddr(in string) (ret []byte) {
	result, ver, err := base58.CheckDecode(in)
	if nil == err {
		ret = append(ret, ver)
		ret = append(ret, result...)
		return
	}
	fmt.Println(err)
	return
}

// Base58EncodeAddr 将地址字节码编码为base58字符串
func Base58EncodeAddr(in []byte) string {
	return base58.CheckEncode(in[1:], in[0]) // first byte is version, reset is data
}

// Base64Decode ...
func Base64Decode(in string) []byte {
	ret, _ := base64.StdEncoding.DecodeString(in)
	return ret
}

// Base64Encode ...
func Base64Encode(in []byte) string {
	return base64.StdEncoding.EncodeToString(in)
}

// HexDecode ...
func HexDecode(in string) []byte {
	ret, _ := hex.DecodeString(in)
	return ret
}

// HexEncode ...
func HexEncode(in []byte) string {
	return hex.EncodeToString(in)
}

// BinaryBigEndianEncodeInt64 ...
func BinaryBigEndianEncodeInt64(num int64) []byte {
	ret := make([]byte, 8)
	binary.BigEndian.PutUint64(ret, uint64(num))
	return ret
}

// BinaryBigEndianDecodeUint64 ...
func BinaryBigEndianDecodeUint64(d []byte) uint64 {
	return binary.BigEndian.Uint64(d)
}

// ConverTimestamp convert ms to yyyy-mm-dd hh24:mi:ss.ms
func ConverTimestamp(ts int64) string {
	str := fmt.Sprintf("%v", ts)
	tsv := ts
	if len(str) == 19 {
		//
	} else if len(str) == 16 {
		tsv *= 1000
	} else if len(str) == 13 {
		tsv *= 1000000
	} else if len(str) == 10 {
		tsv *= 1000000000
	} else {
		return str
	}
	return time.Unix(0, tsv).Format("2006-01-02 15:04:05.000000")
}

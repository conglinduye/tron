package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// TestNet 是否连接 testNet，默认false
var TestNet bool

// GetRandSolidityNodeAddr 随机获取一个solidity node ip
func GetRandSolidityNodeAddr() string {
	if TestNet {
		return TestSolidityNodeList[rand.Int31n(int32(len(TestSolidityNodeList)))]
	}
	return SolidityNodeList[rand.Int31n(int32(len(SolidityNodeList)))]
}

// GetRandFullNodeAddr 随机获取一个full node ip
func GetRandFullNodeAddr() string {
	if TestNet {
		return TestFullNodeList[rand.Int31n(int32(len(TestFullNodeList)))]
	}
	return FullNodeList[rand.Int31n(int32(len(FullNodeList)))]
}

// 地址前缀 测试/主网
const (
	AddressPrefixTest = "a0" //a0 + address
	AddressPrefixMain = "41" //41 + address

	DefaultGrpPort = 50051
	DefaultP2pPort = 18888
)

// Node List info from:
// https://github.com/tronprotocol/Documentation/blob/master/TRX_CN/Official_Public_Node.md

// SolidityNodeList Solidity节点列表
var SolidityNodeList = []string{
	"39.105.66.80",   // good
	"47.254.39.153",  // good
	"47.89.244.227",  // good
	"39.105.118.15",  // good
	"47.75.108.229",  // good
	"34.234.164.105", // good
	"18.221.34.0",    // time out happen, good
	"35.178.11.0",    // good
	"35.180.18.107",  // good
	// "52.63.152.13",   // time out happen +++++
	// "18.231.123.107", // time out happen +++++
}

// FullNodeList Full节点列表
var FullNodeList = []string{
	"54.236.37.243", // not fully implement
	"52.53.189.99",  // not fully implement
	"18.196.99.16",
	"34.253.187.192",
	"52.56.56.149",
	"35.180.51.163",
	"54.252.224.209",
	"18.228.15.36",
	"52.15.93.92",
	"34.220.77.106",
	"13.127.47.162",
	"13.124.62.58",
	"13.229.128.108",
	"35.182.37.246",
	"34.200.228.125",
	"18.220.232.201",
	"13.57.30.186",
	"35.165.103.105",
	"18.184.238.21",
	"34.250.140.143",
	"35.176.192.130",
	"52.47.197.188",
	"52.62.210.100",
	"13.231.4.243",
	"18.231.76.29",
	"35.154.90.144",
	"13.125.210.234",
	"13.250.40.82",
	"35.183.101.48",
	// "47.104.11.194", // grpc connection failed
}

// Faield, error:rpc error: code = Unavailable desc = all SubConns are in TransientFailure, latest connection error: connection error: desc = "transport: Error while dialing dial tcp 47.104.11.194:50051: connect: connection refused"

// TestSolidityNodeList ...
var TestSolidityNodeList = []string{
	"47.75.13.181",
}

// TestFullNodeList ...
var TestFullNodeList = []string{
	"47.90.240.201",
	"47.89.188.246",
	"47.90.208.195",
	"47.89.188.162",
	"47.89.185.110",
	"47.89.183.137",
	"47.90.240.239",
	"47.88.55.186",
	"47.254.75.152",
	"47.254.36.2",
	"47.254.73.154",
	"47.254.20.22",
	"47.254.33.129",
	"47.254.45.208",
	"47.74.159.205",
	"47.74.149.105",
	"47.74.144.205",
	"47.74.159.52",
	"47.88.237.77",
	"47.74.149.180",
	"47.88.229.149",
	"47.74.182.133",
	"47.88.229.123",
	"47.74.152.210",
	"47.75.205.223",
	"47.75.113.95",
	"47.75.57.234",
}

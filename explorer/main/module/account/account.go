package account

import (
	"time"

	"github.com/tronprotocol/grpc-gateway/api"
	"github.com/tronprotocol/grpc-gateway/core"
	"github.com/wlcy/tron/explorer/core/utils"
)

// Account 账户信息
type Account struct {
	Raw            *core.Account
	NetRaw         *api.AccountNetMessage
	Name           string
	Addr           string
	CreateTime     int64
	IsWitness      int8
	Fronzen        string
	AssetIssueName string

	AssetBalance map[string]int64
	Votes        string

	// acccount net info
	freeNetUsed    int64
	freeNetLimit   int64
	netUsed        int64
	netLimit       int64
	totalNetLimit  int64
	totalNetWeight int64
	AssetNetUsed   string
	AssetNetLimit  string

	/*
		`account_name` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'Account name',
		`address` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'Base 58 encoding address',
		`balance` bigint(20) NOT NULL DEFAULT '0' COMMENT 'TRX balance, in sun',
		`create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '账户创建时间',
		`latest_operation_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '账户最后操作时间',
		`is_witness` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否为wintness; 0: 不是，1:是',
		`frozen` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '冻结金额, 投票权',
		`create_unix_time` int(32) NOT NULL DEFAULT '0' COMMENT '账户创建时间unix时间戳，用于分区',
		`allowance` bigint(20) DEFAULT '0',
		`latest_withdraw_time` bigint(20) DEFAULT '0',
		`latest_consume_time` bigint(20) DEFAULT '0',
		`latest_consume_free_time` bigint(20) DEFAULT '0',
		`votes` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '',
	*/
}

var beginTime, _ = time.Parse("2006-01-02 15:03:04.999999", "2018-06-25 00:00:00.000000")

// SetRaw 使用 core.Account 对象设置 Account 信息
func (a *Account) SetRaw(raw *core.Account) {
	a.Raw = raw
	a.Name = string(raw.AccountName)
	a.Addr = utils.Base58EncodeAddr(raw.Address)
	a.AssetIssueName = string(raw.AssetIssuedName)
	a.CreateTime = raw.CreateTime
	if a.CreateTime == 0 {
		a.CreateTime = beginTime.UnixNano()
	}
	a.IsWitness = 0
	if raw.IsWitness {
		a.IsWitness = 1
	}
	if len(raw.Frozen) > 0 {
		a.Fronzen = utils.ToJSONStr(raw.Frozen)

	}
	a.AssetBalance = a.Raw.Asset
	if len(raw.Votes) > 0 {
		a.Votes = utils.ToJSONStr(raw.Votes)
	}
}

//SetNetRaw 使用 api.AccountNetMessage 设置 accountNet 信息
func (a *Account) SetNetRaw(netRaw *api.AccountNetMessage) {
	if nil == netRaw {
		return
	}
	a.NetRaw = netRaw
	a.AssetNetUsed = utils.ToJSONStr(netRaw.AssetNetUsed)
	a.AssetNetLimit = utils.ToJSONStr(netRaw.AssetNetLimit)
	a.freeNetUsed = netRaw.FreeNetUsed
	a.freeNetLimit = netRaw.FreeNetLimit
	a.netLimit = netRaw.NetLimit
	a.netUsed = netRaw.NetUsed
	a.totalNetLimit = netRaw.TotalNetLimit
	a.totalNetWeight = netRaw.TotalNetWeight
}

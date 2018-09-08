/*
tronscan DataBase Version 1.0
*/

/*
账户信息
更新规则：通过分析交易中出现的账户，触发账户更新
*/

CREATE TABLE `account` (
  `account_name` varchar(300) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'Account name',
  `address` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL  DEFAULT '' COMMENT 'Base 58 encoding address',
  `balance` bigint NOT NULL DEFAULT '0' COMMENT 'TRX balance, in sun',
  `create_time` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '账户创建时间',
  `latest_operation_time` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '账户最后操作时间',
  `is_witness` tinyint NOT NULL DEFAULT '0' COMMENT '是否为wintness; 0: 不是，1:是',
  `modified_time` timestamp(6) NOT NULL  DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '记录更新时间',
  `fronze_amount` bigint NOT NULL DEFAULT '0' COMMENT '冻结金额, 投票权',
  `create_unix_time` int(32) NOT NULL DEFAULT '0' COMMENT '账户创建时间unix时间戳，用于分区',
  UNIQUE KEY `uniq_account_address` (`address`,`create_unix_time` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
/*!50100 PARTITION BY RANGE (create_unix_time)
(PARTITION p0 VALUES LESS THAN (1530403200) ENGINE = InnoDB,
 PARTITION p1 VALUES LESS THAN (1533081600) ENGINE = InnoDB,
 PARTITION p2 VALUES LESS THAN (1535760000) ENGINE = InnoDB,
 PARTITION p3 VALUES LESS THAN (1538352000) ENGINE = InnoDB,
 PARTITION p4 VALUES LESS THAN (1541030400) ENGINE = InnoDB,
 PARTITION p5 VALUES LESS THAN (1543622400) ENGINE = InnoDB,
 PARTITION p6 VALUES LESS THAN (1546300800) ENGINE = InnoDB,
 PARTITION p7 VALUES LESS THAN (1548979200) ENGINE = InnoDB,
 PARTITION p8 VALUES LESS THAN (1551398400) ENGINE = InnoDB,
 PARTITION p9 VALUES LESS THAN (1554076800) ENGINE = InnoDB,
 PARTITION p10 VALUES LESS THAN (1556668800) ENGINE = InnoDB,
 PARTITION p11 VALUES LESS THAN (1559347200) ENGINE = InnoDB,
 PARTITION p12 VALUES LESS THAN (1561939200) ENGINE = InnoDB,
 PARTITION p13 VALUES LESS THAN (1564617600) ENGINE = InnoDB,
 PARTITION p14 VALUES LESS THAN (1567296000) ENGINE = InnoDB,
 PARTITION p15 VALUES LESS THAN (1569888000) ENGINE = InnoDB,
 PARTITION p16 VALUES LESS THAN (1572566400) ENGINE = InnoDB,
 PARTITION p17 VALUES LESS THAN (1575158400) ENGINE = InnoDB,
 PARTITION p18 VALUES LESS THAN (1577836800) ENGINE = InnoDB,
 PARTITION p19 VALUES LESS THAN (1580515200) ENGINE = InnoDB,
 PARTITION p20 VALUES LESS THAN (1583020800) ENGINE = InnoDB,
 PARTITION p21 VALUES LESS THAN (1585699200) ENGINE = InnoDB,
 PARTITION p22 VALUES LESS THAN (1588291200) ENGINE = InnoDB,
 PARTITION p23 VALUES LESS THAN (1590969600) ENGINE = InnoDB,
 PARTITION p24 VALUES LESS THAN (1593561600) ENGINE = InnoDB,
 PARTITION p25 VALUES LESS THAN MAXVALUE ENGINE = InnoDB) */;

/*
见证人信息 witness_info 
不需要创建表，直接从主网接口获取，存入缓存，使用缓存数据查询关联信息
*/


/*
账户通证余额信息
更新规则: 账户更新时，调用主网接口获取当前账户信息并更新
*/
CREATE TABLE `account_asset_balance` (
  `address` varchar(45) NOT NULL COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT 'Base 58 encoding address for the token owner',
  `token_name` varchar(300) NOT NULL COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '通证名称',
  `creator_address` varchar(45) NOT NULL COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT 'Token creator address',
  `balance` bigint NOT NULL DEFAULT '0' COMMENT '通证余额',
  KEY `idx_account_asset_balance_addr` (`address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


/*
账户投票结果，只存储有效的结果
更新规则：当前主网没有接口直接获取结果，可能需要分析投票交易进行更新


*/
CREATE TABLE `account_vote_result` (
  `address` varchar(45) NOT NULL COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'voter address',
  `to_address` varchar(45) NOT NULL COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '投票接收人',
  `vote` bigint NOT NULL DEFAULT '0' COMMENT '投票数',
  KEY `idx_account_vote_result_id` (`address`, `vote` DESC),
  KEY `idx_account_vote_result_to` (`to_address`, `vote` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


/*
通证发行信息
直接从主网接口获取

代币的余额在创建账户的余额记录中没有
因此代币的余额需要独立的计算逻辑: -> sum(participate contract) where token_name = 'TOKEN_NAME' and to_address = 'CREATOR_ADDRESS'
*/
CREATE TABLE `asset_info` (
  `token_name` varchar(300) NOT NULL COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '通证名称',
  `creator_address` varchar(45) NOT NULL COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT 'Token creator address',
  `abbr` varchar(100) NOT NULL COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '通证缩写',
  `total_supply` bigint NOT NULL DEFAULT '0' COMMENT '通证发行量',
  `frozen_info` text NULL COMMENT '冻结信息，json',
  `trx_num` bigint NOT NULL DEFAULT '0' COMMENT '通证汇率分子，单位sun',
  `num` bigint NOT NULL DEFAULT '0' COMMENT '通证汇率分母',
  `start_time` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '开始时间',
  `end_time` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '结束时间',
  `token_desc` text NULL COLLATE utf8mb4_unicode_ci COMMENT '通证说明',
  `url` varchar(800) NOT NULL COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '通证主页url'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

/*
区块信息表
按照区块高度的hash分表
更新规则：solityNode为主，fullnode为辅
*/
CREATE TABLE `blocks` (
  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID。高度',
  `block_hash` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '区块hash',
  `parent_hash` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '区块父级hash',
  `witness_address` varchar(300) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '代表节点地址',
  `tx_trie_hash` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '验证数根的hash值',
  `block_size` int(32) DEFAULT '0' COMMENT '区块大小',
  `transaction_num` int(32) DEFAULT '0' COMMENT '交易数',
  `confirmed` tinyint(4) DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块创建时间',
  `modified_time` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '记录更新时间',
  PRIMARY KEY (`block_id`),
  UNIQUE KEY `uniq_blocks_id` (`block_id` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
/*!50100 PARTITION BY HASH (`block_id`)
PARTITIONS 100 */;

/*
所有交易信息表
按照区块高度的hash分表
更新规则：solityNode为主，fullnode为辅
*/
CREATE TABLE `transactions` (
  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID，高度',
  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\nAccountCreateContract = 0;\r\nTransferContract = 1;\r\nTransferAssetContract = 2;\r\nVoteAssetContract = 3;\r\nVoteWitnessContract = 4;\r\nWitnessCreateContract = 5;\r\nAssetIssueContract = 6;\r\nWitnessUpdateContract = 8;\r\nParticipateAssetIssueContract = 9;\r\nAccountUpdateContract = 10;\r\nFreezeBalanceContract = 11;\r\nUnfreezeBalanceContract = 12;\r\nWithdrawBalanceContract = 13;\r\nUnfreezeAssetContract = 14;\r\nUpdateAssetContract = 15;\r\nProposalCreateContract = 16;\r\nProposalApproveContract = 17;\r\nProposalDeleteContract = 18;\r\nSetAccountIdContract = 19;\r\nCustomContract = 20;\r\n// BuyStorageContract = 21;\r\n// BuyStorageBytesContract = 22;\r\n// SellStorageContract = 23;\r\nCreateSmartContract = 30;\r\nTriggerSmartContract = 31;\r\nGetContract = 32;\r\nUpdateSettingContract = 33;\r\nExchangeCreateContract = 41;\r\nExchangeInjectContract = 42;\r\nExchangeWithdrawContract = 43;\r\nExchangeTransactionContract = 44;',
  `contract_data` varchar(5000) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易内容数据,原始数据byte hex encoding',
  `result_data` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易结果对象byte hex encoding',
  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
  `to_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '接收方地址',
  `fee` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易花费 单位 sun',
  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
  `expire_time` bigint(20) DEFAULT '0',
  `modified_time` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '记录更新时间',
  PRIMARY KEY (`trx_hash`,`block_id`),
  KEY `idx_transactions_hash_create_time` (`block_id`,`trx_hash`,`create_time` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
/*!50100 PARTITION BY HASH (`block_id`)
PARTITIONS 100 */;


/*
所有TRX转账表
按照区块高度的hash分表
更新规则：solityNode为主，fullnode为辅
*/
CREATE TABLE `contract_trx_transfer` (
  `trx_hash` varchar(64) NOT NULL DEFAULT '' COMMENT '交易hash',
  `block_id` bigint NOT NULL DEFAULT '0' COMMENT '区块ID',
  `owner_address` varchar(300) NOT NULL DEFAULT '' COMMENT '发起方地址',
  `to_address` varchar(300) NOT NULL DEFAULT '' COMMENT '接收方地址',
  `amount` bigint NOT NULL DEFAULT '0' COMMENT '转账金额。单位sun',
  `token_name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'token名称',
  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\nAccountCreateContract = 0;\\r\\nTransferContract = 1;\\r\\nTransferAssetContract = 2;\\r\\nVoteAssetContract = 3;\\r\\nVoteWitnessContract = 4;\\r\\nWitnessCreateContract = 5;\\r\\nAssetIssueContract = 6;\\r\\nWitnessUpdateContract = 8;\\r\\nParticipateAssetIssueContract = 9;\\r\\nAccountUpdateContract = 10;\\r\\nFreezeBalanceContract = 11;\\r\\nUnfreezeBalanceContract = 12;\\r\\nWithdrawBalanceContract = 13;\\r\\nUnfreezeAssetContract = 14;\\r\\nUpdateAssetContract = 15;\\r\\nProposalCreateContract = 16;\\r\\nProposalApproveContract = 17;\\r\\nProposalDeleteContract = 18;\\r\\nSetAccountIdContract = 19;\\r\\nCustomContract = 20;\\r\\n// BuyStorageContract = 21;\\r\\n// BuyStorageBytesContract = 22;\\r\\n// SellStorageContract = 23;\\r\\nCreateSmartContract = 30;\\r\\nTriggerSmartContract = 31;\\r\\nGetContract = 32;\\r\\nUpdateSettingContract = 33;\\r\\nExchangeCreateContract = 41;\\r\\nExchangeInjectContract = 42;\\r\\nExchangeWithdrawContract = 43;\\r\\nExchangeTransactionContract = 44;',
  `confirmed` tinyint NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
  `modified_time` timestamp(6) NOT NULL  DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '记录更新时间',
  PRIMARY KEY (`trx_hash`,`block_id`),
  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
PARTITION BY HASH(block_id) PARTITIONS 100;

/*
所有通证转账表
按照区块高度的hash分表
更新规则：solityNode为主，fullnode为辅
*/
CREATE TABLE `contract_token_transfer` (
  `trx_hash` varchar(64) NOT NULL DEFAULT '' COMMENT '交易hash',
  `block_id` bigint NOT NULL DEFAULT '0' COMMENT '区块ID',
  `owner_address` varchar(300) NOT NULL DEFAULT '' COMMENT '发起方地址',
  `to_address` varchar(300) NOT NULL DEFAULT '' COMMENT '接收方地址',
  `amount` bigint NOT NULL DEFAULT '0' COMMENT '转账金额。单位sun',
  `token_name` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'token名称',
  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\nAccountCreateContract = 0;\\r\\nTransferContract = 1;\\r\\nTransferAssetContract = 2;\\r\\nVoteAssetContract = 3;\\r\\nVoteWitnessContract = 4;\\r\\nWitnessCreateContract = 5;\\r\\nAssetIssueContract = 6;\\r\\nWitnessUpdateContract = 8;\\r\\nParticipateAssetIssueContract = 9;\\r\\nAccountUpdateContract = 10;\\r\\nFreezeBalanceContract = 11;\\r\\nUnfreezeBalanceContract = 12;\\r\\nWithdrawBalanceContract = 13;\\r\\nUnfreezeAssetContract = 14;\\r\\nUpdateAssetContract = 15;\\r\\nProposalCreateContract = 16;\\r\\nProposalApproveContract = 17;\\r\\nProposalDeleteContract = 18;\\r\\nSetAccountIdContract = 19;\\r\\nCustomContract = 20;\\r\\n// BuyStorageContract = 21;\\r\\n// BuyStorageBytesContract = 22;\\r\\n// SellStorageContract = 23;\\r\\nCreateSmartContract = 30;\\r\\nTriggerSmartContract = 31;\\r\\nGetContract = 32;\\r\\nUpdateSettingContract = 33;\\r\\nExchangeCreateContract = 41;\\r\\nExchangeInjectContract = 42;\\r\\nExchangeWithdrawContract = 43;\\r\\nExchangeTransactionContract = 44;',
  `confirmed` tinyint NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
  `modified_time` timestamp(6) NOT NULL  DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '记录更新时间',
   PRIMARY KEY (`trx_hash`,`block_id`),
   KEY `idx_token_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
PARTITION BY HASH(block_id) PARTITIONS 100;

/*
所有投票信息表
按照区块高度的hash分表
更新规则：solityNode为主，fullnode为辅
*/
CREATE TABLE `contract_vote_witness` (
  `trx_hash` varchar(64) NOT NULL DEFAULT '' COMMENT '交易hash',
  `block_id` bigint NOT NULL DEFAULT '0' COMMENT '区块ID',
  `voter_address` varchar(300) NOT NULL DEFAULT '' COMMENT '投票人地址',
  `candidate_address` varchar(300) NOT NULL DEFAULT '' COMMENT '被投票人地址',
  `vote_num` int(64) NOT NULL DEFAULT '0' COMMENT '投票数',
  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '类型，同blocks类型',
  `confirmed` tinyint NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
  `modified_time` timestamp(6) NOT NULL  DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '记录更新时间',
  PRIMARY KEY (`trx_hash`,`block_id`),
  KEY `idx_vote_witness_hash_create_time` (`block_id`,`trx_hash`,`create_time` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
PARTITION BY HASH(block_id) PARTITIONS 100;

/*
所有认购通证信息表
按照区块高度的hash分表
更新规则：solityNode为主，fullnode为辅
*/
CREATE TABLE `contract_participate_asset` (
  `trx_hash` varchar(64) NOT NULL DEFAULT '' COMMENT '交易hash',
  `block_id` bigint NOT NULL DEFAULT '0' COMMENT '区块ID',
  `owner_address` varchar(300) NOT NULL DEFAULT '' COMMENT '投票人地址',
  `to_address` varchar(300) NOT NULL DEFAULT '' COMMENT '被投票人地址',
  `amount` bigint NOT NULL DEFAULT '0' COMMENT '认购金额',
  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '类型，同blocks类型',
  `token_name` varchar(300) NOT NULL DEFAULT '' COMMENT '通证名称',
  `confirmed` tinyint NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
  `modified_time` timestamp(6) NOT NULL  DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '记录更新时间',
  PRIMARY KEY (`trx_hash`,`block_id`),
  KEY `idx_participate_asset_hash_create_time` (`block_id`,`trx_hash`,`create_time` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
PARTITION BY HASH(block_id) PARTITIONS 100;
/*
创建通证信息表
按照区块高度的hash分表
更新规则：solityNode为主，fullnode为辅
*/
CREATE TABLE `contract_asset_issue` (
  `trx_hash` varchar(64) NOT NULL DEFAULT '' COMMENT '交易hash',
  `block_id` bigint NOT NULL DEFAULT '0' COMMENT '区块ID',
  `token_name` varchar(300) NOT NULL COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '通证名称',
  `creator_address` varchar(300) NOT NULL COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT 'Token creator address',
  `total_supply` bigint NOT NULL DEFAULT '0' COMMENT '通证发行量',
  `trx_num` bigint NOT NULL DEFAULT '0' COMMENT '通证汇率分子，单位sun',
  `num` bigint NOT NULL DEFAULT '0' COMMENT '通证汇率分母',
  `start_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '开始时间',
  `end_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '结束时间',
  `decay_ratio` float4 NOT NULL DEFAULT '0' COMMENT 'decay_ratio',
  `vote_score` int(32) NOT NULL DEFAULT '0' COMMENT 'vote_score',
  `token_desc` text NULL COLLATE utf8mb4_unicode_ci COMMENT '通证说明',
  `url` varchar(800) NOT NULL COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '通证主页url',
  `frozen_info` text NULL COMMENT '冻结信息，json',
  `abbr` varchar(100) NOT NULL COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '通证缩写',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
  `modified_time` timestamp(6) NOT NULL  DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '记录更新时间',
  PRIMARY KEY (`trx_hash`,`block_id`),
  KEY `idx_asset_issue_hash` (`trx_hash`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
PARTITION BY HASH(block_id) PARTITIONS 100;

/*
候选人创建信息表  保存链外信息
*/
CREATE TABLE `wlcy_witness_create_info` (
  `address` varchar(200) NOT NULL DEFAULT '' COMMENT '候选人地址',
  `url` varchar(800) NOT NULL COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '候选人主页url',
  `create_time` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `modified_time` timestamp(6) NOT NULL  DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '记录更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

/*
基金会信息表  保存链外信息
*/
CREATE TABLE `wlcy_funds_info` (
  `id` int(32) NOT NULL DEFAULT '0' COMMENT 'ID',
  `address` varchar(300) NOT NULL DEFAULT '' COMMENT '基金会地址',
  `create_time` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `modified_time` timestamp(6) NOT NULL  DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '记录更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


/*
地域信息表  保存链外信息
*/
CREATE TABLE `wlcy_geo_info` (
  `ip` varchar(64) NOT NULL DEFAULT '' COMMENT 'IP',
  `city` varchar(200) NOT NULL DEFAULT '' COMMENT '城市',
  `country` varchar(200) NOT NULL DEFAULT '' COMMENT '国家',
  `lat` varchar(200) NOT NULL DEFAULT '0' COMMENT 'ip所属经纬度',
  `lng` varchar(200) NOT NULL DEFAULT '0' COMMENT 'ip所属经纬度',
  PRIMARY KEY (`ip`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

/*
超级代表信息表  保存链外信息
*/
CREATE TABLE `wlcy_sr_account` (
  `address` varchar(300) NOT NULL DEFAULT '' COMMENT '超级代表地址',
  `github_link` varchar(300) NOT NULL DEFAULT '' COMMENT '超级代表github',
  `create_time` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `modified_time` timestamp(6) NOT NULL  DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '记录更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

/*
trx请求信息表 
*/
CREATE TABLE `wlcy_trx_request` (
  `address` varchar(300) NOT NULL DEFAULT '' COMMENT '请求地址',
  `ip` varchar(300) NOT NULL DEFAULT '' COMMENT '请求对应的ip',
  `create_time` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `modified_time` timestamp(6) NOT NULL  DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '记录更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
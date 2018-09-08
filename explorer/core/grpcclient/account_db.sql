/*
账户信息
更新规则：通过分析交易中出现的账户，触发账户更新
*/
CREATE TABLE `account` (
  `account_name` varchar(300) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'Account name',
  `address` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL  DEFAULT '' COMMENT 'Base 58 encoding address',
  `balance` bigint NOT NULL DEFAULT '0' COMMENT 'TRX balance, in sun',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '账户创建时间',
  `latest_operation_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '账户最后操作时间',
  `is_witness` tinyint NOT NULL DEFAULT '0' COMMENT '是否为wintness; 0: 不是，1:是',
  `modified_time` timestamp NOT NULL  DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
  `fronze_amount` bigint NOT NULL DEFAULT '0' COMMENT '冻结金额, 投票权',
  UNIQUE KEY `uniq_account_address` (`address`,`create_time` DESC),
  KEY `idx_account_create_time` (`account_name`,`address`,`create_time` DESC),
  KEY `idx_acoount_addr_balance` (`address`,`balance` DESC,`create_time` DESC),
  KEY `idx_account_witness` (`create_time`,`is_witness`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
/*!50100 PARTITION BY RANGE (unix_timestamp(`create_time`))
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
  KEY `idx_account_asset_balance_id` (`address`,`token_name`, `balance` DESC),
  KEY `idx_account_asset_balance_addr` (`address`),
  KEY `idx_account_asset_balance_token_balance` (`token_name`,`balance` DESC),
  KEY `idx_account_asset_balance_creator_token_balance` (`token_name`,`creator_address`, `balance` DESC)
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
  `start_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '开始时间',
  `end_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '结束时间',
  `token_desc` text NULL COLLATE utf8mb4_unicode_ci COMMENT '通证说明',
  `url` varchar(800) NOT NULL COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '通证主页url',
  KEY `idx_asset_info_token_addr` (`token_name`, `creator_address`),
  KEY `idx_asset_info_token_trx` (`token_name`, `trx_num` DESC),
  KEY `idx_asset_info_token_total` (`token_name`, `trx_num` DESC, `total_supply` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


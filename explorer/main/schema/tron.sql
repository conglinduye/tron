-- MySQL dump 10.13  Distrib 5.7.23, for linux-glibc2.12 (x86_64)
--
-- Host: localhost    Database: tron
-- ------------------------------------------------------
-- Server version	5.7.23

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `account_asset_balance`
--

DROP TABLE IF EXISTS `account_asset_balance`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `account_asset_balance` (
  `address` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'Base 58 encoding address for the token owner',
  `asset_name` varchar(300) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '通证名称',
  `creator_address` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'Token creator address',
  `balance` bigint(20) NOT NULL DEFAULT '0' COMMENT '通证余额',
  KEY `idx_account_asset_balance_addr` (`address`),
  KEY `idx_asset_name` (`asset_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `account_vote_result`
--

DROP TABLE IF EXISTS `account_vote_result`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `account_vote_result` (
  `address` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'voter address',
  `to_address` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '投票接收人',
  `vote` bigint(20) NOT NULL DEFAULT '0' COMMENT '投票数',
  KEY `idx_account_vote_result_id` (`address`,`vote`),
  KEY `idx_account_vote_result_to` (`to_address`,`vote`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `asset_issue`
--

DROP TABLE IF EXISTS `asset_issue`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `asset_issue` (
  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
  `asset_name` varchar(200) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT 'asset_name',
  `asset_abbr` varchar(100) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT 'asset_abbr',
  `total_supply` bigint(20) NOT NULL DEFAULT '0' COMMENT '发行量',
  `frozen_supply` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '冻结量',
  `trx_num` bigint(20) NOT NULL DEFAULT '0' COMMENT 'price 分子',
  `num` bigint(20) NOT NULL DEFAULT '0' COMMENT 'price 分母',
  `participated` bigint(20) NOT NULL DEFAULT '0' COMMENT '已筹集资金',
  `start_time` bigint(20) NOT NULL DEFAULT '0',
  `end_time` bigint(20) NOT NULL DEFAULT '0',
  `order_num` bigint(20) NOT NULL DEFAULT '0',
  `vote_score` int(11) NOT NULL DEFAULT '0',
  `asset_desc` text COLLATE utf8mb4_bin NOT NULL COMMENT 'hexEncoding []byte, if you want display text content, hex decoding and cast to string, as there are data not coding with utf-8',
  `url` varchar(500) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `free_asset_net_limit` bigint(20) NOT NULL DEFAULT '0',
  `public_free_asset_net_limit` bigint(20) NOT NULL DEFAULT '0',
  `public_free_asset_net_usage` bigint(20) NOT NULL DEFAULT '0',
  `public_latest_free_net_time` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`owner_address`,`asset_name`),
  KEY `idx_owner_address` (`owner_address`),
  KEY `idx_asset_name` (`asset_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `blocks`
--

DROP TABLE IF EXISTS `blocks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
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
  UNIQUE KEY `uniq_blocks_id` (`block_id`),
  KEY `idx_blocks_create_time` (`create_time`),
  KEY `idx_block_witaddr` (`witness_address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
/*!50100 PARTITION BY HASH (`block_id`)
PARTITIONS 100 */;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `contract_account_create`
--

DROP TABLE IF EXISTS `contract_account_create`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `contract_account_create` (
  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `block_id` bigint(20) NOT NULL DEFAULT '0',
  `contract_type` int(11) NOT NULL DEFAULT '0',
  `create_time` bigint(20) NOT NULL DEFAULT '0',
  `expire_time` bigint(20) NOT NULL DEFAULT '0',
  `confirmed` tinyint(4) NOT NULL DEFAULT '0',
  `owner_address` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `account_address` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `account_type` tinyint(4) NOT NULL DEFAULT '0',
  PRIMARY KEY (`trx_hash`,`block_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
/*!50100 PARTITION BY HASH (`block_id`)
PARTITIONS 100 */;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `contract_account_update`
--

DROP TABLE IF EXISTS `contract_account_update`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `contract_account_update` (
  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
  `expire_time` bigint(20) NOT NULL DEFAULT '0',
  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
  `account_name` varchar(300) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT 'account_name',
  PRIMARY KEY (`trx_hash`,`block_id`),
  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
/*!50100 PARTITION BY HASH (`block_id`)
PARTITIONS 100 */;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `contract_asset_issue`
--

DROP TABLE IF EXISTS `contract_asset_issue`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `contract_asset_issue` (
  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
  `expire_time` bigint(20) NOT NULL DEFAULT '0',
  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
  `asset_name` varchar(200) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT 'asset_name',
  `asset_abbr` varchar(100) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT 'asset_abbr',
  `total_supply` bigint(20) NOT NULL DEFAULT '0' COMMENT '发行量',
  `frozen_supply` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '冻结量',
  `trx_num` bigint(20) NOT NULL DEFAULT '0' COMMENT 'price 分子',
  `num` bigint(20) NOT NULL DEFAULT '0' COMMENT 'price 分母',
  `start_time` bigint(20) NOT NULL DEFAULT '0',
  `end_time` bigint(20) NOT NULL DEFAULT '0',
  `order_num` bigint(20) NOT NULL DEFAULT '0',
  `vote_score` int(11) NOT NULL DEFAULT '0',
  `asset_desc` text COLLATE utf8mb4_bin NOT NULL,
  `url` varchar(500) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `free_asset_net_limit` bigint(20) NOT NULL DEFAULT '0',
  `public_free_asset_net_limit` bigint(20) NOT NULL DEFAULT '0',
  `public_free_asset_net_usage` bigint(20) NOT NULL DEFAULT '0',
  `public_latest_free_net_time` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`trx_hash`,`block_id`),
  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
/*!50100 PARTITION BY HASH (`block_id`)
PARTITIONS 100 */;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `contract_asset_transfer`
--

DROP TABLE IF EXISTS `contract_asset_transfer`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `contract_asset_transfer` (
  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
  `expire_time` bigint(20) NOT NULL DEFAULT '0',
  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
  `to_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '接收方地址',
  `amount` bigint(20) NOT NULL DEFAULT '0' COMMENT '转账金额。单位sun',
  `asset_name` varchar(200) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  PRIMARY KEY (`trx_hash`,`block_id`),
  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time`),
  KEY `idx_asset_name` (`asset_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
/*!50100 PARTITION BY HASH (`block_id`)
PARTITIONS 100 */;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `contract_create_smart`
--

DROP TABLE IF EXISTS `contract_create_smart`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `contract_create_smart` (
  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
  `expire_time` bigint(20) NOT NULL DEFAULT '0',
  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
  `contract_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `abi` text COLLATE utf8mb4_bin NOT NULL,
  `byte_code` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `call_value` bigint(20) NOT NULL DEFAULT '0',
  `consume_user_resource_percent` bigint(20) NOT NULL DEFAULT '0',
  `name` varchar(500) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  PRIMARY KEY (`trx_hash`,`block_id`),
  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
/*!50100 PARTITION BY HASH (`block_id`)
PARTITIONS 100 */;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `contract_freeze_balance`
--

DROP TABLE IF EXISTS `contract_freeze_balance`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `contract_freeze_balance` (
  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
  `expire_time` bigint(20) NOT NULL DEFAULT '0',
  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
  `frozen_balance` bigint(20) NOT NULL DEFAULT '0',
  `frozen_duration` bigint(20) NOT NULL DEFAULT '0',
  `resource` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`trx_hash`,`block_id`),
  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
/*!50100 PARTITION BY HASH (`block_id`)
PARTITIONS 100 */;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `contract_participate_asset`
--

DROP TABLE IF EXISTS `contract_participate_asset`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `contract_participate_asset` (
  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
  `expire_time` bigint(20) NOT NULL DEFAULT '0',
  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
  `to_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
  `asset_name` varchar(200) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT 'token名称',
  `amount` bigint(20) NOT NULL DEFAULT '0' COMMENT '转账金额。单位sun',
  PRIMARY KEY (`trx_hash`,`block_id`),
  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time`),
  KEY `idx_asset_name` (`asset_name`),
  KEY `idx_to_address` (`to_address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
/*!50100 PARTITION BY HASH (`block_id`)
PARTITIONS 100 */;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `contract_set_account_id`
--

DROP TABLE IF EXISTS `contract_set_account_id`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `contract_set_account_id` (
  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
  `expire_time` bigint(20) NOT NULL DEFAULT '0',
  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
  `account_id` varchar(500) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT 'account_id',
  PRIMARY KEY (`trx_hash`,`block_id`),
  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
/*!50100 PARTITION BY HASH (`block_id`)
PARTITIONS 100 */;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `contract_transfer`
--

DROP TABLE IF EXISTS `contract_transfer`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `contract_transfer` (
  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
  `expire_time` bigint(20) NOT NULL DEFAULT '0',
  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
  `to_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '接收方地址',
  `amount` bigint(20) NOT NULL DEFAULT '0' COMMENT '转账金额。单位sun',
  `asset_name` varchar(200) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  PRIMARY KEY (`trx_hash`,`block_id`),
  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time`),
  KEY `idx_trx_owner` (`owner_address`),
  KEY `idx_trx_to_addr` (`to_address`),
  KEY `idx_trx_asset_name` (`asset_name`),
  KEY `idx_trx_transfer_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
/*!50100 PARTITION BY HASH (`block_id`)
PARTITIONS 100 */;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `contract_trigger_smart`
--

DROP TABLE IF EXISTS `contract_trigger_smart`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `contract_trigger_smart` (
  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
  `expire_time` bigint(20) NOT NULL DEFAULT '0',
  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
  `contract_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `call_value` bigint(20) NOT NULL DEFAULT '0',
  `call_data` text COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`trx_hash`,`block_id`),
  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
/*!50100 PARTITION BY HASH (`block_id`)
PARTITIONS 100 */;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `contract_unfreeze_asset`
--

DROP TABLE IF EXISTS `contract_unfreeze_asset`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `contract_unfreeze_asset` (
  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
  `expire_time` bigint(20) NOT NULL DEFAULT '0',
  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
  PRIMARY KEY (`trx_hash`,`block_id`),
  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
/*!50100 PARTITION BY HASH (`block_id`)
PARTITIONS 100 */;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `contract_unfreeze_balance`
--

DROP TABLE IF EXISTS `contract_unfreeze_balance`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `contract_unfreeze_balance` (
  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
  `expire_time` bigint(20) NOT NULL DEFAULT '0',
  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
  `resource` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`trx_hash`,`block_id`),
  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
/*!50100 PARTITION BY HASH (`block_id`)
PARTITIONS 100 */;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `contract_update_asset`
--

DROP TABLE IF EXISTS `contract_update_asset`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `contract_update_asset` (
  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
  `expire_time` bigint(20) NOT NULL DEFAULT '0',
  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
  `asset_desc` text COLLATE utf8mb4_bin NOT NULL,
  `url` text COLLATE utf8mb4_bin NOT NULL,
  `new_limit` bigint(20) NOT NULL DEFAULT '0',
  `new_public_limit` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`trx_hash`,`block_id`),
  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
/*!50100 PARTITION BY HASH (`block_id`)
PARTITIONS 100 */;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `contract_update_setting`
--

DROP TABLE IF EXISTS `contract_update_setting`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `contract_update_setting` (
  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
  `expire_time` bigint(20) NOT NULL DEFAULT '0',
  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
  `contract_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `consume_user_resource_percent` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`trx_hash`,`block_id`),
  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
/*!50100 PARTITION BY HASH (`block_id`)
PARTITIONS 100 */;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `contract_vote_asset`
--

DROP TABLE IF EXISTS `contract_vote_asset`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `contract_vote_asset` (
  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
  `expire_time` bigint(20) NOT NULL DEFAULT '0',
  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
  `vote_address` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '投票地址',
  `support` tinyint(4) NOT NULL DEFAULT '0' COMMENT 'support',
  `count` bigint(20) NOT NULL DEFAULT '0' COMMENT 'count',
  PRIMARY KEY (`trx_hash`,`block_id`),
  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
/*!50100 PARTITION BY HASH (`block_id`)
PARTITIONS 100 */;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `contract_vote_witness`
--

DROP TABLE IF EXISTS `contract_vote_witness`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `contract_vote_witness` (
  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
  `expire_time` bigint(20) NOT NULL DEFAULT '0',
  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
  `votes` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '投票详情，JSON',
  `support` tinyint(4) NOT NULL DEFAULT '0' COMMENT 'support 0:false, 1:true',
  PRIMARY KEY (`trx_hash`,`block_id`),
  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
/*!50100 PARTITION BY HASH (`block_id`)
PARTITIONS 100 */;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `contract_witness_create`
--

DROP TABLE IF EXISTS `contract_witness_create`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `contract_witness_create` (
  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
  `expire_time` bigint(20) NOT NULL DEFAULT '0',
  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
  `url` varchar(500) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT 'url',
  PRIMARY KEY (`trx_hash`,`block_id`),
  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
/*!50100 PARTITION BY HASH (`block_id`)
PARTITIONS 100 */;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `contract_witness_update`
--

DROP TABLE IF EXISTS `contract_witness_update`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `contract_witness_update` (
  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID',
  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\\\\nAccountCreateContract = 0;\\\\r\\\\nTransferContract = 1;\\\\r\\\\nTransferAssetContract = 2;\\\\r\\\\nVoteAssetContract = 3;\\\\r\\\\nVoteWitnessContract = 4;\\\\r\\\\nWitnessCreateContract = 5;\\\\r\\\\nAssetIssueContract = 6;\\\\r\\\\nWitnessUpdateContract = 8;\\\\r\\\\nParticipateAssetIssueContract = 9;\\\\r\\\\nAccountUpdateContract = 10;\\\\r\\\\nFreezeBalanceContract = 11;\\\\r\\\\nUnfreezeBalanceContract = 12;\\\\r\\\\nWithdrawBalanceContract = 13;\\\\r\\\\nUnfreezeAssetContract = 14;\\\\r\\\\nUpdateAssetContract = 15;\\\\r\\\\nProposalCreateContract = 16;\\\\r\\\\nProposalApproveContract = 17;\\\\r\\\\nProposalDeleteContract = 18;\\\\r\\\\nSetAccountIdContract = 19;\\\\r\\\\nCustomContract = 20;\\\\r\\\\n// BuyStorageContract = 21;\\\\r\\\\n// BuyStorageBytesContract = 22;\\\\r\\\\n// SellStorageContract = 23;\\\\r\\\\nCreateSmartContract = 30;\\\\r\\\\nTriggerSmartContract = 31;\\\\r\\\\nGetContract = 32;\\\\r\\\\nUpdateSettingContract = 33;\\\\r\\\\nExchangeCreateContract = 41;\\\\r\\\\nExchangeInjectContract = 42;\\\\r\\\\nExchangeWithdrawContract = 43;\\\\r\\\\nExchangeTransactionContract = 44;',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
  `expire_time` bigint(20) NOT NULL DEFAULT '0',
  `confirmed` tinyint(4) NOT NULL DEFAULT '0' COMMENT '确认状态。0 未确认。1 已确认',
  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
  `update_url` varchar(500) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  PRIMARY KEY (`trx_hash`,`block_id`),
  KEY `idx_trx_transfe_hash_create_time` (`block_id`,`trx_hash`,`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
/*!50100 PARTITION BY HASH (`block_id`)
PARTITIONS 100 */;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nodes`
--

DROP TABLE IF EXISTS `nodes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `nodes` (
  `node_host` varchar(300) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT 'host',
  `node_port` int(11) NOT NULL DEFAULT '0' COMMENT 'port',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`node_host`,`node_port`),
  KEY `idx_node_host` (`node_host`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `transactions`
--

DROP TABLE IF EXISTS `transactions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `transactions` (
  `trx_hash` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易hash',
  `block_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块ID，高度',
  `contract_type` int(8) NOT NULL DEFAULT '0' COMMENT '交易类型\nAccountCreateContract = 0;\r\nTransferContract = 1;\r\nTransferAssetContract = 2;\r\nVoteAssetContract = 3;\r\nVoteWitnessContract = 4;\r\nWitnessCreateContract = 5;\r\nAssetIssueContract = 6;\r\nWitnessUpdateContract = 8;\r\nParticipateAssetIssueContract = 9;\r\nAccountUpdateContract = 10;\r\nFreezeBalanceContract = 11;\r\nUnfreezeBalanceContract = 12;\r\nWithdrawBalanceContract = 13;\r\nUnfreezeAssetContract = 14;\r\nUpdateAssetContract = 15;\r\nProposalCreateContract = 16;\r\nProposalApproveContract = 17;\r\nProposalDeleteContract = 18;\r\nSetAccountIdContract = 19;\r\nCustomContract = 20;\r\n// BuyStorageContract = 21;\r\n// BuyStorageBytesContract = 22;\r\n// SellStorageContract = 23;\r\nCreateSmartContract = 30;\r\nTriggerSmartContract = 31;\r\nGetContract = 32;\r\nUpdateSettingContract = 33;\r\nExchangeCreateContract = 41;\r\nExchangeInjectContract = 42;\r\nExchangeWithdrawContract = 43;\r\nExchangeTransactionContract = 44;',
  `contract_data` varchar(2000) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易内容数据,原始数据byte hex encoding',
  `result_data` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '交易结果对象byte hex encoding',
  `owner_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '发起方地址',
  `to_address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '接收方地址',
  `fee` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易花费 单位 sun',
  `confirmed` tinyint(4) NOT NULL DEFAULT '1',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易创建时间',
  `expire_time` bigint(20) DEFAULT '0',
  `modified_time` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '记录更新时间',
  `real_timestamp` bigint(20) NOT NULL DEFAULT '0' COMMENT 'transaction timestamp',
  PRIMARY KEY (`trx_hash`,`block_id`),
  KEY `idx_transactions_hash_create_time` (`block_id`,`trx_hash`,`create_time`),
  KEY `idx_transactions_create_time` (`create_time`),
  KEY `idx_trx_owner` (`owner_address`),
  KEY `idx_trx_toaddr` (`to_address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
/*!50100 PARTITION BY HASH (`block_id`)
PARTITIONS 100 */;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `tron_account`
--

DROP TABLE IF EXISTS `tron_account`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tron_account` (
  `account_name` varchar(300) COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT 'Account name',
  `address` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'Base 58 encoding address',
  `balance` bigint(20) NOT NULL DEFAULT '0' COMMENT 'TRX balance, in sun',
  `create_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '账户创建时间',
  `latest_operation_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '账户最后操作时间',
  `asset_issue_name` varchar(100) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `is_witness` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否为wintness; 0: 不是，1:是',
  `allowance` bigint(20) NOT NULL DEFAULT '0',
  `latest_withdraw_time` bigint(20) NOT NULL DEFAULT '0',
  `latest_consume_time` bigint(20) NOT NULL DEFAULT '0',
  `latest_consume_free_time` bigint(20) NOT NULL DEFAULT '0',
  `frozen` text COLLATE utf8mb4_bin NOT NULL COMMENT '冻结信息',
  `votes` text COLLATE utf8mb4_bin NOT NULL,
  `net_usage` bigint(20) NOT NULL DEFAULT '0' COMMENT 'bandwidth, get from frozen',
  `free_net_limit` bigint(20) NOT NULL DEFAULT '0',
  `net_used` bigint(20) NOT NULL DEFAULT '0',
  `net_limit` bigint(20) NOT NULL DEFAULT '0',
  `total_net_limit` bigint(20) NOT NULL DEFAULT '0',
  `total_net_weight` bigint(20) NOT NULL DEFAULT '0',
  `asset_net_used` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `asset_net_limit` text COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`address`),
  KEY `idx_tron_account_create_time` (`create_time`),
  KEY `idx_account_name` (`account_name`),
  KEY `idx_account_address` (`address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `witness`
--

DROP TABLE IF EXISTS `witness`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `witness` (
  `address` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '地址',
  `vote_count` bigint(20) NOT NULL DEFAULT '0' COMMENT '得票数',
  `public_key` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `url` varchar(500) COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `total_produced` bigint(20) NOT NULL DEFAULT '0' COMMENT '生产块数',
  `total_missed` bigint(20) NOT NULL DEFAULT '0' COMMENT '丢失块数',
  `latest_block_num` bigint(20) NOT NULL DEFAULT '0',
  `latest_slot_num` bigint(20) NOT NULL DEFAULT '0',
  `is_job` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否为超级候选人 0:false, 1:true',
  PRIMARY KEY (`address`),
  UNIQUE KEY `address_UNIQUE` (`address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `wlcy_asset_info`
--

DROP TABLE IF EXISTS `wlcy_asset_info`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `wlcy_asset_info` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `address` varchar(200) NOT NULL,
  `token_name` varchar(500) COLLATE utf8mb4_bin  NOT NULL,
  `token_id` varchar(500) COLLATE utf8mb4_bin DEFAULT NULL,
  `brief` text,
  `website` varchar(500) COLLATE utf8mb4_bin DEFAULT NULL,
  `white_paper` varchar(500) COLLATE utf8mb4_bin DEFAULT NULL,
  `github` varchar(500) COLLATE utf8mb4_bin DEFAULT NULL,
  `country` varchar(50) DEFAULT NULL,
  `credit` int(11) DEFAULT '0',
  `reddit` varchar(500) DEFAULT NULL,
  `twitter` varchar(500) DEFAULT NULL,
  `facebook` varchar(500) DEFAULT NULL,
  `telegram` varchar(500) DEFAULT NULL,
  `steam` varchar(500) DEFAULT NULL,
  `medium` varchar(500) DEFAULT NULL,
  `webchat` varchar(500) DEFAULT NULL,
  `weibo` varchar(500) DEFAULT NULL,
  `review` int(11) DEFAULT '0',
  `status` int(11) DEFAULT '0',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `address` (`address`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `wlcy_asset_logo`
--

DROP TABLE IF EXISTS `wlcy_asset_logo`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `wlcy_asset_logo` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `address` varchar(100) NOT NULL,
  `logo_url` varchar(400) DEFAULT NULL,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `address` (`address`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `wlcy_funds_info`
--

DROP TABLE IF EXISTS `wlcy_funds_info`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `wlcy_funds_info` (
  `id` int(32) NOT NULL DEFAULT '0' COMMENT 'ID',
  `address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '基金会地址',
  `create_time` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `modified_time` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '记录更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `wlcy_geo_info`
--

DROP TABLE IF EXISTS `wlcy_geo_info`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `wlcy_geo_info` (
  `ip` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'IP',
  `city` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '城市',
  `country` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '国家',
  `lat` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '0' COMMENT 'ip所属经纬度',
  `lng` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '0' COMMENT 'ip所属经纬度',
  PRIMARY KEY (`ip`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `wlcy_opreate_log`
--

DROP TABLE IF EXISTS `wlcy_opreate_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `wlcy_opreate_log` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `address` varchar(100) DEFAULT NULL,
  `operator` varchar(50) DEFAULT NULL,
  `pic_url` varchar(500) DEFAULT NULL,
  `record` text,
  `extra` text,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `wlcy_sr_account`
--

DROP TABLE IF EXISTS `wlcy_sr_account`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `wlcy_sr_account` (
  `address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '超级代表地址',
  `github_link` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '超级代表github',
  `create_time` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `modified_time` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '记录更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `wlcy_statistics`
--

DROP TABLE IF EXISTS `wlcy_statistics`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `wlcy_statistics` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `date` bigint(20) NOT NULL DEFAULT '0' COMMENT '时间',
  `avg_block_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块上链平均时间',
  `avg_block_size` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块平均大小',
  `new_block_seen` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块增长数',
  `new_transaction_seen` bigint(20) NOT NULL DEFAULT '0' COMMENT '交易增长数',
  `new_address_seen` bigint(20) NOT NULL DEFAULT '0' COMMENT '地址增长数',
  `total_block_count` bigint(20) NOT NULL DEFAULT '0' COMMENT '总区块数',
  `total_transaction` bigint(20) NOT NULL DEFAULT '0' COMMENT '总交易数',
  `total_address` bigint(20) NOT NULL DEFAULT '0' COMMENT '总地址数',
  `blockchain_size` bigint(20) NOT NULL DEFAULT '0' COMMENT '区块链总大小',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=88 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `wlcy_trx_request`
--

DROP TABLE IF EXISTS `wlcy_trx_request`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `wlcy_trx_request` (
  `address` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '请求地址',
  `ip` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '请求对应的ip',
  `create_time` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `modified_time` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '记录更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `wlcy_witness_create_info`
--

DROP TABLE IF EXISTS `wlcy_witness_create_info`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `wlcy_witness_create_info` (
  `address` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '候选人地址',
  `url` varchar(800) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '候选人主页url',
  `create_time` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `modified_time` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '记录更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2018-09-20 13:53:24


DROP TABLE IF EXISTS `transactions_index`;
CREATE TABLE `transactions_index` (
  `start_pos` bigint(20) NOT NULL DEFAULT '0',
  `block_id` bigint(20) NOT NULL DEFAULT '0',
  `inner_offset` bigint(20) NOT NULL DEFAULT '0',
  `total_record` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`start_pos`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `contract_transfer_index`;
CREATE TABLE `contract_transfer_index` (
  `start_pos` bigint(20) NOT NULL DEFAULT '0',
  `block_id` bigint(20) NOT NULL DEFAULT '0',
  `inner_offset` bigint(20) NOT NULL DEFAULT '0',
  `total_record` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`start_pos`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

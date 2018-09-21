package service

import (
	"fmt"
	"strings"

	"github.com/wlcy/tron/explorer/lib/log"

	"github.com/wlcy/tron/explorer/web/buffer"

	"github.com/wlcy/tron/explorer/lib/mysql"

	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
)

//QueryBlocksBuffer  从缓存中查询
func QueryBlocksBuffer(req *entity.Blocks) (*entity.BlocksResp, error) {
	var err error
	blockResp := &entity.BlocksResp{}
	blocks := make([]*entity.BlockInfo, 0)
	blockBuffer := buffer.GetBlockBuffer()
	blockResp.Total = blockBuffer.GetMaxBlockID()
	if req.Number != "" {
		block := blockBuffer.GetBlock(mysql.ConvertStringToInt64(req.Number, 0))
		if block == nil {
			log.Debugf("get blocks data in buffer, get them from db instead")
			return QueryBlocks(req)
		}

		blocks = append(blocks, block)
		blockResp.Total = int64(len(blocks))
	} else if req.Producer != "" {
		return QueryBlocks(req)
	} else {
		blocks, err = blockBuffer.GetBlocks(-1, req.Start, req.Limit)
		if err != nil || blocks == nil {
			log.Debugf("get blocks data in buffer, get them from db instead")
			return QueryBlocks(req)
		}
	}
	blockResp.Data = blocks
	return blockResp, nil
}

//QueryBlocks 条件查询  	//?sort=-number&limit=1&count=true&number=2135998
func QueryBlocks(req *entity.Blocks) (*entity.BlocksResp, error) {
	var filterSQL, sortSQL, pageSQL, sortTemp string
	mutiFilter := false

	strSQL := fmt.Sprintf(`
			select block_id,block_hash,block_size,create_time,
			transaction_num,
			tx_trie_hash,parent_hash,witness_address,confirmed
			from tron.blocks
			where 1=1 `)

	//按传入条件拼接sql，很容易错误，需要注意
	if req.Number != "" {
		filterSQL = fmt.Sprintf(" and block_id=%v", req.Number)
	}
	if req.Producer != "" {
		filterSQL = fmt.Sprintf(" and witness_address='%v'", req.Producer)
	}
	sortTemp = "order by"
	if strings.Index(req.Order, "timestamp") > 0 {
		sortTemp = fmt.Sprintf("%v create_time", sortTemp)
		if strings.Index(req.Order, "-") == 0 {
			sortTemp = fmt.Sprintf("%v desc", sortTemp)
		}
		mutiFilter = true
	}

	if strings.Index(req.Sort, "number") > 0 {
		if mutiFilter {
			sortTemp = fmt.Sprintf("%v ,", sortTemp)
		}
		sortTemp = fmt.Sprintf("%v block_id", sortTemp)
		if strings.Index(req.Sort, "-") == 0 {
			sortTemp = fmt.Sprintf("%v desc", sortTemp)
		}
	}
	if sortTemp != "order by" {
		sortSQL = sortTemp
	}
	pageSQL = fmt.Sprintf("limit %v, %v", req.Start, req.Limit)

	return module.QueryBlocksRealize(strSQL, filterSQL, sortSQL, pageSQL)
}

//QueryBlock 精确查询  	//number=2135998
func QueryBlock(req *entity.Blocks) (*entity.BlockInfo, error) {
	var filterSQL string
	strSQL := fmt.Sprintf(`
			select block_id,block_hash,block_size,create_time,
			transaction_num,
			tx_trie_hash,parent_hash,witness_address,confirmed
			from tron.blocks
			where 1=1 `)

	//按传入条件拼接sql，很容易错误，需要注意
	if req.Number != "" {
		filterSQL = fmt.Sprintf(" and block_id=%v", req.Number)
	}
	return module.QueryBlockRealize(strSQL, filterSQL)
}

//QueryBlockBuffer 精确查询  	//number=2135998
func QueryBlockBuffer(req *entity.Blocks) (*entity.BlockInfo, error) {
	block := &entity.BlockInfo{}
	blockBuffer := buffer.GetBlockBuffer()
	if req.Number != "" {
		block = blockBuffer.GetBlock(mysql.ConvertStringToInt64(req.Number, 0))
	}
	return block, nil
}

//QueryBlockLatestBuffer 获取最新块
func QueryBlockLatestBuffer() (*entity.BlockInfo, error) {
	block := &entity.BlockInfo{}
	blockBuffer := buffer.GetBlockBuffer()
	maxNumber := blockBuffer.GetMaxBlockID()
	if maxNumber > 0 {
		block = blockBuffer.GetBlock(maxNumber)
	}
	return block, nil
}

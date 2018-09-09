package service

import (
	"fmt"
	"strings"

	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
)

//QueryBlocks 条件查询  	//?sort=-number&limit=1&count=true&number=2135998
func QueryBlocks(req *entity.Blocks) (*entity.BlocksResp, error) {
	var filterSQL, sortSQL, pageSQL, sortTemp string
	mutiFilter := false

	strSQL := fmt.Sprintf(`
			select block_id,block_hash,block_size,create_time,tx_trie_hash,parent_hash,witness_address,confirmed
			from tron.blocks
			where 1=1 `)

	//按传入条件拼接sql，很容易错误，需要注意
	if req.Number != "" {
		filterSQL = fmt.Sprintf(" and block_id=%v", req.Number)
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
	if req.Limit != "" && req.Start != "" {
		pageSQL = fmt.Sprintf("limit %v, %v", req.Start, req.Limit)
	}
	return module.QueryBlocksRealize(strSQL, filterSQL, sortSQL, pageSQL)
}

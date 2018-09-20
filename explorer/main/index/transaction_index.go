package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// TrxIndex transactions 表索引
type TrxIndex struct {
	StartPosition int64 // 外部定位用的偏移量，起点为 limit 0, 1; index start from 0, 所以 StartPosition== 100000, 代表第 100001 条记录
	BlockID       int64 // block_id
	Offset        int64 // 定位用启始Offset (sql limit 的第一各参数)
	Count         int64 // 记录总数，只有 index == 0 的记录记录
}

var step int64 = 100000

func getIndex() []*TrxIndex {
	ret, err := loadIdx()
	var needStore bool
	if nil != err || 0 == len(ret) {
		needStore = true
		// ret, err = genTransactionIndex()
		ret, err = genTransactionIndex2(0, 0, 0)
	}

	fmt.Printf("needStore:%v, ret:%v, err:%v\n", needStore, len(ret), err)
	if nil != err || 0 == len(ret) {
		return nil
	}

	if needStore {
		storeIdx(ret)
	}
	totalTrn = ret[0].Count
	return ret
}

func genTransactionIndex() ([]*TrxIndex, error) {

	strSQL1 := `select trx_hash, block_id from transactions where block_id >= ? order by block_id asc limit ?, 1`
	strSQL2 := `select trx_hash, block_id from transactions where block_id = ?`

	dbb := getMysqlDB()

	txn, err := dbb.Begin()
	if nil != err {
		fmt.Printf("gen txn failed:%v\n", err)
		return nil, err
	}

	stmt1, err := txn.Prepare(strSQL1)
	if nil != err {
		fmt.Printf("prepare SQL (%v) failed:%v\n", strSQL1, err)
		return nil, err
	}
	defer stmt1.Close()

	stmt2, err := txn.Prepare(strSQL2)
	if nil != err {
		fmt.Printf("prepare SLQ (%v) failed:%v\n", strSQL2, err)
		return nil, err
	}
	defer stmt2.Close()

	index := make([]*TrxIndex, 0, 10000)
	index = append(index, &TrxIndex{
		StartPosition: 0,
		BlockID:       0,
		Offset:        0,
	})
	var blockID, offset, pos, round int64
	round = 1
	var trxHash, trxHash2 string
	for {
		pos = 0

		fmt.Printf("round:%-3v, blockID:%v, related offset:%v, absolute offset:%v\n", round, blockID, offset, round*step)
		row := stmt1.QueryRow(blockID, offset)
		if nil != err || nil == row {
			fmt.Printf("sql1 failed:%v\n", err)
			break
		}

		err = row.Scan(&trxHash, &blockID)
		if nil != err {
			fmt.Printf("sql1 scan failed:%v\n", err)
			break
		}

		idx := &TrxIndex{}
		idx.BlockID = blockID
		idx.Offset = step
		idx.StartPosition = step * round

		rows, err := stmt2.Query(blockID)
		if nil != err {
			fmt.Printf("sql2 failed:%v\n", err)
			break
		}

		for rows.Next() {
			err = rows.Scan(&trxHash2, &blockID)
			if nil != err {
				fmt.Printf("sql2 scan failed:%v\n", err)
				break
			}
			pos++
			if trxHash == trxHash2 {
				idx.Offset += pos - 2
			}
		}
		rows.Close()
		index = append(index, idx)

		offset = idx.Offset + 1 // 计数从 step+1开始， step 属于上一个槽位
		round++
	}

	txn.Commit()
	if nil != err {
		fmt.Printf("commit txn failed:%v\n", err)
		// return err
	}

	printIndex(index)

	return index, nil
}

func storeIdx(index []*TrxIndex) error {
	data, err := json.Marshal(index)
	if nil != err {
		fmt.Printf("gen index data failed:%v\n", err)
		return err
	}

	idxFile, err := os.Create(fmt.Sprintf("%v", *gIndexFile))

	if nil != err {
		fmt.Printf("create file %v failed:%v\n", *gIndexFile, err)
		return err
	}
	defer idxFile.Close()

	w := bufio.NewWriter(idxFile)
	defer w.Flush()

	n, err := w.Write(data)
	if n != len(data) || nil != err {
		fmt.Printf("write file err:%v, byte need to write:%v, actual write:%v\n", err, len(data), n)
		return err
	}

	return nil
}

func loadIdx() ([]*TrxIndex, error) {
	idxFile, err := os.Open(*gIndexFile)

	if nil != err {
		fmt.Printf("open file %v failed:%v\n", *gIndexFile, err)
		return nil, err
	}
	defer idxFile.Close()

	data, err := ioutil.ReadAll(idxFile)
	if nil != err || 0 == len(data) {
		fmt.Printf("read index failed:%v or index is empty!\n", err)
		return nil, err
	}

	index := make([]*TrxIndex, 0)
	err = json.Unmarshal(data, &index)
	if nil != err || 0 == len(index) {
		fmt.Printf("parse index context failed:%v or index is empty!\n", err)
		return nil, err
	}
	return index, nil
}

var totalTrn int64 = 1430347

func searchIdxIF() {

	index := getIndex()
	index, _ = updateIndex(index)

	var offset, count int64
	for {
		fmt.Printf("\n============\ninput offset and count (0 0 quit):")
		n, err := fmt.Scanf("%v %v", &offset, &count)
		if 2 < n || nil != err {
			break
		}
		fmt.Println(searchIdx(offset, count, index))
		fmt.Printf("\n-----------\n")

	}

	fmt.Printf("bye!!!\n")

}

func searchIdx(offset, count int64, index []*TrxIndex) string {
	if offset >= totalTrn {
		fmt.Printf("invalid offset:%v, total count:%v, index range:[0, %v]\n", offset, totalTrn, totalTrn-1)
		return ""
	}

	ascOffset := totalTrn - offset - 1
	ascOffsetIdx := ascOffset / step
	ascInnerOffsetIdx := ascOffset % step

	if ascOffsetIdx >= int64(len(index)) {
		fmt.Printf("invalid offset:%v, err index:%v\n", offset, ascOffset)
		return ""
	}

	fmt.Printf("offset:%v, ascOffset:%v, ascOffsetIdx:%v, ascInnerOffsetIdx:%v\n", offset, ascOffset, ascOffsetIdx, ascInnerOffsetIdx)

	idx := index[ascOffsetIdx]
	return fmt.Sprintf("select trx_hash, block_id from transactions where block_id >= '%v' order by block_id asc limit %v, %v;\n",
		idx.BlockID, idx.Offset+ascInnerOffsetIdx, count)

}

func updateIndex(index []*TrxIndex) ([]*TrxIndex, bool) {
	db := getMysqlDB()

	txn, err := db.Begin()
	if nil != err {
		fmt.Printf("Create db txn failed:%v\n", txn)
		return index, false
	}

	row := txn.QueryRow("select count(*) from transactions")

	var total int64
	err = row.Scan(&total)
	if nil != err {
		fmt.Printf("scan transaction count failed:%v\n", err)
		return index, false
	}

	totalTrn = total

	var maxIdx = total / step
	round := int64(len(index))

	if int64(len(index)) > maxIdx {
		fmt.Printf("updateIndex: current max transaction:%v, max index:%v, current index length:%v, do not need update\n", total, maxIdx, round)
		return index, false
	}
	fmt.Printf("updateIndex: current max transaction:%v, max index:%v, current index length:%v, updating index ......\n", total, maxIdx, round)

	var newIndex []*TrxIndex
	if 0 < round {
		idx := index[round-1]
		newIndex, err = genTransactionIndex2(round, idx.BlockID, idx.Offset)
	} else {
		newIndex, err = genTransactionIndex2(0, 0, 0)
	}
	if nil != err || 0 == len(newIndex) {
		return index, false
	}

	index = append(index, newIndex...)
	if len(index) > 1 {
		index[0].Count = total
		index[1].Count = step
	}
	storeIdx(index)

	return index, true

}

func genTransactionIndex2(round, blockID, offset int64) ([]*TrxIndex, error) {

	fmt.Printf("genTransactionIndex2, start round:%v, block_id:%v, offset:%v\n", round, blockID, offset)

	strSQL1 := `select trx_hash, block_id from transactions where block_id >= ? order by block_id asc limit ?, 1`
	strSQL2 := `select trx_hash, block_id from transactions where block_id = ?`

	dbb := getMysqlDB()

	txn, err := dbb.Begin()
	if nil != err {
		fmt.Printf("gen txn failed:%v\n", err)
		return nil, err
	}

	stmt1, err := txn.Prepare(strSQL1)
	if nil != err {
		fmt.Printf("prepare SQL (%v) failed:%v\n", strSQL1, err)
		return nil, err
	}
	defer stmt1.Close()

	stmt2, err := txn.Prepare(strSQL2)
	if nil != err {
		fmt.Printf("prepare SLQ (%v) failed:%v\n", strSQL2, err)
		return nil, err
	}
	defer stmt2.Close()

	index := make([]*TrxIndex, 0, 10)

	var pos int64
	if 0 == round {
		// add round 0 index
		index = append(index, &TrxIndex{
			StartPosition: 0,
			BlockID:       0,
			Offset:        0,
		})
		round = 1 // the step * round offset start position
	}
	var trxHash, trxHash2 string
	offset += step
	for {
		pos = 0

		fmt.Printf("round:%-3v, blockID:%v, related offset:%v, absolute offset:%v\n", round, blockID, offset, round*step)
		row := stmt1.QueryRow(blockID, offset)
		if nil != err || nil == row {
			fmt.Printf("sql1 failed:%v\n", err)
			break
		}

		err = row.Scan(&trxHash, &blockID)
		if nil != err {
			fmt.Printf("sql1 scan failed:%v\n", err)
			break
		}

		idx := &TrxIndex{}
		idx.BlockID = blockID
		idx.Offset = 0
		idx.StartPosition = step * round

		rows, err := stmt2.Query(blockID)
		if nil != err {
			fmt.Printf("sql2 failed:%v\n", err)
			break
		}

		for rows.Next() {
			err = rows.Scan(&trxHash2, &blockID)
			if nil != err {
				fmt.Printf("sql2 scan failed:%v\n", err)
				break
			}
			if trxHash == trxHash2 {
				idx.Offset = pos
				break
			}
			pos++
		}
		rows.Close()
		index = append(index, idx)

		offset = idx.Offset + step // 获取 next round 相对偏移 step
		round++
	}

	txn.Commit()
	if nil != err {
		fmt.Printf("commit txn failed:%v\n", err)
		// return err
	}

	printIndex(index)

	return index, nil
}

func printIndex(index []*TrxIndex) {
	fmt.Printf("total index:%v\n", len(index))
	for id, idx := range index {
		fmt.Printf("%v-->%#v\n", id, idx)
		fmt.Printf("position:%v --> select trx_hash, block_id from transactions where block_id >= '%v' order by block_id asc limit %v, 1;\n", idx.StartPosition, idx.BlockID, idx.Offset)
	}
}

func storeIdxToDB(index []*TrxIndex) {
	ddb := getMysqlDB()

	txn, err := ddb.Begin()
	if nil != err {
		return
	}

	strSQL1 := "insert into transactions_index (start_pos, block_id, inner_offset, total_record) values (?, ?, ?, ?)"

	stmt1, err := txn.Prepare(strSQL1)
	if nil != err {
		fmt.Printf("prepare SQL (%v) failed:%v\n", strSQL1, err)
		return
	}
	defer stmt1.Close()

	strSQL2 := "delete from transactions_index where start_pos >=0"
	stmt2, err := txn.Prepare(strSQL2)
	if nil != err {
		fmt.Printf("prepare SLQ (%v) failed:%v\n", strSQL2, err)
	}
	defer stmt2.Close()

	stmt2.Exec()

	if len(index) > 1 {
		index[0].Count = totalTrn
		index[1].Count = step
	}

	for _, idx := range index {
		_, err := stmt1.Exec(idx.StartPosition, idx.BlockID, idx.Offset, idx.Count)
		if nil != err {
			fmt.Printf("insert transactions_index %#v failed:%v\n", idx, err)
		}
	}

	txn.Commit()

}

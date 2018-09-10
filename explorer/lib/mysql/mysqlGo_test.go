package mysql

import (
	"fmt"
	"testing"
)

func TestSelect(t *testing.T) {
	var db *TronDB

	db, err := OpenDB("mysql", "budev:budev@tcp(127.0.0.1:3306)/tron?charset=utf8")

	if err != nil {
		fmt.Println("opendb fail")
		return
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Println(err.Error())
		return
	}
	res, err := db.Select(string("select * from blocks"))

	if err != nil {
		fmt.Println("select failed")
		return
	}

	for k, v := range res.Colmns() {
		fmt.Println(k, v)
	}

	for res.NextT() {
		for k, _ := range res.Colmns() {
			fmt.Println(k, ":", res.GetField(k))
		}
	}
}

func TestSelectIgnoreColumnCase(t *testing.T) {
	var db *TronDB

	db, err := OpenDB("mysql", "root:root@tcp(127.0.0.1:3306)/tron?charset=utf8")

	if err != nil {
		fmt.Println("opendb fail")
		return
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Println(err.Error())
		return
	}
	res, err := db.Select(string("select * from blocks"))

	if err != nil {
		fmt.Println("select failed")
		return
	}

	for res.NextT() {
		id := res.GetField("appID")
		name := res.GetField("name")
		desc := res.GetField("description")
		typeid := res.GetField("AppTypeID")
		platform := res.GetField("platform")
		download := res.GetField("downloadUrl")
		pkgname := res.GetField("packageName")
		state := res.GetField("state")
		ctime := res.GetField("cTime")
		createby := res.GetField("cUID")

		fmt.Println(id, name, desc, typeid, platform, download, pkgname, state, ctime, createby)
	}
}

func TestUpdate(t *testing.T) {
	var db *TronDB

	db, err := OpenDB("mysql", "tron:tron@tcp(118.216.57.65:3306)/tron?charset=utf8")

	if err != nil {
		fmt.Println("opendb fail")
		return
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Println(err.Error())
		return
	}

	strSQL := "update blocks set name = 'test1' where id in (2,3)"
	key, rows, err := db.Update(strSQL)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("update success , key=%v, rows=%v", key, rows)
	}

}

func TestDelete(t *testing.T) {
	var db *TronDB

	db, err := OpenDB("mysql", "tron:tron@tcp(118.216.57.65:3306)/tron?charset=utf8")

	if err != nil {
		fmt.Println("opendb fail")
		return
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Println(err.Error())
		return
	}

	strSQL := " delete from blocks  where block_id = 4"
	key, rows, err := db.Update(strSQL)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("delete success , key=%v, rows=%v", key, rows)
	}

}

package grpcclient

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/golang/glog"
	"github.com/tronprotocol/grpc-gateway/core"
	"github.com/wlcy/tron/explorer/core/utils"

	_ "github.com/go-sql-driver/mysql"
)

func TestWalletExt(*testing.T) {

	client := NewWalletExt(fmt.Sprintf("%s:50051", utils.GetRandSolidityNodeAddr()))

	err := client.Connect()
	if nil != err {
		fmt.Println(err)
		return
	}

	fmt.Println(client.GetState(), client.Target())

	addr := "TGo9Me13BSagSHXmKZDbZrLaFW9PXYYs3T"
	addr = "TDPgbSpKrLnaBMF79QUg3aigsG1tsWoxLJ"

	// trans, _ := client.GetTransactionsFromThis2(addr, 0, 100)
	trans, _ := client.GetTransactionsToThisi2(addr, 0, 100)
	// utils.VerifyCall(client.GetTransactionsToThis(addr, 0, 100))

	// a := api.TransactionList{Transaction: trans}
	utils.VerifyCall(trans, nil)

	// serializeTransactionToDB(addr, 0, trans)

}

func serializeTransactionToDB(addr string, beginIdx int64, trans []*core.Transaction) {
	dbb := getMysqlDB()

	txn, err := dbb.Begin()
	if err != nil {
		glog.Fatal(err)
	}
	/*
		'trx_hash','varchar(100)','NO','PRI',NULL,''
		'contract_type','int(11)','YES','',NULL,''
		'owner_address','varchar(45)','YES','',NULL,''
		'signature','varchar(100)','YES','',NULL,''
		'contract_data','varchar(300)','YES','',NULL,''
		'idx','int(11)','YES','',NULL,''

	*/
	sqlstr := "insert into addr_transactions (trx_hash,signature, owner_address, idx, contract_type, contract_data) values (?, ?, ?, ?, ?, ?)"
	stmt, err := txn.Prepare(sqlstr)
	if nil != err {
		glog.Fatal(err)
	}

	for idx, tran := range trans {
		if nil == tran || nil == tran.RawData {
			continue
		}
		if len(tran.RawData.Contract) > 0 {
			utils.VerifyCall(tran, nil)
			trxHash := utils.HexEncode(utils.CalcTransactionHash(tran))
			signature := utils.HexEncode(tran.Signature[0])
			fmt.Printf("idx:%v\n\thash:%v\n\tsign:%v\n", idx, trxHash, signature)
			_, err = stmt.Exec(
				trxHash, signature, addr,
				beginIdx+int64(idx),
				tran.RawData.Contract[0].Type,
				utils.HexEncode(tran.RawData.Contract[0].Parameter.Value))
		} else {
			glog.Errorf("transaction contract is empty!")
		}
		if err != nil {
			glog.Fatal(err)
		}
	}

	err = txn.Commit()
	if err != nil {
		glog.Fatal(err)
	}
	fmt.Println("Program finished successfully")
}

func getMysqlDB() *sql.DB {
	db, err := sql.Open("mysql", "tron:tron@tcp(172.16.21.224:3306)/tron")
	if nil != err {
		return nil
	}
	return db
}

func TestDecode(*testing.T) {
	txid := "3aZ9A3QMD/Tk6WFFA81ZZwcOTJAI5KMKXpQeNkZzw3g="
	fmt.Println(utils.HexEncode(utils.Base64Decode(txid)))
}

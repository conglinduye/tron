package entity

// ReportResp
type ReportResp struct {
	Success 				bool				`json:"success"`			// success
	Data 					[]*ReportOverview	`json:"data"`				// data
}

// Overview
type ReportOverview struct {
	Date					int64				`json:"date"`				// date
	TotalTransaction		int64				`json:"totalTransaction"`	// totalTransaction
	AvgBlockTime			int64				`json:"avgBlockTime"`		// avgBlockTime
	AvgBlockSize			int64				`json:"avgBlockSize"`		// avgBlockSize
	TotalBlockCount			int64				`json:"totalBlockCount"`	// totalBlockCount
	NewAddressSeen			int64				`json:"newAddressSeen"`		// newAddressSeen
	BlockchainSize			int64				`json:"blockchainSize"`		// blockchainSize
	TotalAddress			int64				`json:"totalAddress"`		// totalAddress

}

type ReportBlock struct {
	TotalCount			   	int64				`json:"totalCount"`			// totalCount
	TotalSize				int64				`json:"totalSize"`			// totalSize
}
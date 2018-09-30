package mongo

import (
	"log"
	"testing"

	"github.com/wlcy/tron/explorer/lib/mysql"
	"gopkg.in/mgo.v2/bson"
)

type eventLog struct {
	ID              bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	BlockNum        string        `json:"block_number,omitempty" bson:"block_number,omitempty"`
	BlockTimestamp  int64         `json:"block_timestamp,omitempty" bson:"block_timestamp,omitempty"`
	ContractAddress string        `json:"contract_address,omitempty" bson:"contract_address,omitempty"`
	EventName       string        `json:"event_name,omitempty" bson:"event_name,omitempty"`
	Result          []string      `json:"result,omitempty" bson:"result,omitempty"`
	TransactionID   string        `json:"transaction_id,omitempty" bson:"transaction_id,omitempty"`
}

func TestMongo(t *testing.T) {
	Initialize("47.90.203.178", "18890", "EventLogCenter", "root", "root")
	result, _ := GetMongodbInstance().GetMultiRecord("EventLogCenter", "eventLog", bson.M{"transaction_id": "635603183edb0071edb1e545354e54c2f26e82067d6c5e2115a2436bc453d517"}, bson.M{})
	for _, rr := range result {
		event := &eventLog{}
		bsonBytes, _ := bson.Marshal(rr)
		bson.Unmarshal(bsonBytes, event)
		ss, _ := mysql.JSONObjectToString(event)
		log.Printf("total json:%v\n", ss)
		log.Printf("total:%v", event)
	}

}

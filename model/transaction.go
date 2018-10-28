package model

import (
	"github.com/globalsign/mgo/bson"
	. "time"
)

type Transaction struct {
	ID               bson.ObjectId `bson:"_id" json:"id"`
	AuthorisedAmount GBP           `bson:"authorisedAmount" json:"authorisedAmount"`
	CapturedAmount   GBP           `bson:"capturedAmount" json:"capturedAmount"`
	MerchantName     string        `bson:"merchantName" json:"merchantName"`
	Timestamp        Time          `bson:"timestamp" json:"timestamp"`
}

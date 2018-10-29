package response

import (
	"github.com/globalsign/mgo/bson"
	. "time"
)

type Statement struct {
	Account      AccountDetails         `bson:"accountDetails" json:"accountDetails"`
	FromDate     *Time                  `bson:"fromDate" json:"fromDate"`
	ToDate       *Time                  `bson:"toDate" json:"toDate"`
	Balance      Balance                `bson:"balance" json:"balance"`
	Transactions []StatementTransaction `bson:"transactions" json:"transactions"`
}

type AccountDetails struct {
	FirstName  string `bson:"firstName" json:"firstName"`
	Surname    string `bson:"surname" json:"surname"`
	CardNumber string `bson:"cardNumber" json:"cardNumber"`
}

type StatementTransaction struct {
	ID               bson.ObjectId `bson:"_id" json:"id"`
	AuthorisedAmount float64       `bson:"authorisedAmount" json:"authorisedAmount"`
	CapturedAmount   float64       `bson:"capturedAmount" json:"capturedAmount"`
	MerchantName     string        `bson:"merchantName" json:"merchantName"`
	Timestamp        Time          `bson:"timestamp" json:"timestamp"`
}
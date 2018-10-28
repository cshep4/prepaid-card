package model

import (
	. "github.com/globalsign/mgo/bson"
)

type Account struct {
	ID               ObjectId   `bson:"_id" json:"id"`
	CardNumber       string     `bson:"cardNumber" json:"cardNumber"`
	FirstName        string     `bson:"firstName" json:"firstName"`
	Surname          string     `bson:"surname" json:"surname"`
	TotalBalance     GBP        `bson:"totalBalance" json:"totalBalance"`
	AvailableBalance GBP        `bson:"availableBalance" json:"availableBalance"`
	Transactions     []ObjectId `bson:"transactions" json:"transactions"`
}
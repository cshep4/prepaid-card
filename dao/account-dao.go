package dao

import (
	"github.com/globalsign/mgo/bson"
	. "simple-prepaid-card/constant"
	. "simple-prepaid-card/model"
)

type AccountDAO struct {

}

func (a AccountDAO) Collection() string {
	return ACCOUNT
}

func (a *AccountDAO) FindByCardNumber(cardNumber string) (Account, error) {
	var account Account
	err := db.C(a.Collection()).Find(bson.M{"cardNumber": cardNumber}).One(&account)
	return account, err
}

func (a *AccountDAO) FindByTransactionId(value bson.ObjectId) (Account, error) {
	var account Account
	err := db.C(a.Collection()).Find(bson.M{"transactions": bson.M{"$in": []bson.ObjectId{value}}}).One(&account)
	return account, err
}

func (a *AccountDAO) Insert(account Account) error {
	err := db.C(a.Collection()).Insert(&account)
	return err
}

func (a *AccountDAO) Update(account Account) error {
	err := db.C(a.Collection()).UpdateId(account.ID, &account)
	return err
}

func (a *AccountDAO) Count() (int, error) {
	return db.C(a.Collection()).Count()
}
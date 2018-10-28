package dao

import (
	"github.com/globalsign/mgo/bson"
	. "prepaid-card/constant"
	. "prepaid-card/model"
)

type AccountDAO struct {

}

func (a AccountDAO) Collection() string {
	return ACCOUNT
}

func (a *AccountDAO) FindAll() ([]Account, error) {
	var accounts []Account
	err := db.C(a.Collection()).Find(bson.M{}).All(&accounts)
	return accounts, err
}

func (a *AccountDAO) FindById(id string) (Account, error) {
	var account Account
	err := db.C(a.Collection()).FindId(bson.ObjectIdHex(id)).One(&account)
	return account, err
}

func (a *AccountDAO) FindByCardNumber(cardNumber string) (Account, error) {
	var account Account
	err := db.C(a.Collection()).Find(bson.M{"cardNumber": cardNumber}).One(&account)
	return account, err
}

func (a *AccountDAO) FindByField(field string, value interface{}) (Account, error) {
	var account Account
	err := db.C(a.Collection()).Find(nil).Select(bson.M{field: value}).One(&account)
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

func (a *AccountDAO) Delete(account Account) error {
	err := db.C(a.Collection()).Remove(&account)
	return err
}

func (a *AccountDAO) Update(account Account) error {
	err := db.C(a.Collection()).UpdateId(account.ID, &account)
	return err
}

func (a *AccountDAO) Count() (int, error) {
	return db.C(a.Collection()).Count()
}
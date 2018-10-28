package dao

import (
	"github.com/globalsign/mgo/bson"
	. "simple-prepaid-card/constant"
	. "simple-prepaid-card/model"
)

type TransactionDAO struct {

}

func (t TransactionDAO) Collection() string {
	return TRANSACTION
}

func (t *TransactionDAO) FindAll() ([]Transaction, error) {
	var transactions []Transaction
	err := db.C(t.Collection()).Find(bson.M{}).All(&transactions)
	return transactions, err
}

func (t *TransactionDAO) FindById(id string) (Transaction, error) {
	var transaction Transaction
	err := db.C(t.Collection()).FindId(bson.ObjectIdHex(id)).One(&transaction)
	return transaction, err
}

func (t *TransactionDAO) FindByField(field string, value interface{}) (Transaction, error) {
	var transaction Transaction
	err := db.C(t.Collection()).Find(nil).Select(bson.M{field: value}).One(&transaction)
	return transaction, err
}

func (t *TransactionDAO) FindTransactionByField(field string, value interface{}) (Transaction, error) {
	var transaction Transaction
	err := db.C(t.Collection()).Find(bson.M{"transactions": bson.M{"$elemMatch": bson.M{field: value}}}).One(&transaction)
	return transaction, err
}

func (t *TransactionDAO) Insert(transaction Transaction) error {
	err := db.C(t.Collection()).Insert(&transaction)
	return err
}

func (t *TransactionDAO) Delete(transaction Transaction) error {
	err := db.C(t.Collection()).Remove(&transaction)
	return err
}

func (t *TransactionDAO) Update(transaction Transaction) error {
	err := db.C(t.Collection()).UpdateId(transaction.ID, &transaction)
	return err
}

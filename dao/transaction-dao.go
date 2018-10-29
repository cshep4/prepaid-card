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

func (t *TransactionDAO) FindById(id string) (Transaction, error) {
	var transaction Transaction
	err := db.C(t.Collection()).FindId(bson.ObjectIdHex(id)).One(&transaction)
	return transaction, err
}

func (t *TransactionDAO) FindAllByIds(id []bson.ObjectId) ([]Transaction, error) {
	var transactions []Transaction
	err := db.C(t.Collection()).Find(bson.M{"_id": bson.M{"$in": id}}).All(&transactions)
	return transactions, err
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

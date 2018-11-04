package dao

import (
	"fmt"
	"github.com/globalsign/mgo"
	"os"
	. "simple-prepaid-card/model"
)

type DAO interface {
	Collection() string
}

var db *mgo.Database

func ConnectDb() {
	uri := os.Getenv("MONGODB_URI")
	mongo, err := mgo.ParseURL(uri)

	s, err := mgo.Dial(uri)

	if err != nil {
		panic(err.Error())
	}

	s.SetSafe(&mgo.Safe{})

	db = s.DB(mongo.Database)

	fmt.Println("Database connected")
}

func SaveAccountAndTransaction(account Account,
								transaction Transaction,
								saveAccount func(a Account) error,
								saveTransaction func(t Transaction) error,
								rollbackAccount func(),
								rollbackTransaction func()) error {
	saveErr := make(chan error, 2)
	defer close(saveErr)

	go func() { saveErr <- saveAccount(account) }()
	go func() { saveErr <- saveTransaction(transaction) }()

	accountError := <-saveErr
	transactionError := <-saveErr

	if accountError != nil && transactionError == nil {
		rollbackTransaction()
		return accountError
	} else if accountError == nil && transactionError != nil {
		rollbackAccount()
		return transactionError
	} else if accountError != nil && transactionError != nil {
		return accountError
	}

	return nil
}
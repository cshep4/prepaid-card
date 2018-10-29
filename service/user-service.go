package service

import (
	. "github.com/ahl5esoft/golang-underscore"
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
	"net/http"
	. "simple-prepaid-card/dao"
	. "simple-prepaid-card/model"
	. "simple-prepaid-card/model/request"
	. "simple-prepaid-card/model/response"
	"simple-prepaid-card/response"
	"simple-prepaid-card/util"
	"time"
)

type UserService struct {
	transactionDao TransactionDAO
	accountDao AccountDAO
}

func (u *UserService) CreateAccount(w http.ResponseWriter, r *http.Request) {
	newAccount, err := DecodeNewAccount(r.Body)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	c, err := u.accountDao.Count()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	account := Account{
		ID: bson.NewObjectId(),
		CardNumber: util.CreateUniqueCardNumber(c),
		FirstName: newAccount.FirstName,
		Surname: newAccount.Surname,
		TotalBalance: 0,
		AvailableBalance: 0,
	}

	if err := u.accountDao.Insert(account); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Json(w, http.StatusCreated, account)
}

func (u *UserService) Load(w http.ResponseWriter, r *http.Request) {
	loadMoney, err := DecodeLoadMoney(r.Body)

	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	account, err := u.accountDao.FindByCardNumber(*loadMoney.CardNumber)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid card number")
		return
	}

	account.AvailableBalance += ToGBP(*loadMoney.Amount)
	account.TotalBalance += ToGBP(*loadMoney.Amount)

	if err := u.accountDao.Update(account); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Json(w, http.StatusOK, map[string]string{"result": "success"})
}

func (u *UserService) ViewBalance(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	account, err := u.accountDao.FindByCardNumber(params["cardNumber"])

	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid card number")
		return
	}

	response.Json(w, http.StatusOK, Balance{
		AvailableBalance: account.AvailableBalance.ToPounds(),
		TotalBalance: account.TotalBalance.ToPounds(),
	})
}

func (u *UserService) ViewStatement(w http.ResponseWriter, r *http.Request) {
	fromDate, err1 := util.ParseDateParameter(r.URL.Query().Get("fromDate"))
	toDate, err2 := util.ParseDateParameter(r.URL.Query().Get("toDate"))

	if err1 != nil || err2 != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	cardNumber := mux.Vars(r)["cardNumber"]

	account, err := u.accountDao.FindByCardNumber(cardNumber)

	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid card number")
		return
	}

	transactions, err := u.retrieveTransactionsInRange(account.Transactions, fromDate, toDate)

	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	a := AccountDetails{
		FirstName: account.FirstName,
		Surname: account.Surname,
		CardNumber: account.CardNumber,
	}

	b := Balance{
		AvailableBalance: account.AvailableBalance.ToPounds(),
		TotalBalance: account.TotalBalance.ToPounds(),
	}

	resp := Statement{
		Account:      a,
		FromDate:     fromDate,
		ToDate:       toDate,
		Balance:      b,
		Transactions: transactions,
	}

	response.Json(w, http.StatusOK, resp)
}

func (u *UserService) retrieveTransactionsInRange(ids []bson.ObjectId, fromDate *time.Time, toDate *time.Time) ([]StatementTransaction, error) {
	transactions, err := u.transactionDao.FindAllByIds(ids)

	if err != nil || transactions == nil {
		return []StatementTransaction{}, err
	}

	statementTransactions := Map(transactions, func (t Transaction, _ int) StatementTransaction {
		return StatementTransaction{
			ID:               t.ID,
			AuthorisedAmount: t.AuthorisedAmount.ToPounds(),
			CapturedAmount:   t.CapturedAmount.ToPounds(),
			MerchantName:     t.MerchantName,
			Timestamp:        t.Timestamp,
		}
	})

	if fromDate == nil && toDate == nil {
		return statementTransactions.([]StatementTransaction), err
	}

	predicate := transactionPredicate(fromDate, toDate)

	st := Where(statementTransactions, func(t StatementTransaction, i int) bool { return predicate(t.Timestamp, toDate, fromDate) })

	if st == nil {
		return []StatementTransaction{}, nil
	}

	return st.([]StatementTransaction), nil
}

func transactionPredicate(fromDate *time.Time, toDate *time.Time) (func(timestamp time.Time, toDate *time.Time, fromDate *time.Time) bool) {
	if fromDate == nil && toDate != nil {
		return isTransactionBeforeDate
	} else if fromDate != nil && toDate == nil {
		return isTransactionAfterDate
	}

	return isTransactionInRange
}

func isTransactionInRange(timestamp time.Time, toDate *time.Time, fromDate *time.Time) bool {
	return isTransactionBeforeDate(timestamp, toDate, fromDate) && isTransactionAfterDate(timestamp, toDate, fromDate)
}

func isTransactionBeforeDate(timestamp time.Time, toDate *time.Time, fromDate *time.Time) bool {
	return timestamp.Before(toDate.AddDate(0, 0, 1))
}

func isTransactionAfterDate(timestamp time.Time, toDate *time.Time, fromDate *time.Time) bool {
	return timestamp.After(*fromDate)
}
package service

import (
	"errors"
	"github.com/globalsign/mgo/bson"
	"io"
	"net/http"
	. "prepaid-card/dao"
	. "prepaid-card/model"
	. "prepaid-card/model/request"
	"prepaid-card/response"
	"sync"
	"time"
)

type MerchantService struct {
	transactionDao TransactionDAO
	accountDao     AccountDAO
}

var mu = &sync.Mutex{}

func (m *MerchantService) AuthoriseTransaction(w http.ResponseWriter, r *http.Request) {
	authorisationRequest, err := DecodeAuthorisationRequest(r.Body)

	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	mu.Lock()
	defer mu.Unlock()

	account, err := m.accountDao.FindByCardNumber(*authorisationRequest.CardNumber)

	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid card number")
		return
	}

	transactionAmount := ToGBP(*authorisationRequest.Amount)

	if account.AvailableBalance-transactionAmount < 0 {
		response.Error(w, http.StatusBadRequest, "Insufficient funds")
		return
	}

	transaction := Transaction{
		ID:               bson.NewObjectId(),
		AuthorisedAmount: transactionAmount,
		CapturedAmount:   0,
		MerchantName:     *authorisationRequest.MerchantName,
		Timestamp:        time.Now(),
	}

	account.AvailableBalance -= transactionAmount
	account.Transactions = append(account.Transactions, transaction.ID)

	rollbackAccount := func() {
		account.AvailableBalance += transaction.AuthorisedAmount
		account.Transactions = account.Transactions[:len(account.Transactions)-1]
		m.accountDao.Update(account)
	}

	rollbackTransaction := func() {
		m.transactionDao.Delete(transaction)
	}

	if err := SaveAccountAndTransaction(account, transaction, m.accountDao.Update, m.transactionDao.Insert, rollbackAccount, rollbackTransaction); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
	}

	response.Json(w, http.StatusCreated, transaction)
}

func (m *MerchantService) CaptureTransaction(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	transaction, amountToCapture, err := m.retrieveTransactionAndAmount(r.Body)

	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if transaction.CapturedAmount + *amountToCapture > transaction.AuthorisedAmount {
		response.Error(w, http.StatusBadRequest, "Cannot capture more than the authorised amount")
		return
	}

	account, err := m.accountDao.FindByTransactionId(transaction.ID)

	account.TotalBalance -= *amountToCapture
	transaction.CapturedAmount += *amountToCapture

	rollbackAccount := func() {
		account.TotalBalance += *amountToCapture
		m.accountDao.Update(account)
	}

	rollbackTransaction := func() {
		transaction.CapturedAmount -= *amountToCapture
		m.transactionDao.Update(*transaction)
	}

	if err := SaveAccountAndTransaction(account, *transaction, m.accountDao.Update, m.transactionDao.Update, rollbackAccount, rollbackTransaction); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
	}

	response.Json(w, http.StatusOK, transaction)
}

func (m *MerchantService) ReverseAuthorisation(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	transaction, amountToReverse, err := m.retrieveTransactionAndAmount(r.Body)

	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if *amountToReverse > transaction.AuthorisedAmount - transaction.CapturedAmount {
		response.Error(w, http.StatusBadRequest, "Cannot reverse more than authorised and not already captured")
		return
	}

	account, err := m.accountDao.FindByTransactionId(transaction.ID)

	account.AvailableBalance += *amountToReverse
	transaction.AuthorisedAmount -= *amountToReverse

	rollbackAccount := func() {
		account.AvailableBalance -= *amountToReverse
		m.accountDao.Update(account)
	}

	rollbackTransaction := func() {
		transaction.AuthorisedAmount += *amountToReverse
		m.transactionDao.Update(*transaction)
	}

	if err := SaveAccountAndTransaction(account, *transaction, m.accountDao.Update, m.transactionDao.Update, rollbackAccount, rollbackTransaction); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
	}

	response.Json(w, http.StatusOK, transaction)
}

func (m *MerchantService) Refund(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	transaction, amountToRefund, err := m.retrieveTransactionAndAmount(r.Body)

	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if *amountToRefund > transaction.CapturedAmount {
		response.Error(w, http.StatusBadRequest, "Cannot refund more than captured")
		return
	}

	account, err := m.accountDao.FindByTransactionId(transaction.ID)

	refundTransaction := Transaction{
		ID:               bson.NewObjectId(),
		AuthorisedAmount: -*amountToRefund,
		CapturedAmount:   -*amountToRefund,
		MerchantName:     "REFUND - " + transaction.MerchantName,
		Timestamp:        time.Now(),
	}

	account.AvailableBalance += *amountToRefund
	account.TotalBalance += *amountToRefund
	account.Transactions = append(account.Transactions, refundTransaction.ID)

	rollbackAccount := func() {
		account.AvailableBalance -= refundTransaction.AuthorisedAmount
		account.TotalBalance -= refundTransaction.AuthorisedAmount
		account.Transactions = account.Transactions[:len(account.Transactions)-1]
		m.accountDao.Update(account)
	}

	rollbackTransaction := func() {
		m.transactionDao.Delete(refundTransaction)
	}

	if err := SaveAccountAndTransaction(account, refundTransaction, m.accountDao.Update, m.transactionDao.Insert, rollbackAccount, rollbackTransaction); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
	}

	response.Json(w, http.StatusOK, refundTransaction)
}

func (m *MerchantService) retrieveTransactionAndAmount(body io.ReadCloser) (*Transaction, *GBP, error) {
	transactionRequest, err := DecodeTransactionRequest(body)

	if err != nil {
		return nil, nil, err
	}

	transaction, err := m.transactionDao.FindById(transactionRequest.TransactionId)

	if err != nil {
		return nil, nil, errors.New("Invalid transaction ID")
	}

	amountToCapture := ToGBP(transactionRequest.Amount)

	return &transaction, &amountToCapture, nil
}
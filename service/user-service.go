package service

import (
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
	"net/http"
	. "simple-prepaid-card/dao"
	. "simple-prepaid-card/model"
	. "simple-prepaid-card/model/request"
	. "simple-prepaid-card/model/response"
	"simple-prepaid-card/response"
	"simple-prepaid-card/util"
)

type UserService struct {
	dao AccountDAO
}

func (u *UserService) CreateAccount(w http.ResponseWriter, r *http.Request) {
	newAccount, err := DecodeNewAccount(r.Body)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	c, err := u.dao.Count()
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

	if err := u.dao.Insert(account); err != nil {
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

	account, err := u.dao.FindByCardNumber(*loadMoney.CardNumber)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid card number")
		return
	}

	account.AvailableBalance += ToGBP(*loadMoney.Amount)
	account.TotalBalance += ToGBP(*loadMoney.Amount)

	if err := u.dao.Update(account); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Json(w, http.StatusOK, map[string]string{"result": "success"})
}

func (u *UserService) ViewBalance(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	account, err := u.dao.FindByCardNumber(params["cardNumber"])

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
	panic("implement me")
}

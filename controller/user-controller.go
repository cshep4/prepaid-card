package controller

import (
	"github.com/gorilla/mux"
	. "prepaid-card/constant"
	. "prepaid-card/service"
)

type UserController struct {

}

func (u UserController) Path() string {
	return "/" + ACCOUNT
}

func (u UserController) CreateRoutes(r *mux.Router) {
	userService := NewUserService()

	r.HandleFunc(u.Path() + "/create", userService.CreateAccount).Methods("POST")
	r.HandleFunc(u.Path() + "/load", userService.Load).Methods("PUT")
	r.HandleFunc(u.Path() + "/balance/{cardNumber}", userService.ViewBalance).Methods("GET")
	r.HandleFunc(u.Path() + "/statement/{cardNumber}", userService.ViewStatement).Methods("GET")
}


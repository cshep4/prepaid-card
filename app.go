package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	. "prepaid-card/controller"
	. "prepaid-card/dao"
)

var userController UserController
var merchantController MerchantController

func init() {
	ConnectDb()

	userController = UserController{}
	merchantController = MerchantController{}
}

func main() {
	r := mux.NewRouter()

	userController.CreateRoutes(r)
	merchantController.CreateRoutes(r)

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}

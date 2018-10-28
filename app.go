package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	. "simple-prepaid-card/controller"
	. "simple-prepaid-card/dao"
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

	if err := http.ListenAndServe(os.Getenv("PORT"), r); err != nil {
		log.Fatal(err)
	}
}

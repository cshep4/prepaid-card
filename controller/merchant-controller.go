package controller

import (
	"github.com/gorilla/mux"
	. "prepaid-card/constant"
	. "prepaid-card/service"
)

type MerchantController struct {

}

func (m MerchantController) Path() string {
	return "/" + TRANSACTION
}

func (m MerchantController) CreateRoutes(r *mux.Router) {
	merchantService := NewMerchantService()

	r.HandleFunc(m.Path() + "/authorisation", merchantService.AuthoriseTransaction).Methods("POST")
	r.HandleFunc(m.Path() + "/capture", merchantService.CaptureTransaction).Methods("PUT")
	r.HandleFunc(m.Path() + "/reverse", merchantService.ReverseAuthorisation).Methods("PUT")
	r.HandleFunc(m.Path() + "/refund", merchantService.Refund).Methods("POST")
}
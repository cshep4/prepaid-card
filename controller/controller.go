package controller

import "github.com/gorilla/mux"

type Controller interface {
	CreateRoutes(router *mux.Router)
	Path() string
}
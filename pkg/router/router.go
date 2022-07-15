package router

import (
	"github.com/gorilla/mux"
	"github.com/ragsharan/assignment/pkg/handler"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/answer", handler.GetAnswer).Methods("GET")
	router.HandleFunc("/api/v1/answer", handler.AddAnswer).Methods("POST")
	router.HandleFunc("/api/v1/answer", handler.UpdateAnswer).Methods("PUT")
	router.HandleFunc("/api/v1/answer", handler.RemoveAnswer).Methods("DELETE")

	return router
}

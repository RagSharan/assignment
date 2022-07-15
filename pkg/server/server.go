package server

import (
	"net/http"

	"github.com/ragsharan/assignment/pkg/router"
)

func Run() {

	router := router.NewRouter()

	http.ListenAndServe(":8080", router)

}

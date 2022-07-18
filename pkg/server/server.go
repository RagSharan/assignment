package server

import (
	"context"
	"net/http"

	"github.com/ragsharan/assignment/pkg/config"
	"github.com/ragsharan/assignment/pkg/router"
)

func Run() {
	ctx := context.Background()

	db, client := config.ConnectDB()
	ctx = context.WithValue(ctx, "dbIns", db)
	defer config.CloseConnection(client)

	router := router.NewRouter()

	http.ListenAndServe(":8080", router)

}

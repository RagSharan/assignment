package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ragsharan/assignment/pkg/service/v0/events"
)

var (
	eventServiceIns events.IEventService = events.NewService()
)

func GetEventsById(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	urlPath := strings.Split(req.URL.Path, "/")
	id := strings.TrimSpace(urlPath[4])
	result, err := eventServiceIns.GetEventsById(id)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(err)
		return
	}
	res.WriteHeader(http.StatusAccepted)
	json.NewEncoder(res).Encode(result)

}

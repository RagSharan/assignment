package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ragsharan/assignment/pkg/model"
	service "github.com/ragsharan/assignment/pkg/service/v0"
)

var (
	serviceIns service.IService = service.NewService()
)

func GetAnswer(res http.ResponseWriter, req *http.Request) {
	filter := make(map[string]interface{})
	res.Header().Set("Content-Type", "application/json")
	data := req.URL.Query()
	for k, v := range data {
		if len(v) == 1 {
			filter[k] = v[0]
		}
	}

	result, err := serviceIns.GetAnswer(filter)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(err)
		return
	}
	res.WriteHeader(http.StatusAccepted)
	json.NewEncoder(res).Encode(result)
}

func AddAnswer(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var answer model.Answer
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
	}
	json.Unmarshal(data, &answer)
	result, err := serviceIns.AddAnswer(answer)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(err)
		return
	}
	res.WriteHeader(http.StatusAccepted)
	json.NewEncoder(res).Encode(result)
}

func UpdateAnswer(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var answer model.Answer
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
	}
	json.Unmarshal(data, &answer)
	result, err := serviceIns.UpdateAnswer(answer)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(err)
		return
	}
	res.WriteHeader(http.StatusAccepted)
	json.NewEncoder(res).Encode(result)
}

func RemoveAnswer(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	filter := make(map[string]interface{})
	data := req.URL.Query()
	for k, v := range data {
		if len(v) == 1 {
			filter[k] = v[0]
		}
	}
	result, err := serviceIns.RemoveAnswer(filter)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(err)
		return
	}
	res.WriteHeader(http.StatusAccepted)
	json.NewEncoder(res).Encode(result)
}

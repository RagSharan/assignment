package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ragsharan/assignment/pkg/model"
	"github.com/ragsharan/assignment/pkg/service/v0/answers"
)

var (
	serviceIns answers.IService = answers.NewService()
)

func GetAnswer(res http.ResponseWriter, req *http.Request) {
	filter := make(map[string]interface{})
	res.Header().Set("Content-Type", "application/json")
	queryParams := req.URL.Query()
	fmt.Println("query params", queryParams)
	for k, v := range queryParams {
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

func GetAnswerById(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	urlPath := strings.Split(req.URL.Path, "/")
	id := strings.TrimSpace(urlPath[4])
	result, err := serviceIns.GetAnswerById(id)
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
	urlPath := strings.Split(req.URL.Path, "/")
	id := strings.TrimSpace(urlPath[4])
	result, err := serviceIns.RemoveAnswer(id)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(err)
		return
	}
	res.WriteHeader(http.StatusAccepted)
	json.NewEncoder(res).Encode(result)
}

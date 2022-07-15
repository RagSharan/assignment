package service

import (
	query "github.com/ragsharan/assignment/pkg/db/mongodb"
	"github.com/ragsharan/assignment/pkg/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	queryIns query.Iquery = query.NewQuery()
)

const collection string = "answer"

type service struct{}

type IService interface {
	GetAnswer(filter map[string]interface{}) ([]primitive.M, error)
	AddAnswer(answer model.Answer) (*mongo.InsertOneResult, error)
	UpdateAnswer(answer model.Answer) (*mongo.UpdateResult, error)
	RemoveAnswer(filter map[string]interface{}) (*mongo.DeleteResult, error)
}

func NewService() IService {
	return &service{}
}

func (*service) GetAnswer(params map[string]interface{}) ([]primitive.M, error) {
	result, err := queryIns.FindOne(collection, params)

	return result, err
}

func (*service) AddAnswer(answer model.Answer) (*mongo.InsertOneResult, error) {
	result, err := queryIns.Create(collection, answer)
	return result, err
}

func (*service) UpdateAnswer(answer model.Answer) (*mongo.UpdateResult, error) {
	result, err := queryIns.UpdateById(collection, answer)

	return result, err
}

func (*service) RemoveAnswer(filter map[string]interface{}) (*mongo.DeleteResult, error) {
	result, err := queryIns.DeleteDocument(collection, filter)
	return result, err
}

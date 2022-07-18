package answers

import (
	"errors"

	query "github.com/ragsharan/assignment/pkg/db/mongodb"
	"github.com/ragsharan/assignment/pkg/model"
	"github.com/ragsharan/assignment/pkg/service/v0/events"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	queryIns query.Iquery         = query.NewQuery()
	eventIns events.IEventService = events.NewService()
)

const collection string = "answers"

type service struct{}

type IService interface {
	GetAnswerList(params map[string]interface{}) (answers []primitive.M, err error)
	GetAnswer(filter map[string]interface{}) ([]model.Answer, error)
	AddAnswer(answer model.Answer) (*mongo.InsertOneResult, error)
	UpdateAnswer(answer model.Answer) (*mongo.UpdateResult, error)
	RemoveAnswer(filter map[string]interface{}) (*mongo.DeleteResult, error)
}

func NewService() IService {
	return &service{}
}

func (*service) GetAnswerList(params map[string]interface{}) (result []primitive.M, err error) {

	result, err = queryIns.FindList(collection, params)
	if err != nil {
		return
	}

	return result, err
}

func (*service) GetAnswer(params map[string]interface{}) ([]model.Answer, error) {
	var answers []model.Answer
	filter := genFilter(params)
	answers, err := queryIns.FindAll(collection, filter)
	return answers, err
}

func (*service) AddAnswer(answer model.Answer) (*mongo.InsertOneResult, error) {
	result, err := queryIns.Create(collection, answer)
	if err == nil {
		//str := result.InsertedID.(string)
		answer.Id = result.InsertedID.(primitive.ObjectID)
		_, err = recordEvent("create", answer)
	}
	return result, err
}

func (*service) UpdateAnswer(answer model.Answer) (*mongo.UpdateResult, error) {
	id := answer.Id
	result, err := queryIns.UpdateById(collection, id, answer)
	if err == nil {
		_, err = recordEvent("update", answer)
	}
	return result, err
}

func (*service) RemoveAnswer(params map[string]interface{}) (*mongo.DeleteResult, error) {
	filter := genFilter(params)
	result, err := queryIns.DeleteDocument(collection, filter)
	if err == nil && result.DeletedCount != 0 {
		var answer model.Answer
		for k, v := range params {
			answer = model.Answer{
				Key:   k,
				Value: v,
			}
		}
		_, err = recordEvent("delete", answer)
	}
	return result, err
}

func recordEvent(eventName string, answer model.Answer) (*mongo.InsertOneResult, error) {
	event := model.Event{
		EventName: eventName,
		Data:      answer,
	}
	data, err := eventIns.AddEvents(event)
	if err != nil {
		err = errors.New("events are not recorded")
	}
	return data, err
}

func genFilter(params map[string]interface{}) primitive.M {
	var filter primitive.M
	for k, v := range params {
		if k == "id" {
			k = "_id"
			id, _ := primitive.ObjectIDFromHex(v.(string))
			filter = bson.M{
				k: id,
			}
		} else {
			filter = bson.M{
				k: v,
			}
		}
	}
	return filter
}

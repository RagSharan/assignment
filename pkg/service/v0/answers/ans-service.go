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
	GetAnswerById(id string) (model.Answer, error)
	GetAnswer(filter map[string]interface{}) ([]model.Answer, error)
	AddAnswer(answer model.Answer) (*mongo.InsertOneResult, error)
	UpdateAnswer(answer model.Answer) (*mongo.UpdateResult, error)
	RemoveAnswer(id string) (*mongo.DeleteResult, error)
}

func NewService() IService {
	return &service{}
}
func (*service) GetAnswerById(param string) (model.Answer, error) {
	id, _ := primitive.ObjectIDFromHex(param)
	filter := primitive.M{
		"_id": id,
	}
	answers, err := queryIns.FindAll(collection, filter)
	return answers[0], err
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

func (*service) RemoveAnswer(param string) (*mongo.DeleteResult, error) {
	id, _ := primitive.ObjectIDFromHex(param)
	filter := primitive.M{
		"_id": id,
	}
	result, err := queryIns.DeleteDocument(collection, filter)
	if err == nil && result.DeletedCount != 0 {
		answer := model.Answer{
			Id:    id,
			Key:   "",
			Value: nil,
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

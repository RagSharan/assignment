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
	result, err := queryIns.FindAll(collection, filter)
	data := result[0].(primitive.D)
	answer := convertInStruct(data)

	return answer, err
}
func (*service) GetAnswer(params map[string]interface{}) ([]model.Answer, error) {
	var answers []model.Answer
	filter := genFilter(params)
	result, err := queryIns.FindAll(collection, filter)
	for _, v := range result {
		data := v.(primitive.D)
		answer := convertInStruct(data)
		answers = append(answers, answer)
	}

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
	if err == nil && result.ModifiedCount != 0 {
		_, err = recordEvent("update", answer)
	}
	return result, err
}

func (*service) RemoveAnswer(param string) (*mongo.DeleteResult, error) {
	id, _ := primitive.ObjectIDFromHex(param)
	filter := primitive.M{
		"_id": id,
	}
	results, _ := queryIns.FindAll(collection, filter)
	data := results[0].(primitive.D)
	answer := convertInStruct(data)

	result, err := queryIns.DeleteDocument(collection, filter)
	if err == nil && result.DeletedCount != 0 {
		_, err = recordEvent("delete", answer)
	}
	return result, err
}

func recordEvent(eventName string, answer model.Answer) (*mongo.InsertOneResult, error) {
	x := model.Data{
		Key:   answer.Key,
		Value: answer.Value,
	}
	event := model.Event{
		EventId:   answer.Id,
		EventName: eventName,
		Data:      x,
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

func convertInStruct(data primitive.D) (structObj model.Answer) {
	byteD, _ := bson.Marshal(data)
	bson.Unmarshal(byteD, &structObj)
	return
}

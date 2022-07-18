package events

import (
	query "github.com/ragsharan/assignment/pkg/db/mongodb"
	"github.com/ragsharan/assignment/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	queryIns query.Iquery = query.NewQuery()
)

const collection string = "eventsdata"

type eventService struct{}

type IEventService interface {
	GetEventsById(param string) ([]model.Event, error)
	GetEvents(filter map[string]interface{}) ([]model.Event, error)
	AddEvents(event model.Event) (*mongo.InsertOneResult, error)
	UpdateEvents(event model.Event) (*mongo.UpdateResult, error)
	RemoveEvents(filter map[string]interface{}) (*mongo.DeleteResult, error)
}

func NewService() IEventService {
	return &eventService{}
}
func (*eventService) GetEventsById(param string) ([]model.Event, error) {
	var events []model.Event
	id, _ := primitive.ObjectIDFromHex(param)
	filter := primitive.M{
		"eventId": id,
	}
	result, err := queryIns.FindAll(collection, filter)
	for _, v := range result {
		data := v.(primitive.D)
		event := convertInStruct(data)
		events = append(events, event)
	}
	return events, err
}

func (*eventService) GetEvents(filter map[string]interface{}) ([]model.Event, error) {
	var events []model.Event
	result, err := queryIns.FindAll(collection, filter)

	for _, v := range result {
		data := v.(primitive.D)
		event := convertInStruct(data)
		events = append(events, event)
	}
	return events, err
}

func (*eventService) AddEvents(event model.Event) (*mongo.InsertOneResult, error) {
	result, err := queryIns.Create(collection, event)
	return result, err
}

func (*eventService) UpdateEvents(event model.Event) (*mongo.UpdateResult, error) {
	id := event.Id
	result, err := queryIns.UpdateById(collection, id, event)

	return result, err
}

func (*eventService) RemoveEvents(filter map[string]interface{}) (*mongo.DeleteResult, error) {
	result, err := queryIns.DeleteDocument(collection, filter)
	return result, err
}
func convertInStruct(data primitive.D) (structObj model.Event) {
	byteD, _ := bson.Marshal(data)
	bson.Unmarshal(byteD, &structObj)
	return
}

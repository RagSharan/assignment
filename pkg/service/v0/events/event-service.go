package events

import (
	"errors"
	"fmt"

	query "github.com/ragsharan/assignment/pkg/db/mongodb"
	"github.com/ragsharan/assignment/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	queryIns query.Iquery = query.NewQuery()
)

const collection string = "events"

type eventService struct{}

type IEventService interface {
	GetEvents(filter map[string]interface{}) ([]model.Event, error)
	AddEvents(event model.Event) (*mongo.InsertOneResult, error)
	UpdateEvents(event model.Event) (*mongo.UpdateResult, error)
	RemoveEvents(filter map[string]interface{}) (*mongo.DeleteResult, error)
}

func NewService() IEventService {
	return &eventService{}
}

func (*eventService) GetEvents(params map[string]interface{}) ([]model.Event, error) {
	var events []model.Event
	var filter primitive.M
	for k, v := range params {
		if k == "id" {
			k = "_id"
			id, e := primitive.ObjectIDFromHex(v.(string))
			if e != nil {
				fmt.Println("hexa error", e)
			}
			fmt.Println("id of hexa", id)
			filter = bson.M{
				k: id,
			}
		} else {
			filter = bson.M{
				k: v,
			}
		}
	}
	fmt.Println("filter value", filter)
	// result, err := queryIns.FindAll(collection, filter)
	// for _, v := range result {
	// 	event, _ := v.(model.Event)
	// 	events = append(events, event)
	// }
	fmt.Println("events", events)
	err := errors.New("bad data")
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

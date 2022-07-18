package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	EventId   primitive.ObjectID `json:"eventId" bson:"eventId"`
	EventName string             `json:"eventName" bson:"eventName"`
	Data      Data               `json:"data" bson:"data"`
}

type Answer struct {
	Id    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Key   string             `json:"key" bson:"key,omitempty"`
	Value interface{}        `json:"value" bson:"value"`
}

type Data struct {
	Key   string
	Value interface{}
}

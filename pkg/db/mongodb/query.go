package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type query struct{}

type Iquery interface {
	FindAll(collName string, filter primitive.M) ([]interface{}, error)
	Create(collName string, doc interface{}) (*mongo.InsertOneResult, error)
	UpdateById(collName string, id primitive.ObjectID, doc interface{}) (*mongo.UpdateResult, error)
	DeleteDocument(collName string, filter primitive.M) (*mongo.DeleteResult, error)
}

func NewQuery() Iquery {
	return &query{}
}

func (*query) Create(collName string, doc interface{}) (*mongo.InsertOneResult, error) {
	database, client := connectDB()
	defer closeConnection(client)
	collection := database.Collection(collName)
	result, err := collection.InsertOne(context.Background(), doc)
	if err != nil {
		log.Println(err)
	}
	return result, err
}

func (*query) FindAll(collName string, filter primitive.M) ([]interface{}, error) {
	ctx := context.Background()
	database, client := connectDB()
	defer closeConnection(client)
	collection := database.Collection(collName)

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Println(err)
	}
	var result []interface{}
	cursor.All(ctx, &result)
	defer cursor.Close(ctx)

	return result, err
}

func (*query) UpdateById(collName string, id primitive.ObjectID, doc interface{}) (*mongo.UpdateResult, error) {
	database, client := connectDB()
	defer closeConnection(client)
	collection := database.Collection(collName)

	filter := bson.M{"_id": id}
	update := bson.M{"$set": doc}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Println(err)
	}
	return result, err
}

func (*query) DeleteDocument(collName string, filter primitive.M) (*mongo.DeleteResult, error) {
	database, client := connectDB()
	defer closeConnection(client)
	collection := database.Collection(collName)

	deleteResult, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Deleted %v documents in the collection\n", deleteResult.DeletedCount)
	return deleteResult, err
}

func connectDB() (*mongo.Database, *mongo.Client) {
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println(err)
	}
	database := client.Database("Questions")
	fmt.Println("Connected to MongoDB!")
	return database, client
}
func closeConnection(client *mongo.Client) {
	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Connection to MongoDB instance is closed.")
}

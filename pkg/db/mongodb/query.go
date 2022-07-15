package mongodb

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type query struct{}

type Iquery interface {
	FindOne(collName string, doc map[string]interface{}) ([]bson.M, error)
	Create(collName string, doc interface{}) (*mongo.InsertOneResult, error)
	UpdateById(collName string, id primitive.ObjectID, doc interface{}) (*mongo.UpdateResult, error)
	DeleteDocument(collName string, doc interface{}) (*mongo.DeleteResult, error)
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

func (*query) FindOne(collName string, params map[string]interface{}) ([]bson.M, error) {
	database, client := connectDB()
	defer closeConnection(client)
	collection := database.Collection(collName)
	findOptions := options.Find()
	filter := formateFilter(params)
	cursor, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Println(err)
	}
	var result []bson.M
	for cursor.Next(context.TODO()) {
		var data bson.M
		cursor.Decode(&data)
		result = append(result, data)
	}
	return result, err
}

func (*query) UpdateById(collName string, id primitive.ObjectID, doc interface{}) (*mongo.UpdateResult, error) {
	database, client := connectDB()
	defer closeConnection(client)
	collection := database.Collection(collName)

	filter := bson.M{"_id": id}
	update := bson.M{"$set": doc}
	result, err := collection.UpdateByID(context.TODO(), filter, update)
	if err != nil {
		log.Println(err)
	}
	return result, err
}

func (*query) DeleteDocument(collName string, doc interface{}) (*mongo.DeleteResult, error) {
	database, client := connectDB()
	defer closeConnection(client)
	collection := database.Collection(collName)
	kmap, err := trimObject(doc)
	if err != nil {
		log.Println(err)
	}
	filter := formateFilter(kmap)
	deleteResult, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Deleted %v documents in the collection\n", deleteResult.DeletedCount)
	return deleteResult, err
}

func trimObject(doc interface{}) (map[string]interface{}, error) {
	var kmap map[string]interface{}
	data, err := bson.Marshal(doc)
	if err != nil {
		log.Println(err)
	}
	err = bson.Unmarshal(data, &kmap)
	log.Println("kmap", kmap)
	return kmap, err
}
func formateFilter(kmap map[string]interface{}) primitive.M {
	var filter primitive.M
	keyRecursion(&filter, &kmap)
	log.Println("filter", filter)
	return filter
}

func keyRecursion(filter *primitive.M, kmap *map[string]interface{}, tempKey ...string) {
	for key, value := range *kmap {
		if reflect.TypeOf(value).Kind() != reflect.Map {
			if tempKey != nil {
				key = tempKey[0] + "." + key
			}
			*filter = bson.M{key: value}
		} else {
			tempKey := key
			tempMap := value.(map[string]interface{})
			keyRecursion(filter, &tempMap, tempKey)
		}
	}
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

func formateUpdate(kmap map[string]interface{}) (primitive.M, []primitive.M) {
	var filter primitive.M
	var update []primitive.M
	i := 0
	for key, value := range kmap {
		if i == 0 {
			zmap := map[string]interface{}{key: value}
			filter = formateFilter(zmap)
			i = 1
			continue
		}
		zmap := map[string]interface{}{key: value}
		var temp primitive.M
		keyRecursion(&temp, &zmap)
		final := bson.M{"$set": temp} //$set could be replaced by other methods
		update = append(update, final)
	}
	return filter, update
}

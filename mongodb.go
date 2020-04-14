package main
// Most of this is code from https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

type mongodb struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func (m *mongodb) init() {
	connectionString := os.Getenv("MONGODB_CONNECTION_STRING")
	clientOptions := options.Client().ApplyURI(connectionString)
	var err error
	m.client, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = m.client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	m.collection = m.client.Database("todo").Collection("todo_items")
}

func (m *mongodb) createItem(item Item) (Item, error) {
	insertResult, err := m.collection.InsertOne(context.TODO(), item)
	if err != nil {
		return Item{}, errors.New("Unable to insert item into database.")
	}
	item.Id = insertResult.InsertedID.(primitive.ObjectID).Hex()
	return item, nil
}

func (m *mongodb) deleteItem(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("Invalid ID.")
	}

	filter := bson.D{{"_id", objID}}

	_, err = m.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return &ErrorItemNotFound{Id: id}
	}
	return nil
}

func (m *mongodb) updateItem(id string, td Item) (Item, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Item{}, errors.New("Invalid ID.")
	}

	filter := bson.D{{"_id", objID}}
	update := bson.D{
		{"$set", bson.M{"Description": td.Description, "Completed": td.Completed}},
	}

	_, err = m.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return Item{}, &ErrorItemNotFound{Id: id}
	}
	td.Id = id
	return td, nil
}

func (m *mongodb) getItem(id string) (Item, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Item{}, errors.New("Invalid ID.")
	}

	filter := bson.D{{"_id", objID}}

	var item Item
	err = m.collection.FindOne(context.TODO(), filter).Decode(&item)
	if err != nil {
		return Item{}, &ErrorItemNotFound{Id: id}
	}

	item.Id = id
	return item, nil
}

func (m *mongodb) allItems() ([]Item, error) {
	findOptions := options.Find()
	var results []Item
	var emptyResults []Item

	cur, err := m.collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		return emptyResults, err
	}

	for cur.Next(context.TODO()) {
		var elem Item
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

	    elements, _ := cur.Current.Elements()
	    elem.Id = elements[0].Value().ObjectID().Hex()

		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		return emptyResults, err
	}

	cur.Close(context.TODO())
	return results, err
}

func (m *mongodb) close() {
	m.client.Disconnect(context.TODO())
}

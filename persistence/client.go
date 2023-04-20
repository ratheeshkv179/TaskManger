package persistence

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoClient struct {
	ctx           context.Context
	userName      string
	password      string
	port          string
	hostname      string
	connectionURI string
	Client        *mongo.Client
}

func (d *MongoClient) Init(username, password, hostname, port string) {
	d.userName = username
	d.password = password
	d.port = port
	d.hostname = hostname
	//d.connectionURI = fmt.Sprintf("mongodb://%s:%s@%s:%s", d.userName, d.password, d.hostname, d.port)
	d.connectionURI = fmt.Sprintf("mongodb://%s:%s", d.hostname, d.port)
}

func (d *MongoClient) Connect() {
	var err error
	var cancel context.CancelFunc
	d.ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	clientOpts := options.ClientOptions{}
	clientOpts.ApplyURI(d.connectionURI)
	d.Client, err = mongo.Connect(d.ctx, &clientOpts)
	if err != nil {
		fmt.Println(fmt.Sprintf("MongoDB connection error %s", err.Error()))
	}
}

func (d *MongoClient) Disconnect() {
	err := d.Client.Disconnect(d.ctx)
	if err != nil {
		fmt.Println(fmt.Sprintf("MongoDB disconnection error %s", err.Error()))
	}
}

func (d *MongoClient) Get(dbName, collectionName string) ([]bson.D, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	var response []bson.D
	defer cancel()
	collection := d.Client.Database(dbName).Collection(collectionName)
	cur, err := collection.Find(ctx, bson.D{})
	defer cur.Close(ctx)
	if err != nil {
		fmt.Println("error", err.Error())
	} else {
		for cur.Next(ctx) {
			var result bson.D
			err = cur.Decode(&result)
			response = append(response, result)
			if err != nil {
				fmt.Println("error", err.Error())
				return response, err
			}
		}
	}
	return response, nil
}

func (d *MongoClient) GetOne(dbName, collectionName string, filter interface{}) (bson.D, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	var response bson.D
	defer cancel()
	collection := d.Client.Database(dbName).Collection(collectionName)
	result := collection.FindOne(ctx, filter)
	err := result.Decode(&response)
	if err != nil {
		fmt.Println("error", err.Error())
		return response, err
	}
	return response, nil
}

func (d *MongoClient) UpSert(dbName, collectionName string, filter, updatedDoc interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := d.Client.Database(dbName).Collection(collectionName)
	result, err := collection.UpdateOne(ctx, filter, updatedDoc)
	if err != nil {
		fmt.Println("error", err.Error())
		return err
	}
	fmt.Println("Updated Document in MongoDB", result.MatchedCount, result.UpsertedCount)
	return nil
}

func (d *MongoClient) DeleteOne(dbName, collectionName string, filter interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := d.Client.Database(dbName).Collection(collectionName)
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println("error", err.Error())
		return err
	}
	fmt.Println("Deleted Document from MongoDB", result.DeletedCount)
	return nil
}

func (d *MongoClient) InsertOne(dbName, collectionName string, doc interface{}) (bson.D, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	var res bson.D
	defer cancel()
	collection := d.Client.Database(dbName).Collection(collectionName)
	result, err := collection.InsertOne(ctx, doc)
	if err != nil {
		fmt.Println("error", err.Error())
		return res, err
	}
	fmt.Println("Inserted Document to MongoDB", result.InsertedID)
	return res, nil
}

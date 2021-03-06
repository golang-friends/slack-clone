package models

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-friends/slack-clone/authservice/configs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database represents the information needed to connect to a database
type Database struct {
	DBURI      string `json:"dburi"`
	DBName     string `json:"dbname"`
	DBUser     string `json:"dbuser"`
	DBPassword string `json:"dbpassword"`
}

// Db ...
var Db mongo.Database

// ConnectToDB ...
func ConnectToDB(conf configs.Configuration) {
	db := Database{
		DBURI:      conf.DB.DBURI,
		DBUser:     conf.DB.DBUser,
		DBPassword: conf.DB.DBPass,
		DBName:     conf.DB.DBName,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s", db.DBUser, db.DBPassword, db.DBURI)))
	if err != nil {
		log.Fatal("Error connecting to DB: ", err.Error())
	}
	Db = *client.Database(db.DBName)
}

// ConnectToTestDB takes in a url string for testing purposes
func ConnectToTestDB(url string) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		log.Fatal("Error connecting to DB: ", err.Error())
	}
	Db = *client.Database("auth_service_test")
	collections, _ := Db.ListCollectionNames(ctx, bson.M{})
	for _, collection := range collections {
		Db.Collection(collection).Drop(ctx)
	}
}

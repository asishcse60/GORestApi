package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
	"time"
)
var once sync.Once
var dbClientStorage *MongoDBClientStorage

//need a refactor for ctx instance
func GetDatabaseClientStorage() (*MongoDBClientStorage, error) {
	//Set the client options
	once.Do(func() {
		dbClientStorage = &MongoDBClientStorage{}

		clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
		client, err := mongo.NewClient(clientOptions )
		if err != nil{
			log.Fatal(err)
		}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		err = client.Connect(ctx)
		if err != nil {
			log.Fatal(err)
		}
		//defer client.Disconnect(ctx)
		 dbClientStorage.Client = client
		 dbClientStorage.DatabaseName = "inventory"
		//ping the connection
		err = client.Ping(ctx, nil)
		if err != nil{
			log.Fatal(err)
		} else{
			log.Print("Ping Connection Success!")
		}
		//defer client.Disconnect(ctx)
	})

	fmt.Println("Congratulations! Mongodb is connected!")

	return dbClientStorage, nil
}
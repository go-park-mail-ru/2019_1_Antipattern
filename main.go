package main

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

type User struct {
	Id           primitive.ObjectID `bson:"_id"`
	Mail         string             `bson:"mail, omitempty"`
	Login        string             `bson:"login, omitempty"`
	PasswordHash string             `bson:"passwordhash, omitempty"`
	Score        int                `bson:"score, omitempty"`
}

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://api_db:27017"))

	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}
	options := options.Find()
	options.SetLimit(10)

	options.Sort = bson.M{"score": -1, "login": 1}
	options.SetProjection(bson.M{"passwordhash": 0})
	collection := client.Database("kpacubo").Collection("users")
	//ash := User{"abc@mail.ru", "log", "12345", 103}
	//_, err = collection.InsertOne(ctx, ash)
	if err != nil {
		log.Fatal(err)
	}

	cursor, err := collection.Find(ctx, bson.D{{"$or", bson.A{
		bson.D{{"score", 103}}, bson.D{{"score", 10}}}}}, options)

	if err != nil {
		log.Fatal(err)
	}
	for cursor.Next(ctx) {
		var user User
		err = cursor.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(user)
	}
	//fmt.Println("Inserted a single document: ", insertResult.InsertedID)

}

package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

func main() {
	h := sha256.New()
	users := map[string]string{
		"admin":      "fCRmh4Q2J7Rseqkz",
		"packt":      "RE4zfHB35VPtTkbT",
		"mlabouardy": "L3nSFRcZzNQ67bcc",
	}

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://admin:1596034420ze@localhost:27017/test?authSource=admin"))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	collection := client.Database("demo").Collection("users")

	// remember we can't just change passwords' type with string after they are hashed;
	// because the type maybe not correctly for mongodb to store.
	for username, password := range users {
		h.Write([]byte(password))
		b := h.Sum(nil)
		collection.InsertOne(ctx, bson.M{
			"username": username,
			"password": base64.StdEncoding.EncodeToString(b),
		})
	}
}
// Recipes API
//
// This is a sample recipes API.
// You can find out more about the API at https://github.com/luoshengyue/gin_project.
//
// Schemes: http
// Host: localhost:8080
// BasePath: /
// Version: 1.0.0
// Contact: luoyuehe <bancangbaize@gmail.com> https://sodasweet.cn
//
// Consumes:
// - application/json
// Produces:
// - application/json
// swagger:meta
package main

import (
	"context"
	"encoding/json"
	handlers "gin_project/handlers"
	models "gin_project/models"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"io/ioutil"
	"log"
)

var MONGO_URI ="mongodb://admin:1596034420ze@localhost:27017/test?authSource=admin"
var MONGO_DATABASE = "demo"

var recipesHandler *handlers.RecipesHandler
var recipes []models.Recipe

func init()  {
	recipes := make([]models.Recipe, 0)
	file, _ := ioutil.ReadFile("recipes.json")
	_ = json.Unmarshal([]byte(file), &recipes)

	// 连接到 mongoDB；
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MONGO_URI))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")

	var listOfRecipes []interface{}
	for _, recipe := range recipes {
		// Here we use primitive.NewObjectID() to make sure recipe.ID is the type of "primitive.Object"
		// if we don't do this, the data type of ID in Json file is string instead of "primitive.Object"
		// And we set Bson type of data, it will make id read from file is zero;
		// So we need to do this function, so that we can make listOfRecipes is correctly to Insert DB.
		recipe.ID = primitive.NewObjectID()
		listOfRecipes = append(listOfRecipes, recipe)
	}

	collection := client.Database(MONGO_DATABASE).Collection("recipes")
	insertManyResult, err := collection.InsertMany(ctx, listOfRecipes)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Insert recipes: ", len(insertManyResult.InsertedIDs))

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	status := redisClient.Ping()
	log.Println(status)
	// Here we must execute FlushDB() to clean our Redis cache,
	// If not, the result of response may be not correct.
	redisClient.FlushDB()

	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)
}

func main() {
	router := gin.Default()
	router.POST("/recipes", recipesHandler.NewRecipeHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id",recipesHandler.DeleteRecipeHandler)
	router.Run()
}
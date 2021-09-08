package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"io/ioutil"
	"net/http"
	"time"
)

type Recipe struct {
	ID  string  `json:"id"`
	Name    string  `json:"name"`
	Tags    []string    `json:"tags"`
	Ingredients     []string    `json:"ingredients"`
	Instruction     []string    `json:"instruction"`
	PublishedAt     time.Time   `json:"publishedAt"`
}

func NewRecipeHandler(c *gin.Context)  {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":err.Error(),
		})
		return
	}
	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusOK, recipe)
}

func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

var recipes []Recipe

func init()  {
	recipes = make([]Recipe, 0)
	file, _ := ioutil.ReadFile("F:\\新下载\\Building-Distributed-Applications-in-Gin-main\\chapter02\\recipes.json")
	_ = json.Unmarshal([]byte(file), &recipes)
}

func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.Run()
}
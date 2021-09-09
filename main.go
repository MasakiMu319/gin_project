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

// NewRecipeHandler implemented HTTP POST protocol
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

// ListRecipesHandler implemented HTTP GET protocol
func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

// UpdateRecipeHandler implemented HTTP PUT protocol
func UpdateRecipeHandler(c *gin.Context)  {
	id := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error":err.Error(),
		})
		return
	}
	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipes not found",
		})
		return
	}
	recipes[index] = recipe
	c.JSON(http.StatusOK, recipe)
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
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.Run()
}
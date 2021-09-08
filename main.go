package gin_project

import (
	"github.com/gin-gonic/gin"
	"time"
)

type Recipe struct {
	Name    string  `json:"name"`
	Tags    []string    `json:"tags"`
	Ingredients     []string    `json:"ingredients"`
	Instruction     []string    `json:"instruction"`
	PublishedAt     time.Time   `json:"publishedAt"`
}

func main() {
	router := gin.Default()
	router.Run()
}
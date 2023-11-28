package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func helloworld(c *gin.Context) {
	c.String(http.StatusOK, "Hello World! Time : %s", time.Now().Format(time.RFC3339Nano))
}

func main() {
	engine := gin.New()
	engine.GET("/", helloworld)

	engine.Run("localhost:8080")
}

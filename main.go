package main

import (
	"net/http"
	"zanzibar-dag/config"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.ReadConfig(); err != nil {
		panic(err.Error())
	}

	server := gin.Default()
	server.GET("/healthy", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})

	relationRoute := server.Group("/relation")
	{
		relationRoute.GET("/")
		relationRoute.POST("/")
		relationRoute.DELETE("/")

		relationRoute.POST("/check")
		relationRoute.POST("/get-shortest-path")
		relationRoute.POST("/get-all-paths")
		relationRoute.POST("/get-all-object-relations")
		relationRoute.POST("/clear-all-relations")
	}

	server.Run()
}

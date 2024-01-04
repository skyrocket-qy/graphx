package main

import (
	"net/http"
	"zanzibar-dag/config"
	"zanzibar-dag/internal/delivery"
	"zanzibar-dag/internal/infra/sql"
	"zanzibar-dag/internal/usecase"

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

	db, err := sql.InitDb()
	if err != nil {
		panic(err)
	}

	sqlRepo, err := sql.NewOrmRepository(db)
	if err != nil {
		panic(err)
	}

	usecaseRepo := usecase.NewUsecaseRepository(sqlRepo)

	handlerRepo := delivery.NewHandlerRepository(usecaseRepo)

	relationHandler := handlerRepo.RelationHandler
	relationRoute := server.Group("/relation")
	{
		relationRoute.GET("/", relationHandler.GetAll)
		relationRoute.POST("/", relationHandler.Create)
		relationRoute.DELETE("/", relationHandler.Delete)

		relationRoute.POST("/check", relationHandler.Check)
		relationRoute.POST("/get-shortest-path", relationHandler.GetShortestPath)
		relationRoute.POST("/get-all-paths", relationHandler.GetAllPaths)
		relationRoute.POST("/get-all-object-relations", relationHandler.GetAllObjectRelations)
		relationRoute.POST("/get-all-subbject-relations", relationHandler.GetAllSubjectRelations)
		relationRoute.POST("/clear-all-relations", relationHandler.ClearAllRelations)
	}

	server.Run()
}

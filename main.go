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
	relationRouter := server.Group("/relation")
	{
		relationRouter.GET("/", relationHandler.GetAll)
		relationRouter.POST("/", relationHandler.Create)
		relationRouter.DELETE("/", relationHandler.Delete)

		relationRouter.POST("get-all-namespaces", relationHandler.GetAllNamespaces)
		relationRouter.POST("/check", relationHandler.Check)
		relationRouter.POST("/get-shortest-path", relationHandler.GetShortestPath)
		relationRouter.POST("/get-all-paths", relationHandler.GetAllPaths)
		relationRouter.POST("/get-all-object-relations", relationHandler.GetAllObjectRelations)
		relationRouter.POST("/get-all-subbject-relations", relationHandler.GetAllSubjectRelations)
		relationRouter.POST("/clear-all-relations", relationHandler.ClearAllRelations)
	}

	server.Run()
}

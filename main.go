package main

import (
	"net/http"
	"zanzibar-dag/config"
	"zanzibar-dag/internal/delivery"
	"zanzibar-dag/internal/infra/sql"
	"zanzibar-dag/internal/usecase"

	"zanzibar-dag/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Swagger API
// @version         1.0
// @description
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	docs.SwaggerInfo.Host = "localhost:8080"

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
		relationRouter.GET("/query", relationHandler.Query)
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

	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	server.Run(":8080")
}

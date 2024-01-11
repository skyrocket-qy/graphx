package main

import (
	"log"
	"net"
	"net/http"

	"github.com/skyrocketOoO/zanazibar-dag/config"
	"github.com/skyrocketOoO/zanazibar-dag/docs"
	"github.com/skyrocketOoO/zanazibar-dag/internal/delivery"
	"github.com/skyrocketOoO/zanazibar-dag/internal/delivery/proto"
	"github.com/skyrocketOoO/zanazibar-dag/internal/infra/sql"
	"github.com/skyrocketOoO/zanazibar-dag/internal/usecase"
	"google.golang.org/grpc"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

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
		relationRouter.GET("/", relationHandler.Get)
		relationRouter.POST("/", relationHandler.Create)
		relationRouter.DELETE("/", relationHandler.Delete)

		relationRouter.POST("/delete-by-queries", relationHandler.DeleteByQueries)
		relationRouter.POST("/batch-operation", relationHandler.BatchOperation)

		relationRouter.POST("get-all-namespaces", relationHandler.GetAllNamespaces)
		relationRouter.POST("/check", relationHandler.Check)
		relationRouter.POST("/get-shortest-path", relationHandler.GetShortestPath)
		relationRouter.POST("/get-all-paths", relationHandler.GetAllPaths)
		relationRouter.POST("/get-all-object-relations", relationHandler.GetAllObjectRelations)
		relationRouter.POST("/get-all-subject-relations", relationHandler.GetAllSubjectRelations)
		relationRouter.POST("/clear-all-relations", relationHandler.ClearAllRelations)
	}

	go func() {
		grpcServer := grpc.NewServer()
		proto.RegisterRelationServiceServer(grpcServer, proto.NewRelationHandler(usecaseRepo.RelationUsecase))
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// /swagger/index.html
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	server.Run(":8080")
}

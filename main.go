package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/skyrocketOoO/zanazibar-dag/config"
	"github.com/skyrocketOoO/zanazibar-dag/docs"
	"github.com/skyrocketOoO/zanazibar-dag/internal/delivery"
	"github.com/skyrocketOoO/zanazibar-dag/internal/delivery/proto"
	"github.com/skyrocketOoO/zanazibar-dag/internal/delivery/rest"
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

	var wg sync.WaitGroup

	server := gin.Default()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: server,
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.GET("/healthy", func(c *gin.Context) {
			c.JSON(http.StatusOK, nil)
		})
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

		vd := rest.NewVisualDelivery(*usecase.NewVisualUsecase())
		server.GET("/visual", vd.SeeTree)

		//swagger/index.html
		server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		server.GET("/proto/doc", func(c *gin.Context) {
			c.File("grpc-doc/index.html")
		})

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	wg.Add(1)
	grpcServer := grpc.NewServer()
	go func() {
		defer wg.Done()
		proto.RegisterRelationServiceServer(grpcServer, proto.NewRelationHandler(usecaseRepo.RelationUsecase))
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		defer lis.Close()

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	// Block until a signal is received
	<-quit

	// Graceful shutdown for both servers
	log.Println("Received signal. Shutting down...")

	// Shut down the Gin server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	// Shut down the gRPC server
	grpcServer.GracefulStop()

	// Wait for all goroutines to finish
	wg.Wait()

	log.Println("Graceful shutdown complete.")
}

package api

import (
	"github.com/skyrocketOoO/zanazibar-dag/internal/delivery/rest"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Binding(r *gin.Engine, d *rest.Delivery) {
	r.GET("/ping", d.Ping)
	r.GET("/healthy", d.Healthy)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	relRouter := r.Group("/relation")
	{
		relRouter.GET("/", d.Get)
		relRouter.POST("/", d.Create)
		relRouter.DELETE("/", d.Delete)

		relRouter.DELETE("/delete-by-queries", d.DeleteByQueries)
		relRouter.POST("/batch-operation", d.BatchOperation)

		relRouter.GET("get-all-namespaces", d.GetAllNs)
		relRouter.GET("/check", d.Check)
		relRouter.GET("/get-shortest-path", d.GetShortestPath)
		relRouter.GET("/get-all-paths", d.GetAllPaths)
		relRouter.GET("/get-all-object-relations", d.GetAllObjRels)
		relRouter.GET("/get-all-subject-relations", d.GetAllSbjRels)
		relRouter.DELETE("/all-edges", d.ClearAllEdges)
	}
}

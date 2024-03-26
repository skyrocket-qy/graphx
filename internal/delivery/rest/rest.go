package rest

import (
	"go-server-template/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RestDelivery struct {
	usecase domain.Usecase
}

func NewRestDelivery(usecase domain.Usecase) *RestDelivery {
	return &RestDelivery{
		usecase: usecase,
	}
}

// @Summary Check the server started
// @Accept json
// @Produce json
// @Success 200 {object} domain.Response
// @Router /ping [get]
func (d *RestDelivery) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, domain.Response{Message: "pong"})
}

// @Summary Check the server healthy
// @Accept json
// @Produce json
// @Success 200 {object} domain.Response
// @Failure 503 {object} domain.Response
// @Router /healthy [get]
func (d *RestDelivery) Healthy(c *gin.Context) {
	// do something check
	if err := d.usecase.Healthy(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, domain.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.Response{Message: "healthy"})
}

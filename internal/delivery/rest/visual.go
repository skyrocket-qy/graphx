package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/skyrocketOoO/zanazibar-dag/domain"
	"github.com/skyrocketOoO/zanazibar-dag/internal/usecase"
)

type VisualDelivery struct {
	VisualUsecase usecase.VisualUsecase
}

func NewVisualDelivery(visualUsecase usecase.VisualUsecase) *VisualDelivery {
	return &VisualDelivery{
		VisualUsecase: visualUsecase,
	}
}

func (d *VisualDelivery) SeeTree(c *gin.Context) {
	address, _ := d.VisualUsecase.SeeTree(c.Request.Context(), domain.Node{}, 1)
	c.File(address)
}

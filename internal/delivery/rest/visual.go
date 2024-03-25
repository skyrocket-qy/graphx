package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/skyrocketOoO/zanazibar-dag/internal/usecase"
)

type VisualDelivery struct {
	VisualUsecase usecase.VisualUsecase
}

func NewVisualDelivery() *VisualDelivery {
	return &VisualDelivery{}
}

func (h *VisualDelivery) SeeTree(c *gin.Context) {
	c.File("internal/usecase/html/tree.html")
}

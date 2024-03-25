package usecase

import (
	"context"

	"github.com/skyrocketOoO/zanazibar-dag/domain"
)

type VisualUsecase struct {
}

func NewVisualUsecase() *VisualUsecase {
	return &VisualUsecase{}
}

func (u *VisualUsecase) SeeTree(c context.Context, node domain.Node, maxDepth int) (string, error) {
	return "internal/usecase/html/tree.html", nil
}

package delivery

import (
	"github.com/skyrocketOoO/zanazibar-dag/internal/delivery/rest"
	"github.com/skyrocketOoO/zanazibar-dag/internal/usecase"
)

type HandlerRepository struct {
	RelationHandler rest.RelationHandler
}

func NewHandlerRepository(ucRepo *usecase.UsecaseRepository) *HandlerRepository {
	return &HandlerRepository{
		RelationHandler: *rest.NewRelationHandler(ucRepo.RelationUsecase),
	}
}

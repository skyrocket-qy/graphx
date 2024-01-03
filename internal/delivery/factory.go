package delivery

import (
	"zanzibar-dag/internal/usecase"
)

type HandlerRepository struct {
	RelationHandler RelationHandler
}

func NewHandlerRepository(ucRepo *usecase.UsecaseRepository) *HandlerRepository {
	return &HandlerRepository{
		RelationHandler: *NewRelationHandler(ucRepo.RelationUsecase),
	}
}

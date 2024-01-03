package usecase

import (
	ucdomain "zanzibar-dag/domain/usecase"
	"zanzibar-dag/internal/infra/sql"
)

type UsecaseRepository struct {
	RelationUsecase ucdomain.RelationUsecase
}

func NewUsecaseRepository(sqlRepo *sql.OrmRepository) *UsecaseRepository {
	return &UsecaseRepository{
		RelationUsecase: NewRelationUsecase(&sqlRepo.RelationshipRepo),
	}
}

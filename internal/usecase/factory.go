package usecase

import (
	ucdomain "github.com/skyrocketOoO/zanazibar-dag/domain/usecase"
	"github.com/skyrocketOoO/zanazibar-dag/internal/infra/sql"
)

type UsecaseRepository struct {
	RelationUsecase ucdomain.RelationUsecase
}

func NewUsecaseRepository(sqlRepo *sql.OrmRepository) *UsecaseRepository {
	return &UsecaseRepository{
		RelationUsecase: NewRelationUsecase(&sqlRepo.RelationshipRepo),
	}
}

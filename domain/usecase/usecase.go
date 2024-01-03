package usecasedom

import (
	"zanzibar-dag/domain"
)

type RelationUsecase interface {
	GetAll() ([]string, error)
	Query(relation domain.Relation) ([]domain.Relation, error)
	Create(relation domain.Relation) error
	Delete(relation domain.Relation) error

	Check(from domain.Node, to domain.Node) (bool, error)
	GetShortestPath(from domain.Node, to domain.Node) ([]string, error)
	GetAllPaths(from domain.Node, to domain.Node) ([]string, error)
	GetObjectRelations(from domain.Node) ([]string, error)
	ClearAllRelations() error
}

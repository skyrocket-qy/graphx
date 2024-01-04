package usecasedom

import (
	"zanzibar-dag/domain"
)

type RelationUsecase interface {
	GetAll() ([]domain.Relation, error)
	Query(relation domain.Relation) ([]domain.Relation, error)
	Create(relation domain.Relation) error
	Delete(relation domain.Relation) error

	GetAllNamespaces() ([]string, error)
	Check(from domain.Node, to domain.Node) (bool, error)
	GetShortestPath(from domain.Node, to domain.Node) ([]domain.Relation, error)
	GetAllPaths(from domain.Node, to domain.Node) ([][]domain.Relation, error)
	GetAllObjectRelations(subject domain.Node) ([]domain.Relation, error)
	GetAllSubjectRelations(object domain.Node) ([]domain.Relation, error)

	ClearAllRelations() error
}

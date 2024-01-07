package sqldom

import (
	"zanzibar-dag/domain"
)

type Relation struct {
	ID               uint   `gorm:"primaryKey"`
	ObjectNamespace  string `gorm:"index:idx_object"`
	ObjectName       string `gorm:"index:idx_object"`
	Relation         string `gorm:"index:idx_object"`
	SubjectNamespace string `gorm:"index:idx_subject"`
	SubjectName      string `gorm:"index:idx_subject"`
	SubjectRelation  string `gorm:"index:idx_subject"`
	AllColumns       string `gorm:"unique"`
}

type RelationRepository interface {
	Create(relation domain.Relation) error
	Delete(relation domain.Relation) error
	BatchOperation(operations []domain.Operation) error
	GetAll() ([]domain.Relation, error)
	Query(query domain.Relation) ([]domain.Relation, error)
	GetAllNamespaces() ([]string, error)
	DeleteAll() error
}

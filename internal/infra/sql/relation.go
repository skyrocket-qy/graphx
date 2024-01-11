package sql

import (
	"errors"
	"strings"

	"github.com/skyrocketOoO/go-utility/set"
	"github.com/skyrocketOoO/zanazibar-dag/domain"
	sqldom "github.com/skyrocketOoO/zanazibar-dag/domain/infra/sql"

	"gorm.io/gorm"
)

type RelationRepository struct {
	DB *gorm.DB
}

func NewRelationRepository(db *gorm.DB) *RelationRepository {
	return &RelationRepository{DB: db}
}

func (r *RelationRepository) Create(relation domain.Relation) error {
	sqlRelation := convertToSqlModel(relation)
	return r.DB.Create(&sqlRelation).Error
}

func (r *RelationRepository) Delete(relation domain.Relation) error {
	return r.DB.Where("all_columns = ?", concatAttr(relation)).Delete(&sqldom.Relation{}).Error
}

func (r *RelationRepository) DeleteByQueries(queries []domain.Relation) error {
	operations := set.NewSet[domain.Operation]()
	for _, query := range queries {
		relations, err := r.Query(query)
		if err != nil {
			return err
		}
		for _, relation := range relations {
			operations.Add(domain.Operation{
				Type:     domain.DeleteOperation,
				Relation: relation,
			})
		}
	}

	return r.BatchOperation(operations.ToSlice())
}

func (r *RelationRepository) BatchOperation(operations []domain.Operation) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, operation := range operations {
		switch operation.Type {
		case domain.CreateOperation:
			if err := r.Create(operation.Relation); err != nil {
				tx.Rollback()
				return err
			}
		case domain.DeleteOperation:
			if err := r.Delete(operation.Relation); err != nil {
				tx.Rollback()
				return err
			}
		case domain.CreateIfNotExistOperation:
			if err := r.Create(operation.Relation); err != nil {
				if err != gorm.ErrDuplicatedKey {
					tx.Rollback()
					return err
				}
			}
		default:
			tx.Rollback()
			return errors.New("invalid operation type")
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *RelationRepository) GetAll() ([]domain.Relation, error) {
	// var relations []sqldom.Relation
	// var err error

	// done := make(chan bool)
	// go func() {
	// 	defer close(done)
	// 	err = r.DB.Find(&relations).Error
	// }()
	// <-done

	// if err != nil {
	// 	return nil, err
	// }

	var relations []sqldom.Relation
	if err := r.DB.Find(&relations).Error; err != nil {
		return nil, err
	}

	newRelations := make([]domain.Relation, len(relations))
	for i, relation := range relations {
		newRelations[i] = convertToRelation(relation)
	}
	return newRelations, nil
}

func (r *RelationRepository) Query(query domain.Relation) ([]domain.Relation, error) {
	var relations []sqldom.Relation
	if err := r.DB.Where(&query).Find(&relations).Error; err != nil {
		return nil, err
	}
	newRelations := make([]domain.Relation, len(relations))
	for i, relation := range relations {
		newRelations[i] = convertToRelation(relation)
	}
	return newRelations, nil
}

func (r *RelationRepository) GetAllNamespaces() ([]string, error) {
	sqlQuery := `
		SELECT DISTINCT namespace
		FROM (
			SELECT object_namespace AS namespace FROM relations
			UNION
			SELECT subject_namespace AS namespace FROM relations
		) AS namespaces
	`
	var namespaces []string
	if err := r.DB.Raw(sqlQuery).Scan(&namespaces).Error; err != nil {
		return nil, err
	}

	return namespaces, nil
}

func (r *RelationRepository) DeleteAll() error {
	query := "DELETE FROM relations"
	if err := r.DB.Exec(query).Error; err != nil {
		return err
	}
	return nil
}

func convertToSqlModel(relation domain.Relation) sqldom.Relation {
	return sqldom.Relation{
		ObjectNamespace:  relation.ObjectNamespace,
		ObjectName:       relation.ObjectName,
		Relation:         relation.Relation,
		SubjectNamespace: relation.SubjectNamespace,
		SubjectName:      relation.SubjectName,
		SubjectRelation:  relation.SubjectRelation,
		AllColumns:       concatAttr(relation),
	}
}

func concatAttr(relation domain.Relation) string {
	return strings.Join(
		[]string{
			relation.ObjectNamespace,
			relation.ObjectName,
			relation.Relation,
			relation.SubjectNamespace,
			relation.SubjectName,
			relation.SubjectRelation,
		},
		"%",
	)
}

func convertToRelation(relation sqldom.Relation) domain.Relation {
	return domain.Relation{
		ObjectNamespace:  relation.ObjectNamespace,
		ObjectName:       relation.ObjectName,
		Relation:         relation.Relation,
		SubjectNamespace: relation.SubjectNamespace,
		SubjectName:      relation.SubjectName,
		SubjectRelation:  relation.SubjectRelation,
	}
}

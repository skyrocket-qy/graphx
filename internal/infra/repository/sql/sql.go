package sql

import (
	"context"
	"strings"

	errors "github.com/rotisserie/eris"
	"github.com/skyrocketOoO/go-utility/set"
	"github.com/skyrocketOoO/zanazibar-dag/domain"
	"gorm.io/gorm"
)

type SqlRepository struct {
	db *gorm.DB
}

func NewSqlRepository(db *gorm.DB) (*SqlRepository, error) {
	return &SqlRepository{
		db: db,
	}, nil
}

func (r *SqlRepository) Ping(c context.Context) error {
	db, err := r.db.DB()
	if err != nil {
		return errors.New(err.Error())
	}
	if err := db.PingContext(c); err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func (r *SqlRepository) Create(edge domain.Edge) error {
	sqlRel := convertToSqlModel(edge)
	if err := r.db.Create(&sqlRel).Error; err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func (r *SqlRepository) Delete(edge domain.Edge) error {
	if err := r.db.Where("all_columns = ?", concatAttr(edge)).Delete(&Edge{}).
		Error; err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func (r *SqlRepository) DeleteByQueries(queries []domain.Edge) error {
	operations := set.NewSet[domain.Operation]()
	for _, query := range queries {
		edges, err := r.Query(query)
		if err != nil {
			return err
		}
		for _, edge := range edges {
			operations.Add(domain.Operation{
				Type: domain.DeleteOperation,
				Edge: edge,
			})
		}
	}

	return r.BatchOperation(operations.ToSlice())
}

func (r *SqlRepository) BatchOperation(operations []domain.Operation) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, operation := range operations {
		switch operation.Type {
		case domain.CreateOperation:
			if err := r.Create(operation.Edge); err != nil {
				tx.Rollback()
				return err
			}
		case domain.DeleteOperation:
			if err := r.Delete(operation.Edge); err != nil {
				tx.Rollback()
				return err
			}
		case domain.CreateIfNotExistOperation:
			if err := r.Create(operation.Edge); err != nil {
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

func (r *SqlRepository) GetAll(options ...domain.PageOptions) ([]domain.Edge,
	uint, error) {
	var edges []Edge
	var err error
	if len(options) > 0 {
		options := options[0]
		err = r.db.Where("id > ?", options.LastID).Order("id").
			Limit(options.PageSize).Find(&edges).Error
	} else {
		err = r.db.Find(&edges).Error
	}
	if err != nil {
		return nil, 0, err
	}

	newEdges := make([]domain.Edge, len(edges))
	for i, edge := range edges {
		newEdges[i] = convertToRel(edge)
	}
	return newEdges, edges[len(edges)-1].ID, nil
}

func (r *SqlRepository) Query(query domain.Edge) ([]domain.Edge, error) {
	var edges []Edge
	if err := r.db.Where(&query).Find(&edges).Error; err != nil {
		return nil, err
	}
	newEdges := make([]domain.Edge, len(edges))
	for i, edge := range edges {
		newEdges[i] = convertToRel(edge)
	}
	return newEdges, nil
}

func (r *SqlRepository) GetAllNs() ([]string, error) {
	sqlQuery := `
		SELECT DISTINCT namespace
		FROM (
			SELECT object_namespace AS namespace FROM edges
			UNION
			SELECT subject_namespace AS namespace FROM edges
		) AS namespaces
	`
	var nss []string
	if err := r.db.Raw(sqlQuery).Scan(&nss).Error; err != nil {
		return nil, err
	}

	return nss, nil
}

func (r *SqlRepository) DeleteAll() error {
	query := "DELETE FROM edges"
	if err := r.db.Exec(query).Error; err != nil {
		return err
	}
	return nil
}

func convertToSqlModel(edge domain.Edge) Edge {
	return Edge{
		ObjNs:      edge.ObjNs,
		ObjName:    edge.ObjName,
		ObjRel:     edge.ObjRel,
		SbjNs:      edge.SbjNs,
		SbjName:    edge.SbjName,
		SbjRel:     edge.SbjRel,
		AllColumns: concatAttr(edge),
	}
}

func concatAttr(edge domain.Edge) string {
	return strings.Join(
		[]string{
			edge.ObjNs,
			edge.ObjName,
			edge.ObjRel,
			edge.SbjNs,
			edge.SbjName,
			edge.SbjRel,
		},
		"%",
	)
}

func convertToRel(edge Edge) domain.Edge {
	return domain.Edge{
		ObjNs:   edge.ObjNs,
		ObjName: edge.ObjName,
		ObjRel:  edge.ObjRel,
		SbjNs:   edge.SbjNs,
		SbjName: edge.SbjName,
		SbjRel:  edge.SbjRel,
	}
}

package domain

import "context"

type SqlRepository interface {
	Ping(c context.Context) error
	Create(edge Edge) error
	Delete(edge Edge) error
	DeleteByQueries(queries []Edge) error
	BatchOperation(operations []Operation) error
	GetAll(options ...PageOptions) (edges []Edge, lastID uint, err error)
	Query(query Edge) ([]Edge, error)
	GetAllNs() ([]string, error)
	DeleteAll() error
}

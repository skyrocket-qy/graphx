package sql

import (
	"context"

	"gorm.io/gorm"
)

type OrmRepository struct {
	db *gorm.DB
}

func NewOrmRepository(db *gorm.DB) (*OrmRepository, error) {
	return &OrmRepository{
		db: db,
	}, nil
}

func (r *OrmRepository) Ping(c context.Context) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	if err := db.PingContext(c); err != nil {
		return err
	}

	return nil
}

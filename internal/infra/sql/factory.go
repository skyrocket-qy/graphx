package sql

import (
	sqldomain "zanzibar-dag/domain/infra/sql"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDb() (*gorm.DB, error) {
	return gorm.Open(
		// sqlite.Open("gorm.db"), &gorm.Config{
		// 	Logger: nil,
		// },
		postgres.Open(viper.GetString("database.postgres.dsn")), &gorm.Config{
			Logger: nil,
		},
	)
}

type OrmRepository struct {
	RelationshipRepo RelationRepository
}

func NewOrmRepository(db *gorm.DB) (*OrmRepository, error) {
	if err := db.AutoMigrate(&sqldomain.Relation{}); err != nil {
		return nil, err
	}

	return &OrmRepository{
		RelationshipRepo: *NewRelationRepository(db),
	}, nil
}

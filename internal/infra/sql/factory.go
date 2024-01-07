package sql

import (
	"fmt"
	"os"
	sqldomain "zanzibar-dag/domain/infra/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDb() (*gorm.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
		getEnv("POSTGRES_HOST", "localhost"),
		getEnv("POSTGRES_PORT", "5432"),
		getEnv("POSTGRES_USER", "zanzibar-dag"),
		getEnv("POSTGRES_PASSWORD", "zanzibar-dag"),
		getEnv("POSTGRES_DB", "zanzibar-dag"),
		getEnv("POSTGRES_TIMEZONE", "Asia/Taipei"),
	)
	return gorm.Open(
		postgres.Open(connStr), &gorm.Config{
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

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

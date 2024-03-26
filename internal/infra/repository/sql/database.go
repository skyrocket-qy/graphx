package sql

import (
	"fmt"

	errors "github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDB(database string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	switch database {
	case "pg":
		log.Info().Msg("Connecting to Postgres")
		connStr := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
			viper.GetString("postgres.host"),
			viper.GetString("postgres.port"),
			viper.GetString("postgres.user"),
			viper.GetString("postgres.password"),
			viper.GetString("postgres.db"),
			viper.GetString("postgres.sslmode"),
			viper.GetString("postgres.timezone"),
		)
		db, err = gorm.Open(
			postgres.Open(connStr), &gorm.Config{
				Logger: nil,
			},
		)
	case "sqlite":
		log.Info().Msg("Connecting to Sqlite")
		db, err = gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	default:
		return nil, errors.New("database not supported")
	}
	if err != nil {
		return nil, errors.New(err.Error())
	}

	if err = db.AutoMigrate(&Edge{}); err != nil {
		return nil, errors.New(err.Error())
	}

	return db, nil
}

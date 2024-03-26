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
		return gorm.Open(
			postgres.Open(connStr), &gorm.Config{
				Logger: nil,
			},
		)
	case "sqlite":
		log.Info().Msg("Connecting to Sqlite")
		return gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	}
	return nil, errors.New("database not supported")
}

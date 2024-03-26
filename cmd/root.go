/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"time"

	"github.com/skyrocketOoO/zanazibar-dag/api"
	"github.com/skyrocketOoO/zanazibar-dag/config"
	"github.com/skyrocketOoO/zanazibar-dag/docs"
	"github.com/skyrocketOoO/zanazibar-dag/internal/delivery/rest"
	"github.com/skyrocketOoO/zanazibar-dag/internal/delivery/rest/middleware"
	"github.com/skyrocketOoO/zanazibar-dag/internal/infra/repository/sql"
	"github.com/skyrocketOoO/zanazibar-dag/internal/usecase"

	errors "github.com/rotisserie/eris"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// Here is the enum flag variable declaration
var flagDatabaseEnum DatabaseEnum

func workFunc(cmd *cobra.Command, args []string) {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}) // human-friendly logging without efficiency
	log.Info().Msg("Logger initialized")

	if err := config.ReadConfig(); err != nil {
		log.Fatal().Msg(errors.ToString(err, true))
	}

	docs.SwaggerInfo.Title = "Swagger API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/v2"
	docs.SwaggerInfo.Schemes = []string{"http"}

	sqlDb, err := sql.InitDB(string(flagDatabaseEnum))
	if err != nil {
		log.Fatal().Msg(errors.ToString(err, true))
	}
	defer func() {
		db, _ := sqlDb.DB()
		db.Close()
	}()

	sqlRepo, err := sql.NewSqlRepository(sqlDb)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	usecase := usecase.NewUsecase(sqlRepo)
	delivery := rest.NewDelivery(usecase)

	router := gin.Default()
	router.Use(middleware.CORS())
	api.Binding(router, delivery)

	port, _ := cmd.Flags().GetString("port")
	router.Run(":" + port)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "",
	Short: "A brief description of your application",
	Long:  `The longer description`,
	Run:   workFunc,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.Flags().StringP("port", "p", "8080", "port")
	rootCmd.Flags().Var(&flagDatabaseEnum, "database", `database enum. allowed: "pg", "sqlite"`)
}

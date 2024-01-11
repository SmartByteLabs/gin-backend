package main

import (
	"errors"
	"fmt"

	"github.com/princeparmar/9and9-templeCMS-backend.git/internal/config"
	"github.com/princeparmar/9and9-templeCMS-backend.git/internal/model"
	"github.com/princeparmar/9and9-templeCMS-backend.git/internal/router"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/ginhelper"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/logger"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/rbac/mysql"
	"github.com/spf13/cobra"
)

var rootCMD = &cobra.Command{}
var log = logger.New()

func init() {
	runServerCMD := &cobra.Command{
		Use: "runserver",
		RunE: func(cmd *cobra.Command, args []string) error {
			conf := config.NewConfigFromEnv()

			if conf.App.JWTSecret == "" {
				return errors.New("JWT secrete is required")
			}

			db := database.Connect(log, conf.DatabaseConfig)

			ginRouterHelper := router.NewRouterHandler(log, conf, db)
			log.Info("Starting server...")
			ginhelper.StartServer(log, db, conf.App.Port, ginRouterHelper.RegisterRoute)

			return nil
		},
	}

	rootCMD.AddCommand(runServerCMD)

	setupCommand := &cobra.Command{
		Use: "setup",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Setting up database...")
			conf := config.NewConfigFromEnv()
			db := database.Connect(log, conf.DatabaseConfig)

			err := mysql.CreateAllTables(cmd.Context(), db)
			if err != nil {
				return err
			}

			err = model.CreateAllTables(cmd.Context(), db)
			if err != nil {
				return err
			}

			fmt.Println("Database setup completed")

			return nil
		},
	}

	rootCMD.AddCommand(setupCommand)
}

func main() {
	if err := rootCMD.Execute(); err != nil {
		log.Fatal(err, "Failed to execute command")
	}
}

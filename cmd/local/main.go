package main

import (
	"context"
	"os"
	"watchtower/internal/database"
	"watchtower/internal/logging"
	"watchtower/internal/watchtower"
)

func main() {

	appConfig := watchtower.LoadConfig()
	logger := logging.FromContext(context.Background()).With("service", "local")

	_, db, err := database.NewDBFromProvider(appConfig.AppDir)
	if err != nil {
		logger.Error("Error creating database", "error", err)
		os.Exit(1)
	}

	defer func() {
		if err = db.Close(); err != nil {
			logger.Error("Error closing database", "error", err)
		}
	}()

	migrator := database.NewMigrator(db)

	if err = migrator.Init(); err != nil {
		logger.Error("Error running migrations", "error", err)
		os.Exit(1)
	}
}

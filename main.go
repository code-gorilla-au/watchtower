package main

import (
	"context"
	"embed"
	"os"
	"watchtower/internal/database"
	"watchtower/internal/watchtower"

	"watchtower/internal/logging"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

var (
	Version = "0.0.0-dev"
)

//go:embed all:frontend/build
var assets embed.FS

func main() {
	ctx := context.Background()
	appConfig := LoadConfig()
	logger := logging.New(appConfig.LogLevel)

	logger.Debug("Starting watchtower", "config", appConfig)

	databaseQueries, db, err := database.NewDBFromProvider(appConfig.DbFilePath)
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

	wt := watchtower.NewService(ctx, databaseQueries, db)
	worker := watchtower.NewOrgSyncWorker(wt)

	app := NewApp(worker)

	err = wails.Run(&options.App{
		Title:  "Watchtower",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind: []interface{}{
			app,
			wt,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}

}

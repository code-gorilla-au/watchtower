package main

import (
	"context"
	"embed"
	"os"
	"watchtower/internal/database"
	"watchtower/internal/watchtower"

	"github.com/code-gorilla-au/go-toolbox/logging"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/build
var assets embed.FS

func main() {
	ctx := context.Background()
	appConfig := LoadConfig()
	logger := logging.New(appConfig.LogLevel, logging.LoggerJSON)

	logger.Info("Starting watchtower", "config", appConfig)

	migrator, err := database.NewMigrator(appConfig.DbFilePath)
	if err != nil {
		logger.Error("Error creating database migrator", "error", err)
		os.Exit(1)
	}

	if err = migrator.Init(); err != nil {
		logger.Error("Error running migrations", "error", err)
		os.Exit(1)
	}

	db, err := database.NewDBFromProvider(appConfig.DbFilePath)
	if err != nil {
		logger.Error("Error creating database", "error", err)
		os.Exit(1)
	}

	wt := watchtower.NewService(ctx, db)
	worker := watchtower.NewOrgSyncWorker(wt)

	app := NewApp(worker)

	err = wails.Run(&options.App{
		Title:  "watchtower",
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

package main

import (
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/code-gorilla-au/env"
)

type Config struct {
	Env        string
	AppDir     string
	DbFilePath string
	LogLevel   slog.Level
}

const appDirPath = "watchtower"

func LoadConfig() Config {
	env.LoadEnvFile(".env.local")

	environment := os.Getenv("ENVIRONMENT")
	localDevDBPath := os.Getenv("LOCAL_DEV_DIR")
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	logLevel := env.GetAsIntWithDefault("LOG_LEVEL", 4)

	appDir := path.Join(homeDir, appDirPath)
	if environment == "local" {
		fmt.Print("LOCAL MODE")
		current, _ := os.Getwd()

		appDir = path.Join(current, localDevDBPath)
	} else {
		if err = os.Mkdir(appDir, 0755); err != nil {
			panic(err)
		}
	}

	return Config{
		Env:        environment,
		AppDir:     appDir,
		DbFilePath: appDir,
		LogLevel:   slog.Level(logLevel),
	}
}

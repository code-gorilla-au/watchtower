package watchtower

import (
	"log"
	"log/slog"
	"os"
	"path"

	"github.com/code-gorilla-au/env"
)

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
		log.Println("LOCAL MODE")
		current, _ := os.Getwd()

		appDir = path.Join(current, localDevDBPath)
	} else {
		// folder can already exist
		_ = os.Mkdir(appDir, 0750)
	}

	return Config{
		Env:      environment,
		AppDir:   appDir,
		LogLevel: slog.Level(logLevel),
	}
}

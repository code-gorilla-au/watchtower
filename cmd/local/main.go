package main

import (
	"fmt"
	"watchtower/internal/database"

	"github.com/code-gorilla-au/env"
)

func main() {
	env.LoadEnvFile(".env.local")
	env.WithStrictMode()

	path := env.GetAsString("LOCAL_DEV_DIR")
	db, err := database.NewMigrator(path)
	if err != nil {
		fmt.Print("Error:", err)
		panic(err)
	}

	err = db.Init()
	if err != nil {
		fmt.Print("Error:", err)
		panic(err)
	}
}

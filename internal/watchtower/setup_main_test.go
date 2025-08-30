package watchtower

import (
	"fmt"
	"os"
	"testing"
	"watchtower/internal/database"
)

var _testDB *database.Queries

func TestMain(m *testing.M) {
	q, db, err := database.NewDBFromProvider(":memory:")
	if err != nil {
		fmt.Print("Error:", err)
		os.Exit(1)
	}

	defer func() {
		if err = db.Close(); err != nil {
			fmt.Print("Error:", err)
		}
	}()

	migrator := database.NewMigrator(db)

	if err = migrator.Init(); err != nil {
		fmt.Print("Error running migrations", "error", err)
		os.Exit(1)
	}

	_testDB = q

	exitCode := m.Run()
	os.Exit(exitCode)
}

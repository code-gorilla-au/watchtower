package watchtower

import (
	"watchtower/internal/database"
)

type Service struct {
	db database.Queries
}

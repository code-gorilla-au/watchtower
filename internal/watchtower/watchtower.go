package watchtower

import (
	"context"
	"watchtower/internal/database"

	"github.com/code-gorilla-au/go-toolbox/github"
	"github.com/code-gorilla-au/go-toolbox/logging"
)

// NewService creates and returns a new Service instance with the provided database queries.
func NewService(ctx context.Context, db *database.Queries) *Service {
	return &Service{
		ghClient: github.New(logging.FromContext(ctx)),
		ctx:      ctx,
		db:       db,
	}
}

func (s *Service) Startup(ctx context.Context) {
	s.ctx = ctx
}

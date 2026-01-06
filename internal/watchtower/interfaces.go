package watchtower

import (
	"context"
	"watchtower/internal/database"

	"watchtower/internal/github"
)

//go:generate moq -rm -stub -out mocks.gen.go . ghClient

type ghClient interface {
	SearchRepos(owner string, topic string, token string) (github.QuerySearch[github.Repository], error)
	GetRepoDetails(owner string, repo string, token string) (github.QueryRepository, error)
}

type InsightsStore interface {
	GetPullRequestInsights(ctx context.Context, arg database.GetPullRequestInsightsParams) (database.GetPullRequestInsightsRow, error)
	GetSecuritiesInsights(ctx context.Context, arg database.GetSecuritiesInsightsParams) (database.GetSecuritiesInsightsRow, error)
}

var _ InsightsStore = (*database.Queries)(nil)

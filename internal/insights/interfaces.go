package insights

import (
	"context"
	"watchtower/internal/database"
)

type Store interface {
	GetPullRequestInsightsByOrg(ctx context.Context, arg database.GetPullRequestInsightsByOrgParams) (database.GetPullRequestInsightsByOrgRow, error)
	GetSecuritiesInsightsByOrg(ctx context.Context, arg database.GetSecuritiesInsightsByOrgParams) (database.GetSecuritiesInsightsByOrgRow, error)
}

var _ Store = (*database.Queries)(nil)

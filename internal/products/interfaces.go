package products

import (
	"context"
	"database/sql"
	"watchtower/internal/database"
)

type ProductStore interface {
	ProductBaseStore
	RepoStore
	PullRequestStore
	SecurityStore
}

var _ ProductStore = (*database.Queries)(nil)

type ProductBaseStore interface {
	CreateProduct(ctx context.Context, arg database.CreateProductParams) (database.Product, error)
	GetProductByID(ctx context.Context, id int64) (database.Product, error)
	ListProductsByOrganisation(ctx context.Context, organisationID sql.NullInt64) ([]database.Product, error)
	UpdateProduct(ctx context.Context, arg database.UpdateProductParams) (database.Product, error)
	UpdateProductSync(ctx context.Context, id int64) error
	DeleteProduct(ctx context.Context, id int64) error
}

type RepoStore interface {
	CreateRepo(ctx context.Context, arg database.CreateRepoParams) (database.Repository, error)
	UpdateRepo(ctx context.Context, arg database.UpdateRepoParams) (database.Repository, error)
	GetRepoByName(ctx context.Context, name string) (database.Repository, error)
	GetReposByProductID(ctx context.Context, id int64) ([]database.GetReposByProductIDRow, error)
	DeleteReposByProductID(ctx context.Context, id int64) error
}

type PullRequestStore interface {
	GetPullRequestByProductIDAndState(ctx context.Context, arg database.GetPullRequestByProductIDAndStateParams) ([]database.GetPullRequestByProductIDAndStateRow, error)
	GetPullRequestsByOrganisationAndState(ctx context.Context, arg database.GetPullRequestsByOrganisationAndStateParams) ([]database.GetPullRequestsByOrganisationAndStateRow, error)
	CreatePullRequest(ctx context.Context, arg database.CreatePullRequestParams) (database.PullRequest, error)
	UpdatePullRequest(ctx context.Context, arg database.UpdatePullRequestParams) (database.PullRequest, error)
	GetPullRequestByExternalID(ctx context.Context, externalID string) (database.PullRequest, error)
	GetRecentPullRequests(ctx context.Context) ([]database.GetRecentPullRequestsRow, error)
	DeletePullRequestsByProductID(ctx context.Context, id int64) error
}

type SecurityStore interface {
	GetSecurityByProductIDAndState(ctx context.Context, arg database.GetSecurityByProductIDAndStateParams) ([]database.GetSecurityByProductIDAndStateRow, error)
	GetSecurityByOrganisationAndState(ctx context.Context, arg database.GetSecurityByOrganisationAndStateParams) ([]database.GetSecurityByOrganisationAndStateRow, error)
	CreateSecurity(ctx context.Context, arg database.CreateSecurityParams) (database.Security, error)
	GetRecentSecurity(ctx context.Context) ([]database.GetRecentSecurityRow, error)
	UpdateSecurity(ctx context.Context, arg database.UpdateSecurityParams) (database.Security, error)
	GetSecurityByExternalID(ctx context.Context, externalID string) (database.Security, error)
	DeleteSecurityByProductID(ctx context.Context, id int64) error
}

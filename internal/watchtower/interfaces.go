package watchtower

import (
	"context"
	"database/sql"
	"watchtower/internal/database"

	"watchtower/internal/github"
)

//go:generate moq -rm -stub -out mocks.gen.go . ghClient RepoStore ProductStore OrgStore

type ghClient interface {
	SearchRepos(owner string, topic string, token string) (github.QuerySearch[github.Repository], error)
	GetRepoDetails(owner string, repo string, token string) (github.QueryRepository, error)
}

type RepoStore interface {
	CreateRepo(ctx context.Context, arg database.CreateRepoParams) (database.Repository, error)
	GetReposByProductID(ctx context.Context, productID int64) ([]database.Repository, error)
	GetPullRequestByProductIDAndState(ctx context.Context, arg database.GetPullRequestByProductIDAndStateParams) ([]database.PullRequest, error)
	GetPullRequestsByOrganisationAndState(ctx context.Context, arg database.GetPullRequestsByOrganisationAndStateParams) ([]database.PullRequest, error)
	CreatePullRequest(ctx context.Context, arg database.CreatePullRequestParams) (database.PullRequest, error)
	GetSecurityByProductIDAndState(ctx context.Context, arg database.GetSecurityByProductIDAndStateParams) ([]database.Security, error)
	GetSecurityByOrganisationAndState(ctx context.Context, arg database.GetSecurityByOrganisationAndStateParams) ([]database.Security, error)
	CreateSecurity(ctx context.Context, arg database.CreateSecurityParams) (database.Security, error)
}

type ProductStore interface {
	CreateProduct(ctx context.Context, arg database.CreateProductParams) (database.Product, error)
	GetProductByID(ctx context.Context, id int64) (database.Product, error)
	ListProductsByOrganisation(ctx context.Context, organisationID sql.NullInt64) ([]database.Product, error)
	UpdateProduct(ctx context.Context, arg database.UpdateProductParams) (database.Product, error)
	UpdateProductSync(ctx context.Context, id int64) error
	DeleteSecurityByProductID(ctx context.Context, id int64) error
	DeletePullRequestsByProductID(ctx context.Context, id int64) error
	DeleteReposByProductID(ctx context.Context, id int64) error
	DeleteProduct(ctx context.Context, id int64) error
}

type OrgStore interface {
	SetOrgsDefaultFalse(ctx context.Context) error
	CreateOrganisation(ctx context.Context, arg database.CreateOrganisationParams) (database.Organisation, error)
	GetOrganisationByID(ctx context.Context, id int64) (database.Organisation, error)
	GetDefaultOrganisation(ctx context.Context) (database.Organisation, error)
	SetDefaultOrg(ctx context.Context, id int64) (database.Organisation, error)
	ListOrganisations(ctx context.Context) ([]database.Organisation, error)
	ListOrgsOlderThanUpdatedAt(ctx context.Context, updatedAt int64) ([]database.Organisation, error)
	GetOrganisationForProduct(ctx context.Context, productID sql.NullInt64) (database.Organisation, error)
	DeleteOrg(ctx context.Context, id int64) error
	DeleteProductOrganisationByOrgID(ctx context.Context, organisationID sql.NullInt64) error
	UpdateOrganisation(ctx context.Context, arg database.UpdateOrganisationParams) (database.Organisation, error)
	UpdateProductSync(ctx context.Context, id int64) error
	AddProductToOrganisation(ctx context.Context, arg database.AddProductToOrganisationParams) error
}

var _ RepoStore = (*database.Queries)(nil)
var _ ProductStore = (*database.Queries)(nil)
var _ OrgStore = (*database.Queries)(nil)

package products

import (
	"context"
	"database/sql"
	"watchtower/internal/database"
)

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

	CreateRepo(ctx context.Context, arg database.CreateRepoParams) (database.Repository, error)
	GetReposByProductID(ctx context.Context, id int64) ([]database.GetReposByProductIDRow, error)
	GetPullRequestByProductIDAndState(ctx context.Context, arg database.GetPullRequestByProductIDAndStateParams) ([]database.GetPullRequestByProductIDAndStateRow, error)
	GetPullRequestsByOrganisationAndState(ctx context.Context, arg database.GetPullRequestsByOrganisationAndStateParams) ([]database.GetPullRequestsByOrganisationAndStateRow, error)
	CreatePullRequest(ctx context.Context, arg database.CreatePullRequestParams) (database.PullRequest, error)
	GetSecurityByProductIDAndState(ctx context.Context, arg database.GetSecurityByProductIDAndStateParams) ([]database.GetSecurityByProductIDAndStateRow, error)
	GetSecurityByOrganisationAndState(ctx context.Context, arg database.GetSecurityByOrganisationAndStateParams) ([]database.GetSecurityByOrganisationAndStateRow, error)
	CreateSecurity(ctx context.Context, arg database.CreateSecurityParams) (database.Security, error)
}

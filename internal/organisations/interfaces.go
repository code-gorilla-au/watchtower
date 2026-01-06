package organisations

import (
	"context"
	"database/sql"
	"watchtower/internal/database"
)

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
	AddProductToOrganisation(ctx context.Context, arg database.AddProductToOrganisationParams) error
}

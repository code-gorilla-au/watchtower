package watchtower

import (
	"context"
	"database/sql"
	"time"
	"watchtower/internal/database"

	"github.com/code-gorilla-au/go-toolbox/logging"
)

type CreateOrgParams struct {
	FriendlyName string
	Namespace    string
	Token        string
	Description  string
}

func (o organisationService) Create(ctx context.Context, params CreateOrgParams) (OrganisationDTO, error) {
	logger := logging.FromContext(ctx)
	logger.Debug("Creating organisation")

	var orgModel database.Organisation
	var err error

	err = database.WithTxnContext(ctx, o.txnDB, func(tx *sql.Tx) error {
		txnStore := o.txnFunc(tx)
		if err = txnStore.SetOrgsDefaultFalse(ctx); err != nil {
			logger.Error("Error setting default org", "error", err)

			return err
		}

		orgModel, err = txnStore.CreateOrganisation(ctx, database.CreateOrganisationParams{
			FriendlyName: params.FriendlyName,
			Namespace:    params.Namespace,
			Token:        params.Token,
			Description:  params.Description,
		})

		if err != nil {
			logger.Error("Error creating organisation", "error", err)

			return err
		}

		return nil
	})

	return ToOrganisationDTO(orgModel), err
}

func (o organisationService) Get(ctx context.Context, id int64) (OrganisationDTO, error) {
	logger := logging.FromContext(ctx)
	logger.Debug("Fetching organisation", "id", id)
	model, err := o.store.GetOrganisationByID(ctx, id)
	if err != nil {
		logger.Error("Error fetching organisation", "error", err)
		return OrganisationDTO{}, err
	}

	return ToOrganisationDTO(model), nil
}

func (o organisationService) GetDefault(ctx context.Context) (OrganisationDTO, error) {
	logger := logging.FromContext(ctx)
	logger.Debug("Fetching default organisation")
	model, err := o.store.GetDefaultOrganisation(ctx)
	if err != nil {
		logger.Error("Error fetching default organisation", "error", err)
		return OrganisationDTO{}, err
	}

	return ToOrganisationDTO(model), nil
}

func (o organisationService) SetDefault(ctx context.Context, id int64) (OrganisationDTO, error) {
	logger := logging.FromContext(ctx)
	logger.Debug("setting default org", "org", id)

	if err := o.store.SetOrgsDefaultFalse(ctx); err != nil {
		logger.Error("Error setting default org", "error", err)

		return OrganisationDTO{}, err
	}

	model, err := o.store.SetDefaultOrg(ctx, id)
	if err != nil {
		logger.Error("Error setting default org", "error", err)

		return OrganisationDTO{}, err
	}

	return ToOrganisationDTO(model), nil
}

func (o organisationService) GetAll(ctx context.Context) ([]OrganisationDTO, error) {
	logger := logging.FromContext(ctx)
	logger.Debug("Listing all organisations")

	models, err := o.store.ListOrganisations(ctx)
	if err != nil {
		logger.Error("Error listing organisations", "error", err)

		return nil, err
	}

	return ToOrganisationDTOs(models), nil
}

func (o organisationService) GetStaleOrgs(ctx context.Context) ([]OrganisationDTO, error) {
	logger := logging.FromContext(ctx)
	logger.Debug("Fetching stale organisations")

	fiveMinutesAgo := time.Now().Add(-5 * time.Minute).Unix()

	models, err := o.store.ListOrgsOlderThanUpdatedAt(ctx, fiveMinutesAgo)
	if err != nil {
		logger.Error("Error fetching stale organisations", "error", err)
		return nil, err
	}

	return ToOrganisationDTOs(models), nil
}

func (o organisationService) GetOrgAssociatedToProduct(ctx context.Context, productID int64) (InternalOrganisation, error) {
	logger := logging.FromContext(ctx)
	logger.Debug("Fetching organisations associated to product", "product", productID)

	model, err := o.store.GetOrganisationForProduct(ctx, sql.NullInt64{Int64: productID, Valid: true})
	if err != nil {
		logger.Error("Error fetching organisations associated to product", "error", err)
		return InternalOrganisation{}, err
	}

	return ToInternalOrganisation(model), err
}

func (o organisationService) Delete(ctx context.Context, id int64) error {
	logger := logging.FromContext(ctx)
	logger.Debug("Fetching organisation", "id", id)

	if err := o.store.DeleteProductOrganisationByOrgID(ctx, sql.NullInt64{
		Int64: id,
		Valid: true,
	}); err != nil {
		logger.Error("Error deleting organisation", "error", err)
	}

	if err := o.store.DeleteOrg(ctx, id); err != nil {
		logger.Error("Error deleting organisation", "error", err)

		return err
	}

	return nil
}

type UpdateOrgParams struct {
	ID           int64
	DefaultOrg   bool
	FriendlyName string
	Namespace    string
	Description  string
}

func (o organisationService) Update(ctx context.Context, params UpdateOrgParams) (OrganisationDTO, error) {
	logger := logging.FromContext(ctx)
	logger.Debug("Updating organisation", "id", params.ID)

	var model database.Organisation
	var err error

	err = database.WithTxnContext(ctx, o.txnDB, func(tx *sql.Tx) error {
		txnStore := o.txnFunc(tx)
		if params.DefaultOrg {
			if err = txnStore.SetOrgsDefaultFalse(ctx); err != nil {
				logger.Error("Error setting default org", "error", err)

				return err
			}
		}

		model, err = txnStore.UpdateOrganisation(ctx, database.UpdateOrganisationParams{
			ID:           params.ID,
			DefaultOrg:   params.DefaultOrg,
			FriendlyName: params.FriendlyName,
			Namespace:    params.Namespace,
			Description:  params.Description,
		})
		if err != nil {
			logger.Error("Error updating organisation", "error", err)

			return err
		}

		return nil
	})

	return ToOrganisationDTO(model), err
}

func (o organisationService) UpdateSyncDateNow(ctx context.Context, id int64) error {
	logger := logging.FromContext(ctx)
	logger.Debug("Updating sync date")
	return o.store.UpdateProductSync(ctx, id)
}

func (o organisationService) AssociateProductToOrg(ctx context.Context, orgID int64, productID int64) error {
	logger := logging.FromContext(ctx)
	logger.Debug("Associating product to org", "org", orgID, "product", productID)

	if err := o.store.AddProductToOrganisation(ctx, database.AddProductToOrganisationParams{
		ProductID: sql.NullInt64{
			Int64: productID,
			Valid: true,
		},
		OrganisationID: sql.NullInt64{
			Int64: orgID,
			Valid: true,
		},
	}); err != nil {
		logger.Error("Error linking product to organisation", "error", err)
		return err
	}

	return nil
}

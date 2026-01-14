// Package organisations provides a service for managing and retrieving GitHub organizations and their configurations.
package organisations

import (
	"context"
	"database/sql"
	"time"
	"watchtower/internal/database"
	"watchtower/internal/logging"
)

// Service handles organization-related business logic.
type Service struct {
	store   OrgStore
	txnDB   *sql.DB
	txnFunc func(tx *sql.Tx) OrgStore
}

// New creates a new Service instance with the provided store, transaction database and transaction function.
func New(store OrgStore, txnDB *sql.DB, txnFunc func(tx *sql.Tx) OrgStore) *Service {
	return &Service{
		store:   store,
		txnDB:   txnDB,
		txnFunc: txnFunc,
	}
}

// Create creates a new organization with the given parameters and returns the created organization or an error.
func (s Service) Create(ctx context.Context, params CreateOrgParams) (OrganisationDTO, error) {
	logger := logging.FromContext(ctx).With("service", "organisations")

	logger.Debug("Creating organisation")

	var orgModel database.Organisation
	var err error

	err = database.WithTxnContext(ctx, s.txnDB, func(tx *sql.Tx) error {
		txnStore := s.txnFunc(tx)
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

	return toOrganisationDTO(orgModel), err
}

// Get retrieves an organisation by its ID and returns its DTO representation or an error if the operation fails.
func (s Service) Get(ctx context.Context, id int64) (OrganisationDTO, error) {
	logger := logging.FromContext(ctx).With("service", "organisations")

	logger.Debug("Fetching organisation", "id", id)
	model, err := s.store.GetOrganisationByID(ctx, id)
	if err != nil {
		logger.Error("Error fetching organisation", "error", err)
		return OrganisationDTO{}, err
	}

	return toOrganisationDTO(model), nil
}

// GetDefault retrieves the default organisation and returns its DTO representation or an error if the operation fails.
func (s Service) GetDefault(ctx context.Context) (OrganisationDTO, error) {
	logger := logging.FromContext(ctx).With("service", "organisations")

	logger.Debug("Fetching default organisation")
	model, err := s.store.GetDefaultOrganisation(ctx)
	if err != nil {
		logger.Error("Error fetching default organisation", "error", err)
		return OrganisationDTO{}, err
	}

	return toOrganisationDTO(model), nil
}

// SetDefault sets the organisation with the given ID as the default organisation and updates the database accordingly.
func (s Service) SetDefault(ctx context.Context, id int64) (OrganisationDTO, error) {
	logger := logging.FromContext(ctx).With("service", "organisations")

	logger.Debug("Setting default organisation", "id", id)

	var model database.Organisation
	var err error

	err = database.WithTxnContext(ctx, s.txnDB, func(tx *sql.Tx) error {
		txnStore := s.txnFunc(tx)
		if err = txnStore.SetOrgsDefaultFalse(ctx); err != nil {
			return err
		}

		if model, err = txnStore.SetDefaultOrg(ctx, id); err != nil {
			return err
		}

		return nil
	})

	return toOrganisationDTO(model), err
}

// GetAll retrieves all organisations and returns their DTO representations or an error if the operation fails.
func (s Service) GetAll(ctx context.Context) ([]OrganisationDTO, error) {
	logger := logging.FromContext(ctx).With("service", "organisations")

	logger.Debug("Fetching all organisations")
	models, err := s.store.ListOrganisations(ctx)
	if err != nil {
		logger.Error("Error fetching all organisations", "error", err)
		return nil, err
	}

	return toOrganisationDTOs(models), nil
}

// GetStaleOrgs retrieves organisations that haven't been updated for 1 hour.
func (s Service) GetStaleOrgs(ctx context.Context) ([]OrganisationDTO, error) {
	logger := logging.FromContext(ctx).With("service", "organisations")

	logger.Debug("Fetching stale organisations")

	oneHourAgo := time.Now().Add(-1 * time.Hour).Unix()

	models, err := s.store.ListOrgsOlderThanUpdatedAt(ctx, oneHourAgo)
	if err != nil {
		logger.Error("Error fetching stale organisations", "error", err)
		return nil, err
	}

	return toOrganisationDTOs(models), nil
}

// GetOrgAssociatedToProduct retrieves the organisation associated with a product.
func (s Service) GetOrgAssociatedToProduct(ctx context.Context, productID int64) (InternalOrganisation, error) {
	logger := logging.FromContext(ctx).With("service", "organisations")

	logger.Debug("Fetching organisation for product", "productID", productID)
	model, err := s.store.GetOrganisationForProduct(ctx, sql.NullInt64{Int64: productID, Valid: true})
	if err != nil {
		logger.Error("Error fetching organisation for product", "error", err)
		return InternalOrganisation{}, err
	}

	return toInternalOrganisation(model), nil
}

// Delete removes an organisation and its associated product-organisation mapping.
func (s Service) Delete(ctx context.Context, id int64) error {
	logger := logging.FromContext(ctx).With("service", "organisations")

	logger.Debug("Deleting organisation", "id", id)

	err := database.WithTxnContext(ctx, s.txnDB, func(tx *sql.Tx) error {
		txnStore := s.txnFunc(tx)
		if err := txnStore.DeleteProductOrganisationByOrgID(ctx, sql.NullInt64{Int64: id, Valid: true}); err != nil {
			return err
		}

		if err := txnStore.DeleteOrg(ctx, id); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Error("Error deleting organisation", "error", err)
		return err
	}

	return nil
}

// Update updates an existing organisation's details.
func (s Service) Update(ctx context.Context, params UpdateOrgParams) (OrganisationDTO, error) {
	logger := logging.FromContext(ctx).With("service", "organisations")

	logger.Debug("Updating organisation", "id", params.ID)

	var model database.Organisation
	var err error

	err = database.WithTxnContext(ctx, s.txnDB, func(tx *sql.Tx) error {
		txnStore := s.txnFunc(tx)
		if params.DefaultOrg {
			if err = txnStore.SetOrgsDefaultFalse(ctx); err != nil {
				return err
			}
		}

		model, err = txnStore.UpdateOrganisation(ctx, database.UpdateOrganisationParams{
			ID:           params.ID,
			FriendlyName: params.FriendlyName,
			Namespace:    params.Namespace,
			Description:  params.Description,
			DefaultOrg:   params.DefaultOrg,
		})

		return err
	})

	if err != nil {
		logger.Error("Error updating organisation", "error", err)
		return OrganisationDTO{}, err
	}

	return toOrganisationDTO(model), nil
}

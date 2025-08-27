package watchtower

import (
	"context"
	"watchtower/internal/database"

	"github.com/code-gorilla-au/go-toolbox/logging"
)

type CreateOrgParams struct {
	FriendlyName string
	Namespace    string
	Token        string
	Description  string
}

func (o OrganisationService) Create(ctx context.Context, params CreateOrgParams) (OrganisationDTO, error) {
	logger := logging.FromContext(ctx)
	logger.Info("Creating organisation")

	if err := o.db.SetOrgsDefaultFalse(ctx); err != nil {
		logger.Error("Error setting default org", "error", err)

		return OrganisationDTO{}, err
	}

	model, err := o.db.CreateOrganisation(ctx, database.CreateOrganisationParams{
		FriendlyName: params.FriendlyName,
		Namespace:    params.Namespace,
		Token:        params.Token,
		Description:  params.Description,
	})

	if err != nil {
		logger.Error("Error creating organisation", "error", err)

		return OrganisationDTO{}, err
	}

	return ToOrganisationDTO(model), nil
}

func (o OrganisationService) Get(ctx context.Context, id int64) (OrganisationDTO, error) {
	logger := logging.FromContext(ctx)
	logger.Info("Fetching organisation", "id", id)
	model, err := o.db.GetOrganisationByID(ctx, id)
	if err != nil {
		logger.Error("Error fetching organisation", "error", err)
		return OrganisationDTO{}, err
	}

	return ToOrganisationDTO(model), nil
}

func (o OrganisationService) GetDefault(ctx context.Context) (OrganisationDTO, error) {
	logger := logging.FromContext(ctx)
	logger.Info("Fetching default organisation")
	model, err := o.db.GetDefaultOrganisation(ctx)
	if err != nil {
		logger.Error("Error fetching default organisation", "error", err)
		return OrganisationDTO{}, err
	}

	return ToOrganisationDTO(model), nil
}

func (o OrganisationService) SetDefault(ctx context.Context, id int64) (OrganisationDTO, error) {
	logger := logging.FromContext(ctx)
	logger.Info("setting default org", "org", id)

	if err := o.db.SetOrgsDefaultFalse(ctx); err != nil {
		logger.Error("Error setting default org", "error", err)

		return OrganisationDTO{}, err
	}

	model, err := o.db.SetDefaultOrg(ctx, id)
	if err != nil {
		logger.Error("Error setting default org", "error", err)

		return OrganisationDTO{}, err
	}

	return ToOrganisationDTO(model), nil
}

func (o OrganisationService) GetAll(ctx context.Context) ([]OrganisationDTO, error) {
	logger := logging.FromContext(ctx)
	logger.Info("Listing all organisations")

	models, err := o.db.ListOrganisations(ctx)
	if err != nil {
		logger.Error("Error listing organisations", "error", err)

		return nil, err
	}

	return ToOrganisationDTOs(models), nil
}

func (o OrganisationService) Delete(ctx context.Context, id int64) error {
	logger := logging.FromContext(ctx)
	logger.Info("Fetching organisation", "id", id)

	if err := o.db.DeleteOrg(ctx, id); err != nil {
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

func (o OrganisationService) Update(ctx context.Context, params UpdateOrgParams) (OrganisationDTO, error) {
	logger := logging.FromContext(ctx)
	logger.Info("Updating organisation", "id", params.ID)

	if params.DefaultOrg {
		if err := o.db.SetOrgsDefaultFalse(ctx); err != nil {
			logger.Error("Error setting default org", "error", err)

			return OrganisationDTO{}, err
		}
	}

	model, err := o.db.UpdateOrganisation(ctx, database.UpdateOrganisationParams{
		ID:           params.ID,
		DefaultOrg:   params.DefaultOrg,
		FriendlyName: params.FriendlyName,
		Namespace:    params.Namespace,
		Description:  params.Description,
	})
	if err != nil {
		logger.Error("Error updating organisation", "error", err)

		return OrganisationDTO{}, err
	}

	return ToOrganisationDTO(model), nil
}

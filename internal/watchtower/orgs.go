package watchtower

import (
	"watchtower/internal/database"

	"github.com/code-gorilla-au/go-toolbox/logging"
)

// CreateOrganisation creates a new organisation in the database using the specified friendly name and namespace.
// It logs the creation process and returns the created organisation DTO or an error if the operation fails.
func (s *Service) CreateOrganisation(friendlyName string, namespace string, token string, description string) (OrganisationDTO, error) {
	logger := logging.FromContext(s.ctx)
	logger.Info("Creating organisation")

	if err := s.db.SetOrgsDefaultFalse(s.ctx); err != nil {
		logger.Error("Error setting default org", "error", err)
		return OrganisationDTO{}, err
	}

	model, err := s.db.CreateOrganisation(s.ctx, database.CreateOrganisationParams{
		FriendlyName: friendlyName,
		Namespace:    namespace,
		Token:        token,
		Description:  description,
	})

	if err != nil {
		logger.Error("Error creating organisation", "error", err)
		return OrganisationDTO{}, err
	}

	return ToOrganisationDTO(model), nil
}

// GetDefaultOrganisation returns the default organisation (marked as default_org = 1)
func (s *Service) GetDefaultOrganisation() (OrganisationDTO, error) {
	logger := logging.FromContext(s.ctx)
	logger.Info("Fetching default organisation")

	model, err := s.db.GetDefaultOrganisation(s.ctx)
	if err != nil {
		logger.Error("Error fetching default organisation", "error", err)
		return OrganisationDTO{}, err
	}
	return ToOrganisationDTO(model), nil
}

func (s *Service) SetDefaultOrg(id int64) (OrganisationDTO, error) {
	logger := logging.FromContext(s.ctx)
	logger.Info("setting default org", "org", id)

	if err := s.db.SetOrgsDefaultFalse(s.ctx); err != nil {
		logger.Error("Error setting default org", "error", err)
		return OrganisationDTO{}, err
	}

	model, err := s.db.SetDefaultOrg(s.ctx, id)
	if err != nil {
		logger.Error("Error setting default org", "error", err)
		return OrganisationDTO{}, err
	}

	return ToOrganisationDTO(model), nil
}

func (s *Service) GetOrganisationByID(id int64) (OrganisationDTO, error) {
	logger := logging.FromContext(s.ctx)
	logger.Info("Fetching organisation", "id", id)
	model, err := s.db.GetOrganisationByID(s.ctx, id)

	if err != nil {
		logger.Error("Error fetching organisation", "error", err)
		return OrganisationDTO{}, err
	}
	return ToOrganisationDTO(model), nil
}

// GetAllOrganisations returns all organisations ordered by friendly_name
func (s *Service) GetAllOrganisations() ([]OrganisationDTO, error) {
	logger := logging.FromContext(s.ctx)
	logger.Info("Listing all organisations")

	models, err := s.db.ListOrganisations(s.ctx)
	if err != nil {
		logger.Error("Error listing organisations", "error", err)
		return nil, err
	}

	result := make([]OrganisationDTO, 0, len(models))
	for _, m := range models {
		result = append(result, ToOrganisationDTO(m))
	}

	return result, nil
}

// DeleteOrganisation deletes an organisation from the database by its ID. Returns an error if the operation fails.
func (s *Service) DeleteOrganisation(id int64) error {
	logger := logging.FromContext(s.ctx)
	logger.Info("Fetching organisation", "id", id)

	if err := s.db.DeleteOrg(s.ctx, id); err != nil {
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
}

// UpdateOrganisation updates the details of an organisation based on the provided parameters.
// It sets other organisations' default status to false if a new default organisation is specified.
// Returns the updated organisation details as OrganisationDTO or an error if the update fails.
func (s *Service) UpdateOrganisation(params UpdateOrgParams) (OrganisationDTO, error) {
	logger := logging.FromContext(s.ctx)
	logger.Info("Updating organisation", "id", params.ID)

	if params.DefaultOrg {
		if err := s.db.SetOrgsDefaultFalse(s.ctx); err != nil {
			logger.Error("Error setting default org", "error", err)
			return OrganisationDTO{}, err
		}
	}

	model, err := s.db.UpdateOrganisation(s.ctx, database.UpdateOrganisationParams{
		ID:           params.ID,
		DefaultOrg:   params.DefaultOrg,
		FriendlyName: params.FriendlyName,
		Namespace:    params.Namespace,
	})
	if err != nil {
		logger.Error("Error updating organisation", "error", err)
		return OrganisationDTO{}, err
	}

	return ToOrganisationDTO(model), nil
}

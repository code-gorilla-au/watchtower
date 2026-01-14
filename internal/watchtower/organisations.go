package watchtower

import "watchtower/internal/organisations"

// CreateOrganisation creates a new organisation in the database using the specified friendly name and namespace.
// It logs the creation process and returns the created organisation DTO or an error if the operation fails.
func (s *Service) CreateOrganisation(friendlyName string, namespace string, token string, description string) (organisations.OrganisationDTO, error) {
	return s.orgSvc.Create(s.ctx, organisations.CreateOrgParams{
		FriendlyName: friendlyName,
		Namespace:    namespace,
		Token:        token,
		Description:  description,
	})
}

// GetDefaultOrganisation returns the default organisation (marked as default_org = 1).
func (s *Service) GetDefaultOrganisation() (organisations.OrganisationDTO, error) {
	return s.orgSvc.GetDefault(s.ctx)
}

// SetDefaultOrg marks the organization with the specified ID as the default organization.
func (s *Service) SetDefaultOrg(id int64) (organisations.OrganisationDTO, error) {
	return s.orgSvc.SetDefault(s.ctx, id)
}

// GetOrganisationByID retrieves an organization's details by its ID.
func (s *Service) GetOrganisationByID(id int64) (organisations.OrganisationDTO, error) {
	return s.orgSvc.Get(s.ctx, id)
}

// GetAllOrganisations returns all organisations ordered by friendly_name.
func (s *Service) GetAllOrganisations() ([]organisations.OrganisationDTO, error) {
	return s.orgSvc.GetAll(s.ctx)
}

// DeleteAllOrgs deletes all organizations and their associated products.
func (s *Service) DeleteAllOrgs() error {
	list, err := s.orgSvc.GetAll(s.ctx)
	if err != nil {
		return err
	}

	for _, org := range list {
		_ = s.DeleteOrganisation(org.ID)
	}

	return nil
}

// DeleteOrganisation deletes an organisation from the database by its ID. Returns an error if the operation fails.
func (s *Service) DeleteOrganisation(id int64) error {
	list, err := s.productSvc.GetByOrg(s.ctx, id)
	if err != nil {
		return err
	}

	for _, p := range list {
		_ = s.productSvc.DeleteProduct(s.ctx, p.ID)
	}

	return s.orgSvc.Delete(s.ctx, id)
}

// UpdateOrganisation updates the details of an organisation based on the provided parameters.
// It sets other organisations' default status to false if a new default organisation is specified.
// Returns the updated organisation details as OrganisationDTO or an error if the update fails.
func (s *Service) UpdateOrganisation(params organisations.UpdateOrgParams) (organisations.OrganisationDTO, error) {
	return s.orgSvc.Update(s.ctx, params)
}

package watchtower

import (
	"context"
	"watchtower/internal/database"

	"github.com/code-gorilla-au/go-toolbox/github"
	"github.com/code-gorilla-au/go-toolbox/logging"
)

// NewService creates and returns a new Service instance with the provided database queries.
func NewService(ctx context.Context, db *database.Queries) *Service {
	return &Service{
		ghClient: github.New(logging.FromContext(ctx)),
		ctx:      ctx,
		db:       db,
		orgSvc: &organisationService{
			db: db,
		},
		productSvc: &productsService{
			db: db,
		},
	}
}

func (s *Service) Startup(ctx context.Context) {
	s.ctx = ctx
}

// CreateOrganisation creates a new organisation in the database using the specified friendly name and namespace.
// It logs the creation process and returns the created organisation DTO or an error if the operation fails.
func (s *Service) CreateOrganisation(friendlyName string, namespace string, token string, description string) (OrganisationDTO, error) {
	return s.orgSvc.Create(s.ctx, CreateOrgParams{
		FriendlyName: friendlyName,
		Namespace:    namespace,
		Token:        token,
		Description:  description,
	})
}

// GetDefaultOrganisation returns the default organisation (marked as default_org = 1).
func (s *Service) GetDefaultOrganisation() (OrganisationDTO, error) {
	return s.orgSvc.GetDefault(s.ctx)
}

func (s *Service) SetDefaultOrg(id int64) (OrganisationDTO, error) {
	return s.orgSvc.SetDefault(s.ctx, id)
}

func (s *Service) GetOrganisationByID(id int64) (OrganisationDTO, error) {
	return s.orgSvc.Get(s.ctx, id)
}

// GetAllOrganisations returns all organisations ordered by friendly_name.
func (s *Service) GetAllOrganisations() ([]OrganisationDTO, error) {
	return s.orgSvc.GetAll(s.ctx)
}

// DeleteOrganisation deletes an organisation from the database by its ID. Returns an error if the operation fails.
func (s *Service) DeleteOrganisation(id int64) error {
	return s.orgSvc.Delete(s.ctx, id)
}

// UpdateOrganisation updates the details of an organisation based on the provided parameters.
// It sets other organisations' default status to false if a new default organisation is specified.
// Returns the updated organisation details as OrganisationDTO or an error if the update fails.
func (s *Service) UpdateOrganisation(params UpdateOrgParams) (OrganisationDTO, error) {
	return s.orgSvc.Update(s.ctx, params)
}

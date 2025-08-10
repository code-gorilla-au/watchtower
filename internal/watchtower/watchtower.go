package watchtower

import (
	"context"
	"watchtower/internal/database"

	"github.com/code-gorilla-au/go-toolbox/logging"
)

// NewService creates and returns a new Service instance with the provided database queries.
func NewService(db database.Queries) *Service {
	return &Service{
		db: db,
	}
}

// CreateOrganisation creates a new organisation in the database using the specified friendly name and namespace.
// It logs the creation process and returns the created organisation DTO or an error if the operation fails.
func (s *Service) CreateOrganisation(ctx context.Context, friendlyName string, namespace string) (OrganisationDTO, error) {
	logger := logging.FromContext(ctx)
	logger.Info("Creating organisation")

	model, err := s.db.CreateOrganisation(ctx, database.CreateOrganisationParams{
		FriendlyName: friendlyName,
		Namespace:    namespace,
	})

	if err != nil {
		logger.Error("Error creating organisation", err)
		return OrganisationDTO{}, err
	}

	return ToOrganisationDTO(model), nil
}

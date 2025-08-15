package watchtower

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
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
	}
}

func (s *Service) Startup(ctx context.Context) {
	s.ctx = ctx
}

// CreateOrganisation creates a new organisation in the database using the specified friendly name and namespace.
// It logs the creation process and returns the created organisation DTO or an error if the operation fails.
func (s *Service) CreateOrganisation(friendlyName string, namespace string, token string) (OrganisationDTO, error) {
	logger := logging.FromContext(s.ctx)
	logger.Info("Creating organisation")

	if err := s.db.SetOrgsDefaultFalse(s.ctx); err != nil {
		logger.Error("Error setting default org", err)
		return OrganisationDTO{}, err
	}

	model, err := s.db.CreateOrganisation(s.ctx, database.CreateOrganisationParams{
		FriendlyName: friendlyName,
		Namespace:    namespace,
		Token:        token,
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
		logger.Error("Error fetching default organisation", err)
		return OrganisationDTO{}, err
	}
	return ToOrganisationDTO(model), nil
}

func (s *Service) SetDefaultOrg(id int64) (OrganisationDTO, error) {
	logger := logging.FromContext(s.ctx)
	logger.Info("setting default org", "org", id)

	if err := s.db.SetOrgsDefaultFalse(s.ctx); err != nil {
		logger.Error("Error setting default org", err)
		return OrganisationDTO{}, err
	}

	model, err := s.db.SetDefaultOrg(s.ctx, id)
	if err != nil {
		logger.Error("Error setting default org", err)
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
		logger.Error("Error listing organisations", err)
		return nil, err
	}

	result := make([]OrganisationDTO, 0, len(models))
	for _, m := range models {
		result = append(result, ToOrganisationDTO(m))
	}

	return result, nil
}

// CreateProduct creates a new product and associates it with an organisation
func (s *Service) CreateProduct(name string, tags []string, organisationID int64) (ProductDTO, error) {
	logger := logging.FromContext(s.ctx)
	logger.Info("Creating product")

	var tagsNS sql.NullString

	tagJson, err := json.Marshal(tags)
	if err != nil {
		logger.Error("Error marshalling tags", err)
		return ProductDTO{}, err
	}

	tagsNS = sql.NullString{String: string(tagJson), Valid: true}

	prod, err := s.db.CreateProduct(s.ctx, database.CreateProductParams{
		Name: name,
		Tags: tagsNS,
	})
	if err != nil {
		logger.Error("Error creating product", err)
		return ProductDTO{}, err
	}

	err = s.db.AddProductToOrganisation(s.ctx, database.AddProductToOrganisationParams{
		ProductID:      sql.NullInt64{Int64: prod.ID, Valid: true},
		OrganisationID: sql.NullInt64{Int64: organisationID, Valid: true},
	})
	if err != nil {
		logger.Error("Error linking product to organisation", err)
		return ProductDTO{}, err
	}

	return ToProductDTO(prod), nil
}

// GetProductByID fetches a product by its ID
func (s *Service) GetProductByID(id int64) (ProductDTO, error) {
	logger := logging.FromContext(s.ctx)
	logger.Info("Fetching product by ID")

	prod, err := s.db.GetProductByID(s.ctx, id)
	if err != nil {
		logger.Error("Error fetching product by ID", err)
		return ProductDTO{}, err
	}

	return ToProductDTO(prod), nil
}

// GetAllProductsForOrganisation lists products linked to the given organisation
func (s *Service) GetAllProductsForOrganisation(organisationID int64) ([]ProductDTO, error) {
	logger := logging.FromContext(s.ctx)
	logger.Info("Listing products for organisation")

	models, err := s.db.ListProductsByOrganisation(s.ctx, sql.NullInt64{Int64: organisationID, Valid: true})
	if err != nil {
		logger.Error("Error listing products for organisation", err)
		return nil, err
	}

	result := make([]ProductDTO, 0, len(models))
	for _, m := range models {
		result = append(result, ToProductDTO(m))
	}

	return result, nil
}

// UpdateProduct updates a product and returns the updated entity
func (s *Service) UpdateProduct(id int64, name string, tags *string) (ProductDTO, error) {
	logger := logging.FromContext(s.ctx)
	logger.Info("Updating product")

	var tagsNS sql.NullString
	if tags != nil {
		tagsNS = sql.NullString{String: *tags, Valid: true}
	}

	err := s.db.UpdateProduct(s.ctx, database.UpdateProductParams{
		Name: name,
		Tags: tagsNS,
		ID:   id,
	})
	if err != nil {
		logger.Error("Error updating product", err)
		return ProductDTO{}, err
	}

	prod, err := s.db.GetProductByID(s.ctx, id)
	if err != nil {
		logger.Error("Error fetching updated product", err)
		return ProductDTO{}, err
	}

	return ToProductDTO(prod), nil
}

func (s *Service) DeleteProduct(id int64) error {
	logger := logging.FromContext(s.ctx)
	logger.Info("Deleting product")
	return s.db.DeleteProduct(s.ctx, id)
}

func (s *Service) SyncProduct(id int64) error {
	logger := logging.FromContext(s.ctx)

	product, err := s.GetProductByID(id)
	if err != nil {
		logger.Error("Error fetching product", err)
		return err
	}

	org, err := s.db.GetOrganisationForProduct(s.ctx, sql.NullInt64{Int64: product.ID, Valid: true})
	if err != nil {
		logger.Error("Error fetching organisation for product", err)
		return err
	}

	for _, tag := range product.Tags {
		logger.Info("Searching for repo with tag", "tag", tag)

		repos, apiErr := s.ghClient.SearchRepos(org.Namespace, strings.TrimSpace(tag), org.Token)
		if apiErr != nil {
			logger.Error("Error searching for repos", "error", apiErr)
			return apiErr
		}

		for _, repo := range repos.Data.Search.Edges {
			_, err = s.db.CreateRepo(s.ctx, database.CreateRepoParams{
				Name:  repo.Node.Name,
				Url:   repo.Node.Url,
				Topic: tag,
				Owner: repo.Node.Owner.Login,
			})
			if err != nil {
				logger.Error("Error creating repo", "error", err)
				return err
			}
		}

	}

	return nil
}

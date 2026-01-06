package watchtower

import "watchtower/internal/products"

// CreateProduct creates a new product and associates it with an organisation.
func (s *Service) CreateProduct(name string, description string, tags []string, organisationID int64) (products.ProductDTO, error) {

	prod, err := s.productSvc.Create(s.ctx, products.CreateProductParams{
		Name: name,
		Tags: tags,
		Desc: description,
	})
	if err != nil {
		return products.ProductDTO{}, err
	}

	err = s.orgSvc.AssociateProductToOrg(s.ctx, organisationID, prod.ID)
	if err != nil {
		return products.ProductDTO{}, err
	}

	return prod, nil
}

// GetProductByID fetches a product by its ID.
func (s *Service) GetProductByID(id int64) (products.ProductDTO, error) {
	return s.productSvc.Get(s.ctx, id)
}

// GetAllProductsForOrganisation lists products linked to the given organisation.
func (s *Service) GetAllProductsForOrganisation(organisationID int64) ([]products.ProductDTO, error) {
	return s.productSvc.GetByOrg(s.ctx, organisationID)
}

// UpdateProduct updates a product and returns the updated entity.
func (s *Service) UpdateProduct(id int64, name string, description string, tags []string) (products.ProductDTO, error) {
	return s.productSvc.Update(s.ctx, products.UpdateProductParams{
		ID:          id,
		Name:        name,
		Tags:        tags,
		Description: description,
	})
}

// DeleteProduct removes a product from the database by its ID. Returns an error if the operation fails.
func (s *Service) DeleteProduct(id int64) error {
	return s.productSvc.DeleteProduct(s.ctx, id)
}

// GetProductRepos retrieves a list of repositories associated with the product identified by the given ID.
func (s *Service) GetProductRepos(id int64) ([]products.RepositoryDTO, error) {
	return s.productSvc.GetRepos(s.ctx, id)
}

// GetProductPullRequests retrieves a list of pull requests for a specified product ID.
func (s *Service) GetProductPullRequests(id int64) ([]products.PullRequestDTO, error) {
	return s.productSvc.GetPullRequests(s.ctx, id)
}

// GetPullRequestByOrganisation retrieves a list of pull requests for a given organization by its ID.
func (s *Service) GetPullRequestByOrganisation(id int64) ([]products.PullRequestDTO, error) {
	return s.productSvc.GetPullRequestByOrg(s.ctx, id)
}

// GetSecurityByProductID retrieves a list of security details for a given product ID. Returns an error if the operation fails.
func (s *Service) GetSecurityByProductID(productID int64) ([]products.SecurityDTO, error) {
	return s.productSvc.GetSecurity(s.ctx, productID)
}

// GetSecurityByOrganisation retrieves a list of security details associated with a specific organization by its ID.
func (s *Service) GetSecurityByOrganisation(id int64) ([]products.SecurityDTO, error) {
	return s.productSvc.GetSecurityByOrg(s.ctx, id)
}

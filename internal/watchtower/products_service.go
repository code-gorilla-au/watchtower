package watchtower

import (
	"context"
	"database/sql"
	"encoding/json"
	"watchtower/internal/database"

	"github.com/code-gorilla-au/go-toolbox/github"
	"github.com/code-gorilla-au/go-toolbox/logging"
)

type CreateProductParams struct {
	Name string
	Desc string
	Tags []string
}

func (p *productsService) Create(ctx context.Context, params CreateProductParams) (ProductDTO, error) {
	logger := logging.FromContext(ctx)
	logger.Info("Creating product")

	var tagsNS sql.NullString

	tagJson, err := json.Marshal(params.Tags)
	if err != nil {
		logger.Error("Error marshalling tags", "error", err)

		return ProductDTO{}, err
	}

	tagsNS = sql.NullString{String: string(tagJson), Valid: true}

	prod, err := p.db.CreateProduct(ctx, database.CreateProductParams{
		Name:        params.Name,
		Tags:        tagsNS,
		Description: params.Desc,
	})
	if err != nil {
		logger.Error("Error creating product", "error", err)

		return ProductDTO{}, err
	}

	return ToProductDTO(prod), nil
}

func (p *productsService) Get(ctx context.Context, id int64) (ProductDTO, error) {
	logger := logging.FromContext(ctx)
	logger.Info("Fetching product by ID")

	prod, err := p.db.GetProductByID(ctx, id)
	if err != nil {
		logger.Error("Error fetching product by ID", "error", err)

		return ProductDTO{}, err
	}

	return ToProductDTO(prod), nil
}

func (p *productsService) GetByOrg(ctx context.Context, orgID int64) ([]ProductDTO, error) {
	logger := logging.FromContext(ctx)
	logger.Info("Listing products for organisation")

	models, err := p.db.ListProductsByOrganisation(ctx, sql.NullInt64{Int64: orgID, Valid: true})
	if err != nil {
		logger.Error("Error listing products for organisation", "error", err)

		return nil, err
	}

	return ToProductDTOs(models), nil
}

type UpdateProductParams struct {
	ID   int64
	Name string
	Tags []string
}

func (p *productsService) Update(ctx context.Context, params UpdateProductParams) (ProductDTO, error) {
	logger := logging.FromContext(ctx)
	logger.Debug("Updating product")

	data, err := json.Marshal(params.Tags)
	if err != nil {
		logger.Error("Error marshalling tags", "error", err)
	}

	var tagsNS sql.NullString
	if params.Tags != nil {
		tagsNS = sql.NullString{String: string(data), Valid: true}
	}

	model, err := p.db.UpdateProduct(ctx, database.UpdateProductParams{
		Name: params.Name,
		Tags: tagsNS,
		ID:   params.ID,
	})
	if err != nil {
		logger.Error("Error updating product", "error", err)

		return ProductDTO{}, err
	}

	return ToProductDTO(model), nil
}

func (p *productsService) GetRepos(ctx context.Context, id int64) ([]RepositoryDTO, error) {
	return p.repoService.GetRepos(ctx, id)
}

func (p *productsService) GetPullRequests(ctx context.Context, id int64) ([]PullRequestDTO, error) {
	return p.repoService.GetPullRequests(ctx, id)
}

func (p *productsService) GetPullRequestByOrg(ctx context.Context, orgID int64) ([]PullRequestDTO, error) {
	return p.repoService.GetPullRequestByOrg(ctx, orgID)
}

func (p *productsService) GetSecurity(ctx context.Context, id int64) ([]SecurityDTO, error) {
	return p.repoService.GetSecurity(ctx, id)
}

func (p *productsService) GetSecurityByOrg(ctx context.Context, orgID int64) ([]SecurityDTO, error) {
	return p.repoService.GetSecurityByOrg(ctx, orgID)
}

func (p *productsService) BulkInsertRepos(ctx context.Context, reposList []github.Node[github.Repository], tag string) error {

	params := ToCreateRepoFromGithub(reposList, tag)

	return p.repoService.BulkCreateRepos(ctx, params)

}

func (p *productsService) BulkInsertRepoDetails(ctx context.Context, repoDetails github.QueryRepository) error {
	logger := logging.FromContext(ctx)

	prs := ToCreatePRsFromGithubRepos(repoDetails.Data.Repository.PullRequests, repoDetails.Data.Repository.Name)
	secs := ToSecParamsFromGithubVulnerabilities(repoDetails.Data.Repository.VulnerabilityAlerts, repoDetails.Data.Repository.Name)

	if err := p.repoService.BulkCreatePullRequest(ctx, prs); err != nil {
		logger.Error("Error creating pull requests", "error", err)
		return err
	}

	if err := p.repoService.BulkCreateSecurity(ctx, secs); err != nil {
		logger.Error("Error creating security", "error", err)
		return err
	}

	return nil
}

func (p *productsService) DeleteProduct(ctx context.Context, id int64) error {
	logger := logging.FromContext(ctx)
	logger.Debug("Deleting product")

	if err := p.db.DeleteSecurityByProductID(ctx, id); err != nil {
		logger.Error("Error deleting security for product", "error", err)
	}

	if err := p.db.DeletePullRequestsByProductID(ctx, id); err != nil {
		logger.Error("Error deleting pull requests for product", "error", err)
	}

	if err := p.db.DeleteReposByProductID(ctx, id); err != nil {
		logger.Error("Error deleting repos for product", "error", err)
	}

	return p.db.DeleteProduct(ctx, id)
}

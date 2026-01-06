package watchtower

import (
	"context"
	"time"
	"watchtower/internal/organisations"
)

type Service struct {
	ctx        context.Context
	orgSvc     *organisations.Service
	productSvc *productsService
	ghClient   ghClient
}

type productsService struct {
	db          ProductStore
	repoService *repoService
}

type repoService struct {
	db RepoStore
}

type ProductDTO struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductOrganisationDTO struct {
	ProductID      *int64 `json:"product_id,omitempty"`
	OrganisationID *int64 `json:"organisation_id,omitempty"`
}

type PullRequestDTO struct {
	ID             int64     `json:"id"`
	ExternalID     string    `json:"external_id"`
	Title          string    `json:"title"`
	RepositoryName string    `json:"repository_name"`
	URL            string    `json:"url"`
	State          string    `json:"state"`
	Author         string    `json:"author"`
	Tag            string    `json:"tag"`
	ProductName    string    `json:"product_name"`
	MergedAt       time.Time `json:"merged_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type RepositoryDTO struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	Topic       string    `json:"topic"`
	Owner       string    `json:"owner"`
	ProductName string    `json:"product_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SecurityDTO struct {
	ID             int64     `json:"id"`
	ExternalID     string    `json:"external_id"`
	RepositoryName string    `json:"repository_name"`
	PackageName    string    `json:"package_name"`
	State          string    `json:"state"`
	Severity       string    `json:"severity"`
	PatchedVersion string    `json:"patched_version"`
	Tag            string    `json:"tag"`
	ProductName    string    `json:"product_name"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

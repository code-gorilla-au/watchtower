package products

import (
	"time"
)

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

type CreateProductParams struct {
	Name string
	Desc string
	Tags []string
}

type UpdateProductParams struct {
	ID          int64
	Name        string
	Description string
	Tags        []string
}

type CreateRepoParams struct {
	Name  string
	Url   string
	Topic string
	Owner string
}

type UpdateRepoParams struct {
	ID    int64
	Name  string
	Url   string
	Topic string
	Owner string
}

type CreatePRParams struct {
	ExternalID     string
	Title          string
	RepositoryName string
	Url            string
	State          string
	Author         string
	MergedAt       *time.Time
	CreatedAt      time.Time
}

type CreateSecurityParams struct {
	ExternalID     string
	RepositoryName string
	PackageName    string
	State          string
	Severity       string
	PatchedVersion string
	FixedAt        *time.Time
	CreatedAt      time.Time
}

package watchtower

import (
	"context"
	"time"
	"watchtower/internal/database"
)

type Service struct {
	ctx      context.Context
	db       *database.Queries
	ghClient ghClient
}

// DTOs with snake_case JSON tags and time.Time timestamps

type OrganisationDTO struct {
	ID           int64     `json:"id"`
	FriendlyName string    `json:"friendly_name"`
	Description  string    `json:"description"`
	Namespace    string    `json:"namespace"`
	DefaultOrg   bool      `json:"default_org"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
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
	MergedAt       time.Time `json:"merged_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type RepositoryDTO struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Topic     string    `json:"topic"`
	Owner     string    `json:"owner"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SecurityDTO struct {
	ID             int64     `json:"id"`
	ExternalID     string    `json:"external_id"`
	RepositoryName string    `json:"repository_name"`
	PackageName    string    `json:"package_name"`
	State          string    `json:"state"`
	Severity       string    `json:"severity"`
	PatchedVersion string    `json:"patched_version"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

package watchtower

import (
	"time"
	"watchtower/internal/database"
)

type Service struct {
	db database.Queries
}

// DTOs with snake_case JSON tags and time.Time timestamps

type OrganisationDTO struct {
	ID           int64     `json:"id"`
	FriendlyName string    `json:"friendly_name"`
	Namespace    string    `json:"namespace"`
	DefaultOrg   bool      `json:"default_org"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ProductDTO struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Tags      *string   `json:"tags,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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

// Conversion helpers from database models to DTOs

func toTime(ts int64) time.Time {
	return time.Unix(ts, 0).UTC()
}

func ToOrganisationDTO(m database.Organisation) OrganisationDTO {
	return OrganisationDTO{
		ID:           m.ID,
		FriendlyName: m.FriendlyName,
		Namespace:    m.Namespace,
		DefaultOrg:   m.DefaultOrg,
		CreatedAt:    toTime(m.CreatedAt),
		UpdatedAt:    toTime(m.UpdatedAt),
	}
}

func ToProductDTO(m database.Product) ProductDTO {
	var tagsPtr *string
	if m.Tags.Valid {
		val := m.Tags.String
		tagsPtr = &val
	}
	return ProductDTO{
		ID:        m.ID,
		Name:      m.Name,
		Tags:      tagsPtr,
		CreatedAt: toTime(m.CreatedAt),
		UpdatedAt: toTime(m.UpdatedAt),
	}
}

func ToProductOrganisationDTO(m database.ProductOrganisation) ProductOrganisationDTO {
	var pidPtr *int64
	var oidPtr *int64
	if m.ProductID.Valid {
		v := m.ProductID.Int64
		pidPtr = &v
	}
	if m.OrganisationID.Valid {
		v := m.OrganisationID.Int64
		oidPtr = &v
	}
	return ProductOrganisationDTO{
		ProductID:      pidPtr,
		OrganisationID: oidPtr,
	}
}

func ToPullRequestDTO(m database.PullRequest) PullRequestDTO {
	return PullRequestDTO{
		ID:             m.ID,
		ExternalID:     m.ExternalID,
		Title:          m.Title,
		RepositoryName: m.RepositoryName,
		URL:            m.Url,
		State:          m.State,
		Author:         m.Author,
		MergedAt:       toTime(m.MergedAt),
		CreatedAt:      toTime(m.CreatedAt),
		UpdatedAt:      toTime(m.UpdatedAt),
	}
}

func ToRepositoryDTO(m database.Repository) RepositoryDTO {
	return RepositoryDTO{
		ID:        m.ID,
		Name:      m.Name,
		URL:       m.Url,
		Topic:     m.Topic,
		Owner:     m.Owner,
		CreatedAt: toTime(m.CreatedAt),
		UpdatedAt: toTime(m.UpdatedAt),
	}
}

func ToSecurityDTO(m database.Security) SecurityDTO {
	return SecurityDTO{
		ID:             m.ID,
		ExternalID:     m.ExternalID,
		RepositoryName: m.RepositoryName,
		PackageName:    m.PackageName,
		State:          m.State,
		Severity:       m.Severity,
		PatchedVersion: m.PatchedVersion,
		CreatedAt:      toTime(m.CreatedAt),
		UpdatedAt:      toTime(m.UpdatedAt),
	}
}

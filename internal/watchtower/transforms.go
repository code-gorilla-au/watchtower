package watchtower

import (
	"encoding/json"
	"time"
	"watchtower/internal/database"
)

// Conversion helpers from database models to DTOs

func toTime(ts int64) time.Time {
	return time.Unix(ts, 0).UTC()
}

func ToOrganisationDTO(m database.Organisation) OrganisationDTO {
	return OrganisationDTO{
		ID:           m.ID,
		FriendlyName: m.FriendlyName,
		Description:  m.Description,
		Namespace:    m.Namespace,
		DefaultOrg:   m.DefaultOrg,
		CreatedAt:    toTime(m.CreatedAt),
		UpdatedAt:    toTime(m.UpdatedAt),
	}
}

func ToProductDTOs(models []database.Product) []ProductDTO {
	result := make([]ProductDTO, 0, len(models))
	for _, m := range models {
		result = append(result, ToProductDTO(m))
	}

	return result
}

func ToProductDTO(m database.Product) ProductDTO {
	var tagList []string

	if m.Tags.Valid {
		val := m.Tags.String

		_ = json.Unmarshal([]byte(val), &tagList)
	}

	return ProductDTO{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Tags:        tagList,
		CreatedAt:   toTime(m.CreatedAt),
		UpdatedAt:   toTime(m.UpdatedAt),
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
		MergedAt:       toTime(m.MergedAt.Int64),
		CreatedAt:      toTime(m.CreatedAt),
		UpdatedAt:      toTime(m.UpdatedAt),
	}
}

func toPullRequestDTOs(models []database.PullRequest) []PullRequestDTO {
	result := make([]PullRequestDTO, 0, len(models))
	for _, m := range models {
		result = append(result, ToPullRequestDTO(m))
	}

	return result
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

func ToSecurityDTOs(models []database.Security) []SecurityDTO {
	result := make([]SecurityDTO, 0, len(models))
	for _, m := range models {
		result = append(result, ToSecurityDTO(m))
	}

	return result
}

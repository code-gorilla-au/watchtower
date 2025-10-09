package watchtower

import (
	"encoding/json"
	"time"
	"watchtower/internal/database"

	"watchtower/internal/github"
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

func ToOrganisationDTOs(models []database.Organisation) []OrganisationDTO {
	result := make([]OrganisationDTO, 0, len(models))
	for _, m := range models {
		result = append(result, ToOrganisationDTO(m))
	}

	return result
}

func ToInternalOrganisation(m database.Organisation) InternalOrganisation {
	return InternalOrganisation{
		OrganisationDTO: OrganisationDTO{
			ID:           m.ID,
			FriendlyName: m.FriendlyName,
			Description:  m.Description,
			Namespace:    m.Namespace,
			DefaultOrg:   m.DefaultOrg,
			CreatedAt:    toTime(m.CreatedAt),
			UpdatedAt:    toTime(m.UpdatedAt),
		},
		Token: m.Token,
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

func ProductToPullRequestDTO(m database.GetPullRequestByProductIDAndStateRow) PullRequestDTO {
	return PullRequestDTO{
		ID:             m.ID,
		ExternalID:     m.ExternalID,
		Title:          m.Title,
		RepositoryName: m.RepositoryName,
		URL:            m.Url,
		State:          m.State,
		Author:         m.Author,
		Tag:            m.Tag,
		ProductName:    m.ProductName,
		MergedAt:       toTime(m.MergedAt.Int64),
		CreatedAt:      toTime(m.CreatedAt),
		UpdatedAt:      toTime(m.UpdatedAt),
	}
}

func toPullRequestDTOs(models []database.GetPullRequestByProductIDAndStateRow) []PullRequestDTO {
	result := make([]PullRequestDTO, 0, len(models))
	for _, m := range models {
		result = append(result, ProductToPullRequestDTO(m))
	}

	return result
}

func OrgToPullRequestDTO(m database.GetPullRequestsByOrganisationAndStateRow) PullRequestDTO {
	return PullRequestDTO{
		ID:             m.ID,
		ExternalID:     m.ExternalID,
		Title:          m.Title,
		RepositoryName: m.RepositoryName,
		URL:            m.Url,
		State:          m.State,
		Author:         m.Author,
		Tag:            m.Tag,
		ProductName:    m.ProductName,
		MergedAt:       toTime(m.MergedAt.Int64),
		CreatedAt:      toTime(m.CreatedAt),
		UpdatedAt:      toTime(m.UpdatedAt),
	}
}

func orgToPullRequestDTOs(models []database.GetPullRequestsByOrganisationAndStateRow) []PullRequestDTO {
	result := make([]PullRequestDTO, 0, len(models))
	for _, m := range models {
		result = append(result, OrgToPullRequestDTO(m))
	}

	return result
}

func ToRepositoryDTO(m database.GetReposByProductIDRow) RepositoryDTO {
	return RepositoryDTO{
		ID:          m.ID,
		Name:        m.Name,
		URL:         m.Url,
		Topic:       m.Topic,
		Owner:       m.Owner,
		ProductName: m.ProductName,
		CreatedAt:   toTime(m.CreatedAt),
		UpdatedAt:   toTime(m.UpdatedAt),
	}
}

func ToSecurityDTO(m database.GetSecurityByProductIDAndStateRow) SecurityDTO {
	return SecurityDTO{
		ID:             m.ID,
		ExternalID:     m.ExternalID,
		RepositoryName: m.RepositoryName,
		PackageName:    m.PackageName,
		State:          m.State,
		Severity:       m.Severity,
		PatchedVersion: m.PatchedVersion,
		Tag:            m.Tag,
		ProductName:    m.ProductName,
		CreatedAt:      toTime(m.CreatedAt),
		UpdatedAt:      toTime(m.UpdatedAt),
	}
}

func ToSecurityDTOs(models []database.GetSecurityByProductIDAndStateRow) []SecurityDTO {
	result := make([]SecurityDTO, 0, len(models))
	for _, m := range models {
		result = append(result, ToSecurityDTO(m))
	}

	return result
}

func OrgToSecurityDTO(m database.GetSecurityByOrganisationAndStateRow) SecurityDTO {
	return SecurityDTO{
		ID:             m.ID,
		ExternalID:     m.ExternalID,
		RepositoryName: m.RepositoryName,
		PackageName:    m.PackageName,
		State:          m.State,
		Severity:       m.Severity,
		PatchedVersion: m.PatchedVersion,
		Tag:            m.Tag,
		ProductName:    m.ProductName,
		CreatedAt:      toTime(m.CreatedAt),
		UpdatedAt:      toTime(m.UpdatedAt),
	}
}

func OrgToSecurityDTOs(models []database.GetSecurityByOrganisationAndStateRow) []SecurityDTO {
	result := make([]SecurityDTO, 0, len(models))
	for _, m := range models {
		result = append(result, OrgToSecurityDTO(m))
	}

	return result
}

func ToCreateRepoFromGithub(repos []github.Node[github.Repository], tag string) []CreateRepoParams {
	result := make([]CreateRepoParams, 0, len(repos))

	for _, repo := range repos {
		result = append(result, CreateRepoParams{
			Name:  repo.Node.Name,
			Url:   repo.Node.Url,
			Topic: tag,
			Owner: repo.Node.Owner.Login,
		})
	}

	return result
}

func ToCreatePRsFromGithubRepos(prs github.RootNode[github.PullRequest], repoName string) []CreatePRParams {
	result := make([]CreatePRParams, len(prs.Nodes))

	for _, pr := range prs.Nodes {
		result = append(result, CreatePRParams{
			ExternalID:     pr.ID,
			Title:          pr.Title,
			RepositoryName: repoName,
			Url:            pr.Permalink,
			State:          string(pr.State),
			Author:         pr.Author.Login,
			MergedAt:       pr.MergedAt,
			CreatedAt:      pr.CreatedAt,
		})
	}

	return result
}

func ToSecParamsFromGithubVulnerabilities(secs github.RootNode[github.VulnerabilityAlerts], repoName string) []CreateSecurityParams {
	result := make([]CreateSecurityParams, len(secs.Nodes))

	for _, sec := range secs.Nodes {
		result = append(result, CreateSecurityParams{
			ExternalID:     sec.ID,
			RepositoryName: repoName,
			PackageName:    sec.SecurityVulnerability.Package.Name,
			State:          string(sec.State),
			Severity:       string(sec.SecurityVulnerability.Advisory.Severity),
			PatchedVersion: sec.SecurityVulnerability.FirstPatchedVersion.Identifier,
			CreatedAt:      sec.CreatedAt,
			FixedAt:        sec.FixedAt,
		})
	}

	return result
}

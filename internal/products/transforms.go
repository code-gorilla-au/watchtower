package products

import (
	"encoding/json"
	"time"
	"watchtower/internal/database"
	"watchtower/internal/github"
)

// toTime converts a Unix timestamp to time.Time.
func toTime(ts int64) time.Time {
	return time.Unix(ts, 0).UTC()
}

func toProductDTOs(models []database.Product) []ProductDTO {
	result := make([]ProductDTO, 0, len(models))
	for _, m := range models {
		result = append(result, toProductDTO(m))
	}

	return result
}

func toProductDTO(m database.Product) ProductDTO {
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

func productToPullRequestDTO(m database.GetPullRequestByProductIDAndStateRow) PullRequestDTO {
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
		result = append(result, productToPullRequestDTO(m))
	}

	return result
}

func orgToPullRequestDTO(m database.GetPullRequestsByOrganisationAndStateRow) PullRequestDTO {
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
		result = append(result, orgToPullRequestDTO(m))
	}

	return result
}

func toRepositoryDTO(m database.GetReposByProductIDRow) RepositoryDTO {
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

func toSecurityDTO(m database.GetSecurityByProductIDAndStateRow) SecurityDTO {
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

func toSecurityDTOs(models []database.GetSecurityByProductIDAndStateRow) []SecurityDTO {
	result := make([]SecurityDTO, 0, len(models))
	for _, m := range models {
		result = append(result, toSecurityDTO(m))
	}

	return result
}

func orgToSecurityDTO(m database.GetSecurityByOrganisationAndStateRow) SecurityDTO {
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

func orgToSecurityDTOs(models []database.GetSecurityByOrganisationAndStateRow) []SecurityDTO {
	result := make([]SecurityDTO, 0, len(models))
	for _, m := range models {
		result = append(result, orgToSecurityDTO(m))
	}

	return result
}

func toCreateRepoFromGithub(repos []github.Node[github.Repository], tag string) []CreateRepoParams {
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

func toCreatePRsFromGithubRepos(prs github.RootNode[github.PullRequest], repoName string) []CreatePRParams {
	result := make([]CreatePRParams, 0, len(prs.Nodes))

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

func toSecParamsFromGithubVulnerabilities(secs github.RootNode[github.VulnerabilityAlerts], repoName string) []CreateSecurityParams {
	result := make([]CreateSecurityParams, 0, len(secs.Nodes))

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

func fromRecentlyChangedPRModel(model database.GetRecentPullRequestsRow) RecentlyChangedEntity {
	return RecentlyChangedEntity{
		ExternalID:     model.ExternalID,
		OrganisationID: model.OrganisationID.Int64,
		RepositoryName: model.RepositoryName,
	}
}

func fromRecentlyChangedPRModels(models []database.GetRecentPullRequestsRow) []RecentlyChangedEntity {
	result := make([]RecentlyChangedEntity, 0, len(models))
	for _, m := range models {
		result = append(result, fromRecentlyChangedPRModel(m))
	}

	return result
}

func fromRecentlyChangedSecurityModel(model database.GetRecentSecurityRow) RecentlyChangedEntity {
	return RecentlyChangedEntity{
		ExternalID:     model.ExternalID,
		OrganisationID: model.OrganisationID.Int64,
		RepositoryName: model.RepositoryName,
	}
}

func fromRecentlyChangedSecurityModels(models []database.GetRecentSecurityRow) []RecentlyChangedEntity {
	result := make([]RecentlyChangedEntity, 0, len(models))
	for _, m := range models {
		result = append(result, fromRecentlyChangedSecurityModel(m))
	}

	return result
}

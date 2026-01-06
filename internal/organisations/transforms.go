package organisations

import (
	"time"
	"watchtower/internal/database"
)

// toTime converts a Unix timestamp to time.Time.
func toTime(ts int64) time.Time {
	return time.Unix(ts, 0).UTC()
}

// toOrganisationDTO converts a database Organisation model to an OrganisationDTO.
func toOrganisationDTO(m database.Organisation) OrganisationDTO {
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

// toOrganisationDTOs converts a slice of database Organisation models to a slice of OrganisationDTOs.
func toOrganisationDTOs(models []database.Organisation) []OrganisationDTO {
	result := make([]OrganisationDTO, 0, len(models))
	for _, m := range models {
		result = append(result, toOrganisationDTO(m))
	}

	return result
}

// toInternalOrganisation converts a database Organisation model to an InternalOrganisation.
func toInternalOrganisation(m database.Organisation) InternalOrganisation {
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

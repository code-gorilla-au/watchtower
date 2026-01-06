package organisations

import (
	"time"
)

type CreateOrgParams struct {
	FriendlyName string
	Namespace    string
	Token        string
	Description  string
}

type UpdateOrgParams struct {
	ID           int64
	FriendlyName string
	Namespace    string
	Token        string
	Description  string
	DefaultOrg   bool
}

type OrganisationDTO struct {
	ID           int64     `json:"id"`
	FriendlyName string    `json:"friendly_name"`
	Description  string    `json:"description"`
	Namespace    string    `json:"namespace"`
	DefaultOrg   bool      `json:"default_org"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type InternalOrganisation struct {
	OrganisationDTO
	Token string
}

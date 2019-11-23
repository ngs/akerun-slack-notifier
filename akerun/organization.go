package akerun

// Organization .
type Organization struct {
	ID string `json:"id"`
}

// OrganizationsResponse .
type OrganizationsResponse struct {
	Organizations []Organization `json:"organizations"`
}

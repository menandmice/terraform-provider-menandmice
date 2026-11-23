package menandmice

import (
	"errors"
	"fmt"
	"strings"
)

// Organization represents a Micetro organization object.
type Organization struct {
	Ref string `json:"ref,omitempty"`
	OrganizationProperties
}

// OrganizationProperties are the properties of an organization object.
// They are used when listing organizations, e.g. to resolve an organization
// name to its reference.
type OrganizationProperties struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type findOrganizationsResponse struct {
	Organizations []Organization `json:"organizations"`
	TotalResults  int            `json:"totalResults"`
}

// FindOrganizations returns organizations matching the given filter.
// The filter is a map of property name to value, serialized into a
// Micetro query string.
func (c Mmclient) FindOrganizations(filter map[string]interface{}) ([]Organization, error) {
	var re findOrganizationsResponse
	var query map[string]interface{}
	if filter != nil && len(filter) > 0 {
		query = map[string]interface{}{"filter": map2filter(filter)}
	}
	err := c.Get(&re, "organizations/", query)
	return re.Organizations, err
}

// login authenticates against the Micetro session API and stores the
// returned session token as a Bearer token on the client. Operating within a
// specific organization (see SetCurrentOrganization) requires this session
// based authentication.
func (c *Mmclient) login(username, password string) error {
	var loginResponse struct {
		Session string `json:"session"`
	}
	var errorResponse ErrorResponse
	r, err := c.R().
		SetBody(map[string]string{"loginName": username, "password": password}).
		SetResult(&loginResponse).
		SetError(&errorResponse).
		Post("micetro/sessions")
	if err != nil {
		return err
	}
	if err := ResponseError(r, errorResponse); err != nil {
		return err
	}
	if loginResponse.Session == "" {
		return errors.New("login succeeded but no session token was returned")
	}
	c.SetAuthToken(loginResponse.Session)
	return nil
}

// resolveOrganization turns the user supplied organization value into an
// organization object ref. If the value already looks like a ref (contains a
// '/') it is returned unchanged, otherwise it is treated as an organization
// name and looked up.
func (c *Mmclient) resolveOrganization(organization string) (string, error) {
	if strings.Contains(organization, "/") {
		return organization, nil
	}
	organizations, err := c.FindOrganizations(map[string]interface{}{"name": organization})
	if err != nil {
		return "", err
	}
	switch len(organizations) {
	case 0:
		return "", fmt.Errorf("organization %q not found", organization)
	case 1:
		return organizations[0].Ref, nil
	default:
		return "", fmt.Errorf("multiple organizations found for %q, please use the organization ref instead", organization)
	}
}

// SetCurrentOrganization switches the current session to the given
// organization. After this call, operations such as creating DNS records or
// DHCP scopes are performed within that organization.
func (c *Mmclient) SetCurrentOrganization(ref string) error {
	var errorResponse ErrorResponse
	r, err := c.R().
		SetBody(map[string]string{"organizationRef": ref}).
		SetError(&errorResponse).
		Post("command/SetCurrentOrganization")
	if err != nil {
		return err
	}
	return ResponseError(r, errorResponse)
}

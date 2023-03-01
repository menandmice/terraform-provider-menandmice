package menandmice

type DHCPReservation struct {
	Ref               string   `json:"ref,omitempty"`
	OwnerRef          string   `json:"ownerRef,omitempty"`
	ReservationMethod string   `json:"reservationMethod,omitempty"`
	Addresses         []string `json:"addresses,omitempty"`
	DHCPReservationPropertie
}

type DHCPReservationPropertie struct {
	Name             string `json:"name"`
	Type             string `json:"type,omitempty"`
	Description      string `json:"description,omitempty"`
	ClientIdentifier string `json:"clientIdentifier,omitempty"`
	DDNSHostName     string `json:"ddnsHostName,omitempty"`
	Filename         string `json:"filename,omitempty"`
	ServerName       string `json:"serverName,omitempty"`
	NextServer       string `json:"nextServer,omitempty"`
}

type readDHCPReservationResponse struct {
	Result struct {
		DHCPReservation `json:"dhcpReservation"`
	} `json:"result"`
}

func (c *Mmclient) ReadDHCPReservation(ref string) (*DHCPReservation, error) {
	var re readDHCPReservationResponse
	err := c.Get(&re, "DHCPReservations/"+ref, nil)
	if reqError, ok := err.(*RequestError); ok && reqError.StatusCode == ResourceNotFound {
		//DHCPReservationNotFound not found
		return nil, nil
	}
	return &re.Result.DHCPReservation, err
}

type createDHCPReservationRequest struct {
	DHCPReservation DHCPReservation `json:"dhcpReservation"`
	SaveComment     string          `json:"saveComment"`
}

func (c *Mmclient) CreateDHCPReservation(dhcpReservation DHCPReservation, owner string) (string, error) {

	var objRef string

	postcreate := createDHCPReservationRequest{
		DHCPReservation: dhcpReservation,
		SaveComment:     "created by terraform",
	}
	var re RefResponse
	err := c.Post(postcreate, &re, "DHCPServers/"+owner+"/DHCPReservations")

	if err != nil {
		return objRef, err
	}

	return re.Result.Ref, err

}

func (c *Mmclient) DeleteDHCPReservation(ref string) error {

	err := c.Delete(deleteRequest("DHCPReservation"), "DHCPReservations/"+ref)
	if reqError, ok := err.(*RequestError); ok && reqError.StatusCode == ResourceNotFound {
		//DHCPReservationNotFound not found, so nothing to delete
		return nil
	}
	return err
}

type updateDHCPReservationRequest struct {
	Ref               string `json:"ref"`
	ObjType           string `json:"objType"`
	SaveComment       string `json:"saveComment"`
	DeleteUnspecified bool   `json:"deleteUnspecified"`

	Properties DHCPReservationPropertie `json:"properties"`
}

func (c *Mmclient) UpdateDHCPReservation(properties DHCPReservationPropertie, ref string) error {

	update := updateDHCPReservationRequest{
		Ref:               ref,
		ObjType:           "DHCPReservation",
		SaveComment:       "updated by terraform",
		DeleteUnspecified: true,
		Properties:        properties,
	}

	return c.Put(update, "DHCPReservations/"+ref)
}

type DHCPScope struct {
	Ref           string `json:"ref"`
	Name          string `json:"name"`
	RangeRef      string `json:"rangeRef"`
	DHCPServerRef string `json:"dhcpServerRef"`
	Superscope    string `json:"superscope"`
	Description   string `json:"description"`
	Available     int    `json:"available"`
	Enabled       bool   `json:"enabled"`
}

type findDHCPScopeResponse struct {
	Result struct {
		DHCPScopes   []DHCPScope `json:"dhcpScopes"`
		TotalResults int         `json:"totalResults"`
	} `json:"result"`
}

// function ReadDHCPScope is never needed

func (c Mmclient) FindDHCPScope(filter map[string]interface{}) ([]DHCPScope, error) {
	var re findDHCPScopeResponse

	query := map[string]interface{}{"filter": map2filter(filter)}

	err := c.Get(&re, "DHCPScopes/", query)
	return re.Result.DHCPScopes, err
}

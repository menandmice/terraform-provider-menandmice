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

type ReadDHCPReservationResponse struct {
	Result struct {
		DHCPReservation `json:"dhcpReservation"`
	} `json:"result"`
}

func (c *Mmclient) ReadDHCPReservation(ref string) (DHCPReservation, error) {
	var re ReadDHCPReservationResponse
	err := c.Get(&re, "DHCPReservations/"+ref, nil, nil)
	return re.Result.DHCPReservation, err
}

type CreateDHCPReservationRequest struct {
	DHCPReservation DHCPReservation `json:"dhcpReservation"`
	SaveComment     string          `json:"saveComment"`
}

func (c *Mmclient) CreateDHCPReservation(dhcpReservation DHCPReservation, owner string) (string, error) {

	var objRef string

	postcreate := CreateDHCPReservationRequest{
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

	return c.Delete(deleteRequest("DHCPReservation"), "DHCPReservations/"+ref)
}

type UpdateDHCPReservationRequest struct {
	Ref               string `json:"ref"`
	ObjType           string `json:"objType"`
	SaveComment       string `json:"saveComment"`
	DeleteUnspecified bool   `json:"deleteUnspecified"`

	Properties DHCPReservationPropertie `json:"properties"`
}

func (c *Mmclient) UpdateDHCPReservation(properties DHCPReservationPropertie, ref string) error {

	update := UpdateDHCPReservationRequest{
		Ref:               ref,
		ObjType:           "DHCPReservation",
		SaveComment:       "updated by terraform",
		DeleteUnspecified: true,
		Properties:        properties,
	}

	return c.Put(update, "DHCPReservations/"+ref)
}

package menandmice

import "fmt"

type IPAMRecord struct {
	Ref     string    `json:"addrRef,omitempty"`
	Address string    `json:"address"`
	DNSHost []DNSHost `json:"dnsHosts,omitempty"`
	// DHCPReservations []DHCPReservation `json:"dhcpReservations,omitempty"`
	// DHCPLeases []???  "dhcpLeases,omitempty"`
	DiscoveryType             string    `json:"discoveryType,omitempty"`
	PTRStatus                 string    `json:"ptrStatus,omitempty"`
	LastSeenDate              string    `json:"lastSeenDate,omitempty"`
	LastDiscoveryDate         string    `json:"lastDiscoveryDate,omitempty"`
	LastKnownClientIdentifier string    `json:"lastKnownClientIdentifier,omitempty"`
	ExtraneousPTR             bool      `json:"extraneousPTR,omitempty"`
	Device                    string    `json:"device,omitempty"`
	State                     string    `json:"state,omitempty"`
	Usage                     int       `json:"usage,omitempty"`
	HoldInfo                  *HoldInfo `json:"holdInfo,omitempty"`
	IPAMProperties
}
type IPAMProperties struct {
	Claimed          bool              `json:"claimed"`
	Interface        string            `json:"interace,omitempty"`
	CustomProperties map[string]string `json:"customProperties,omitempty"`

	// CloudDeviceInfo []string `json:"cloudDeviceInfo,omitempty"`
}

type HoldInfo struct {
	ExpiryTime string `json:"expiryTime,omitempty"`
	Username   string `json:"username,omitempty"`
}

type DNSHost struct {
	DNSRecord      DNSRecord   `json:"dnsRecord"`
	PTRStatus      bool        `json:"ptrStatus"`
	RelatedRecords []DNSRecord `json:"relatedRecords"`
}

type ReadIPAMRECResponse struct {
	Result struct {
		IPAMRecord `json:"ipamRecord"`
	} `json:"result"`
}

func (c *Mmclient) ReadIPAMRec(ref string) (IPAMRecord, error) {
	var re ReadIPAMRECResponse
	err := c.Get(&re, "IPAMRecords/"+ref, nil, nil)
	return re.Result.IPAMRecord, err
}

// TODO because this will only set IPAMProperties. and ignore other. maybe change to :
// func (c *Mmclient) CreateIPAMRec(ipamProperites IPAMProperties,rec string) error {
func (c *Mmclient) CreateIPAMRec(ipamRecord IPAMRecord) error {

	// TODO this function will query ipamRecord bassed on ip. But this is not unque
	//		you have to use SetCurrentAddressSpace in the client initalisation. or here en prevent race conditions

	// we need to check if IPAMRecord not already exist. because creation is done via update/PUT
	existingIPAMRecord, err := c.ReadIPAMRec(ipamRecord.Address)

	if err != nil {
		return err
	}
	if existingIPAMRecord.Claimed == true {
		return fmt.Errorf("There already exist a DHCPReservations for: %v", existingIPAMRecord.Address)
	}
	return c.UpdateIPAMRec(ipamRecord.IPAMProperties, ipamRecord.Address)
}

func (c *Mmclient) DeleteIPAMRec(ref string) error {

	return c.Delete(deleteRequest("IPAddress"), "IPAMRecords/"+ref)
}

type UpdateIPAMRecRequest struct {
	Ref               string `json:"ref"`
	ObjType           string `json:"objType"`
	SaveComment       string `json:"saveComment"`
	DeleteUnspecified bool   `json:"deleteUnspecified"`

	// we cant use IPAMProperties for this because CustomProperties should be flattend first
	Properties map[string]interface{} `json:"properties"`
}

func (c *Mmclient) UpdateIPAMRec(ipamProperties IPAMProperties, ref string) error {

	// A work around to create properties with same fields as DNSZoneProperties but with flattend CustomProperties
	// first mask CustomProperties in DNSZoneProperties
	// Then convert to map considerting `json:"omitempty"`
	// Then add CustomProperties 1 by 1

	customProperties := ipamProperties.CustomProperties
	ipamProperties.CustomProperties = nil

	properties, err := toMap(ipamProperties)

	if err != nil {
		return err
	}

	for key, value := range customProperties {
		properties[key] = value
	}

	update := UpdateIPAMRecRequest{
		Ref:               ref,
		ObjType:           "IPAddress",
		SaveComment:       "updated by terraform",
		DeleteUnspecified: true,
		Properties:        properties,
	}

	return c.Put(update, "IPAMRecords/"+ref)
}

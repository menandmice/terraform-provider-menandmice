package menandmice

type IPAMRecord struct {
	Ref     string    `json:"addrRef,omitempty"`
	Address string    `json:"address"`
	DNSHost []DNSHost `json:"dnsHosts,omitempty"`
	// DHCPReservations          []??? `json:"dhcpReservations,omitempty"`
	// DHCPLeases []???  "dhcpLeases,omitempty"`
	//TODO how to set DiscoveryType
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

func (c *Mmclient) CreateIPAMRec(ipamRecord IPAMRecord) error {
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

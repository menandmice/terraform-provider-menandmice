package menandmice

import "fmt"

type IPAMRecord struct {
	Ref     string `json:"addrRef,omitempty"`
	Address string `json:"address"`
	// DNSHost []DNSHost `json:"dnsHosts,omitempty"`	// works not used for now
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
	PTRStatus      string      `json:"ptrStatus"`
	RelatedRecords []DNSRecord `json:"relatedRecords"`
}

type readIPAMRECResponse struct {
	Result struct {
		IPAMRecord `json:"ipamRecord"`
	} `json:"result"`
}

func (c *Mmclient) ReadIPAMRec(ref string) (IPAMRecord, error) {
	var re readIPAMRECResponse
	err := c.Get(&re, "IPAMRecords/"+ref, nil)
	return re.Result.IPAMRecord, err
}

// TODO because this will only set IPAMProperties and ignore others. Maybe change to:
// func (c *Mmclient) CreateIPAMRec(ipamProperites IPAMProperties,rec string) error {
func (c *Mmclient) CreateIPAMRec(ipamRecord IPAMRecord) error {

	// TODO this function will query ipamRecord bassed on IP address. But this is not unique
	//		you have to use SetCurrentAddressSpace in the client initalisation or here and prevent race conditions

	// we need to check if IPAMRecord already exists, because creation is done via update/PUT
	existingIPAMRecord, err := c.ReadIPAMRec(ipamRecord.Address)

	if err != nil {
		return err
	}
	if existingIPAMRecord.Claimed {
		return fmt.Errorf("DHCPReservations already exists for: %v", existingIPAMRecord.Address)
	}
	return c.UpdateIPAMRec(ipamRecord.IPAMProperties, ipamRecord.Address)
}

func (c *Mmclient) DeleteIPAMRec(ref string) error {

	return c.Delete(deleteRequest("IPAddress"), "IPAMRecords/"+ref)
}

type updateIPAMRecRequest struct {
	Ref               string `json:"ref"`
	ObjType           string `json:"objType"`
	SaveComment       string `json:"saveComment"`
	DeleteUnspecified bool   `json:"deleteUnspecified"`

	// We can't use IPAMProperties for this because CustomProperties should be flattend first
	Properties map[string]interface{} `json:"properties"`
}

func (c *Mmclient) UpdateIPAMRec(ipamProperties IPAMProperties, ref string) error {

	// A workaround to create IPAMproperties with same fields as IPAMPropertiesproperties but with flattend CustomProperties
	// First mask CustomProperties in IPAMProperties
	// Then convert to map considering `json:"omitempty"`
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

	update := updateIPAMRecRequest{
		Ref:               ref,
		ObjType:           "IPAddress",
		SaveComment:       "updated by terraform",
		DeleteUnspecified: true,
		Properties:        properties,
	}

	return c.Put(update, "IPAMRecords/"+ref)
}

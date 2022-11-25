package menandmice

type DNSZone struct {
	Ref          string   `json:"ref,omitempty"`
	Name         string   `json:"name"`
	AdIntegrated bool     `json:"adIntegrated"`
	DNSViewRef   string   `json:"dnsViewRef,omitempty"`
	DNSViewRefs  []string `json:"dnsViewRefs,omitempty"`
	Authority    string   `json:"authority,omitempty"`

	Created      string `json:"created,omitempty"`
	LastModified string `json:"lastModified,omitempty"`

	// Properties that can be updated
	DNSZoneProperties
}

type DNSZoneProperties struct {
	Dynamic           bool              `json:"dynamic,omitempty"`
	ZoneType          string            `json:"type,omitempty"`
	DnssecSigned      bool              `json:"dnssecSigned,omitempty"`
	KskIDs            string            `json:"kskIDs,omitempty"`
	ZskIDs            string            `json:"zskIDs,omitempty"`
	CustomProperties  map[string]string `json:"customProperties,omitempty"`
	AdReplicationType string            `json:"adReplicationType,omitempty"`
	AdPartition       string            `json:"adPartition,omitempty"`
	DisplayName       string            `json:"displyaName,omitempty"`
}

type FindDNSZoneResponse struct {
	Result struct {
		DNSZones     []DNSZone `json:"dnsZones"`
		TotalResults int       `json:"totalResults"`
	} `json:"result"`
}

func (c Mmclient) FindDNSZones(limit int, filter map[string]string) ([]DNSZone, error) {
	var re FindDNSZoneResponse
	query := map[string]interface{}{
		"limit": limit,
	}

	if folderRef, ok := filter["folderRef"]; ok {
		query["folderRef"] = folderRef
		delete(filter, "folderRef")
	}

	if dnsViewRef, ok := filter["dnsViewRef"]; ok {
		query["dnsViewRef"] = dnsViewRef
		delete(filter, "dnsViewRef")
	}

	if dnsServerRef, ok := filter["dnsServerRef"]; ok {
		query["dnsServerRef"] = dnsServerRef
		delete(filter, "dnsServerRef")
	}

	err := c.Get(&re, "dnszones/", query, filter)
	// TODO return empyt list if you get error view server etz does not exist
	return re.Result.DNSZones, err
}

type ReadDNSZoneResponse struct {
	Result struct {
		DNSZone `json:"dnsZone"`
	} `json:"result"`
}

func (c Mmclient) ReadDNSZone(ref string) (*DNSZone, error) {
	var re ReadDNSZoneResponse
	err := c.Get(&re, "DNSZones/"+ref, nil, nil)
	if reqError, ok := err.(*RequestError); ok && reqError.StatusCode == ResourceNotFound {
		return nil, nil
	}

	return &re.Result.DNSZone, err
}

type CreateDNSZoneRequest struct {
	DNSZone     DNSZone  `json:"dnsZone"`
	SaveComment string   `json:"saveComment"`
	Masters     []string `json:"masters,omitempty"`
}

func (c *Mmclient) CreateDNSZone(dnszone DNSZone, masters []string) (string, error) {
	var objRef string
	postcreate := CreateDNSZoneRequest{
		DNSZone:     dnszone,
		SaveComment: "created by terraform",
		Masters:     masters,
	}
	var re RefResponse
	err := c.Post(postcreate, &re, "DNSZones")

	if err != nil {
		return objRef, err
	}

	return re.Result.Ref, err
}

func (c *Mmclient) DeleteDNSZone(ref string) error {
	err := c.Delete(deleteRequest("DNSZone"), "DNSZones/"+ref)

	if reqError, ok := err.(*RequestError); ok && reqError.StatusCode == ResourceNotFound {
		//DNS Zone not found, so nothing to delete
		return nil
	}
	return err
}

type UpdateDNSZoneRequest struct {
	Ref               string `json:"ref"`
	ObjType           string `json:"objType"`
	SaveComment       string `json:"saveComment"`
	DeleteUnspecified bool   `json:"deleteUnspecified"`

	// we can`t use DNSZoneProperties for this because CustomProperties should be flattend first
	Properties map[string]interface{} `json:"properties"`
}

func (c *Mmclient) UpdateDNSZone(dnsZoneProperties DNSZoneProperties, ref string) error {

	// A workaround to create properties with same fields as DNSZoneProperties but with flattend CustomProperties
	// first mask CustomProperties in DNSZoneProperties
	// Then convert to map considering `json:"omitempty"`
	// Then add CustomProperties 1 by 1

	customProperties := dnsZoneProperties.CustomProperties
	dnsZoneProperties.CustomProperties = nil

	properties, err := toMap(dnsZoneProperties)

	if err != nil {
		return err
	}
	for key, value := range customProperties {
		properties[key] = value
	}

	update := UpdateDNSZoneRequest{
		Ref:               ref,
		ObjType:           "DNSZone",
		SaveComment:       "updated by terraform",
		DeleteUnspecified: true,
		Properties:        properties,
	}

	return c.Put(update, "DNSZones/"+ref)
}

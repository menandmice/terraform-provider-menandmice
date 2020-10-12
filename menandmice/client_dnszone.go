package menandmice

type DNSZone struct {
	Ref          string   `json:"ref,omitempty"`
	AdIntegrated bool     `json:"adIntegrated"`
	DnsViewRef   string   `json:"dnsViewRef,omitempty"`
	DnsViewRefs  []string `json:"dnsViewRefs,omitempty"`
	Authority    string   `json:"authority,omitempty"`
	DNSZoneProperties
}

type DNSZoneProperties struct {
	Name         string `json:"name"`
	Dynamic      bool   `json:"dynamic,omitempty"`
	ZoneType     string `json:"type,omitempty"`
	DnssecSigned bool   `json:"dnssecSigned,omitempty"`
	KskIDs       string `json:"kskIDs,omitempty"`
	ZskIDs       string `json:"zskIDs,omitempty"`
	// TODO CustomProperties map[string]string `json:"customProperties,omitempty"`
	AdReplicationType string `json:"adReplicationType,omitempty"`
	AdPartition       string `json:"adPartition,omitempty"`
	Created           string `json:"created,omitempty"`
	LastModified      string `json:"lastModified,omitempty"`

	DisplayName string `json:"displyaName,omitempty"`
}

type FindDNSZoneResponse struct {
	Result struct {
		DNSZones     []DNSZone `json:"dnsZones"`
		TotalResults int       `json:"totalResults"`
	} `json:"result"`
}

func (c Mmclient) FindDNSZone(filter map[string]string) (error, []DNSZone) {
	var re FindDNSZoneResponse
	err := c.Get(&re, "dnszones/", filter)
	return err, re.Result.DNSZones
}

type ReadDNSZoneResponse struct {
	Result struct {
		DNSZone `json:"dnsZone"`
	} `json:"result"`
}

func (c Mmclient) ReadDNSZone(ref string) (error, DNSZone) {
	var re ReadDNSZoneResponse
	err := c.Get(&re, "dnszones/"+ref, nil)
	return err, re.Result.DNSZone
}

type CreateDNSZoneRequest struct {
	DNSZone     DNSZone  `json:"dnsZone"`
	SaveComment string   `json:"saveComment"`
	Masters     []string `json:"masters,omitempty"`
}

func (c *Mmclient) CreateDNSZone(dnszone DNSZone, masters []string) (error, string) {
	var objRef string
	postcreate := CreateDNSZoneRequest{
		DNSZone:     dnszone,
		SaveComment: "created by terraform",
		Masters:     masters,
	}
	var re RefResponse
	err := c.Post(postcreate, &re, "DNSZones")

	if err != nil {
		return err, objRef
	}

	return err, re.Result.Ref
}

// TODO this could be shared between all delete
type DeleteDNSZoneRequest struct {
	SaveComment  string `json:"saveComment"`
	ForceRemoval bool   `json:"forceRemoval"`
	// objType string

}

func (c *Mmclient) DeleteDNSZone(ref string) error {

	del := DeleteDNSZoneRequest{
		ForceRemoval: true,
		SaveComment:  "deleted by terraform",
	}
	return c.Delete(del, "DNSZones/"+ref)
}

type UpdateDNSZoneRequest struct {
	Ref string `json:"ref"`
	// objType Unknown
	SaveComment       string            `json:"saveComment"`
	DeleteUnspecified bool              `json:"deleteUnspecified"`
	Properties        DNSZoneProperties `json:"properties"`
}

func (c *Mmclient) UpdateDNSZone(dnsZoneProperties DNSZoneProperties, ref string) error {

	update := UpdateDNSZoneRequest{
		Ref:               ref,
		SaveComment:       "updated by terraform",
		DeleteUnspecified: true,
		Properties:        dnsZoneProperties,
	}
	return c.Put(update, "DNSZones/"+ref)
}

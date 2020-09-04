package menandmice

type DNSzone struct {
	Ref string `json:"ref,omitempty"`
	DNSZoneProperties
}

type DNSZoneProperties struct {
	Name         string   `json:"name"`
	Dynamic      bool     `json:"dynamic,omitempty"`
	AdIntegrated bool     `json:"adIntegrated"`
	DnsViewRef   string   `json:"dnsViewRef,omitempty"`
	DnsViewRefs  []string `json:"dnsViewRefs,omitempty"`
	Authority    string   `json:"authority,omitempty"`
	ZoneType     string   `json:"type,omitempty"`
	DnssecSigned bool     `json:"dnssecSigned,omitempty"`
	KskIDs       string   `json:"kskIDs,omitempty"`
	ZskIDs       string   `json:"zskIDs,omitempty"`
	// TODO CustomProperties map[string]string `json:"customProperties,omitempty"`
	AdReplicationType string `json:"adReplicationType,omitempty"`
	AdPartition       string `json:"adPartition,omitempty"`
	Created           string `json:"created,omitempty"`
	LastModified      string `json:"lastModified,omitempty"`

	DisplayName string `json:"displyaName,omitempty"`
}

type ReadDNSzoneResponse struct {
	Result struct {
		DNSzone `json:"dnsZone"`
	} `json:"result"`
}

func (c Mmclient) ReadDNSzone(ref string) (error, DNSzone) {
	var re ReadDNSzoneResponse
	//TODO fix ref
	err := c.Get(&re, "dnszones/"+ref)
	return err, re.Result.DNSzone
}

type CreateDNSzoneRequest struct {
	DNSzone     DNSzone  `json:"dnsZone"`
	SaveComment string   `json:"saveComment"`
	Master      []string `json:"master,omitempty"`
}

type CreateDNSzoneResponse struct {
	Result struct {
		Ref string `json:"ref"`
	} `json:"result"`
}

func (c *Mmclient) CreateDNSzone(dnszone DNSzone) (error, string) {
	var objRef string
	postcreate := CreateDNSzoneRequest{
		DNSzone:     dnszone,
		SaveComment: "created by terraform",
		// TODO Master : ,
	}
	var re CreateDNSzoneResponse
	err := c.Post(postcreate, &re, "DNSzones")

	if err != nil {
		return err, objRef
	}

	return err, re.Result.Ref
}

// TODO this could be shared between all delete
type DeleteDNSzoneRequest struct {
	SaveComment  string `json:"saveComment"`
	ForceRemoval bool   `json:"forceRemoval"`
	// objType string

}

func (c *Mmclient) DeleteDNSZone(ref string) error {

	del := DeleteDNSzoneRequest{
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

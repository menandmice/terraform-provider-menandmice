package menandmice

import "encoding/json"

type DNSZone struct {
	Ref          string   `json:"ref,omitempty"`
	AdIntegrated bool     `json:"adIntegrated"`
	DNSViewRef   string   `json:"dnsViewRef,omitempty"`
	DNSViewRefs  []string `json:"dnsViewRefs,omitempty"`
	Authority    string   `json:"authority,omitempty"`

	Created      string `json:"created,omitempty"`
	LastModified string `json:"lastModified,omitempty"`

	// properties that can be updated
	DNSZoneProperties
}

type DNSZoneProperties struct {
	Name              string            `json:"name"`
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

func (c Mmclient) FindDNSZone(filter map[string]string) ([]DNSZone, error) {
	var re FindDNSZoneResponse
	err := c.Get(&re, "dnszones/", filter)
	return re.Result.DNSZones, err
}

type ReadDNSZoneResponse struct {
	Result struct {
		DNSZone `json:"dnsZone"`
	} `json:"result"`
}

func (c Mmclient) ReadDNSZone(ref string) (DNSZone, error) {
	var re ReadDNSZoneResponse
	err := c.Get(&re, "dnszones/"+ref, nil)
	return re.Result.DNSZone, err
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
	Ref               string `json:"ref"`
	ObjType           string `json:"objType"`
	SaveComment       string `json:"saveComment"`
	DeleteUnspecified bool   `json:"deleteUnspecified"`

	// we cant use DNSZoneProperties for this because CustomProperties should be flattend first
	Properties map[string]interface{} `json:"properties"`
}

func (c *Mmclient) UpdateDNSZone(dnsZoneProperties DNSZoneProperties, ref string) error {

	// A work around to create properties with same fields as DNSZoneProperties but with flattend CustomProperties
	// first mask CustomProperties in DNSZoneProperties
	// Then marshal and Unmarshal in Map[string]interface
	// Then add CustomProperties 1 by 1

	customProperties := dnsZoneProperties.CustomProperties
	dnsZoneProperties.CustomProperties = nil

	var properties map[string]interface{}
	serialized, err := json.Marshal(dnsZoneProperties)

	if err != nil {
		return err
	}

	json.Unmarshal(serialized, &properties)

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

	dnsZoneProperties.CustomProperties = nil

	return c.Put(update, "DNSZones/"+ref)
}

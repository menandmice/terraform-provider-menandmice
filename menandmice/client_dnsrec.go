package menandmice

import (
	"errors"
)

type DNSRecord struct {
	Ref        string `json:"ref,omitempty"`
	DNSZoneRef string `json:"dnsZoneRef"`
	Rectype    string `json:"type"`
	DNSProperties
}

type DNSProperties struct {
	Name    string  `json:"name"`
	Rectype string  `json:"type"`
	Ttl     *string `json:"ttl,omitempty"`
	Data    string  `json:"data"`
	Comment string  `json:"comment,omitempty"`
	Aging   int     `json:"aging,omitempty"`
	Enabled bool    `json:"enabled,omitempty"`
}

type FindDNSRecResponse struct {
	Result struct {
		DNSRecords   []DNSRecord `json:"dnsRecords"`
		TotalResults int         `json:"totalResults"`
	} `json:"result"`
}

func (c Mmclient) FindDNSRec(zone string, filter map[string]string) ([]DNSRecord, error) {
	var re FindDNSRecResponse
	err := c.Get(&re, "DNSZones/"+zone+"/DNSRecords", nil, filter)
	return re.Result.DNSRecords, err
}

type ReadDNSRecResponse struct {
	Result struct {
		DNSRecord `json:"dnsRecord"`
	} `json:"result"`
}

func (c *Mmclient) ReadDNSRec(ref string) (DNSRecord, error) {
	var re ReadDNSRecResponse
	err := c.Get(&re, "dnsrecords/"+ref, nil, nil)
	return re.Result.DNSRecord, err
}

type CreateDNSRecResponse struct {
	Result struct {
		ObjRef []string `json:"objRefs"`
		Error  []string `json:"errors"`
	} `json:"result"`
}

type CreateDNSRecRequest struct {
	DNSRecords  []DNSRecord `json:"dnsRecords"`
	SaveComment string      `json:"saveComment"`
	// TODO autoAssignRangeRef string
	// TODO dnsZoneRef string
	ForceOverrideOfNamingConflictCheck bool `json:"forceOverrideOfNamingConflictCheck"`
}

func (c *Mmclient) CreateDNSRec(dnsrec DNSRecord) (string, error) {
	var objRef string
	postcreate := CreateDNSRecRequest{
		DNSRecords:                         []DNSRecord{dnsrec},
		SaveComment:                        "created by terraform",
		ForceOverrideOfNamingConflictCheck: false,
	}
	var re CreateDNSRecResponse
	err := c.Post(postcreate, &re, "DNSRecords")

	// TODO if dnsZoneRef does not exit you can confusing error "Missing object reference." give better messages

	if err != nil {
		return objRef, err
	}

	if len(re.Result.Error) > 0 {
		return objRef, errors.New(re.Result.Error[0])
	}

	if len(re.Result.ObjRef) != 1 {
		return objRef, errors.New("faild to create dns_record")
	}

	return re.Result.ObjRef[0], err
}

func (c *Mmclient) DeleteDNSRec(ref string) error {

	return c.Delete(deleteRequest("DNSRecord"), "DNSRecords/"+ref)
}

type UpdateDNSRecRequest struct {
	Ref string `json:"ref"`
	// objType Unknown
	SaveComment       string        `json:"saveComment"`
	DeleteUnspecified bool          `json:"deleteUnspecified"`
	Properties        DNSProperties `json:"properties"`
}

func (c *Mmclient) UpdateDNSRec(dnsProperties DNSProperties, ref string) error {

	update := UpdateDNSRecRequest{
		Ref:               ref,
		SaveComment:       "updated by terraform",
		DeleteUnspecified: true,
		Properties:        dnsProperties,
	}
	return c.Put(update, "DNSRecords/"+ref)
}

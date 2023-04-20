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
	Name    string `json:"name"`
	TTL     string `json:"ttl,omitempty"`
	Data    string `json:"data"`
	Comment string `json:"comment,omitempty"`
	Aging   int    `json:"aging,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}

type findDNSRecResponse struct {
	Result struct {
		DNSRecords   []DNSRecord `json:"dnsRecords"`
		TotalResults int         `json:"totalResults"`
	} `json:"result"`
}

func (c Mmclient) FindDNSRec(zone string, filter map[string]interface{}) ([]DNSRecord, error) {
	var re findDNSRecResponse

	query := map[string]interface{}{"filter": map2filter(filter)}

	err := c.Get(&re, "DNSZones/"+zone+"/DNSRecords", query)
	return re.Result.DNSRecords, err
}

type readDNSRecResponse struct {
	Result struct {
		DNSRecord `json:"dnsRecord"`
	} `json:"result"`
}

func (c *Mmclient) ReadDNSRec(ref string) (*DNSRecord, error) {
	var re readDNSRecResponse
	err := c.Get(&re, "dnsrecords/"+ref, nil)

	if reqError, ok := err.(*RequestError); ok && reqError.StatusCode == ResourceNotFound {
		return nil, nil
	}
	return &re.Result.DNSRecord, err
}

type createDNSRecResponse struct {
	Result struct {
		ObjRef []string `json:"objRefs"`
		Error  []string `json:"errors"`
	} `json:"result"`
}

type createDNSRecRequest struct {
	DNSRecords  []DNSRecord `json:"dnsRecords"`
	SaveComment string      `json:"saveComment"`
	// autoAssignRangeRef string
	// dnsZoneRef string
	ForceOverrideOfNamingConflictCheck bool `json:"forceOverrideOfNamingConflictCheck"`
}

func (c *Mmclient) CreateDNSRec(dnsrec DNSRecord) (string, error) {
	var objRef string
	postcreate := createDNSRecRequest{
		DNSRecords:                         []DNSRecord{dnsrec},
		SaveComment:                        "created by terraform",
		ForceOverrideOfNamingConflictCheck: false,
	}
	var re createDNSRecResponse
	err := c.Post(postcreate, &re, "DNSRecords")

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

	err := c.Delete(deleteRequest("DNSRecord"), "DNSRecords/"+ref)
	if reqError, ok := err.(*RequestError); ok && reqError.StatusCode == ResourceNotFound {
		return nil
	}
	return err
}

type updateDNSRecRequest struct {
	Ref string `json:"ref"`
	// objType Unknown
	SaveComment       string        `json:"saveComment"`
	DeleteUnspecified bool          `json:"deleteUnspecified"`
	Properties        DNSProperties `json:"properties"`
}

func (c *Mmclient) UpdateDNSRec(dnsProperties DNSProperties, ref string) error {

	update := updateDNSRecRequest{
		Ref:               ref,
		SaveComment:       "updated by terraform",
		DeleteUnspecified: true,
		Properties:        dnsProperties,
	}
	return c.Put(update, "DNSRecords/"+ref)
}

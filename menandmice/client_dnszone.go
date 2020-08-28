package menandmice

type DNSzone struct {
	Ref          *string `json:"ref,omitempty"`
	Name         string  `json:"name"`
	Dynamic      bool    `json:"dynamic,omitempty"`
	AdIntegrated bool    `json:"AdIntergrated,omitempty"`
	DnsViewRef   string  `json:"dnsViewRef":omitempty`
	Authority    string  `json:"authority,omitempty"`
	ZoneType     string  `json:"type,omitempty"`
	DnssecSigned bool    `json:"dnssecSigned,omitempty"`
	KskIDs       string  `json:"kskIDs,omitempty"`
	ZskIDs       string  `json:"zskIDs,omitempty"`
	// TODO CustomProperties map[string]string `json:"customProperties,omitempty"`
	// TODO adReplicationType
	// TOOD adPartition

	Created       string `json:"created,omitempty"`
	LastModiefied string `json:"lastModiefied,omitempty"`
	DisplayName   string `json:"displyaName,omitempty"`
}

type ReadDNSzoneResponse struct {
	Result struct {
		DNSzone `json:"dnsZone"`
	} `json:"result"`
}

func (c Mmclient) ReadDNSzone(ref string) (error, DNSzone) {
	var re ReadDNSzoneResponse
	err := c.Get(&re, "dnszones/"+ref)
	return err, re.Result.DNSzone
}

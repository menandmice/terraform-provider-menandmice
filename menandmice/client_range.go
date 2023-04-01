package menandmice

import "errors"

type Range struct {
	Ref               string     `json:"ref,omitempty"`
	Name              string     `json:"name"`
	From              string     `json:"from,omitempty"`
	To                string     `json:"to,omitempty"`
	ParentRef         string     `json:"parentRef,omitempty"`
	AdSiteRef         string     `json:"adSiteRef,omitempty"`
	AdSiteDisplayName string     `json:"adSiteDisplayName,omitempty"`
	ChildRanges       []NamedRef `json:"childRanges"`
	// seems redundant
	// IsLeaf            bool       `json:"isLeaf"`
	// NumChildren int        `json:"numchildren"`
	DhcpScopes  []NamedRef `json:"dhcpScopes"`
	Authority   *Authority `json:"authority,omitempty"`
	HasSchedule bool       `json:"hasSchedule"`
	HasMonitor  bool       `json:"hasMonitor"`

	UtilizationPercentage int                   `json:"utilizationPercentage,omitempty"`
	HasRogueAddresses     bool                  `json:"hasRogueAddresses,omitempty"`
	CloudNetworkRef       string                `json:"cloudNetworkRef,omitempty"`
	CloudAllocationPools  []CloudAllocationPool `json:"cloudAllocationPools,omitempty"`

	InheritAccess        bool                  `json:"inheritAccess"`
	DiscoveredProperties []DiscoveryProperties `json:"discoveredProperties,omitempty"`
	Created              string                `json:"created,omitempty"`
	LastModified         string                `json:"lastModified,omitempty"`
	FolderRef            string                `json:"folderRef,omitempty"`
	RangeProperties
}

type RangeProperties struct {
	IsContainer bool `json:"isContainer"`
	Locked      bool `json:"locked"`
	AutoAssign  bool `json:"autoAssign"`
	Subnet      bool `json:"subnet"`
	// TODO should be CustomProperties map[string]interface{} `json:"customProperties"`
	CustomProperties map[string]string `json:"customProperties,omitempty"`
}
type Authority struct {
	Name    string   `json:"name,omitempty"`
	Type    string   `json:"type,omitempty"`
	SubType string   `json:"subType,omitempty"`
	Sources []Source `json:"source,omitempty"`
}

type Source struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Ref     string `json:"ref"`
	Enabled bool   `json:"enabled"`
}
type CloudAllocationPool struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Discovery struct {
	Interval int    `json:"interval,omitempty"`
	Unit     string `json:"unit,omitempty"` // TODO make enum Minutes , Hours , Days , Weeks , Months
	Enabled  bool   `json:"enabled"`
	// StartTime string `json:"startTime,omitempty"` // TODO better time format
}

type NamedRef struct {
	Ref     string `json:"ref"`
	ObjType string `json:"objType"`
	Name    string `json:"name"`
}

type DiscoveryProperties struct {
	RouterName           string `json:"routerName"`
	Gateway              string `json:"gateway"`
	InterfaceID          int    `json:"interfaceID"`
	InterfaceName        string `json:"interfaceName"`
	VLANID               int
	InterfaceDescription string `json:"interfaceDescription"`
	VRFName              string
}

type findRangesResponse struct {
	Result struct {
		Ranges       []Range `json:"ranges"`
		TotalResults int     `json:"totalResults"`
	} `json:"result"`
}

func (c Mmclient) FindRanges(limit int, filter map[string]interface{}) ([]Range, error) {
	var re findRangesResponse
	query := map[string]interface{}{
		"limit": limit,
	}

	if folderRef, ok := filter["folderRef"]; ok {
		query["folderRef"] = folderRef
		delete(filter, "folderRef")
	}

	if rawFilter, ok := filter["filter"]; ok {
		query["filter"] = rawFilter.(string)
	} else {
		query["filter"] = map2filter(filter)
	}
	err := c.Get(&re, "ranges/", query)
	if reqError, ok := err.(*RequestError); ok && reqError.StatusCode == ObjectNotFoundForReference {
		// folder, view, or server where not found, is not error but return empty list
		return []Range{}, nil
	}
	return re.Result.Ranges, err
}

type readRangeResponse struct {
	Result struct {
		Range `json:"range"`
	} `json:"result"`
}

func (c Mmclient) ReadRange(ref string) (*Range, error) {
	var re readRangeResponse
	err := c.Get(&re, "Ranges/"+ref, nil)
	if reqError, ok := err.(*RequestError); ok && reqError.StatusCode == ResourceNotFound {
		return nil, nil
	}

	return &re.Result.Range, err
}

type createRangeRequest struct {
	Range       Range     `json:"range"`
	SaveComment string    `json:"saveComment"`
	Discovery   Discovery `json:"discovery"`
}

func (c *Mmclient) CreateRange(iprange Range, discovery Discovery) (string, error) {
	var objRef string
	postcreate := createRangeRequest{
		Range:       iprange,
		SaveComment: "created by terraform",
		Discovery:   discovery,
	}
	var re RefResponse
	err := c.Post(postcreate, &re, "Ranges")

	// TODO better error messages required custom property
	if err != nil {
		return objRef, err
	}

	return re.Result.Ref, err
}

func (c *Mmclient) DeleteRange(ref string) error {

	err := c.Delete(deleteRequest("Range"), "Ranges/"+ref)
	if reqError, ok := err.(*RequestError); ok && reqError.StatusCode == ResourceNotFound {
		return nil
	}
	return err
}

type updateRangeRequest struct {
	Ref               string `json:"ref"`
	ObjType           string `json:"objType"`
	SaveComment       string `json:"saveComment"`
	DeleteUnspecified bool   `json:"deleteUnspecified"`

	Properties map[string]interface{} `json:"properties"`
}

// TODO Add test for this. there was bug in this, but it was never detected.
// Are there properties can be update in place
func (c *Mmclient) UpdateRange(rangeProperties RangeProperties, ref string) error {

	// A workaround to create properties with same fields as DNSZoneProperties but with flattend CustomProperties
	// first mask CustomProperties in DNSZoneProperties
	// Then convert to map considering `json:"omitempty"`
	// Then add CustomProperties 1 by 1

	customProperties := rangeProperties.CustomProperties
	rangeProperties.CustomProperties = nil

	properties, err := toMap(rangeProperties)

	if err != nil {
		return err
	}
	for key, value := range customProperties {
		properties[key] = value
	}

	update := updateRangeRequest{
		Ref:     ref,
		ObjType: "Range",
		// TODO  reuse same constant everywhere for comment
		SaveComment:       "updated by terraform",
		DeleteUnspecified: true, // TODO false
		Properties:        properties,
	}

	return c.Put(update, "Ranges/"+ref)
}

type nextFreeAddressRespons struct {
	Result struct {
		Address string `json:"address"`
	} `json:"result"`
}

type NextFreeAddressRequest struct {
	RangeRef           string
	StartAddress       string
	Ping               bool
	ExcludeDHCP        bool
	TemporaryClaimTime int
}

func (c Mmclient) NextFreeAddress(request NextFreeAddressRequest) (string, error) {
	var re nextFreeAddressRespons
	query := map[string]interface{}{
		"ping":               request.Ping,
		"excludeDHCP":        request.ExcludeDHCP,
		"temporaryClaimTime": request.TemporaryClaimTime,
	}
	if request.StartAddress != "" {
		query["startAddress"] = request.StartAddress
	}
	err := c.Get(&re, "Ranges/"+request.RangeRef+"/NextFreeAddress", query)
	return re.Result.Address, err
}

type availableAddressBlocksRespons struct {
	Result struct {
		AddressBlocks []AddressBlock `json:"addressBlocks"`
	} `json:"result"`
}

type AddressBlock struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type AvailableAddressBlocksRequest struct {
	RangeRef           string
	StartAddress       string
	Size               int
	Mask               int
	Limit              int
	IgnoreSubnetFlag   bool
	TemporaryClaimTime int
}

func (c Mmclient) AvailableAddressBlocks(request AvailableAddressBlocksRequest) ([]AddressBlock, error) {

	var re availableAddressBlocksRespons
	query := map[string]interface{}{
		"limit":              request.Limit,
		"ignoreSubnetFlag":   request.IgnoreSubnetFlag,
		"size":               request.Size,
		"temporaryClaimTime": request.TemporaryClaimTime,
	}
	if request.Mask != 0 {
		query["mask"] = request.Mask
	}

	if request.Size != 0 {
		query["size"] = request.Size
	}

	if request.StartAddress != "" {
		query["startAddress"] = request.StartAddress
	}
	if request.RangeRef == "" {
		return nil, errors.New("Range not specified")
	}
	err := c.Get(&re, "Ranges/"+request.RangeRef+"/AvailableAddressBlocks", query)
	return re.Result.AddressBlocks, err

}

package menandmice

type Range struct {
	Ref  string `json:"ref,omitempty"`
	Name string `json:"name"`
	From string `json:"from"`
	To   string `json:"to"`
	// ParentRef string `json:"parentRef"`
	// AdSiteRef string         `json: "adSiteRef"`
	// AdSiteDisplayName string `json: "adSiteDisplayName"`
	// ChildRanges: [] ChildRanges `json: "childRanges"`
	// IsLeaf bool `json:isLeaf`
	// NumChildren int
	// dhcpScopes
	// authority
	// subnet:  bool
	Locked           bool                   `json:"locked"`
	AutoAssign       bool                   `json:"autoAssign"`
	HasSchedule      bool                   `json:"hasSchedule"`
	HasMonitor       bool                   `json:"hasMonitor"`
	CustomProperties map[string]interface{} `json:"customProperties"`
	// inheritAcces bool
	//  isContainer bool
	//       "utilizationPercentage": 0,                        "utilizationPercentage": 0,
	//       "hasRogueAddresses": false,                        "hasRogueAddresses": false,
	//       "cloudNetworkRef": "string",                       "cloudNetworkRef": "string",
	//       "cloudAllocationPools": [                          "cloudAllocationPools": [

	// "discoveredProperties"
	//       "created": "2022-09-02T12:37:53.585Z",             "created": "2022-09-02T12:37:53.585Z",
	//       "lastModified": "2022-09-02T12:37:53.585Z",        "lastModified": "2022-09-02T12:37:53.585Z",
	//       "folderRef": "string"                              "folderRef": "string"
	RangeProperties
}

type RangeProperties struct {
}

type Discovery struct {
	Interval int    `json:"interval,omitempty"`
	Unit     string `json:"unit,omitempty"` // TODO make enum Minutes , Hours , Days , Weeks , Months
	Enabled  bool   `json:"Enabled,omitempty"`
	// StartTime string `json:"startTime,omitempty"` // TODO better time format
}

//       "discoveredProperties": [                          "discoveredProperties": [
//         {                                                  {
//           "routerName": "string",                            "routerName": "string",
//           "gateway": "string",                               "gateway": "string",
//           "interfaceID": 0,                                  "interfaceID": 0,
//           "interfaceName": "string",                         "interfaceName": "string",
//           "VLANID": 0,                                       "VLANID": 0,
//           "interfaceDescription": "string",                  "interfaceDescription": "string",
//           "VRFName": "string"                                "VRFName": "string"
//         }                                                  }
//       ],                                                 ],

//       "authority": {                                     "authority": {
//         "name": "string",                                  "name": "string",
//         "type": "string",                                  "type": "string",
//         "subtype": "Scope",                                "subtype": "Scope",
//         "sources": [                                       "sources": [
//           {                                                  {
//             "name": "string",                                  "name": "string",
//             "type": "string",                                  "type": "string",
//             "ref": "string",                                   "ref": "string",
//             "enabled": false                                   "enabled": false
//           }                                                  }
//         ]                                                  ]
//       },                                                 },

//       "dhcpScopes": [                                    "dhcpScopes": [
//         {                                                  {
//           "ref": "string",                                   "ref": "string",
//           "objType": "Unknown",                              "objType": "Unknown",
//           "name": "string"                                   "name": "string"
//         }                                                  }
//       ],                                                 ],

type ReadRangeResponse struct {
	Result struct {
		Range `json:"range"`
	} `json:"result"`
}

func (c Mmclient) ReadRange(ref string) (*Range, error) {
	var re ReadRangeResponse
	err := c.Get(&re, "Ranges/"+ref, nil, nil)
	if reqError, ok := err.(*RequestError); ok && reqError.StatusCode == ResourceNotFound {
		return nil, nil
	}

	return &re.Result.Range, err
}

type CreateRangeRequest struct {
	DNSZone     Range     `json:"range"`
	SaveComment string    `json:"saveComment"`
	Discovery   Discovery `json:"discovery,omitempty"`
}

func (c *Mmclient) CreateRangeZone(iprange Range, discovery Discovery) (string, error) {
	var objRef string
	postcreate := CreateRangeRequest{
		DNSZone:     iprange,
		SaveComment: "created by terraform",
		Discovery:   discovery,
	}
	var re RefResponse
	err := c.Post(postcreate, &re, "Ranges")

	if err != nil {
		return objRef, err
	}

	return re.Result.Ref, err
}

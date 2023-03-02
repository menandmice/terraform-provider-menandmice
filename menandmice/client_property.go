package menandmice

type PropertieDefinietion struct {
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	System       bool     `json:"system"`
	ReadOnly     bool     `json:"readOnly"`
	MultiLine    bool     `json:"multiLine"`
	DefaultValue string   `json:"defaultValue"`
	ListItems    []string `json:"listItems"`
	//     "parentProperty": "string",
	//     "cloudTags": [
	//       {
	//         "name": "string",
	//         "cloudRef": "string"
	//       }
	//     ]
	//   },
}

type readPropertyResponse struct {
	Result struct {
		PropertyDefinitions []PropertieDefinietion `json:"propertyDefinitions"`
	} `json:"result"`
}

// ReadProperty get map op properties from resouce
func (c *Mmclient) ReadProperty(resource, ref string) (map[string]PropertieDefinietion, error) {

	var propertiesDefinitions = map[string]PropertieDefinietion{}
	var re readPropertyResponse
	err := c.Get(&re, resource+"/"+ref+"/PropertyDefinitions", nil)
	if err != nil {
		return propertiesDefinitions, err
	}
	for _, property := range re.Result.PropertyDefinitions {
		propertiesDefinitions[property.Name] = property
	}
	return propertiesDefinitions, nil
}

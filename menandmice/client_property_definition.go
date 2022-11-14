package menandmice

import (
	"errors"
	"strings"
)

type PropertyDefinition struct {
	// ref            string
	Name           string     `json:"Name"`
	Type           string     `json:"type"`
	System         bool       `json:"system,omitempty"`
	Mandatory      bool       `json:"mandatory,omitempty"`
	ReadOnly       bool       `json:"readOnly,omitempty"`
	MultiLine      bool       `json:"multiLine,omitempty"`
	DefaultValue   string     `json:"defaultValue,omitempty"`
	ListItems      []string   `json:"listItems,omitempty"`
	ParentProperty string     `json:"parentProperty,omitempty"`
	CloudTags      []CloudTag `json:"cloudTags,omitempty"`
}

type CloudTag struct {
	Name     string `json:"Name"`
	CloudRef string `json:"CloudRef"`
}

type PropertyDefinitionResponse struct {
	Result struct {
		PropertyDefinition []PropertyDefinition `json:"propertyDefinitioni`
	} `json:"result"`
}

var RESOURCES_WITH_CUSTOMPROPERTIES = []string{"dnszones", "ranges"}

// go lang does not have a contain/ in in std libery
func arryContainsString(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

//TODO whould enum type for resource work better here?
func (c *Mmclient) FindPropertyDefinitions(resource string) ([]PropertyDefinition, error) {
	var re PropertyDefinitionResponse
	resource = strings.ToLower(resource)
	if !arryContainsString(RESOURCES_WITH_CUSTOMPROPERTIES, resource) {
		return nil, errors.New("can't get resource definition for this resource")
	}
	path := resource + "/dummy/PropertyDefinitions"
	err := c.Get(&re, path, nil, nil)
	if reqError, ok := err.(*RequestError); ok && reqError.StatusCode == ResourceNotFound {
		//DHCPReservationNotFound not found
		return nil, nil
	}
	return re.Result.PropertyDefinition, err
}

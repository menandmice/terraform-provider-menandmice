package menandmice

type Group struct {
	Ref          string   `json:"ref,omitempty"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	AdIntegrated bool     `json:"adIntegrated"`
	BuiltIn      bool     `json:"builtin"`
	GroupMembers []Member `json:"groupMembers"`
	Roles        []Member `json:"roles"`
}

type Member struct {
	Ref     string `json:"ref"`
	ObjType string `json:"objType"`
	Name    string `json:"name"`
}

type ReadGroupResponse struct {
	Result struct {
		Group `json:"group"`
	} `json:"result"`
}

func (c *Mmclient) readGroup(ref string) (Group, error) {
	var re ReadGroupResponse
	err := c.Get(&re, "Groups/"+ref, nil, nil)
	return re.Result.Group, err

}

type CreateGroupRequest struct {
	Group       Group  ` json:"group"`
	SaveComment string `json:"saveComment"`
}

func (c *Mmclient) CreatGroup(group Group) (string, error) {
	var objRef string
	postcreate := CreateGroupRequest{
		Group:       group,
		SaveComment: "created by terraform",
	}
	var re RefResponse
	err := c.Post(postcreate, &re, "Groups/")

	if err != nil {
		return objRef, err
	}

	return re.Result.Ref, err
}

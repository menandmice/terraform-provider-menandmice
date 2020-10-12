package menandmice

type Group struct {
	Ref          string   `json:"ref,omitempty"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	AdIntegrated bool     `json:"adIntegrated"`
	BuiltIn      bool     `json:builtin`
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

func (c *Mmclient) readGroup(ref string) (error, Group) {
	var re ReadGroupResponse
	err := c.Get(&re, "Groups/"+ref, nil)
	return err, re.Result.Group

}

type CreateGroupRequest struct {
	Group       Group  ` json:"group"`
	SaveComment string `json:"saveComment"`
}

func (c *Mmclient) CreatGroup(group Group) (error, string) {
	var objRef string
	postcreate := CreateGroupRequest{
		Group:       group,
		SaveComment: "created by terraform",
	}
	var re RefResponse
	err := c.Post(postcreate, &re, "Groups/")

	if err != nil {
		return err, objRef
	}

	return err, re.Result.Ref
}

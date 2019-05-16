package auth

const PlanGrantType = "plan"

type Grant struct {
	UserName string
	Access   []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"access"`
}

func (grant *Grant) IsAuthorizedForPlan(name string) bool {
	return grant.IsAuthorized(PlanGrantType, name)
}

func (grant *Grant) IsAuthorized(accessType string, name string) bool {
	for _, access := range grant.Access {
		if access.Type == accessType {
			if access.Name == name {
				return true
			}
		}
	}

	return false
}

package auth

type Grant struct {
	UserName string
	Access []struct{
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"access"`
}

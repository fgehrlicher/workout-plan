package auth

type Grant struct {
	UserName string `json:"sub"`
	Access []struct{
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"access"`
}

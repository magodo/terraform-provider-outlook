package client

type User struct {
	ID            string `json:"id"`
	DisplayName   string `json:"displayName"`
	PrincipalName string `json:"userPrincipalName"`
}

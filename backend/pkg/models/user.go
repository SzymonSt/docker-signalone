package models

type User struct {
	UserId           string `json:"userId"`
	UserName         string `json:"userName"`
	IsPro            bool   `json:"isPro"`
	AgentBearerToken string `json:"agentBearerToken"`
}

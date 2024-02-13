package models

type User struct {
	UserId           string `json:"userId" bson:"userId"`
	UserName         string `json:"userName" bson:"userName"`
	IsPro            bool   `json:"isPro" bson:"isPro"`
	AgentBearerToken string `json:"agentBearerToken" bson:"agentBearerToken"`
	Counter          int32  `json:"counter" bson:"counter"`
}

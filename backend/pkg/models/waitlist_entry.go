package models

type WaitlistEntry struct {
	Email       string `json:"email" bson:"email"`
	CompanyName string `json:"companyName" bson:"companyName"`
}

package models

import "github.com/golang-jwt/jwt/v5"

type User struct {
	UserId           string `json:"userId" bson:"userId"`
	UserName         string `json:"userName" bson:"userName"`
	IsPro            bool   `json:"isPro" bson:"isPro"`
	AgentBearerToken string `json:"agentBearerToken" bson:"agentBearerToken"`
	Counter          int32  `json:"counter" bson:"counter"`
	Type             string `json:"type" bson:"type"`
}

type GoogleClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	jwt.RegisteredClaims
}

type GoogleTokenRequest struct {
	IdToken string `json:"idToken"`
}

type JWTClaimsWithId struct {
	Id string `json:"id"`
	jwt.RegisteredClaims
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

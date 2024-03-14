package models

import "github.com/golang-jwt/jwt/v5"

type User struct {
	UserId           string `json:"userId" bson:"userId"`
	UserName         string `json:"userName" bson:"userName"`
	IsPro            bool   `json:"isPro" bson:"isPro"`
	AgentBearerToken string `json:"agentBearerToken" bson:"agentBearerToken"`
	Counter          int32  `json:"counter" bson:"counter"`
	Type             string `json:"type" bson:"type"`

	//If user type is signalone
	PasswordHash          string `json:"passwordHash" bson:"passwordHash"`
	EmailConfirmed        bool   `json:"emailConfirmed" bson:"emailConfirmed"`
	EmailConfirmationCode string `json:"emailConfirmationCode" bson:"emailConfirmationCode"`
}

type SignalAccountPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type EmailConfirmationPayload struct {
	Email             string `json:"email"`
	ConfirmationToken string `json:"confirmationToken"`
}

type GithubUserData struct {
	AvatarUrl         string `json:"avatar_url"`
	Bio               string `json:"bio"`
	Blog              string `json:"blog"`
	Company           string `json:"company"`
	CreatedAt         string `json:"created_at"`
	Email             string `json:"email"`
	EventsUrl         string `json:"events_url"`
	Followers         int    `json:"followers"`
	FollowersUrl      string `json:"followers_url"`
	Following         int    `json:"following"`
	FollowingUrl      string `json:"following_url"`
	GistsUrl          string `json:"gists_url"`
	Gravatar_id       string `json:"gravatar_id"`
	Hireable          bool   `json:"hireable"`
	HtmlUrl           string `json:"html_url"`
	Id                int    `json:"id"`
	Location          string `json:"location"`
	Login             string `json:"login"`
	Name              string `json:"name"`
	NodeId            string `json:"node_id"`
	OrganizationsUrl  string `json:"organizations_url"`
	PublicGists       int    `json:"public_gists"`
	PublicRepos       int    `json:"public_repos"`
	ReceivedEventsUrl string `json:"received_events_url"`
	ReposUrl          string `json:"repos_url"`
	SiteAdmin         bool   `json:"site_admin"`
	StarredUrl        string `json:"starred_url"`
	SubscriptionsUrl  string `json:"subscriptions_url"`
	TwitterUsername   string `json:"twitter_username"`
	Type              string `json:"type"`
	Url               string `json:"url"`
}

type GithubTokenRequest struct {
	Code string `json:"code"`
}

type GithubTokenResponse struct {
	AccessToken string `json:"access_token"`
	Token_type  string `json:"token_type"`
	Scope       string `json:"scope"`
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

type JWTClaimsWithUserData struct {
	Id       string `json:"id"`
	UserName string `json:"userName"`
	jwt.RegisteredClaims
}

type AgentTokenClaimsWithUserData struct {
	Id       string `json:"id"`
	UserName string `json:"userName"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type Email struct {
	Email          string `json:"email" binding:"required"`
	MessageContent string `json:"messageContent" binding:"required"`
	MessageTitle   string `json:"messageTitle" binding:"required"`
}

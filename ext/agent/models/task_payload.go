package models

type TaskPayload struct {
	BearerToken string
	BackendUrl  string
	UserId      string
}

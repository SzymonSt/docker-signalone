package models

import "time"

type ContainerData struct {
	LastLog chan time.Time `json:"last_log"`
}

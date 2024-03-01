package models

import "time"

type ContainerData struct {
	LastLog time.Time `json:"last_log"`
}

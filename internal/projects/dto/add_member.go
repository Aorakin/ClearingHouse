package dto

import "github.com/google/uuid"

type Members struct {
	Members []uuid.UUID `json:"members"`
}

package dtos

import "github.com/google/uuid"

type Usage struct {
	TypeID uuid.UUID `json:"type_id"`
	Type   string    `json:"type"`
	Quota  float64   `json:"quota"`
	Usage  float64   `json:"usage"`
}

type UsageResponse struct {
	Usage []Usage `json:"usage"`
}

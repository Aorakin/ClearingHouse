package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateTicketRequest struct {
	PoolName  string `json:"pool_name" validate:"required" `
	GPU       string `json:"gpu" `
	Ram       string `json:"ram"`
	VRam      string `json:"vram"`
	CPU       string `json:"cpu"`
	Storage   string `json:"storage"`
	UsageTime string `json:"usage_time" validate:"required"`
}

type Ticket struct {
	ID                uuid.UUID    `json:"id" validate:"required"`
	NamespaceURN      string       `json:"namespace_urn" validate:"required"`
	GlideletURN       string       `json:"glidelet_urn" validate:"required"`
	Spec              []GliderSpec `json:"spec" validate:"required"`
	ReferenceTicketID string       `json:"reference_ticket_id"`
	RedeemTimeout     time.Time    `json:"redeem_timeout" validate:"required"`
	Lease             time.Time    `json:"lease" validate:"required"`
	Signature         []byte       `json:"signature" validate:"required"`
}

type GliderSpec struct {
	ID        uuid.UUID      `json:"id" validate:"required"`
	Type      string         `json:"type" validate:"required"`
	PoolID    uuid.UUID      `json:"pool_id" validate:"required"`
	Resources []SpecResource `json:"resource" validate:"required"`
}

type SpecResource struct {
	Name     string `json:"name" validate:"required"`
	Quantity string `json:"quantity" validate:"required"`
	Unit     string `json:"unit" validate:"required"`
}

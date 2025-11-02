package dtos

import (
	"time"

	"github.com/google/uuid"
)

type GliderTicketResponse struct {
	Ticket    GliderTicket `json:"ticket"`
	Signature string       `json:"signature"`
}

type GliderTicket struct {
	ID                uuid.UUID  `json:"id"`
	NamespaceURN      string     `json:"namespace_urn"` // namespace.id
	NamespaceName     string     `json:"namespace_name"`
	ProjectURN        string     `json:"project_urn"` // project.id
	ProjectName       string     `json:"project_name"`
	GlideletURN       string     `json:"glidelet_urn"` // resource_pool.URN
	GlideletName      string     `json:"glidelet_name"`
	OrganizationName  string     `json:"organization_name"`
	Spec              GliderSpec `json:"spec"`
	ReferenceTicketID string     `json:"reference_ticket_id"`
	RedeemTimeout     uint       `json:"redeem_timeout"` // in seconds
	Lease             uint       `json:"lease"`          // in seconds
	CreatedAt         time.Time  `json:"created_at"`
}

type GliderSpec struct {
	Type      ResourceUnitType `json:"type"`
	PoolID    string           `json:"pool_id"`
	Resources []SpecResource   `json:"resource"`
}

type SpecResource struct {
	ResourceID string `json:"resource_id"`
	Name       string `json:"name"`
	Quantity   uint   `json:"quantity"`
	Unit       string `json:"unit"`
}

type StatusTicket string

const (
	StatusReady    StatusTicket = "ready"
	StatusRedeemed StatusTicket = "redeemed"
	StatusCanceled StatusTicket = "canceled"
	StatusExpired  StatusTicket = "expired"
)

type ResourceUnitType string

const (
	ResourceUnitTypeCPU ResourceUnitType = "compute"
)

type TicketTransaction struct {
	ID        uuid.UUID    `json:"id"`
	Status    StatusTicket `json:"status"`
	Ticket    string       `json:"ticket"`
	Signature string       `json:"signature"`
}

type Task struct {
	ID      uuid.UUID           `json:"id"`
	Tickets []TicketTransaction `json:"tickets"`
}

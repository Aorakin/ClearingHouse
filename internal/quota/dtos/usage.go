package dtos

type Usage struct {
	TypeID string  `json:"type_id"`
	Type   string  `json:"type"`
	Quota  float64 `json:"quota"`
	Usage  float64 `json:"usage"`
}

type UsageResponse struct {
	Usage []Usage `json:"usage"`
}

package dtos

type IDUri struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type ResourcePoolResponse struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	Resources []ResourceResponse `json:"resources"`
}

type ResourceResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ResourceType string `json:"resource_type"`
	Unit         string `json:"unit"`
	Quantity     uint   `json:"quantity"`
}



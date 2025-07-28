package dtos

import "github.com/ClearingHouse/internal/models"

func NewResourceResponse(r models.Resource) ResourceResponse {
	return ResourceResponse{
		ID:           r.ID.String(),
		Name:         r.Name,
		ResourceType: r.ResourceType.Name,
		Unit:         r.ResourceType.Unit,
		Quantity:     r.Quantity,
	}
}

func NewResourcePoolResponse(pool models.ResourcePool) ResourcePoolResponse {
	response := ResourcePoolResponse{
		ID:        pool.ID.String(),
		Name:      pool.Name,
		Resources: []ResourceResponse{},
	}
	for _, res := range pool.Resources {
		response.Resources = append(response.Resources, NewResourceResponse(res))
	}
	return response
}

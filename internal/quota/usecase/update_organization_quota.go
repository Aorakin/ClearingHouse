package usecase

import (
	"fmt"

	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/quota/dtos"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/google/uuid"
)

func (u *QuotaUsecase) UpdateOrganizationQuota(quotaID uuid.UUID, request *dtos.UpdateOrganizationQuotaRequest, userID uuid.UUID) (*models.OrganizationQuota, error) {
	orgQuota, err := u.quotaRepo.GetOrgQuotaByID(quotaID)
	if err != nil {
		return nil, apiError.NewNotFoundError(fmt.Errorf("organization quota not found: %w", err))
	}

	if err := u.isOrgAdmin(orgQuota.FromOrgID, userID); err != nil {
		return nil, err
	}

	// Update name/description
	if err := u.quotaRepo.UpdateOrganizationQuota(quotaID, request.Name, request.Description); err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to update organization quota: %w", err))
	}

	// Update resources if provided
	if len(request.Resources) > 0 {
		// Build map of existing org quota resources by resource_id
		existingQuantities, err := u.quotaRepo.GetResourceQuantitiesByOrgQuotaID(quotaID)
		if err != nil {
			return nil, apiError.NewInternalServerError(fmt.Errorf("failed to get existing resource quantities: %w", err))
		}

		existingMap := make(map[uuid.UUID]models.ResourceQuantity)
		for _, qty := range existingQuantities {
			existingMap[qty.ResourceProp.ResourceID] = qty
		}

		// Validate resources exist in the node
		resources, err := u.resourceRepo.GetResourcesByNodeID(orgQuota.NodeID)
		if err != nil {
			return nil, apiError.NewInternalServerError(fmt.Errorf("failed to get resources by node ID: %w", err))
		}
		nodeResourceMap := make(map[uuid.UUID]models.Resource)
		for _, r := range resources {
			nodeResourceMap[r.ID] = r
		}

		seenResources := make(map[uuid.UUID]struct{})
		for _, r := range request.Resources {
			if _, duplicate := seenResources[r.ResourceID]; duplicate {
				return nil, apiError.NewBadRequestError(fmt.Errorf("duplicate resource ID: %s", r.ResourceID))
			}
			seenResources[r.ResourceID] = struct{}{}

			nodeResource, exists := nodeResourceMap[r.ResourceID]
			if !exists {
				return nil, apiError.NewNotFoundError(fmt.Errorf("resource %s not found in the specified node", r.ResourceID))
			}

			if r.Quantity > nodeResource.Quantity {
				return nil, apiError.NewBadRequestError(fmt.Errorf("requested quantity %d exceeds available quantity %d for resource %s", r.Quantity, nodeResource.Quantity, r.ResourceID))
			}

			if existingQty, found := existingMap[r.ResourceID]; found {
				// Update existing quantity
				if err := u.quotaRepo.UpdateResourceQuantity(existingQty.ID, r.Quantity); err != nil {
					return nil, apiError.NewInternalServerError(fmt.Errorf("failed to update resource quantity: %w", err))
				}
				// Update resource property (price/duration)
				if err := u.quotaRepo.UpdateResourceProperty(existingQty.ResourcePropID, r.Price, r.Duration); err != nil {
					return nil, apiError.NewInternalServerError(fmt.Errorf("failed to update resource property: %w", err))
				}

				// Cascade down: clamp child project quotas
				if err := u.cascadeOrgQuotaToProjectQuotas(quotaID, r.ResourceID, r.Quantity); err != nil {
					return nil, err
				}
			} else {
				// Create new resource property + quantity
				resourceProperty := models.ResourceProperty{
					ResourceID:  r.ResourceID,
					Price:       r.Price,
					MaxDuration: r.Duration,
				}
				if err := u.quotaRepo.CreateResourceProperty(&resourceProperty); err != nil {
					return nil, apiError.NewInternalServerError(fmt.Errorf("failed to create resource property: %w", err))
				}
				resourceQuantity := models.ResourceQuantity{
					OrganizationQuotaID: &quotaID,
					Quantity:            r.Quantity,
					ResourcePropID:      resourceProperty.ID,
				}
				if err := u.quotaRepo.CreateResourceQuantity(&resourceQuantity); err != nil {
					return nil, apiError.NewInternalServerError(fmt.Errorf("failed to create resource quantity: %w", err))
				}
			}
		}

		// Zero out resources no longer in the request
		for resourceID, existingQty := range existingMap {
			if _, found := seenResources[resourceID]; !found {
				if err := u.quotaRepo.UpdateResourceQuantity(existingQty.ID, 0); err != nil {
					return nil, apiError.NewInternalServerError(fmt.Errorf("failed to zero resource quantity: %w", err))
				}
				// Cascade down: clamp to 0
				if err := u.cascadeOrgQuotaToProjectQuotas(quotaID, resourceID, 0); err != nil {
					return nil, err
				}
			}
		}
	}

	return u.quotaRepo.GetOrgQuotaByID(quotaID)
}

// cascadeOrgQuotaToProjectQuotas clamps all child project quota resources to
// not exceed the new org quota limit, and recursively cascades to namespace quotas.
func (u *QuotaUsecase) cascadeOrgQuotaToProjectQuotas(orgQuotaID uuid.UUID, resourceID uuid.UUID, newLimit uint) error {
	projectQuotas, err := u.quotaRepo.GetProjectQuotasByOrgQuotaID(orgQuotaID)
	if err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to get project quotas for cascade: %w", err))
	}

	for _, pq := range projectQuotas {
		for _, rq := range pq.Resources {
			if rq.ResourceProp.ResourceID == resourceID && rq.Quantity > newLimit {
				if err := u.quotaRepo.UpdateResourceQuantity(rq.ID, newLimit); err != nil {
					return apiError.NewInternalServerError(fmt.Errorf("failed to cascade update project quota resource: %w", err))
				}
				// Continue cascading down to namespace quotas
				if err := u.cascadeProjectQuotaToNamespaceQuotas(pq.ID, resourceID, newLimit); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// cascadeProjectQuotaToNamespaceQuotas clamps all child namespace quota resources
// to not exceed the new project quota limit.
func (u *QuotaUsecase) cascadeProjectQuotaToNamespaceQuotas(projectQuotaID uuid.UUID, resourceID uuid.UUID, newLimit uint) error {
	namespaceQuotas, err := u.quotaRepo.GetNamespaceQuotasByProjectQuotaID(projectQuotaID)
	if err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to get namespace quotas for cascade: %w", err))
	}

	for _, nq := range namespaceQuotas {
		for _, rq := range nq.Resources {
			if rq.ResourceProp.ResourceID == resourceID && rq.Quantity > newLimit {
				if err := u.quotaRepo.UpdateResourceQuantity(rq.ID, newLimit); err != nil {
					return apiError.NewInternalServerError(fmt.Errorf("failed to cascade update namespace quota resource: %w", err))
				}
			}
		}
	}
	return nil
}

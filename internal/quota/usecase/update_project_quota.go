package usecase

import (
	"fmt"

	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/quota/dtos"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/google/uuid"
)

func (u *QuotaUsecase) UpdateProjectQuota(quotaID uuid.UUID, request *dtos.UpdateProjectQuotaRequest, userID uuid.UUID) (*models.ProjectQuota, error) {
	projectQuota, err := u.quotaRepo.GetProjectQuotaByID(quotaID)
	if err != nil {
		return nil, apiError.NewNotFoundError(fmt.Errorf("project quota not found: %w", err))
	}

	if err := u.isOrgAdmin(projectQuota.OrganizationID, userID); err != nil {
		return nil, err
	}

	// Update name/description
	if err := u.quotaRepo.UpdateProjectQuota(quotaID, request.Name, request.Description); err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to update project quota: %w", err))
	}

	// Update resources if provided
	if len(request.Resources) > 0 {
		// Determine the parent limit map based on whether this is inter-org or internal
		parentLimitMap, err := u.getProjectQuotaParentLimits(projectQuota)
		if err != nil {
			return nil, err
		}

		existingQuantities, err := u.quotaRepo.GetResourceQuantitiesByProjectQuotaID(quotaID)
		if err != nil {
			return nil, apiError.NewInternalServerError(fmt.Errorf("failed to get existing resource quantities: %w", err))
		}

		existingMap := make(map[uuid.UUID]models.ResourceQuantity)
		for _, qty := range existingQuantities {
			existingMap[qty.ResourceProp.ResourceID] = qty
		}

		seenResources := make(map[uuid.UUID]struct{})
		for _, r := range request.Resources {
			if _, duplicate := seenResources[r.ResourceID]; duplicate {
				return nil, apiError.NewBadRequestError(fmt.Errorf("duplicate resource ID: %s", r.ResourceID))
			}
			seenResources[r.ResourceID] = struct{}{}

			// Validate against parent limit
			if parentLimit, exists := parentLimitMap[r.ResourceID]; exists {
				if r.Quantity > parentLimit {
					return nil, apiError.NewBadRequestError(fmt.Errorf("requested quantity %d exceeds parent quota %d for resource %s", r.Quantity, parentLimit, r.ResourceID))
				}
			} else {
				return nil, apiError.NewBadRequestError(fmt.Errorf("resource %s not found in parent quota", r.ResourceID))
			}

			if existingQty, found := existingMap[r.ResourceID]; found {
				if err := u.quotaRepo.UpdateResourceQuantity(existingQty.ID, r.Quantity); err != nil {
					return nil, apiError.NewInternalServerError(fmt.Errorf("failed to update resource quantity: %w", err))
				}

				// Cascade down to namespace quotas
				if err := u.cascadeProjectQuotaToNamespaceQuotas(quotaID, r.ResourceID, r.Quantity); err != nil {
					return nil, err
				}
			} else {
				// Resource not yet allocated at project level — need to find the resource prop
				if projectQuota.OrganizationQuotaID != nil {
					// Inter-org: reuse the org quota's ResourceProperty
					orgQuota, err := u.quotaRepo.GetOrgQuotaByID(*projectQuota.OrganizationQuotaID)
					if err != nil {
						return nil, apiError.NewInternalServerError(fmt.Errorf("failed to get org quota: %w", err))
					}
					var propID uuid.UUID
					for _, oqr := range orgQuota.Resources {
						if oqr.ResourceProp.ResourceID == r.ResourceID {
							propID = oqr.ResourcePropID
							break
						}
					}
					if propID == uuid.Nil {
						return nil, apiError.NewBadRequestError(fmt.Errorf("resource %s not found in organization quota", r.ResourceID))
					}
					resourceQuantity := models.ResourceQuantity{
						ProjectQuotaID: &quotaID,
						ResourcePropID: propID,
						Quantity:       r.Quantity,
					}
					if err := u.quotaRepo.CreateResourceQuantity(&resourceQuantity); err != nil {
						return nil, apiError.NewInternalServerError(fmt.Errorf("failed to create resource quantity: %w", err))
					}
				} else {
					return nil, apiError.NewBadRequestError(fmt.Errorf("resource %s does not exist in this project quota", r.ResourceID))
				}
			}
		}

		// Zero out resources no longer in the request
		for resourceID, existingQty := range existingMap {
			if _, found := seenResources[resourceID]; !found {
				if err := u.quotaRepo.UpdateResourceQuantity(existingQty.ID, 0); err != nil {
					return nil, apiError.NewInternalServerError(fmt.Errorf("failed to zero resource quantity: %w", err))
				}
				if err := u.cascadeProjectQuotaToNamespaceQuotas(quotaID, resourceID, 0); err != nil {
					return nil, err
				}
			}
		}
	}

	return u.quotaRepo.GetProjectQuotaByID(quotaID)
}

func (u *QuotaUsecase) UpdateInternalProjectQuota(quotaID uuid.UUID, request *dtos.UpdateInternalProjectQuotaRequest, userID uuid.UUID) (*models.ProjectQuota, error) {
	projectQuota, err := u.quotaRepo.GetProjectQuotaByID(quotaID)
	if err != nil {
		return nil, apiError.NewNotFoundError(fmt.Errorf("project quota not found: %w", err))
	}

	if err := u.isOrgAdmin(projectQuota.OrganizationID, userID); err != nil {
		return nil, err
	}

	// Update name/description
	if err := u.quotaRepo.UpdateProjectQuota(quotaID, request.Name, request.Description); err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to update project quota: %w", err))
	}

	if len(request.Resources) > 0 {
		// Get node resources for upper bound validation
		node, err := u.resourceRepo.GetResourceNodeByID(projectQuota.NodeID)
		if err != nil {
			return nil, apiError.NewInternalServerError(fmt.Errorf("failed to get resource node: %w", err))
		}
		nodeResourceMap := make(map[uuid.UUID]models.Resource)
		for _, r := range node.Resources {
			nodeResourceMap[r.ID] = r
		}

		existingQuantities, err := u.quotaRepo.GetResourceQuantitiesByProjectQuotaID(quotaID)
		if err != nil {
			return nil, apiError.NewInternalServerError(fmt.Errorf("failed to get existing resource quantities: %w", err))
		}

		existingMap := make(map[uuid.UUID]models.ResourceQuantity)
		for _, qty := range existingQuantities {
			existingMap[qty.ResourceProp.ResourceID] = qty
		}

		seenResources := make(map[uuid.UUID]struct{})
		for _, r := range request.Resources {
			if _, duplicate := seenResources[r.ResourceID]; duplicate {
				return nil, apiError.NewBadRequestError(fmt.Errorf("duplicate resource ID: %s", r.ResourceID))
			}
			seenResources[r.ResourceID] = struct{}{}

			nodeRes, exists := nodeResourceMap[r.ResourceID]
			if !exists {
				return nil, apiError.NewBadRequestError(fmt.Errorf("resource %s not found in the resource node", r.ResourceID))
			}
			if r.Quantity > nodeRes.Quantity {
				return nil, apiError.NewBadRequestError(fmt.Errorf("requested quantity %d exceeds available quantity %d for resource %s", r.Quantity, nodeRes.Quantity, r.ResourceID))
			}

			if existingQty, found := existingMap[r.ResourceID]; found {
				if err := u.quotaRepo.UpdateResourceQuantity(existingQty.ID, r.Quantity); err != nil {
					return nil, apiError.NewInternalServerError(fmt.Errorf("failed to update resource quantity: %w", err))
				}
				// Internal projects can update price/maxtime
				if err := u.quotaRepo.UpdateResourceProperty(existingQty.ResourcePropID, r.Price, r.Duration); err != nil {
					return nil, apiError.NewInternalServerError(fmt.Errorf("failed to update resource property: %w", err))
				}

				// Cascade down to namespace quotas
				if err := u.cascadeProjectQuotaToNamespaceQuotas(quotaID, r.ResourceID, r.Quantity); err != nil {
					return nil, err
				}
			} else {
				// Create new resource property + quantity for internal project
				resourceProperty := models.ResourceProperty{
					ResourceID:  r.ResourceID,
					Price:       r.Price,
					MaxDuration: r.Duration,
				}
				if err := u.quotaRepo.CreateResourceProperty(&resourceProperty); err != nil {
					return nil, apiError.NewInternalServerError(fmt.Errorf("failed to create resource property: %w", err))
				}
				resourceQuantity := models.ResourceQuantity{
					ProjectQuotaID: &quotaID,
					Quantity:       r.Quantity,
					ResourcePropID: resourceProperty.ID,
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
				if err := u.cascadeProjectQuotaToNamespaceQuotas(quotaID, resourceID, 0); err != nil {
					return nil, err
				}
			}
		}
	}

	return u.quotaRepo.GetProjectQuotaByID(quotaID)
}

// getProjectQuotaParentLimits returns a map of resource_id -> max quantity
// based on the parent (org quota for inter-org, or node resources for internal).
func (u *QuotaUsecase) getProjectQuotaParentLimits(projectQuota *models.ProjectQuota) (map[uuid.UUID]uint, error) {
	limitMap := make(map[uuid.UUID]uint)

	if projectQuota.OrganizationQuotaID != nil {
		// Inter-org: parent is the org quota
		orgQuota, err := u.quotaRepo.GetOrgQuotaByID(*projectQuota.OrganizationQuotaID)
		if err != nil {
			return nil, apiError.NewInternalServerError(fmt.Errorf("failed to get parent org quota: %w", err))
		}
		for _, rq := range orgQuota.Resources {
			limitMap[rq.ResourceProp.ResourceID] = rq.Quantity
		}
	} else {
		// Internal project: parent is the node resource pool
		node, err := u.resourceRepo.GetResourceNodeByID(projectQuota.NodeID)
		if err != nil {
			return nil, apiError.NewInternalServerError(fmt.Errorf("failed to get resource node: %w", err))
		}
		for _, r := range node.Resources {
			limitMap[r.ID] = r.Quantity
		}
	}

	return limitMap, nil
}

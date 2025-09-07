package usecase

import (
	"errors"
	"fmt"
	"log"

	"github.com/ClearingHouse/helper"
	"github.com/ClearingHouse/internal/models"
	namespaceInterfaces "github.com/ClearingHouse/internal/namespaces/interfaces"
	orgInterfaces "github.com/ClearingHouse/internal/organizations/interfaces"
	projInterfaces "github.com/ClearingHouse/internal/projects/interfaces"
	"github.com/ClearingHouse/internal/quota/dtos"
	"github.com/ClearingHouse/internal/quota/interfaces"
	resourcesInterfaces "github.com/ClearingHouse/internal/resources/interfaces"
	userInterfaces "github.com/ClearingHouse/internal/users/interfaces"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/google/uuid"
)

type QuotaUsecase struct {
	quotaRepo     interfaces.QuotaRepository
	resourceRepo  resourcesInterfaces.ResourceRepository
	namespaceRepo namespaceInterfaces.NamespaceRepository
	orgRepo       orgInterfaces.OrganizationRepository
	projRepo      projInterfaces.ProjectRepository
	userRepo      userInterfaces.UsersRepository
}

func NewQuotaUsecase(quotaRepo interfaces.QuotaRepository, resourceRepo resourcesInterfaces.ResourceRepository, namespaceRepo namespaceInterfaces.NamespaceRepository, orgRepo orgInterfaces.OrganizationRepository, projRepo projInterfaces.ProjectRepository, userRepo userInterfaces.UsersRepository) interfaces.QuotaUsecase {
	return &QuotaUsecase{
		quotaRepo:     quotaRepo,
		resourceRepo:  resourceRepo,
		namespaceRepo: namespaceRepo,
		orgRepo:       orgRepo,
		projRepo:      projRepo,
		userRepo:      userRepo,
	}
}

func (u *QuotaUsecase) CreateOrganizationQuota(request *dtos.CreateOrganizationQuotaRequest, userID uuid.UUID) (*models.OrganizationQuota, error) {
	if len(request.Resources) == 0 {
		return nil, apiError.NewBadRequestError(errors.New("at least one resource quota is required"))
	}

	isOrgAdmin, err := u.isOrgAdmin(request.FromOrganizationID, userID)
	if err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to check organization admin: %w", err))
	}
	if !isOrgAdmin {
		return nil, apiError.NewForbiddenError(errors.New("user is not an admin of the organization"))
	}

	_, err = u.orgRepo.GetOrganizationByID(request.ToOrganizationID)
	if err != nil {
		return nil, apiError.NewNotFoundError(fmt.Errorf("failed to find target organization: %w", err))
	}

	var seenResource = map[uuid.UUID]struct{}{}
	var poolID uuid.UUID

	for i, r := range request.Resources {
		// check duplicate resource
		if _, exists := seenResource[r.ResourceID]; exists {
			return nil, fmt.Errorf("duplicate resource ID: %s", r.ResourceID)
		}
		seenResource[r.ResourceID] = struct{}{}

		// fetch resource
		resource, err := u.resourceRepo.GetResourceByID(r.ResourceID)
		if err != nil {
			return nil, apiError.NewNotFoundError(fmt.Errorf("resource not found: %w", err))
		}

		// check resource ownership
		if resource.ResourcePool.OrganizationID != request.FromOrganizationID {
			return nil, apiError.NewForbiddenError(errors.New("unauthorized: resource not owned by organization"))
		}

		if r.Quantity > resource.Quantity {
			return nil, apiError.NewBadRequestError(fmt.Errorf("requested quantity %d exceeds available quantity %d for resource %s", r.Quantity, resource.Quantity, r.ResourceID))
		}

		// check resource pool consistency
		if i == 0 {
			poolID = resource.ResourcePoolID
		} else if resource.ResourcePoolID != poolID {
			return nil, apiError.NewBadRequestError(errors.New("all resources must be from the same pool"))
		}
	}

	IsOrgQuotaExist, err := u.quotaRepo.IsOrgQuotaExist(request.FromOrganizationID, request.ToOrganizationID, poolID)
	if err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to check existing quota: %w", err))
	}

	if IsOrgQuotaExist {
		return nil, apiError.NewConflictError(errors.New("quota already exists between the organizations for this pool"))
	}

	quotaGroup := &models.OrganizationQuota{
		Name:           request.Name,
		Description:    request.Description,
		ResourcePoolID: poolID,
		FromOrgID:      request.FromOrganizationID,
		ToOrgID:        request.ToOrganizationID,
	}

	err = u.quotaRepo.CreateOrgQuota(quotaGroup)
	if err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to create organization quota: %w", err))
	}

	for _, r := range request.Resources {
		resourceProperty := models.ResourceProperty{
			ResourceID:  r.ResourceID,
			Price:       r.Price,
			MaxDuration: float32(r.Duration),
		}

		err := u.quotaRepo.CreateResourceProperty(&resourceProperty)
		if err != nil {
			return nil, apiError.NewInternalServerError(fmt.Errorf("failed to create resource property: %w", err))
		}

		resourceQuantity := models.ResourceQuantity{
			OrganizationQuotaID: &quotaGroup.ID,
			Quantity:            r.Quantity,
			ResourcePropID:      resourceProperty.ID,
		}

		err = u.quotaRepo.CreateResourceQuantity(&resourceQuantity)
		if err != nil {
			return nil, apiError.NewInternalServerError(fmt.Errorf("failed to create resource quantity: %w", err))
		}
	}

	return quotaGroup, nil
}

func (u *QuotaUsecase) GetOrganizationQuota(fromOrgID uuid.UUID, toOrgID uuid.UUID) ([]models.OrganizationQuota, error) {
	quotas, err := u.quotaRepo.GetOrganizationByRelationship(fromOrgID, toOrgID)
	if err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to get organization quotas: %w", err))
	}
	return quotas, nil
}

func (u *QuotaUsecase) CreateProjectQuota(request *dtos.CreateProjectQuotaRequest, userID uuid.UUID) (*models.ProjectQuota, error) {
	isOrgAdmin, err := u.isOrgAdmin(request.OrgID, userID)
	if err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to check organization admin: %w", err))
	}
	if !isOrgAdmin {
		return nil, apiError.NewForbiddenError(errors.New("user is not an admin of the organization"))
	}

	quota, err := u.quotaRepo.GetOrgQuotaByID(request.OrgQuotaID)
	if err != nil {
		return nil, apiError.NewNotFoundError(fmt.Errorf("failed to find organization quota: %w", err))
	}
	if quota.ToOrgID != request.OrgID {
		return nil, apiError.NewForbiddenError(errors.New("quota does not belong to the organization"))
	}

	quotaResources := make(map[uuid.UUID]models.ResourceQuantity)
	for _, resource := range quota.Resources {
		quotaResources[resource.ResourceProp.ResourceID] = resource
	}
	var seenResource = map[uuid.UUID]struct{}{}

	for _, r := range request.Resources {
		if _, exists := quotaResources[r.ResourceID]; !exists {
			return nil, apiError.NewBadRequestError(fmt.Errorf("resource %s not found in organization quota", r.ResourceID))
		}

		if _, exists := seenResource[r.ResourceID]; exists {
			return nil, apiError.NewConflictError(errors.New("duplicate resource found"))
		}
		seenResource[r.ResourceID] = struct{}{}

		currentUsage, err := u.quotaRepo.GetOrgUsage(quota.ID, r.ResourceID)
		if err != nil {
			return nil, fmt.Errorf("failed to get current usage for resource %s: %w", r.ResourceID, err)
		}

		maximumQuota, err := u.quotaRepo.GetOrgQuotaQuantity(quota.ID, r.ResourceID)
		if err != nil {
			return nil, fmt.Errorf("failed to get maximum quota for resource %s: %w", r.ResourceID, err)
		}

		log.Println("current usage", currentUsage)
		log.Println("maximum quota", maximumQuota)
		if r.Quantity+currentUsage > maximumQuota {
			return nil, errors.New("too much quota requested")
		}
	}

	projectQuota := &models.ProjectQuota{
		Name:                request.Name,
		Description:         request.Description,
		OrganizationID:      request.OrgID,
		ProjectID:           request.ProjectID,
		OrganizationQuotaID: request.OrgQuotaID,
		ResourcePoolID:      quota.ResourcePoolID,
	}
	err = u.quotaRepo.CreateProjectQuota(projectQuota)
	if err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to create project quota: %w", err))
	}

	for _, r := range request.Resources {
		resourceQuantity := models.ResourceQuantity{
			ProjectQuotaID: &projectQuota.ID,
			ResourcePropID: quotaResources[r.ResourceID].ResourceProp.ID,
			Quantity:       r.Quantity,
		}

		err = u.quotaRepo.CreateResourceQuantity(&resourceQuantity)
		if err != nil {
			return nil, apiError.NewInternalServerError(fmt.Errorf("failed to create resource quantity: %w", err))
		}
	}

	return projectQuota, nil
}

func (u *QuotaUsecase) GetProjectQuota(projectID uuid.UUID) ([]models.ProjectQuota, error) {
	quotas, err := u.quotaRepo.GetProjectQuotaByProjectID(projectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to get project quotas: %w", err))
	}
	return quotas, nil
}

func (u *QuotaUsecase) CreateNamespaceQuota(request *dtos.CreateNamespaceQuotaRequest, userID uuid.UUID) (*models.NamespaceQuota, error) {
	if len(request.Resources) == 0 {
		return nil, apiError.NewBadRequestError(errors.New("at least one resource quota is required"))
	}

	project, err := u.projRepo.GetProjectByID(request.ProjectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to find project: %w", err))
	}

	isProjAdmin, err := u.isProjAdmin(project.ID, userID)
	if err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to check project admin status: %w", err))
	}
	if !isProjAdmin {
		return nil, apiError.NewForbiddenError(errors.New("user is not a project admin"))
	}

	quota, err := u.quotaRepo.GetProjectQuotaByID(request.ProjectQuotaID)
	if err != nil {
		return nil, apiError.NewNotFoundError(fmt.Errorf("failed to find project quota: %w", err))
	}
	if quota.ProjectID != request.ProjectID {
		return nil, apiError.NewForbiddenError(errors.New("quota does not belong to the project"))
	}

	quotaResources := make(map[uuid.UUID]models.ResourceQuantity)
	for _, resource := range quota.Resources {
		quotaResources[resource.ResourceProp.ResourceID] = resource
	}
	var seenResource = map[uuid.UUID]struct{}{}

	for _, r := range request.Resources {
		_, err := u.resourceRepo.GetResourceByID(r.ResourceID)
		if err != nil {
			return nil, apiError.NewNotFoundError(fmt.Errorf("resource not found: %w", err))
		}

		if _, exists := quotaResources[r.ResourceID]; !exists {
			return nil, apiError.NewBadRequestError(fmt.Errorf("resource %s not found in project quota", r.ResourceID))
		}

		if _, exists := seenResource[r.ResourceID]; exists {
			return nil, apiError.NewConflictError(errors.New("duplicate resource found"))
		}
		seenResource[r.ResourceID] = struct{}{}

		if r.Quantity > quotaResources[r.ResourceID].Quantity {
			return nil, apiError.NewBadRequestError(fmt.Errorf("requested quantity %d exceeds available quantity %d for resource %s", r.Quantity, quotaResources[r.ResourceID].Quantity, r.ResourceID))
		}
	}

	namespaceQuota := &models.NamespaceQuota{
		Name:           request.Name,
		Description:    request.Description,
		ProjectID:      request.ProjectID,
		ProjectQuotaID: request.ProjectQuotaID,
		ResourcePoolID: quota.ResourcePoolID,
	}
	err = u.quotaRepo.CreateNamespaceQuota(namespaceQuota)
	if err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to create namespace quota: %w", err))
	}

	for _, r := range request.Resources {
		resourceQuantity := models.ResourceQuantity{
			NamespaceQuotaID: &namespaceQuota.ID,
			ResourcePropID:   quotaResources[r.ResourceID].ResourceProp.ID,
			Quantity:         r.Quantity,
		}

		err = u.quotaRepo.CreateResourceQuantity(&resourceQuantity)
		if err != nil {
			return nil, apiError.NewInternalServerError(fmt.Errorf("failed to create resource quantity: %w", err))
		}
	}

	return nil, nil
}

func (u *QuotaUsecase) GetNamespaceQuota(namespaceID uuid.UUID) ([]models.NamespaceQuota, error) {
	quotas, err := u.quotaRepo.GetNamespaceQuotaByNamespaceID(namespaceID)
	if err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to get namespace quota: %w", err))
	}
	return quotas, nil
}

func (u *QuotaUsecase) AssignQuotaToNamespace(request *dtos.AssignQuotaToNamespaceRequest, userID uuid.UUID) error {
	project, err := u.projRepo.GetProjectByID(request.ProjectID)
	if err != nil {
		return apiError.NewInternalServerError(err)
	}

	isProjAdmin, err := u.isProjAdmin(project.ID, userID)
	if err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to check project admin status: %w", err))
	}
	if !isProjAdmin {
		return apiError.NewForbiddenError(errors.New("user is not a project admin"))
	}

	namespaceQuota, err := u.quotaRepo.GetNamespaceQuotaByID(request.QuotaID)
	if err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to find namespace quota group: %w", err))
	}

	if namespaceQuota.ProjectID != request.ProjectID {
		return apiError.NewBadRequestError(errors.New("quota group does not belong to the project"))
	}

	allNamespaceInProj, err := u.namespaceRepo.GetAllNamespacesByProjectID(namespaceQuota.ProjectID)
	if err != nil {
		return apiError.NewInternalServerError(fmt.Errorf("failed to find namespaces by project ID: %w", err))
	}

	nsMap := make(map[uuid.UUID]struct{})
	for _, ns := range allNamespaceInProj {
		nsMap[ns.ID] = struct{}{}
	}

	// validate namespaces
	for _, namespaceID := range request.Namespaces {
		if _, ok := nsMap[namespaceID]; !ok {
			return apiError.NewBadRequestError(fmt.Errorf("namespace %s does not belong to the project %s", namespaceID, namespaceQuota.ProjectID))
		}

	}

	// insert into database
	for _, namespaceID := range request.Namespaces {
		err = u.quotaRepo.AssignQuotaToNamespace(namespaceID, request.QuotaID)
		if err != nil {
			return apiError.NewInternalServerError(fmt.Errorf("failed to assign quota to namespace %s: %w", namespaceID, err))
		}
	}

	return nil
}

func (u *QuotaUsecase) isOrgAdmin(orgID uuid.UUID, userID uuid.UUID) (bool, error) {
	org, err := u.orgRepo.GetOrganizationByID(orgID)
	if err != nil {
		return false, err
	}

	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return false, err
	}

	if !helper.ContainsUserID(org.Admins, user.ID) {
		return false, nil
	}

	return true, nil
}

func (u *QuotaUsecase) isProjAdmin(projID uuid.UUID, userID uuid.UUID) (bool, error) {
	proj, err := u.projRepo.GetProjectByID(projID)
	if err != nil {
		return false, err
	}

	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return false, err
	}

	if !helper.ContainsUserID(proj.Admins, user.ID) {
		return false, nil
	}

	return true, nil
}

// // SERIOUS

// func (u *QuotaUsecase) CreateOrganizationQuota(request *dtos.CreateOrganizationQuotaRequest, userID uuid.UUID) (*models.OrganizationQuota, error) {
// 	if len(request.Resources) == 0 {
// 		return nil, apiError.NewBadRequestError(fmt.Errorf("at least one resource quota is required"))
// 	}

// 	err := u.isOrgAdmin(request.FromOrganizationID, userID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	_, err = u.orgRepo.GetOrganizationByID(request.ToOrganizationID)
// 	if err != nil {
// 		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to find target organization: %w", err))
// 	}

// 	return nil, nil
// }

// func (u *QuotaUsecase) isOrgAdmin(orgID uuid.UUID, userID uuid.UUID) error {
// 	org, err := u.orgRepo.GetOrganizationByID(orgID)
// 	if err != nil {
// 		return apiError.NewInternalServerError(fmt.Errorf("failed to find target organization: %w", err))
// 	}

// 	user, err := u.userRepo.GetByID(userID)
// 	if err != nil {
// 		return apiError.NewInternalServerError(fmt.Errorf("failed to find user: %w", err))
// 	}

// 	if !helper.ContainsUserID(org.Admins, user.ID) {
// 		return apiError.NewForbiddenError(errors.New("user is not an admin of the organization"))
// 	}

// 	return nil
// }

// func (u *QuotaUsecase) CreateOrganizationQuotaGroup(request *dtos.CreateOrganizationQuotaRequest) (*models.OrganizationQuotaGroup, error) {
// 	var resources []models.Resource
// 	var resourceTypeCheck = map[string]bool{}
// 	var poolID uuid.UUID

// 	for i, r := range request.Resources {
// 		resource, err := u.resourceRepo.GetResourceByID(r.ResourceID)
// 		if err != nil {
// 			return nil, fmt.Errorf("resource not found: %w", err)
// 		}

// 		if resource.ResourcePool.OrganizationID != request.FromOrganizationID {
// 			return nil, errors.New("unauthorized: resource not owned by organization")
// 		}

// 		if i == 0 {
// 			poolID = resource.ResourcePoolID
// 		} else if resource.ResourcePoolID != poolID {
// 			return nil, errors.New("all resources must be from the same pool")
// 		}

// 		rtName := resource.ResourceType.Name
// 		if resourceTypeCheck[rtName] {
// 			return nil, fmt.Errorf("duplicate resource type: %s", rtName)
// 		}
// 		resourceTypeCheck[rtName] = true

// 		resources = append(resources, *resource)
// 	}

// 	existingQuotaGroup, err := u.quotaRepo.FindExistingOrganizationQuotaGroup(request.FromOrganizationID, request.ToOrganizationID, poolID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to check existing quota: %w", err)
// 	}

// 	if existingQuotaGroup != nil {
// 		return nil, fmt.Errorf("quota already exists between the organizations for this pool")
// 	}

// 	quotaGroup := &models.OrganizationQuotaGroup{
// 		Name:               request.Name,
// 		Description:        request.Description,
// 		FromOrganizationID: request.FromOrganizationID,
// 		ToOrganizationID:   request.ToOrganizationID,
// 	}

// 	err = u.quotaRepo.CreateOrganizationQuotaGroup(quotaGroup)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create quota group: %w", err)
// 	}

// 	for _, resource := range request.Resources {
// 		resourceProperty := models.ResourceProperty{
// 			ResourceID: resource.ResourceID,
// 			Price:      resource.Price,
// 		}

// 		err := u.quotaRepo.CreateResourceProperty(&resourceProperty)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to create resource property: %w", err)
// 		}

// 		resourceQuantity := models.ResourceQuantity{
// 			OrganizationQuotaGroupID: &quotaGroup.BaseModel.ID,
// 			ResourcePropertyID:       resourceProperty.ID,
// 			Quantity:                 resource.Quantity,
// 		}

// 		err = u.quotaRepo.CreateResourceQuantity(&resourceQuantity)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to create resource quantity: %w", err)
// 		}
// 	}

// 	return quotaGroup, nil
// }

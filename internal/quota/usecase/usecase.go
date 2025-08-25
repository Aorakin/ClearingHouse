package usecase

import (
	"errors"
	"fmt"
	"log"
	"slices"

	"github.com/ClearingHouse/helper"
	"github.com/ClearingHouse/internal/models"
	namespaceInterfaces "github.com/ClearingHouse/internal/namespaces/interfaces"
	orgInterfaces "github.com/ClearingHouse/internal/organizations/interfaces"
	projInterfaces "github.com/ClearingHouse/internal/projects/interfaces"
	"github.com/ClearingHouse/internal/quota/dtos"
	"github.com/ClearingHouse/internal/quota/interfaces"
	resourcesInterfaces "github.com/ClearingHouse/internal/resources/interfaces"
	userInterfaces "github.com/ClearingHouse/internal/users/interfaces"
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

func (u *QuotaUsecase) CreateOrganizationQuotaGroup(request *dtos.CreateOrganizationQuotaRequest) (*models.OrganizationQuotaGroup, error) {
	if len(request.Resources) == 0 {
		return nil, errors.New("at least one resource quota is required")
	}

	isOrgAdmin, err := u.isOrgAdmin(request.Creator, request.FromOrganizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to check organization admin: %w", err)
	}
	if !isOrgAdmin {
		return nil, errors.New("user is not an admin of the organization")
	}

	_, err = u.orgRepo.GetOrganizationByID(request.ToOrganizationID)
	if err != nil {
		return nil, fmt.Errorf("target organization not found: %w", err)
	}

	var resources []models.Resource
	var resourceTypeCheck = map[string]bool{}
	var poolID uuid.UUID

	for i, r := range request.Resources {
		resource, err := u.resourceRepo.GetResourceByID(r.ResourceID)
		if err != nil {
			return nil, fmt.Errorf("resource not found: %w", err)
		}

		if resource.ResourcePool.OrganizationID != request.FromOrganizationID {
			return nil, errors.New("unauthorized: resource not owned by organization")
		}

		if i == 0 {
			poolID = resource.ResourcePoolID
		} else if resource.ResourcePoolID != poolID {
			return nil, errors.New("all resources must be from the same pool")
		}

		rtName := resource.ResourceType.Name
		if resourceTypeCheck[rtName] {
			return nil, fmt.Errorf("duplicate resource type: %s", rtName)
		}
		resourceTypeCheck[rtName] = true

		resources = append(resources, *resource)
	}

	existingQuotaGroup, err := u.quotaRepo.FindExistingOrganizationQuotaGroup(request.FromOrganizationID, request.ToOrganizationID, poolID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing quota: %w", err)
	}

	if existingQuotaGroup != nil {
		return nil, fmt.Errorf("quota already exists between the organizations for this pool")
	}

	quotaGroup := &models.OrganizationQuotaGroup{
		Name:               request.Name,
		Description:        request.Description,
		FromOrganizationID: request.FromOrganizationID,
		ToOrganizationID:   request.ToOrganizationID,
	}

	err = u.quotaRepo.CreateOrganizationQuotaGroup(quotaGroup)
	if err != nil {
		return nil, fmt.Errorf("failed to create quota group: %w", err)
	}

	for _, resource := range request.Resources {
		resourceProperty := models.ResourceProperty{
			ResourceID: resource.ResourceID,
			Price:      resource.Price,
		}

		err := u.quotaRepo.CreateResourceProperty(&resourceProperty)
		if err != nil {
			return nil, fmt.Errorf("failed to create resource property: %w", err)
		}

		resourceQuantity := models.ResourceQuantity{
			OrganizationQuotaGroupID: &quotaGroup.BaseModel.ID,
			ResourcePropertyID:       resourceProperty.ID,
			Quantity:                 resource.Quantity,
		}

		err = u.quotaRepo.CreateResourceQuantity(&resourceQuantity)
		if err != nil {
			return nil, fmt.Errorf("failed to create resource quantity: %w", err)
		}
	}

	return quotaGroup, nil
}

func (u *QuotaUsecase) FindOrganizationQuotaGroup(fromOrgID uuid.UUID, toOrgID uuid.UUID) ([]models.OrganizationQuotaGroup, error) {
	quotaGroups, err := u.quotaRepo.FindOrganizationQuotaGroup(fromOrgID, toOrgID)
	if err != nil {
		return nil, fmt.Errorf("failed to find organization quota groups: %w", err)
	}

	return quotaGroups, nil
}

func (u *QuotaUsecase) CreateProjectQuotaGroup(request *dtos.CreateProjectQuotaRequest) (*models.ProjectQuotaGroup, error) {
	isOrgAdmin, err := u.isOrgAdmin(request.Creator, request.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to check organization admin: %w", err)
	}
	if !isOrgAdmin {
		return nil, errors.New("user is not an admin of the organization")
	}

	for _, resource_pool := range request.ResourcePools {
		quotaGroup, err := u.quotaRepo.FindOrganizationQuotaGroupByID(resource_pool.QuotaGroupID)
		if err != nil {
			return nil, fmt.Errorf("failed to get organization quota group: %w", err)
		}
		if quotaGroup.ToOrganizationID != request.OrganizationID {
			return nil, errors.New("quota does not belong to the organization")
		}

		var resources uuid.UUIDs
		for _, resource := range quotaGroup.Resources {
			resources = append(resources, resource.ResourceProperty.ResourceID)
		}

		for _, resource := range resource_pool.Resources {
			if !slices.Contains(resources, resource.ResourceID) {
				return nil, fmt.Errorf("resource %s not found in organization quota group", resource.ResourceID)
			}

			currentUsage, err := u.quotaRepo.GetOrgUsage(quotaGroup.BaseModel.ID, resource.ResourceID)
			if err != nil {
				return nil, fmt.Errorf("failed to get current usage for resource %s: %w", resource.ResourceID, err)
			}

			maximumQuota, err := u.quotaRepo.GetOrgQuotaQuantity(quotaGroup.BaseModel.ID, resource.ResourceID)
			if err != nil {
				return nil, fmt.Errorf("failed to get maximum quota for resource %s: %w", resource.ResourceID, err)
			}

			log.Println("current usage", currentUsage)
			log.Println("maximum quota", maximumQuota)
			if resource.Quantity+currentUsage > maximumQuota {
				return nil, errors.New("too much quota requested")
			}
		}
	}

	for _, resourcePool := range request.ResourcePools {
		projectQuotaGroup := &models.ProjectQuotaGroup{
			Name:                     request.Name,
			Description:              request.Description,
			OrganizationID:           request.OrganizationID,
			OrganizationQuotaGroupID: resourcePool.QuotaGroupID,
		}

		err := u.quotaRepo.CreateProjectQuotaGroup(projectQuotaGroup)
		if err != nil {
			return nil, fmt.Errorf("failed to create project quota group: %w", err)
		}

		for _, resource := range resourcePool.Resources {
			resourceProperty, err := u.quotaRepo.GetResourcePropertyByOrg(resourcePool.QuotaGroupID, resource.ResourceID)
			if err != nil {
				return nil, fmt.Errorf("failed to get resource property: %w", err)
			}

			resourceQuantity := models.ResourceQuantity{
				ProjectQuotaGroupID: &projectQuotaGroup.BaseModel.ID,
				ResourcePropertyID:  resourceProperty.ID,
				Quantity:            resource.Quantity,
			}

			err = u.quotaRepo.CreateResourceQuantity(&resourceQuantity)
			if err != nil {
				return nil, fmt.Errorf("failed to create resource quantity: %w", err)
			}
		}
	}

	return &models.ProjectQuotaGroup{}, nil
}

func (u *QuotaUsecase) FindProjectQuotaGroup(projectID uuid.UUID) ([]models.ProjectQuotaGroup, error) {
	return u.quotaRepo.FindProjectQuotaGroupByProjectID(projectID)
}

func (u *QuotaUsecase) CreateNamespaceQuotaGroup(request *dtos.CreateNamespaceQuotaRequest) (*models.NamespaceQuotaGroup, error) {
	// Validate the request: make sure ProjectID is valid and resources belong to the right org

	// CHECK USER AUTHORIZATION

	for _, resourcePool := range request.ResourcePools {
		projectQuotaGroup, err := u.quotaRepo.FindProjectQuotaGroupByID(resourcePool.QuotaGroupID)
		if err != nil {
			return nil, fmt.Errorf("failed to get project quota group: %w", err)
		}

		var allowedResources uuid.UUIDs
		for _, resource := range projectQuotaGroup.Resources {
			log.Printf("quota group id %s resource id %s", projectQuotaGroup.BaseModel.ID, resource.ResourceProperty.ResourceID)
			allowedResources = append(allowedResources, resource.ResourceProperty.ResourceID)
		}

		for _, resource := range resourcePool.Resources {
			if !slices.Contains(allowedResources, resource.ResourceID) {
				return nil, fmt.Errorf("resource %s not found in project quota group", resource.ResourceID)
			}

			maxQuota, err := u.quotaRepo.GetProjectQuotaQuantity(projectQuotaGroup.BaseModel.ID, resource.ResourceID)
			if err != nil {
				return nil, fmt.Errorf("failed to get max quota for resource %s: %w", resource.ResourceID, err)
			}

			if resource.Quantity > maxQuota {
				return nil, errors.New("too much quota requested")
			}
		}
	}

	namespaceQuotaGroup := &models.NamespaceQuotaGroup{
		Name:        request.Name,
		Description: request.Description,
		ProjectID:   request.ProjectID,
	}

	err := u.quotaRepo.CreateNamespaceQuotaGroup(namespaceQuotaGroup)
	if err != nil {
		return nil, fmt.Errorf("failed to create namespace quota group: %w", err)
	}

	for _, resourcePool := range request.ResourcePools {
		for _, resource := range resourcePool.Resources {
			resourceProperty, err := u.quotaRepo.GetResourcePropertyByProj(resourcePool.QuotaGroupID, resource.ResourceID)
			if err != nil {
				return nil, fmt.Errorf("failed to get resource property: %w", err)
			}

			resourceQuantity := models.ResourceQuantity{
				NamespaceQuotaGroupID: &namespaceQuotaGroup.BaseModel.ID,
				ResourcePropertyID:    resourceProperty.ID,
				Quantity:              resource.Quantity,
			}

			err = u.quotaRepo.CreateResourceQuantity(&resourceQuantity)
			if err != nil {
				return nil, fmt.Errorf("failed to create resource quantity: %w", err)
			}
		}
	}

	return namespaceQuotaGroup, nil
}

func (u *QuotaUsecase) AssignQuotaToNamespace(request *dtos.AssignQuotaToNamespaceRequest) error {
	project, err := u.projRepo.FindProjectByID(request.ProjectID)
	if err != nil {
		return fmt.Errorf("failed to find project: %w", err)
	}

	isProjAdmin, err := u.IsProjAdmin(request.Creator, project.ID)
	if err != nil {
		return fmt.Errorf("failed to check project admin status: %w", err)
	}
	if !isProjAdmin {
		return fmt.Errorf("user is not a project admin")
	}

	namespaceQuotaGroup, err := u.quotaRepo.FindNamespaceQuotaGroupByID(request.QuotaGroupID)
	if err != nil {
		return fmt.Errorf("failed to find namespace quota group: %w", err)
	}

	if namespaceQuotaGroup.ProjectID != request.ProjectID {
		return fmt.Errorf("quota group does not belong to the project")
	}

	allNamespaceInProj, err := u.namespaceRepo.FindAllNamespacesByProjectID(namespaceQuotaGroup.ProjectID)
	if err != nil {
		return fmt.Errorf("failed to find namespaces by project ID: %w", err)
	}

	nsMap := make(map[uuid.UUID]struct{})
	for _, ns := range allNamespaceInProj {
		nsMap[ns.ID] = struct{}{}
	}

	// validate namespaces
	for _, namespaceID := range request.Namespaces {
		if _, ok := nsMap[namespaceID]; !ok {
			return fmt.Errorf("namespace %s does not belong to the project %s", namespaceID, namespaceQuotaGroup.ProjectID.String())
		}

	}

	// insert into database
	for _, namespaceID := range request.Namespaces {
		err = u.quotaRepo.AssignQuotaToNamespace(namespaceID, request.QuotaGroupID)
		if err != nil {
			return fmt.Errorf("failed to assign quota to namespace %s: %w", namespaceID, err)
		}
	}

	return nil
}

func (u *QuotaUsecase) isOrgAdmin(userID uuid.UUID, orgID uuid.UUID) (bool, error) {
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

func (u *QuotaUsecase) IsProjAdmin(userID uuid.UUID, projID uuid.UUID) (bool, error) {
	proj, err := u.projRepo.FindProjectByID(projID)
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

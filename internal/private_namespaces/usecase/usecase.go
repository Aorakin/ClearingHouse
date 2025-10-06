package usecase

import (
	"errors"
	"fmt"

	"github.com/ClearingHouse/helper"
	"github.com/ClearingHouse/internal/models"
	namespaceInterfaces "github.com/ClearingHouse/internal/namespaces/interfaces"
	orgInterfaces "github.com/ClearingHouse/internal/organizations/interfaces"
	"github.com/ClearingHouse/internal/private_namespaces/dtos"
	"github.com/ClearingHouse/internal/private_namespaces/interfaces"
	quotaInterfaces "github.com/ClearingHouse/internal/quota/interfaces"
	resourceInterfaces "github.com/ClearingHouse/internal/resources/interfaces"
	userInterfaces "github.com/ClearingHouse/internal/users/interfaces"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PrivateNamespaceUsecase struct {
	privNamespaceRepo interfaces.PrivateNamespaceRepository
	namespaceRepo     namespaceInterfaces.NamespaceRepository
	orgRepo           orgInterfaces.OrganizationRepository
	userRepo          userInterfaces.UsersRepository
	quotaRepo         quotaInterfaces.QuotaRepository
	resourceRepo      resourceInterfaces.ResourceRepository
}

func NewPrivateNamespaceUsecase(privNamespaceRepository interfaces.PrivateNamespaceRepository, namespaceRepository namespaceInterfaces.NamespaceRepository, orgRepository orgInterfaces.OrganizationRepository, userRepository userInterfaces.UsersRepository, quotaRepository quotaInterfaces.QuotaRepository, resourceRepository resourceInterfaces.ResourceRepository) interfaces.PrivateNamespaceUsecase {
	return &PrivateNamespaceUsecase{
		privNamespaceRepo: privNamespaceRepository,
		namespaceRepo:     namespaceRepository,
		orgRepo:           orgRepository,
		userRepo:          userRepository,
		quotaRepo:         quotaRepository,
		resourceRepo:      resourceRepository,
	}
}

func (u *PrivateNamespaceUsecase) GetPrivateNamespaceByOwnerID(ownerID uuid.UUID) (*models.Namespace, error) {
	namespace, err := u.privNamespaceRepo.GetPrivateNamespaceByOwnerID(ownerID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apiError.NewNotFoundError("namespace not found")
		}
		return nil, err
	}
	return namespace, nil
}

func (u *PrivateNamespaceUsecase) CreatePrivateNamespace(request *dtos.CreatePrivateNamespaceRequest, userID uuid.UUID) (*models.Namespace, error) {
	isAdmin, err := u.isOrgAdmin(request.OrganizationID, userID)
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, apiError.NewUnauthorizedError("user is not organization admin")
	}

	namespace := &models.Namespace{
		Description: request.Description,
		Credit:      request.Credit,
		OrgID:       request.OrganizationID,
		Name:        request.Name,
	}

	if err := u.privNamespaceRepo.CreatePrivateNamespace(namespace); err != nil {
		return nil, err
	}

	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to fetch user: %w", err))
	}

	user.NamespaceID = &namespace.ID
	if err := u.userRepo.Update(user); err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to update user with namespace ID: %w", err))
	}

	return namespace, nil
}

func (u *PrivateNamespaceUsecase) CreateNamespaceQuota(request *dtos.CreateNamespaceQuotaRequest, userID uuid.UUID) (*models.NamespaceQuota, error) {
	if len(request.Resources) == 0 {
		return nil, apiError.NewBadRequestError(errors.New("at least one resource quota is required"))
	}

	namespace, err := u.namespaceRepo.GetNamespaceByID(request.NamespaceID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apiError.NewNotFoundError(errors.New("namespace not found"))
		}
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to find namespace: %w", err))
	}

	isAdmin, err := u.isOrgAdmin(namespace.OrgID, userID)
	if err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to check organization admin status: %w", err))
	}
	if !isAdmin {
		return nil, apiError.NewForbiddenError(errors.New("user is not organization admin"))
	}

	var seenResource = map[uuid.UUID]struct{}{}

	isQuotaExists, err := u.quotaRepo.IsNamespaceQuotaExists(request.NamespaceID, request.PoolID)
	if err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to check namespace quota existence: %w", err))
	}
	if isQuotaExists {
		return nil, apiError.NewConflictError(errors.New("namespace quota already exists"))
	}

	for _, r := range request.Resources {
		// check duplicate resource
		if _, exists := seenResource[r.ResourceID]; exists {
			return nil, apiError.NewBadRequestError(fmt.Errorf("duplicate resource ID: %s", r.ResourceID))
		}
		seenResource[r.ResourceID] = struct{}{}

		// fetch resource
		resource, err := u.resourceRepo.GetResourceByID(r.ResourceID)
		if err != nil {
			return nil, apiError.NewNotFoundError(fmt.Errorf("resource not found: %w", err))
		}

		if resource.ResourcePoolID != request.PoolID {
			return nil, apiError.NewBadRequestError(errors.New("all resources must belong to the same resource pool"))
		}

		// check resource ownership
		if resource.ResourcePool.OrganizationID != namespace.OrgID {
			return nil, apiError.NewForbiddenError(errors.New("unauthorized: resource not owned by organization"))
		}

		if r.Quantity > resource.Quantity {
			return nil, apiError.NewBadRequestError(fmt.Errorf("requested quantity %d exceeds available quantity %d for resource %s", r.Quantity, resource.Quantity, r.ResourceID))
		}

	}

	namespaceQuota := &models.NamespaceQuota{
		Name:           request.Name,
		Description:    request.Description,
		Namespaces:     []models.Namespace{*namespace},
		ResourcePoolID: request.PoolID,
	}
	err = u.quotaRepo.CreateNamespaceQuota(namespaceQuota)
	if err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to create namespace quota: %w", err))
	}

	for _, r := range request.Resources {
		resourceProperty := models.ResourceProperty{
			ResourceID:  r.ResourceID,
			Price:       r.Price,
			MaxDuration: r.Duration,
		}

		err := u.quotaRepo.CreateResourceProperty(&resourceProperty)
		if err != nil {
			return nil, apiError.NewInternalServerError(fmt.Errorf("failed to create resource property: %w", err))
		}

		resourceQuantity := models.ResourceQuantity{
			NamespaceQuotaID: &namespaceQuota.ID,
			Quantity:         r.Quantity,
			ResourcePropID:   resourceProperty.ID,
		}

		err = u.quotaRepo.CreateResourceQuantity(&resourceQuantity)
		if err != nil {
			return nil, apiError.NewInternalServerError(fmt.Errorf("failed to create resource quantity: %w", err))
		}
	}

	return namespaceQuota, nil
}

func (u *PrivateNamespaceUsecase) isOrgAdmin(orgID uuid.UUID, userID uuid.UUID) (bool, error) {
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

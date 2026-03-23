package usecase

import (
	"errors"
	"fmt"

	namespaceInterfaces "github.com/ClearingHouse/internal/namespaces/interfaces"
	orgInterfaces "github.com/ClearingHouse/internal/organizations/interfaces"
	projInterfaces "github.com/ClearingHouse/internal/projects/interfaces"
	"github.com/ClearingHouse/internal/quota/dtos"
	"github.com/ClearingHouse/internal/quota/interfaces"
	resourcesInterfaces "github.com/ClearingHouse/internal/resources/interfaces"
	userInterfaces "github.com/ClearingHouse/internal/users/interfaces"
	apiError "github.com/ClearingHouse/pkg/api_error"
	helper "github.com/ClearingHouse/pkg/membership"
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

func (u *QuotaUsecase) isOrgAdmin(orgID uuid.UUID, userID uuid.UUID) error {
	org, err := u.orgRepo.GetOrganizationByID(orgID)
	if err != nil {
		return apiError.NewNotFoundError(fmt.Errorf("failed to get organization: %w", err))
	}

	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return apiError.NewNotFoundError(fmt.Errorf("failed to get user: %w", err))
	}

	if !helper.ContainsUserID(org.Admins, user.ID) {
		return apiError.NewForbiddenError(fmt.Errorf("user is not an admin of the organization"))
	}

	return nil
}

func (u *QuotaUsecase) isProjAdmin(projID uuid.UUID, userID uuid.UUID) error {
	proj, err := u.projRepo.GetProjectByID(projID)
	if err != nil {
		return apiError.NewNotFoundError(fmt.Errorf("failed to get project: %w", err))
	}

	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return apiError.NewNotFoundError(fmt.Errorf("failed to get user: %w", err))
	}

	if helper.ContainsUserID(proj.Admins, user.ID) {
		return nil
	}

	// Also allow org admin of the project's organization
	org, err := u.orgRepo.GetOrganizationByID(proj.OrganizationID)
	if err != nil {
		return apiError.NewNotFoundError(fmt.Errorf("failed to get organization: %w", err))
	}

	if helper.ContainsUserID(org.Admins, user.ID) {
		return nil
	}

	return apiError.NewForbiddenError(fmt.Errorf("user is not an admin of the project or its organization"))
}

func (u *QuotaUsecase) isProjMember(projID uuid.UUID, userID uuid.UUID) error {
	proj, err := u.projRepo.GetProjectByID(projID)
	if err != nil {
		return apiError.NewNotFoundError(fmt.Errorf("failed to get project: %w", err))
	}

	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return apiError.NewNotFoundError(fmt.Errorf("failed to get user: %w", err))
	}

	if !helper.ContainsUserID(proj.Members, user.ID) {
		return apiError.NewForbiddenError(fmt.Errorf("user is not a member of the project"))
	}

	return nil
}

func (u *QuotaUsecase) isNamespaceMember(namespaceID uuid.UUID, userID uuid.UUID) error {
	namespace, err := u.namespaceRepo.GetNamespaceByID(namespaceID)
	if err != nil {
		return apiError.NewNotFoundError(fmt.Errorf("failed to get namespace: %w", err))
	}

	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return apiError.NewNotFoundError(fmt.Errorf("failed to get user: %w", err))
	}

	// Allow both owners and members
	if (namespace.OwnerID == nil || *namespace.OwnerID != user.ID) && !helper.ContainsUserID(namespace.Members, userID) {
		return apiError.NewForbiddenError(fmt.Errorf("user is not a member or owner of the namespace"))
	}

	return nil
}

func (u *QuotaUsecase) GetUsage(quotaID uuid.UUID, namespaceID uuid.UUID, userID uuid.UUID) (interface{}, error) {
	// Check if user is a namespace member
	namespaceErr := u.isNamespaceMember(namespaceID, userID)

	// If not a namespace member, check if they're a project member or admin
	if namespaceErr != nil {
		namespace, err := u.namespaceRepo.GetNamespaceByID(namespaceID)
		if err != nil {
			return nil, apiError.NewNotFoundError(fmt.Errorf("failed to get namespace: %w", err))
		}

		// If namespace belongs to a project, check project membership or admin status
		if namespace.ProjectID != nil {
			memberErr := u.isProjMember(*namespace.ProjectID, userID)
			if memberErr != nil {
				// If not a member, check if they're an admin
				if err := u.isProjAdmin(*namespace.ProjectID, userID); err != nil {
					return nil, apiError.NewForbiddenError(fmt.Errorf("user is not a member or admin of the namespace or its project"))
				}
			}
		} else {
			// No project associated, return original namespace error
			return nil, namespaceErr
		}
	}

	isAssigned, err := u.quotaRepo.IsAssigned(namespaceID, quotaID)
	if err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to check quota assignment: %w", err))
	}
	if !isAssigned {
		return nil, apiError.NewBadRequestError(errors.New("quota is not assigned to the namespace"))
	}

	usage, err := u.quotaRepo.GetNamespaceUsageByType(namespaceID, quotaID)
	if err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to fetch usage data: %w", err))
	}

	quota, err := u.quotaRepo.GetQuotaByType(quotaID)
	if err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to fetch quota data: %w", err))
	}

	var response dtos.UsageResponse

	usageMap := make(map[uuid.UUID]float64)
	for _, resourceUsage := range usage.ResourceUsages {
		usageMap[resourceUsage.TypeID] = resourceUsage.Usage
	}

	for _, resourceQuota := range quota.ResourceQuotas {
		usageValue := usageMap[resourceQuota.TypeID]
		response.Usage = append(response.Usage, dtos.Usage{
			TypeID: resourceQuota.TypeID,
			Type:   resourceQuota.Type,
			Quota:  resourceQuota.Quota,
			Usage:  usageValue,
		})
	}

	return response, nil
}

func (u *QuotaUsecase) GetProjectQuotaTotal(projectID uuid.UUID, userID uuid.UUID) (interface{}, error) {
	if err := u.isProjMember(projectID, userID); err != nil {
		if err := u.isProjAdmin(projectID, userID); err != nil {
			return nil, apiError.NewForbiddenError(fmt.Errorf("user is not a member or admin of the project"))
		}
	}

	total, err := u.quotaRepo.GetProjectQuotaTotalByType(projectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(fmt.Errorf("failed to fetch project quota: %w", err))
	}

	return total, nil
}

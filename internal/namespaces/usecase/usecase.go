package usecase

import (
	"fmt"
	"log"

	"github.com/ClearingHouse/helper"
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/namespaces/dtos"
	"github.com/ClearingHouse/internal/namespaces/interfaces"
	projInterfaces "github.com/ClearingHouse/internal/projects/interfaces"
	quotaInterfaces "github.com/ClearingHouse/internal/quota/interfaces"
	userInterfaces "github.com/ClearingHouse/internal/users/interfaces"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/google/uuid"
)

type NamespaceUsecase struct {
	namespaceRepo interfaces.NamespaceRepository
	userRepo      userInterfaces.UsersRepository
	projRepo      projInterfaces.ProjectRepository
	quotaRepo     quotaInterfaces.QuotaRepository
}

func NewNamespaceUsecase(namespaceRepo interfaces.NamespaceRepository, userRepo userInterfaces.UsersRepository, projRepo projInterfaces.ProjectRepository, quotaRepo quotaInterfaces.QuotaRepository) interfaces.NamespaceUsecase {
	return &NamespaceUsecase{
		namespaceRepo: namespaceRepo,
		userRepo:      userRepo,
		projRepo:      projRepo,
		quotaRepo:     quotaRepo,
	}
}

func (u *NamespaceUsecase) CreateNamespace(request *dtos.CreateNamespaceRequest, userID uuid.UUID) (*models.Namespace, error) {
	proj, err := u.projRepo.GetProjectByID(request.ProjectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	if !helper.ContainsUserID(proj.Admins, userID) {
		return nil, apiError.NewUnauthorizedError("user is not project admin")
	}

	namespace := models.Namespace{
		Name:        request.Name,
		Description: request.Description,
		Credit:      request.Credit,
		ProjectID:   &request.ProjectID,
		OrgID:       proj.OrganizationID,
		OwnerID:     userID,
	}

	err = u.namespaceRepo.Create(&namespace)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	return &namespace, nil
}

func (u *NamespaceUsecase) GetAllNamespaces() ([]models.Namespace, error) {
	return u.namespaceRepo.GetAll()
}

func (u *NamespaceUsecase) AddMembers(req *dtos.AddMembersRequest, userID uuid.UUID) (*models.Namespace, error) {
	namespace, err := u.namespaceRepo.GetNamespaceByID(req.NamespaceID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	proj, err := u.projRepo.GetProjectByID(*namespace.ProjectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	if !helper.ContainsUserID(proj.Admins, userID) {
		return nil, apiError.NewUnauthorizedError("user is not project admin")
	}

	existing := make(map[uuid.UUID]struct{})
	for _, m := range namespace.Members {
		existing[m.ID] = struct{}{}
	}

	seenReq := make(map[uuid.UUID]struct{})
	for _, memberID := range req.Members {
		if !helper.ContainsUserID(proj.Members, memberID) {
			return nil, apiError.NewBadRequestError(fmt.Sprintf("user %s is not a member of the project", memberID))
		}
		if _, found := existing[memberID]; found {
			return nil, apiError.NewConflictError(fmt.Sprintf("user %s is already in the namespace", memberID))
		}
		if _, found := seenReq[memberID]; found {
			return nil, apiError.NewBadRequestError(fmt.Sprintf("duplicate user %s in request", memberID))
		}
		seenReq[memberID] = struct{}{}
		if _, err := u.userRepo.GetByID(memberID); err != nil {
			return nil, apiError.NewNotFoundError(fmt.Sprintf("user %s not found", memberID))
		}
	}

	users, err := u.userRepo.GetByIDs(req.Members)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}
	namespace.Members = append(namespace.Members, users...)

	if err := u.namespaceRepo.UpdateMembers(namespace); err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	return namespace, nil
}

func (u *NamespaceUsecase) RemoveMembers(req *dtos.RemoveMembersRequest, userID uuid.UUID) (*models.Namespace, error) {
	namespace, err := u.namespaceRepo.GetNamespaceByID(req.NamespaceID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	proj, err := u.projRepo.GetProjectByID(*namespace.ProjectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	if !helper.ContainsUserID(proj.Admins, userID) {
		return nil, apiError.NewUnauthorizedError("user is not project admin")
	}

	// Create a map of members to remove for quick lookup
	removeMap := make(map[uuid.UUID]struct{})
	seenReq := make(map[uuid.UUID]struct{})
	for _, memberID := range req.Members {
		if _, found := seenReq[memberID]; found {
			return nil, apiError.NewBadRequestError(fmt.Sprintf("duplicate user %s in request", memberID))
		}
		seenReq[memberID] = struct{}{}
		removeMap[memberID] = struct{}{}
	}

	// Check if members to remove exist in namespace and prevent removing owner
	existing := make(map[uuid.UUID]struct{})
	for _, m := range namespace.Members {
		existing[m.ID] = struct{}{}
	}

	for _, memberID := range req.Members {
		if _, found := existing[memberID]; !found {
			return nil, apiError.NewNotFoundError(fmt.Sprintf("user %s is not a member of this namespace", memberID))
		}
		// Prevent removing namespace owner
		if namespace.OwnerID == memberID {
			return nil, apiError.NewBadRequestError(fmt.Sprintf("user %s is the namespace owner and cannot be removed", memberID))
		}
	}

	// Filter out members to remove
	newMembers := []models.User{}
	for _, member := range namespace.Members {
		if _, shouldRemove := removeMap[member.ID]; !shouldRemove {
			newMembers = append(newMembers, member)
		}
	}

	namespace.Members = newMembers

	if err := u.namespaceRepo.UpdateMembers(namespace); err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	return namespace, nil
}

func (u *NamespaceUsecase) GetAllUserNamespaces(projID uuid.UUID, userID uuid.UUID) ([]models.Namespace, error) {
	namespaces, err := u.namespaceRepo.GetAllNamespacesByProjectAndUserID(projID, userID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}
	return namespaces, nil
}

func (u *NamespaceUsecase) GetNamespace(namespaceID uuid.UUID, userID uuid.UUID) (*models.Namespace, error) {
	namespace, err := u.namespaceRepo.GetNamespaceByID(namespaceID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	// Check if user is owner or member
	if namespace.OwnerID != userID && !helper.ContainsUserID(namespace.Members, userID) {
		return nil, apiError.NewUnauthorizedError("user is not namespace owner or member")
	}

	log.Printf("Namespace members: %+v", namespace.Members)
	return namespace, nil
}

func (u *NamespaceUsecase) GetNamespaceQuotas(namespaceID uuid.UUID, userID uuid.UUID) ([]models.NamespaceQuota, error) {
	namespace, err := u.namespaceRepo.GetNamespaceByID(namespaceID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}
	if !helper.ContainsUserID(namespace.Members, userID) {
		return nil, apiError.NewUnauthorizedError("user is not namespace member")
	}

	quotas, err := u.quotaRepo.GetNamespaceQuotaByNamespaceID(namespaceID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	return quotas, nil
}

func (u *NamespaceUsecase) GetNamespaceUsages(namespaceID uuid.UUID, userID uuid.UUID) (*dtos.NamespaceUsageResponse, error) {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}
	namespace, err := u.namespaceRepo.GetNamespaceByID(namespaceID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}
	if !helper.ContainsUserID(namespace.Members, userID) && (*user.NamespaceID != namespace.ID) {
		return nil, apiError.NewUnauthorizedError("user is not namespace member")
	}

	quotas, err := u.namespaceRepo.GetNamespaceQuotaByType(namespace.ID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	usages, err := u.namespaceRepo.GetNamespaceUsageByType(namespace.ID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	var namespaceUsage dtos.NamespaceUsageResponse
	usageMap := make(map[uuid.UUID]float64)
	for _, u := range usages.ResourceUsages {
		usageMap[u.TypeID] = u.Usage
	}

	for _, q := range quotas.ResourceQuotas {
		uVal := usageMap[q.TypeID]
		namespaceUsage.Usage = append(namespaceUsage.Usage, dtos.NamespaceUsage{
			TypeID: q.TypeID,
			Type:   q.Type,
			Quota:  q.Quota,
			Usage:  uVal,
		})
	}

	return &namespaceUsage, nil
}

func (u *NamespaceUsecase) GetAllPrivateNamespaces(userID uuid.UUID) ([]models.Namespace, error) {
	namespaces, err := u.namespaceRepo.GetPrivateNamespaceByUserID(userID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}
	return namespaces, nil
}

func (u *NamespaceUsecase) GetNamespacesByProjectID(projectID uuid.UUID, userID uuid.UUID) ([]models.Namespace, error) {
	// Check if project exists and user has access
	proj, err := u.projRepo.GetProjectByID(projectID)
	if err != nil {
		return nil, apiError.NewNotFoundError("project not found")
	}

	// Check if user is a member or admin of the project
	if !helper.ContainsUserID(proj.Members, userID) && !helper.ContainsUserID(proj.Admins, userID) {
		return nil, apiError.NewUnauthorizedError("user does not have access to this project")
	}

	namespaces, err := u.namespaceRepo.GetAllNamespacesByProjectID(projectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	return namespaces, nil
}

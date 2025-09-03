package usecase

import (
	"fmt"

	"github.com/ClearingHouse/helper"
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/namespaces/dtos"
	"github.com/ClearingHouse/internal/namespaces/interfaces"
	projInterfaces "github.com/ClearingHouse/internal/projects/interfaces"
	userInterfaces "github.com/ClearingHouse/internal/users/interfaces"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/google/uuid"
)

type NamespaceUsecase struct {
	namespaceRepo interfaces.NamespaceRepository
	userRepo      userInterfaces.UsersRepository
	projRepo      projInterfaces.ProjectRepository
}

func NewNamespaceUsecase(namespaceRepo interfaces.NamespaceRepository, userRepo userInterfaces.UsersRepository, projRepo projInterfaces.ProjectRepository) interfaces.NamespaceUsecase {
	return &NamespaceUsecase{
		namespaceRepo: namespaceRepo,
		userRepo:      userRepo,
		projRepo:      projRepo,
	}
}

func (u *NamespaceUsecase) CreateNamespace(request *dtos.CreateNamespaceRequest, userID uuid.UUID) (*models.Namespace, error) {
	proj, err := u.projRepo.GetProjectByID(request.ProjectID)
	if err != nil {
		return nil, apiError.NewInternalServerError(err)
	}

	if !helper.ContainsUserID(proj.Admins, request.Creator) {
		return nil, apiError.NewUnauthorizedError("user is not project admin")
	}

	namespace := models.Namespace{
		Name:        request.Name,
		Description: request.Description,
		Credit:      request.Credit,
		ProjectID:   request.ProjectID,
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

	proj, err := u.projRepo.GetProjectByID(namespace.ProjectID)
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

func (u *NamespaceUsecase) GetAllUserNamespaces(userID uuid.UUID) ([]models.Namespace, error) {
	namespaces, err := u.namespaceRepo.FindAllNamespacesByUserID(userID)
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
	if !helper.ContainsUserID(namespace.Members, userID) {
		return nil, apiError.NewUnauthorizedError("user is not namespace member")
	}

	return namespace, nil
}

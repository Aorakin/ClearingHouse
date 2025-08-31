package usecase

import (
	"fmt"

	"github.com/ClearingHouse/helper"
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/namespaces/dtos"
	"github.com/ClearingHouse/internal/namespaces/interfaces"
	projInterfaces "github.com/ClearingHouse/internal/projects/interfaces"
	userInterfaces "github.com/ClearingHouse/internal/users/interfaces"
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

func (u *NamespaceUsecase) CreateNamespace(request *dtos.CreateNamespaceRequest) (*models.Namespace, error) {
	proj, err := u.projRepo.FindProjectByID(request.ProjectID)
	if err != nil {
		return nil, err
	}

	if !helper.ContainsUserID(proj.Admins, request.Creator) {
		return nil, fmt.Errorf("only project admins can create namespace")
	}

	namespace := models.Namespace{
		Name:        request.Name,
		Description: request.Description,
		Credit:      request.Credit,
		ProjectID:   request.ProjectID,
	}

	err = u.namespaceRepo.Create(&namespace)
	return &namespace, err
}

func (u *NamespaceUsecase) GetAllNamespaces() ([]models.Namespace, error) {
	return u.namespaceRepo.GetAll()
}

func (u *NamespaceUsecase) AddMembers(req *dtos.AddMembersRequest) (*models.Namespace, error) {
	namespace, err := u.namespaceRepo.GetByID(req.NamespaceID)
	if err != nil {
		return nil, err
	}

	proj, err := u.projRepo.FindProjectByID(namespace.ProjectID)
	if err != nil {
		return nil, err
	}

	if !helper.ContainsUserID(proj.Admins, req.Creator) {
		return nil, fmt.Errorf("only project admins can add members to a namespace")
	}

	existing := make(map[uuid.UUID]struct{})
	for _, m := range namespace.Members {
		existing[m.ID] = struct{}{}
	}

	seenReq := make(map[uuid.UUID]struct{})
	for _, memberID := range req.Members {
		if !helper.ContainsUserID(proj.Members, memberID) {
			return nil, fmt.Errorf("user %s is not a member of the project", memberID)
		}
		if _, found := existing[memberID]; found {
			return nil, fmt.Errorf("user %s is already in the namespace", memberID)
		}
		if _, found := seenReq[memberID]; found {
			return nil, fmt.Errorf("duplicate user %s in request", memberID)
		}
		seenReq[memberID] = struct{}{}
		if _, err := u.userRepo.GetByID(memberID); err != nil {
			return nil, fmt.Errorf("user %s not found", memberID)
		}
	}

	users, err := u.userRepo.GetByIDs(req.Members)
	if err != nil {
		return nil, err
	}
	namespace.Members = append(namespace.Members, users...)

	if err := u.namespaceRepo.UpdateMembers(namespace); err != nil {
		return nil, err
	}

	return namespace, nil
}

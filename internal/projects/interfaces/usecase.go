package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/projects/dto"
	"github.com/google/uuid"
)

type ProjectsUsecase interface {
	GetAll(userID uuid.UUID) ([]models.Project, error)
	GetByID(userID uuid.UUID, projectID uuid.UUID) (*models.Project, error)
	Create(userID uuid.UUID, projectRequest *dto.CreateProjectRequest) (*models.Project, error)
	AddMembers(userID uuid.UUID, projectID uuid.UUID, members *dto.Members) (*models.Project, error)
	// Update(userID uuid.UUID, projectID uuid.UUID, projectRequest *dto.CreateProjectRequest) (*models.Project, error)
	// Delete(userID uuid.UUID, projectID uuid.UUID) error
	// RemoveMembers(userID uuid.UUID, projectID uuid.UUID, members *dto.Members) (*models.Project, error)
}

package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/projects/dtos"
	"github.com/google/uuid"
)

type ProjectRepository interface {
	CreateProject(project *models.Project) error
	GetAllProjects() ([]models.Project, error)
	GetProjectsByUserAdminScope(userID uuid.UUID) ([]models.Project, error)

	GetProjectByID(id uuid.UUID) (*models.Project, error)
	UpdateProject(project *models.Project) error
	UpdateMembers(project *models.Project) error
	UpdateAdmins(project *models.Project) error
	DeleteProject(id uuid.UUID) error
	HasProjectQuotas(projectID uuid.UUID) (bool, error)

	GetAllProjectsByUserID(userID uuid.UUID) ([]models.Project, error)
	GetProjectsByOrganizationID(orgID uuid.UUID) ([]models.Project, error)
	GetProjectMembers(projectID uuid.UUID) ([]models.User, error)
	GetProjectQuotaByType(projectID uuid.UUID, userID uuid.UUID) (*dtos.ResourceQuotaResponse, error)
	GetProjectUsageByType(projectID uuid.UUID, userID uuid.UUID) (*dtos.ResourceUsageResponse, error)
}

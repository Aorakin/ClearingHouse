package interfaces

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/projects/dtos"
)

type ProjectUsecase interface {
	CreateProject(request *dtos.CreateProjectRequest) error
	GetAllProjects() ([]models.Project, error)
}

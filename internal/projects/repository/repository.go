package repository

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/projects/interfaces"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) interfaces.ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) CreateProject(project *models.Project) error {
	return r.db.Create(project).Error
}

func (r *ProjectRepository) GetAllProjects() ([]models.Project, error) {
	var projects []models.Project
	err := r.db.Find(&projects).Error
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *ProjectRepository) GetProjectByID(id uuid.UUID) (*models.Project, error) {
	var project models.Project
	err := r.db.Preload("Admins").Preload("Members").First(&project, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) UpdateMembers(project *models.Project) error {
	return r.db.Model(project).Association("Members").Replace(project.Members)
}

func (r *ProjectRepository) GetAllProjectsByUserID(userID uuid.UUID) ([]models.Project, error) {
	var user models.User
	if err := r.db.Preload("MemberProjects").First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return user.MemberProjects, nil
}

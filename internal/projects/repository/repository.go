package repository

import (
	"errors"

	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/projects/interfaces"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectsRepository struct {
	db *gorm.DB
}

func NewProjectsRepository(db *gorm.DB) interfaces.ProjectsRepository {
	return &ProjectsRepository{db: db}
}

func (r *ProjectsRepository) GetAll(userID uuid.UUID) ([]models.Project, error) {
	var projects []models.Project

	if err := r.db.
		Joins("JOIN project_owners ON projects.id = project_owners.project_id").
		Where("project_owners.user_id = ?", userID).
		Find(&projects).Error; err != nil {
		return nil, err
	}

	return projects, nil

}

func (r *ProjectsRepository) GetByID(id uuid.UUID) (*models.Project, error) {
	var project models.Project
	if err := r.db.Preload("Owners").Preload("Members").First(&project, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectsRepository) Create(project *models.Project) (*models.Project, error) {
	if err := r.db.Create(project).Error; err != nil {
		return nil, err
	}
	return project, nil
}

func (r *ProjectsRepository) Update(id uuid.UUID, project *models.Project) (*models.Project, error) {
	var existingProject models.Project
	if err := r.db.First(&existingProject, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("project not found")
		}
		return nil, err
	}

	project.ID = id
	if err := r.db.Model(&existingProject).Updates(project).Error; err != nil {
		return nil, err
	}

	return &existingProject, nil
}

func (r *ProjectsRepository) Delete(id uuid.UUID) error {
	if err := r.db.Delete(&models.Project{}, "id = ?", id).Error; err != nil {
		return err
	}
	if r.db.RowsAffected == 0 {
		return errors.New("announcement not found")
	}
	return nil
}

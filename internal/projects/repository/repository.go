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

type ResourceUsage struct {
	TypeID string  `json:"type_id"`
	Type   string  `json:"type"`
	Usage  float64 `json:"usage"`
}

type ResourceUsageResponse struct {
	ResourceUsage []ResourceUsage `json:"resource_usage"`
}

func (r *ProjectRepository) GetProjectQuotaByType(projectID uuid.UUID, userID uuid.UUID) (interface{}, error) {
	var quotas []models.NamespaceQuota

	// Fetch quotas with resources and resource types
	if err := r.db.
		Joins("JOIN namespace_quotas nqj ON nqj.namespace_quota_id = namespace_quota.id").
		Joins("JOIN namespaces ns ON ns.id = nqj.namespace_id").
		Joins("JOIN namespace_members nm ON nm.namespace_id = ns.id").
		Where("ns.project_id = ? AND nm.user_id = ?", projectID, userID).
		Preload("Resources.ResourceProp.Resource.ResourceType").
		Find(&quotas).Error; err != nil {
		return nil, err
	}

	// Aggregate usage by type
	typeAgg := make(map[string]ResourceUsage)
	for _, quota := range quotas {
		for _, res := range quota.Resources {
			rt := res.ResourceProp.Resource.ResourceType
			rtID := rt.ID.String()

			if _, ok := typeAgg[rtID]; !ok {
				typeAgg[rtID] = ResourceUsage{
					TypeID: rtID,
					Type:   rt.Name,
					Usage:  0,
				}
			}
			tmp := typeAgg[rtID]
			tmp.Usage += float64(res.Quantity)
			typeAgg[rtID] = tmp
		}
	}

	// Convert to slice
	var result []ResourceUsage
	for _, v := range typeAgg {
		result = append(result, v)
	}

	return &ResourceUsageResponse{ResourceUsage: result}, nil
}

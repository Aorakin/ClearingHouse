package repository

import (
	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/projects/dtos"
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

func (r *ProjectRepository) GetProjectQuotaByType(projectID uuid.UUID, userID uuid.UUID) (*dtos.ResourceQuotaResponse, error) {
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
	typeAgg := make(map[string]dtos.ResourceQuota)
	for _, quota := range quotas {
		for _, res := range quota.Resources {
			rt := res.ResourceProp.Resource.ResourceType
			rtID := rt.ID.String()

			if _, ok := typeAgg[rtID]; !ok {
				typeAgg[rtID] = dtos.ResourceQuota{
					TypeID: rtID,
					Type:   rt.Name,
					Quota:  0,
				}
			}
			tmp := typeAgg[rtID]
			tmp.Quota += float64(res.Quantity)
			typeAgg[rtID] = tmp
		}
	}

	var result []dtos.ResourceQuota
	for _, v := range typeAgg {
		result = append(result, v)
	}

	return &dtos.ResourceQuotaResponse{ResourceQuotas: result}, nil
}

func (r *ProjectRepository) GetProjectUsageByType(projectID uuid.UUID, userID uuid.UUID) (*dtos.ResourceUsageResponse, error) {
	var tickets []models.Ticket

	if err := r.db.
		Joins("JOIN namespaces ns ON ns.id = tickets.namespace_id").
		Joins("JOIN namespace_members nm ON nm.namespace_id = ns.id").
		Where("ns.project_id = ? AND nm.user_id = ?", projectID, userID).
		Preload("Resources.Resource.ResourceType").
		Find(&tickets).Error; err != nil {
		return nil, err
	}

	typeAgg := make(map[string]dtos.ResourceUsage)
	for _, t := range tickets {
		for _, tr := range t.Resources {
			rt := tr.Resource.ResourceType
			rtID := rt.ID.String()

			if _, ok := typeAgg[rtID]; !ok {
				typeAgg[rtID] = dtos.ResourceUsage{
					TypeID: rtID,
					Type:   rt.Name,
					Usage:  0,
				}
			}
			tmp := typeAgg[rtID]
			tmp.Usage += float64(tr.Quantity)
			typeAgg[rtID] = tmp
		}
	}

	var result []dtos.ResourceUsage
	for _, v := range typeAgg {
		result = append(result, v)
	}

	return &dtos.ResourceUsageResponse{ResourceUsages: result}, nil
}

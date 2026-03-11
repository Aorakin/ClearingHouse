package repository

import (
	"log"
	"sort"

	"github.com/ClearingHouse/internal/models"
	"github.com/ClearingHouse/internal/projects/dtos"
	"github.com/ClearingHouse/internal/projects/interfaces"
	"github.com/ClearingHouse/pkg/enum"
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
	err := r.db.
		Preload("Admins", func(db *gorm.DB) *gorm.DB { return db.Order("email") }).
		Preload("Members", func(db *gorm.DB) *gorm.DB { return db.Order("email") }).
		Order("name").Find(&projects).Error
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *ProjectRepository) GetProjectsByUserAdminScope(userID uuid.UUID) ([]models.Project, error) {
	var projects []models.Project

	subqueryOrgAdmin := r.db.Table("organization_admins").
		Select("organization_id").
		Where("user_id = ?", userID)

	subqueryProjAdmin := r.db.Table("project_admins").
		Joins("JOIN projects p ON p.id = project_admins.project_id").
		Select("DISTINCT p.organization_id").
		Where("project_admins.user_id = ? AND p.deleted_at IS NULL", userID)

	if err := r.db.
		Preload("Admins", func(db *gorm.DB) *gorm.DB { return db.Order("email") }).
		Preload("Members", func(db *gorm.DB) *gorm.DB { return db.Order("email") }).
		Where("organization_id IN (?) OR organization_id IN (?)", subqueryOrgAdmin, subqueryProjAdmin).
		Order("name").
		Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *ProjectRepository) GetProjectByID(id uuid.UUID) (*models.Project, error) {
	var project models.Project
	err := r.db.
		Preload("Admins", func(db *gorm.DB) *gorm.DB { return db.Order("email") }).
		Preload("Members", func(db *gorm.DB) *gorm.DB { return db.Order("email") }).
		First(&project, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) UpdateProject(project *models.Project) error {
	return r.db.Model(project).Updates(map[string]interface{}{
		"name":        project.Name,
		"description": project.Description,
	}).Error
}

func (r *ProjectRepository) UpdateMembers(project *models.Project) error {
	return r.db.Model(project).Association("Members").Replace(project.Members)
}

func (r *ProjectRepository) UpdateAdmins(project *models.Project) error {
	return r.db.Model(project).Association("Admins").Replace(project.Admins)
}

func (r *ProjectRepository) DeleteProject(id uuid.UUID) error {
	return r.db.Delete(&models.Project{}, "id = ?", id).Error
}

func (r *ProjectRepository) HasProjectQuotas(projectID uuid.UUID) (bool, error) {
	var count int64

	err := r.db.Table("project_quota").
		Where("project_id = ? AND deleted_at IS NULL", projectID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *ProjectRepository) GetAllProjectsByUserID(userID uuid.UUID) ([]models.Project, error) {
	var user models.User
	if err := r.db.Preload("MemberProjects", func(db *gorm.DB) *gorm.DB { return db.Order("name") }).First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return user.MemberProjects, nil
}

func (r *ProjectRepository) GetProjectsByOrganizationID(orgID uuid.UUID) ([]models.Project, error) {
	var projects []models.Project
	err := r.db.
		Preload("Admins", func(db *gorm.DB) *gorm.DB { return db.Order("email") }).
		Preload("Members", func(db *gorm.DB) *gorm.DB { return db.Order("email") }).
		Where("organization_id = ?", orgID).
		Order("name").
		Find(&projects).Error
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *ProjectRepository) GetProjectMembers(projectID uuid.UUID) ([]models.User, error) {
	var project models.Project
	if err := r.db.Preload("Members", func(db *gorm.DB) *gorm.DB { return db.Order("email") }).First(&project, "id = ?", projectID).Error; err != nil {
		return nil, err
	}
	return project.Members, nil
}

func (r *ProjectRepository) GetProjectQuotaByType(projectID uuid.UUID, userID uuid.UUID) (*dtos.ResourceQuotaResponse, error) {
	var quotas []models.NamespaceQuota

	// Fetch quotas with resources and resource types
	if err := r.db.
		Joins("JOIN namespace_quota_template_relations nqtr ON nqtr.namespace_quota_id = namespace_quota.id").
		Joins("JOIN namespace_quota_templates nqt ON nqt.id = nqtr.namespace_quota_template_id").
		Joins("JOIN namespaces ns ON ns.quota_template_id = nqt.id").
		Joins("JOIN namespace_members nm ON nm.namespace_id = ns.id").
		Where("ns.project_id = ? AND nm.user_id = ?", projectID, userID).
		Preload("Resources.ResourceProp.Resource.ResourceType").
		Find(&quotas).Error; err != nil {
		return nil, err
	}

	log.Println(len(quotas))
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

	sort.Slice(result, func(i, j int) bool {
		return enum.ResourceTypeOrder(result[i].Type) < enum.ResourceTypeOrder(result[j].Type)
	})

	return &dtos.ResourceQuotaResponse{ResourceQuotas: result}, nil
}

func (r *ProjectRepository) GetProjectUsageByType(projectID uuid.UUID, userID uuid.UUID) (*dtos.ResourceUsageResponse, error) {
	var tickets []models.Ticket

	if err := r.db.
		Joins("JOIN namespaces ns ON ns.id = tickets.namespace_id").
		Joins("JOIN namespace_members nm ON nm.namespace_id = ns.id").
		Where("ns.project_id = ? AND nm.user_id = ? AND status IN ?", projectID, userID, enum.UsingStatuses).
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

	sort.Slice(result, func(i, j int) bool {
		return enum.ResourceTypeOrder(result[i].Type) < enum.ResourceTypeOrder(result[j].Type)
	})

	return &dtos.ResourceUsageResponse{ResourceUsages: result}, nil
}

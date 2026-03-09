package http

import (
	"net/http"

	"github.com/ClearingHouse/internal/quota/dtos"
	"github.com/ClearingHouse/internal/quota/interfaces"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/ClearingHouse/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type QuotaHandler struct {
	quotaUsecase interfaces.QuotaUsecase
}

func NewQuotaHandler(quotaUsecase interfaces.QuotaUsecase) interfaces.QuotaHandler {
	return &QuotaHandler{
		quotaUsecase: quotaUsecase,
	}
}

func (h *QuotaHandler) CreateOrganizationQuota() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		var request dtos.CreateOrganizationQuotaRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		quota, err := h.quotaUsecase.CreateOrganizationQuota(&request, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusCreated, quota)
	}
}

func (h *QuotaHandler) GetOrganizationQuota() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dtos.FindOrganizationQuotaGroupRequest
		if err := c.ShouldBindQuery(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		fromUUID, err := uuid.Parse(request.FromOrganizationID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}
		toUUID, err := uuid.Parse(request.ToOrganizationID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		quotas, err := h.quotaUsecase.GetOrganizationQuota(fromUUID, toUUID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, quotas)
	}
}

func (h *QuotaHandler) GetOrganizationQuotasByOrgID() gin.HandlerFunc {
	return func(c *gin.Context) {
		orgID := c.Param("org_id")
		orgUUID, err := uuid.Parse(orgID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		quotas, err := h.quotaUsecase.GetOrganizationQuotasByOrgID(orgUUID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, quotas)
	}
}

func (h *QuotaHandler) DeleteOrganizationQuota() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		quotaID := c.Param("quota_id")
		if quotaID == "" {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("Quota ID is required")))
			return
		}

		quotaUUID, err := uuid.Parse(quotaID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		if err := h.quotaUsecase.DeleteOrganizationQuota(quotaUUID, userID); err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Organization quota deleted successfully"})
	}
}

func (h *QuotaHandler) CreateProjectQuota() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		var request dtos.CreateProjectQuotaRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		quota, err := h.quotaUsecase.CreateProjectQuota(&request, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusCreated, quota)
	}
}

func (h *QuotaHandler) GetProjectQuotas() gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("project_id")
		if projectID == "" {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("Project ID is required")))
			return
		}

		projectUUID, err := uuid.Parse(projectID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		quotas, err := h.quotaUsecase.GetProjectQuotas(projectUUID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, quotas)
	}
}

func (h *QuotaHandler) CreateNamespaceQuota() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		var request dtos.CreateNamespaceQuotaRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		quota, err := h.quotaUsecase.CreateNamespaceQuota(&request, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusCreated, quota)
	}
}

func (h *QuotaHandler) GetNamespaceQuota() gin.HandlerFunc {
	return func(c *gin.Context) {
		namespaceID := c.Param("namespace_id")
		if namespaceID == "" {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("Namespace ID is required")))
			return
		}

		namespaceUUID, err := uuid.Parse(namespaceID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		quota, err := h.quotaUsecase.GetNamespaceQuota(namespaceUUID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, quota)
	}
}

func (h *QuotaHandler) UpdateNamespaceQuota() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		quotaID := c.Param("quota_id")
		if quotaID == "" {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("Quota ID is required")))
			return
		}

		quotaUUID, err := uuid.Parse(quotaID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		var request dtos.UpdateNamespaceQuotaRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		quota, err := h.quotaUsecase.UpdateNamespaceQuota(quotaUUID, &request, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, quota)
	}
}

func (h *QuotaHandler) CreateNamespaceQuotaTemplate() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		var request dtos.CreateNamespaceQuotaTemplateRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		template, err := h.quotaUsecase.CreateNamespaceQuotaTemplate(&request, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusCreated, template)
	}
}

func (h *QuotaHandler) GetNamespaceQuotaTemplate() gin.HandlerFunc {
	return func(c *gin.Context) {
		quotaTemplateID := c.Param("quota_template_id")
		if quotaTemplateID == "" {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("Quota Template ID is required")))
			return
		}

		quotaTemplateUUID, err := uuid.Parse(quotaTemplateID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		template, err := h.quotaUsecase.GetNamespaceQuotaTemplate(quotaTemplateUUID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, template)
	}
}

func (h *QuotaHandler) UpdateNamespaceQuotaTemplate() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		quotaTemplateID := c.Param("quota_template_id")
		if quotaTemplateID == "" {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("Quota Template ID is required")))
			return
		}

		quotaTemplateUUID, err := uuid.Parse(quotaTemplateID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		var request dtos.UpdateNamespaceQuotaTemplateRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		template, err := h.quotaUsecase.UpdateNamespaceQuotaTemplate(quotaTemplateUUID, &request, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, template)
	}
}

func (h *QuotaHandler) GetNamespaceQuotaTemplatesByProjectID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		projectID := c.Param("project_id")
		if projectID == "" {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("Project ID is required")))
			return
		}

		projectUUID, err := uuid.Parse(projectID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		templates, err := h.quotaUsecase.GetNamespaceQuotaTemplatesByProjectID(projectUUID, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, templates)
	}
}

func (h *QuotaHandler) AssignQuotaTemplateToNamespace() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		var request dtos.AssignQuotaToNamespaceRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		if err := h.quotaUsecase.AssignQuotaTemplateToNamespace(&request, userID); err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

func (h *QuotaHandler) CreateInternalProjectQuota() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		var request dtos.CreateInternalProjectQuotaRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		quota, err := h.quotaUsecase.CreateInternalProjectQuota(&request, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusCreated, quota)
	}
}

func (h *QuotaHandler) GetUsage() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		quotaID := c.Param("quota_id")
		if quotaID == "" {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("Quota ID is required")))
			return
		}

		namespaceID := c.Param("namespace_id")
		if namespaceID == "" {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("Namespace ID is required")))
			return
		}

		namespaceUUID, err := uuid.Parse(namespaceID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}
		quotaUUID, err := uuid.Parse(quotaID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		usage, err := h.quotaUsecase.GetUsage(quotaUUID, namespaceUUID, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, usage)
	}
}

func (h *QuotaHandler) GetProjectQuotaTotal() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		projectID := c.Param("project_id")
		if projectID == "" {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("Project ID is required")))
			return
		}

		projectUUID, err := uuid.Parse(projectID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		total, err := h.quotaUsecase.GetProjectQuotaTotal(projectUUID, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, total)
	}
}

func (h *QuotaHandler) GetNamespaceQuotaInProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		projectID := c.Param("project_id")
		if projectID == "" {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("Project ID is required")))
			return
		}

		projectUUID, err := uuid.Parse(projectID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		quotas, err := h.quotaUsecase.GetNamespaceQuotaInProject(userID, projectUUID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, quotas)
	}
}

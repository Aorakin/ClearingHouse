package http

import (
	"net/http"

	"github.com/ClearingHouse/internal/organizations/dtos"
	"github.com/ClearingHouse/internal/organizations/interfaces"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/ClearingHouse/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrganizationHandler struct {
	organizationUsecase interfaces.OrganizationUsecase
}

func NewOrganizationHandler(organizationUsecase interfaces.OrganizationUsecase) interfaces.OrganizationHandler {
	return &OrganizationHandler{
		organizationUsecase: organizationUsecase,
	}
}

func (h *OrganizationHandler) CreateOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		var dto dtos.CreateOrganization
		if err := c.ShouldBindJSON(&dto); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		org, err := h.organizationUsecase.CreateOrganization(&dto, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusCreated, org)
	}
}

func (h *OrganizationHandler) GetAllOrganizations() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		isSuperAdmin := c.MustGet("isSuperAdmin").(bool)

		orgs, err := h.organizationUsecase.GetAllOrganizations(userID, isSuperAdmin)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, orgs)
	}
}

func (h *OrganizationHandler) GetOrganizationByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		isSuperAdmin := c.MustGet("isSuperAdmin").(bool)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		var uri dtos.OrganizationURI
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		orgID, err := uuid.Parse(uri.OrgID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		org, uErr := h.organizationUsecase.GetOrganizationByID(orgID, userID, isSuperAdmin)
		if uErr != nil {
			c.JSON(response.ErrorResponseBuilder(uErr))
			return
		}

		c.JSON(http.StatusOK, org)
	}
}

func (h *OrganizationHandler) AddMembers() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		var request dtos.AddMembersRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		org, err := h.organizationUsecase.AddMembers(&request, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, org)
	}
}

func (h *OrganizationHandler) RemoveMembers() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		var request dtos.RemoveMembersRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		org, err := h.organizationUsecase.RemoveMembers(&request, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, org)
	}
}

func (h *OrganizationHandler) DeleteOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		var uri dtos.OrganizationURI
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		orgID, err := uuid.Parse(uri.OrgID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		if err := h.organizationUsecase.DeleteOrganization(orgID, userID); err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "organization deleted successfully"})
	}
}

func (h *OrganizationHandler) UpdateOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		var uri dtos.OrganizationURI
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		orgID, err := uuid.Parse(uri.OrgID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		var dto dtos.UpdateOrganization
		if err := c.ShouldBindJSON(&dto); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		org, uErr := h.organizationUsecase.UpdateOrganization(orgID, &dto, userID)
		if uErr != nil {
			c.JSON(response.ErrorResponseBuilder(uErr))
			return
		}

		c.JSON(http.StatusOK, org)
	}
}

func (h *OrganizationHandler) GetMembers() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}
		members, err := h.organizationUsecase.GetMembers()
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, members)
	}
}
func (h *OrganizationHandler) GetOrganizationMembers() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		isSuperAdmin := c.MustGet("isSuperAdmin").(bool)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}
		var uri dtos.OrganizationURI
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}
		orgID, err := uuid.Parse(uri.OrgID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}
		members, err := h.organizationUsecase.GetOrganizationMembers(orgID, userID, isSuperAdmin)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}
		c.JSON(http.StatusOK, members)
	}
}

func (h *OrganizationHandler) AddAdmins() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		isSuperAdmin := c.MustGet("isSuperAdmin").(bool)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		var request dtos.AddAdminsRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		org, err := h.organizationUsecase.AddAdmins(&request, userID, isSuperAdmin)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, org)
	}
}

func (h *OrganizationHandler) RemoveAdmins() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		isSuperAdmin := c.MustGet("isSuperAdmin").(bool)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		var request dtos.RemoveAdminsRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		org, err := h.organizationUsecase.RemoveAdmins(&request, userID, isSuperAdmin)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, org)
	}
}

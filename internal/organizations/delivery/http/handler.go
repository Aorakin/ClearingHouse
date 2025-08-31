package http

import (
	"net/http"

	"github.com/ClearingHouse/internal/organizations/dtos"
	"github.com/ClearingHouse/internal/organizations/interfaces"
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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var dto dtos.CreateOrganization
		if err := c.ShouldBindJSON(&dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		dto.Creator = userID

		org, err := h.organizationUsecase.CreateOrganization(&dto)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, org)
	}
}
func (h *OrganizationHandler) GetOrganizationByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var uri dtos.OrganizationURI
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		orgID, err := uuid.Parse(uri.OrgID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		org, err := h.organizationUsecase.GetOrganizationByID(orgID, userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
			return
		}

		c.JSON(http.StatusOK, org)
	}
}
func (h *OrganizationHandler) UpdateOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		var uri dtos.OrganizationURI
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		orgID, err := uuid.Parse(uri.OrgID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var body dtos.UpdateOrganization
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		org, err := h.organizationUsecase.UpdateOrganization(orgID, &body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, org)
	}
}
func (h *OrganizationHandler) DeleteOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		var uri dtos.OrganizationURI
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		orgID, err := uuid.Parse(uri.OrgID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := h.organizationUsecase.DeleteOrganization(orgID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

func (h *OrganizationHandler) GetOrganizations() gin.HandlerFunc {
	return func(c *gin.Context) {
		orgs, err := h.organizationUsecase.GetOrganizations()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, orgs)
	}
}

func (h *OrganizationHandler) AddMembers() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var request dtos.AddMembersRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		request.Creator = userID

		org, err := h.organizationUsecase.AddMembers(&request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, org)
	}
}

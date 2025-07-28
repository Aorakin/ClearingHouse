package http

import (
	"log"
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
		var dto dtos.CreateOrganization
		if err := c.ShouldBindJSON(&dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

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
		var uri dtos.OrganizationURI
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		orgID, err := uuid.Parse(uri.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		org, err := h.organizationUsecase.GetOrganizationByID(orgID)
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

		orgID, err := uuid.Parse(uri.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 2. Bind JSON body
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

		orgID, err := uuid.Parse(uri.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Printf("%#v", uri)

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

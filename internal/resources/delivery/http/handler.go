package http

import (
	"net/http"

	"github.com/ClearingHouse/internal/resources/dtos"
	"github.com/ClearingHouse/internal/resources/interfaces"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/ClearingHouse/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ResourceHandler struct {
	ResourceUsecase interfaces.ResourceUsecase
}

func NewResourceHandler(resourceUsecase interfaces.ResourceUsecase) interfaces.ResourceHandler {
	return &ResourceHandler{
		ResourceUsecase: resourceUsecase,
	}
}

func (h *ResourceHandler) GetResource() gin.HandlerFunc {
	return func(c *gin.Context) {
		var uri dtos.IDUri
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		orgID, err := uuid.Parse(uri.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resource, err := h.ResourceUsecase.GetResources(orgID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resource)
	}
}

func (h *ResourceHandler) CreateResourcePool() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dtos.CreateResourcePoolRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resourcePool, err := h.ResourceUsecase.CreateResourcePool(&request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, resourcePool)
	}
}

func (h *ResourceHandler) CreateResourceType() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dtos.CreateResourceTypeRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resourceType, err := h.ResourceUsecase.CreateResourceType(&request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, resourceType)
	}
}

func (h *ResourceHandler) CreateResource() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dtos.CreateResourceRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resource, err := h.ResourceUsecase.CreateResource(&request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, resource)
	}
}

func (h *ResourceHandler) GetResourceTypes() gin.HandlerFunc {
	return func(c *gin.Context) {
		resourceTypes, err := h.ResourceUsecase.GetResourceTypes()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resourceTypes)
	}
}

func (h *ResourceHandler) UpdateResource() gin.HandlerFunc {
	return func(c *gin.Context) {
		var uri dtos.IDUri
		if err := c.ShouldBindUri(&uri); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resourceID, err := uuid.Parse(uri.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var request dtos.UpdateResourceRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resource, err := h.ResourceUsecase.UpdateResource(resourceID, &request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resource)
	}
}

func (h *ResourceHandler) GetResourceProperty() gin.HandlerFunc {
	return func(c *gin.Context) {
		resourceID := c.Param("id")
		if resourceID == "" {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("Resource ID is required")))
			return
		}

		resourceUUID, err := uuid.Parse(resourceID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		resource, err := h.ResourceUsecase.GetResourceProperty(resourceUUID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}
		c.JSON(http.StatusOK, resource)
	}
}

func (h *ResourceHandler) GetResourcePool() gin.HandlerFunc {
	return func(c *gin.Context) {
		resourcePoolID := c.Param("id")
		if resourcePoolID == "" {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("Resource Pool ID is required")))
			return
		}

		resourcePoolUUID, err := uuid.Parse(resourcePoolID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		resourcePool, err := h.ResourceUsecase.GetResourcePool(&resourcePoolUUID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}
		c.JSON(http.StatusOK, resourcePool)
	}
}

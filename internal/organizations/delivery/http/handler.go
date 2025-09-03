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
		orgs, err := h.organizationUsecase.GetAllOrganizations()
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

		org, uErr := h.organizationUsecase.GetOrganizationByID(orgID, userID)
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

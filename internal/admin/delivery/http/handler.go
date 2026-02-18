package http

import (
	"net/http"
	// "os"

	"github.com/ClearingHouse/internal/admin/dtos"
	"github.com/ClearingHouse/internal/admin/interfaces"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/ClearingHouse/pkg/response"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminUsecase interfaces.AdminUsecase
}

func NewAdminHandler(adminUsecase interfaces.AdminUsecase) interfaces.AdminHandler {
	return &AdminHandler{
		adminUsecase: adminUsecase,
	}
}

func (h *AdminHandler) AssignSuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verify admin secret key
		// adminSecret := c.GetHeader("X-Admin-Secret")
		// expectedSecret := os.Getenv("ADMIN_SECRET")

		// if adminSecret == "" || expectedSecret == "" || adminSecret != expectedSecret {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized: invalid admin secret"})
		// 	return
		// }

		var dto dtos.AssignSuperAdminRequest
		if err := c.ShouldBindJSON(&dto); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		err := h.adminUsecase.AssignSuperAdmin(dto.Email)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		c.JSON(http.StatusOK, dtos.SuperAdminResponse{
			Message: "super admin assigned successfully",
			Email:   dto.Email,
		})
	}
}

func (h *AdminHandler) RevokeSuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verify admin secret key
		// adminSecret := c.GetHeader("X-Admin-Secret")
		// expectedSecret := os.Getenv("ADMIN_SECRET")

		// if adminSecret == "" || expectedSecret == "" || adminSecret != expectedSecret {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized: invalid admin secret"})
		// 	return
		// }

		var dto dtos.RevokeSuperAdminRequest
		if err := c.ShouldBindJSON(&dto); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		err := h.adminUsecase.RevokeSuperAdmin(dto.Email)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		c.JSON(http.StatusOK, dtos.SuperAdminResponse{
			Message: "super admin revoked successfully",
			Email:   dto.Email,
		})
	}
}

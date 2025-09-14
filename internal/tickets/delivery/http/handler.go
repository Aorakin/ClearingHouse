package http

import (
	"net/http"

	"github.com/ClearingHouse/internal/tickets/dtos"
	"github.com/ClearingHouse/internal/tickets/interfaces"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/ClearingHouse/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TicketHandler struct {
	ticketUsecase interfaces.TicketUsecase
}

func NewTicketHandler(ticketUsecase interfaces.TicketUsecase) interfaces.TicketHandler {
	return &TicketHandler{
		ticketUsecase: ticketUsecase,
	}
}

func (h *TicketHandler) CreateTicket() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}
		var request dtos.CreateTicketRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		ticket, err := h.ticketUsecase.CreateTicket(&request, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusCreated, ticket)
	}
}

func (h *TicketHandler) GetNamespaceTickets() gin.HandlerFunc {
	return func(c *gin.Context) {
		namespaceIDParam := c.Param("namespace_id")
		namespaceID, err := uuid.Parse(namespaceIDParam)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("invalid namespace ID")))
			return
		}

		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		tickets, err := h.ticketUsecase.GetNamespaceTickets(namespaceID, userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, tickets)
	}
}

func (h *TicketHandler) StartTicket() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dtos.StartTicketsRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		tickets, err := h.ticketUsecase.StartTicket(&request)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, tickets)
	}
}

func (h *TicketHandler) StopTicket() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dtos.StopTicketsRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err)))
			return
		}

		tickets, err := h.ticketUsecase.StopTicket(&request)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, tickets)
	}
}

func (h *TicketHandler) GetUserTickets() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		tickets, err := h.ticketUsecase.GetUserTickets(userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(err))
			return
		}

		c.JSON(http.StatusOK, tickets)
	}
}

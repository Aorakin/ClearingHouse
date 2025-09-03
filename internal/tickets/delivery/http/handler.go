package http

import (
	"net/http"

	"github.com/ClearingHouse/internal/tickets/dtos"
	"github.com/ClearingHouse/internal/tickets/interfaces"
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
			c.JSON(400, gin.H{"error": "unauthorized"})
			return
		}
		var request dtos.CreateTicketRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		request.Creator = userID

		ticket, err := h.ticketUsecase.CreateTicket(&request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, ticket)
	}
}

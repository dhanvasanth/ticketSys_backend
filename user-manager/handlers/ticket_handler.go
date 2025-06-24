package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"ticket-service/models"
	"ticket-service/services"
)

type TicketHandler struct {
	ticketService *services.TicketService
}

func NewTicketHandler() *TicketHandler {
	return &TicketHandler{
		ticketService: services.NewTicketService(),
	}
}

func (h *TicketHandler) GetTickets(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	
	filters := map[string]string{
		"status":          c.Query("status"),
		"priority_id":     c.Query("priority_id"),
		"assignee_id":     c.Query("assignee_id"),
		"organization_id": c.Query("organization_id"),
	}
	
	tickets, total, err := h.ticketService.GetTickets(page, limit, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	response := models.PaginatedResponse{
		Data:  tickets,
		Total: total,
		Page:  page,
		Limit: limit,
	}
	
	c.JSON(http.StatusOK, response)
}

func (h *TicketHandler) GetTicket(c *gin.Context) {
	id := c.Param("id")
	
	ticket, comments, err := h.ticketService.GetTicketByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	
	response := models.TicketResponse{
		Ticket:   ticket,
		Comments: comments,
	}
	
	c.JSON(http.StatusOK, response)
}

func (h *TicketHandler) CreateTicket(c *gin.Context) {
	var req models.CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	ticket, err := h.ticketService.CreateTicket(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"ticket":        ticket,
		"ticket_number": ticket.TicketNumber,
	})
}

func (h *TicketHandler) UpdateTicket(c *gin.Context) {
	id := c.Param("id")
	
	var req models.UpdateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	err := h.ticketService.UpdateTicket(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Ticket updated successfully"})
}

func (h *TicketHandler) DeleteTicket(c *gin.Context) {
	id := c.Param("id")
	
	err := h.ticketService.DeleteTicket(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Ticket deleted successfully"})
}

func (h *TicketHandler) CreateComment(c *gin.Context) {
	ticketID := c.Param("ticket_id")
	
	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	comment, err := h.ticketService.CreateComment(ticketID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"comment": comment})
}

func (h *TicketHandler) GetDashboardStats(c *gin.Context) {
	organizationID := c.Query("organization_id")
	
	stats, err := h.ticketService.GetDashboardStats(organizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, stats)
}


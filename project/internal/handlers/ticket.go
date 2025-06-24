package handlers

import (
    "net/http"
    "strconv"
    "ticket-service/internal/models"
    "ticket-service/internal/services"
    
    "github.com/gin-gonic/gin"
)

type TicketHandler struct {
    ticketService services.TicketService
}

func NewTicketHandler(ticketService services.TicketService) *TicketHandler {
    return &TicketHandler{ticketService: ticketService}
}

func (h *TicketHandler) CreateTicket(c *gin.Context) {
    userID := c.GetUint("user_id")
    
    var req models.CreateTicketRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    ticket, err := h.ticketService.CreateTicket(userID, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, gin.H{
        "message": "Ticket created successfully",
        "ticket":  ticket,
    })
}

func (h *TicketHandler) GetTickets(c *gin.Context) {
    userID := c.GetUint("user_id")
    roleID := c.GetUint("role_id")
    
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    
    tickets, total, err := h.ticketService.GetTickets(userID, roleID, page, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "tickets": tickets,
        "total":   total,
        "page":    page,
        "limit":   limit,
    })
}

func (h *TicketHandler) GetTicket(c *gin.Context) {
    userID := c.GetUint("user_id")
    roleID := c.GetUint("role_id")
    
    ticketID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
        return
    }
    
    ticket, err := h.ticketService.GetTicket(uint(ticketID), userID, roleID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"ticket": ticket})
}

func (h *TicketHandler) UpdateTicket(c *gin.Context) {
    userID := c.GetUint("user_id")
    
    ticketID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
        return
    }
    
    var req models.UpdateTicketRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    ticket, err := h.ticketService.UpdateTicket(uint(ticketID), userID, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message": "Ticket updated successfully",
        "ticket":  ticket,
    })
}

func (h *TicketHandler) AddComment(c *gin.Context) {
    userID := c.GetUint("user_id")
    
    ticketID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
        return
    }
    
    var req models.CreateCommentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    comment, err := h.ticketService.AddComment(uint(ticketID), userID, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, gin.H{
        "message": "Comment added successfully",
        "comment": comment,
    })
}

func (h *TicketHandler) GetAllTickets(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    
    tickets, total, err := h.ticketService.GetTickets(0, 1, page, limit) // Admin view
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "tickets": tickets,
        "total":   total,
        "page":    page,
        "limit":   limit,
    })
}
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ticket-service/models"
	"ticket-service/services"
)

type EmailHandler struct {
	emailService *services.EmailService
}

func NewEmailHandler() *EmailHandler {
	return &EmailHandler{
		emailService: services.NewEmailService(),
	}
}

func (h *EmailHandler) GetEmailTemplates(c *gin.Context) {
	filters := map[string]string{
		"template_type":   c.Query("template_type"),
		"organization_id": c.Query("organization_id"),
		"is_active":       c.Query("is_active"),
	}
	
	templates, err := h.emailService.GetEmailTemplates(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

func (h *EmailHandler) CreateEmailTemplate(c *gin.Context) {
	var req models.CreateEmailTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	template, err := h.emailService.CreateEmailTemplate(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"template": template})
}

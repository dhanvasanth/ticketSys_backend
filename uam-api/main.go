package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"strings"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// Models
type Ticket struct {
	ID              int64           `json:"id" db:"id"`
	TicketNumber    string          `json:"ticket_number" db:"ticket_number"`
	Subject         string          `json:"subject" db:"subject"`
	Description     string          `json:"description" db:"description"`
	Status          string          `json:"status" db:"status"`
	PriorityID      *int64          `json:"priority_id" db:"priority_id"`
	CategoryID      *int64          `json:"category_id" db:"category_id"`
	RequesterID     int64           `json:"requester_id" db:"requester_id"`
	RequesterEmail  string          `json:"requester_email" db:"requester_email"`
	RequesterName   string          `json:"requester_name" db:"requester_name"`
	AssigneeID      *int64          `json:"assignee_id" db:"assignee_id"`
	OrganizationID  *int64          `json:"organization_id" db:"organization_id"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
	DueDate         *time.Time      `json:"due_date" db:"due_date"`
	ResolvedAt      *time.Time      `json:"resolved_at" db:"resolved_at"`
	ClosedAt        *time.Time      `json:"closed_at" db:"closed_at"`
	Source          string          `json:"source" db:"source"`
	Tags            json.RawMessage `json:"tags" db:"tags"`
	CustomFields    json.RawMessage `json:"custom_fields" db:"custom_fields"`
}

type TicketComment struct {
	ID             int64     `json:"id" db:"id"`
	TicketID       int64     `json:"ticket_id" db:"ticket_id"`
	UserID         *int64    `json:"user_id" db:"user_id"`
	Content        string    `json:"content" db:"content"`
	IsPublic       bool      `json:"is_public" db:"is_public"`
	IsFromCustomer bool      `json:"is_from_customer" db:"is_from_customer"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type EmailTemplate struct {
	ID             int64     `json:"id" db:"id"`
	Name           string    `json:"name" db:"name"`
	Subject        string    `json:"subject" db:"subject"`
	Body           string    `json:"body" db:"body"`
	TemplateType   string    `json:"template_type" db:"template_type"`
	OrganizationID *int64    `json:"organization_id" db:"organization_id"`
	IsActive       bool      `json:"is_active" db:"is_active"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// Request/Response DTOs
type CreateTicketRequest struct {
	Subject         string          `json:"subject" binding:"required"`
	Description     string          `json:"description"`
	PriorityID      *int64          `json:"priority_id"`
	CategoryID      *int64          `json:"category_id"`
	RequesterID     int64           `json:"requester_id" binding:"required"`
	RequesterEmail  string          `json:"requester_email" binding:"required"`
	RequesterName   string          `json:"requester_name" binding:"required"`
	AssigneeID      *int64          `json:"assignee_id"`
	OrganizationID  *int64          `json:"organization_id"`
	DueDate         *time.Time      `json:"due_date"`
	Source          string          `json:"source"`
	Tags            json.RawMessage `json:"tags"`
	CustomFields    json.RawMessage `json:"custom_fields"`
}

type UpdateTicketRequest struct {
	Subject         *string         `json:"subject"`
	Description     *string         `json:"description"`
	Status          *string         `json:"status"`
	PriorityID      *int64          `json:"priority_id"`
	CategoryID      *int64          `json:"category_id"`
	AssigneeID      *int64          `json:"assignee_id"`
	DueDate         *time.Time      `json:"due_date"`
	Tags            json.RawMessage `json:"tags"`
	CustomFields    json.RawMessage `json:"custom_fields"`
}

type CreateCommentRequest struct {
	Content        string `json:"content" binding:"required"`
	UserID         *int64 `json:"user_id"`
	IsPublic       *bool  `json:"is_public"`
	IsFromCustomer *bool  `json:"is_from_customer"`
}

type CreateEmailTemplateRequest struct {
	Name           string `json:"name" binding:"required"`
	Subject        string `json:"subject" binding:"required"`
	Body           string `json:"body" binding:"required"`
	TemplateType   string `json:"template_type" binding:"required"`
	OrganizationID *int64 `json:"organization_id"`
	IsActive       *bool  `json:"is_active"`
}

type TicketResponse struct {
	Ticket   *Ticket          `json:"ticket"`
	Comments []*TicketComment `json:"comments,omitempty"`
}

type PaginatedResponse struct {
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
}

type DashboardStats struct {
	TotalTickets        int64 `json:"total_tickets"`
	OpenTickets         int64 `json:"open_tickets"`
	ResolvedTickets     int64 `json:"resolved_tickets"`
	AvgResolutionTime   int64 `json:"avg_resolution_time"`
	TicketsByStatus     map[string]int64 `json:"tickets_by_status"`
	TicketsByPriority   map[string]int64 `json:"tickets_by_priority"`
}

// Database connection
var db *sql.DB

// Initialize database connection
func initDB() {
	var err error
	db, err = sql.Open("mysql", "mariadb:mariadb@tcp(localhost:3306)/ticket?parseTime=true")
	if err != nil {
		panic(err)
	}
	
	if err = db.Ping(); err != nil {
		panic(err)
	}
}

// Ticket Handlers
func getTickets(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit
	
	status := c.Query("status")
	priorityID := c.Query("priority_id")
	assigneeID := c.Query("assignee_id")
	organizationID := c.Query("organization_id")
	
	query := "SELECT * FROM tickets WHERE 1=1"
	args := []interface{}{}
	
	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}
	if priorityID != "" {
		query += " AND priority_id = ?"
		args = append(args, priorityID)
	}
	if assigneeID != "" {
		query += " AND assignee_id = ?"
		args = append(args, assigneeID)
	}
	if organizationID != "" {
		query += " AND organization_id = ?"
		args = append(args, organizationID)
	}
	
	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)
	
	rows, err := db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	
	var tickets []*Ticket
	for rows.Next() {
		ticket := &Ticket{}
		err := rows.Scan(&ticket.ID, &ticket.TicketNumber, &ticket.Subject, &ticket.Description,
			&ticket.Status, &ticket.PriorityID, &ticket.CategoryID, &ticket.RequesterID,
			&ticket.RequesterEmail, &ticket.RequesterName, &ticket.AssigneeID,
			&ticket.OrganizationID, &ticket.CreatedAt, &ticket.UpdatedAt,
			&ticket.DueDate, &ticket.ResolvedAt, &ticket.ClosedAt,
			&ticket.Source, &ticket.Tags, &ticket.CustomFields)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		tickets = append(tickets, ticket)
	}
	
	// Get total count
	var total int64
	countQuery := "SELECT COUNT(*) FROM tickets WHERE 1=1"
	countArgs := []interface{}{}
	
	if status != "" {
		countQuery += " AND status = ?"
		countArgs = append(countArgs, status)
	}
	if priorityID != "" {
		countQuery += " AND priority_id = ?"
		countArgs = append(countArgs, priorityID)
	}
	if assigneeID != "" {
		countQuery += " AND assignee_id = ?"
		countArgs = append(countArgs, assigneeID)
	}
	if organizationID != "" {
		countQuery += " AND organization_id = ?"
		countArgs = append(countArgs, organizationID)
	}
	
	db.QueryRow(countQuery, countArgs...).Scan(&total)
	
	response := PaginatedResponse{
		Data:  tickets,
		Total: total,
		Page:  page,
		Limit: limit,
	}
	
	c.JSON(http.StatusOK, response)
}

func getTicket(c *gin.Context) {
	id := c.Param("id")
	
	ticket := &Ticket{}
	query := "SELECT * FROM tickets WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&ticket.ID, &ticket.TicketNumber, &ticket.Subject,
		&ticket.Description, &ticket.Status, &ticket.PriorityID, &ticket.CategoryID,
		&ticket.RequesterID, &ticket.RequesterEmail, &ticket.RequesterName,
		&ticket.AssigneeID, &ticket.OrganizationID, &ticket.CreatedAt,
		&ticket.UpdatedAt, &ticket.DueDate, &ticket.ResolvedAt, &ticket.ClosedAt,
		&ticket.Source, &ticket.Tags, &ticket.CustomFields)
	
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	
	// Get comments
	comments, err := getTicketComments(ticket.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	response := TicketResponse{
		Ticket:   ticket,
		Comments: comments,
	}
	
	c.JSON(http.StatusOK, response)
}

func createTicket(c *gin.Context) {
	var req CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Generate ticket number
	ticketNumber := generateTicketNumber()
	
	query := `INSERT INTO tickets (ticket_number, subject, description, priority_id, category_id, 
		requester_id, requester_email, requester_name, assignee_id, organization_id, 
		due_date, source, tags, custom_fields) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	result, err := db.Exec(query, ticketNumber, req.Subject, req.Description,
		req.PriorityID, req.CategoryID, req.RequesterID, req.RequesterEmail,
		req.RequesterName, req.AssigneeID, req.OrganizationID, req.DueDate,
		req.Source, req.Tags, req.CustomFields)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	id, _ := result.LastInsertId()
	
	// Get created ticket
	ticket := &Ticket{}
	selectQuery := "SELECT * FROM tickets WHERE id = ?"
	db.QueryRow(selectQuery, id).Scan(&ticket.ID, &ticket.TicketNumber, &ticket.Subject,
		&ticket.Description, &ticket.Status, &ticket.PriorityID, &ticket.CategoryID,
		&ticket.RequesterID, &ticket.RequesterEmail, &ticket.RequesterName,
		&ticket.AssigneeID, &ticket.OrganizationID, &ticket.CreatedAt,
		&ticket.UpdatedAt, &ticket.DueDate, &ticket.ResolvedAt, &ticket.ClosedAt,
		&ticket.Source, &ticket.Tags, &ticket.CustomFields)
	
	c.JSON(http.StatusCreated, gin.H{
		"ticket":        ticket,
		"ticket_number": ticketNumber,
	})
}

func updateTicket(c *gin.Context) {
	id := c.Param("id")
	
	var req UpdateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	query := "UPDATE tickets SET "
	args := []interface{}{}
	setParts := []string{}
	
	if req.Subject != nil {
		setParts = append(setParts, "subject = ?")
		args = append(args, *req.Subject)
	}
	if req.Description != nil {
		setParts = append(setParts, "description = ?")
		args = append(args, *req.Description)
	}
	if req.Status != nil {
		setParts = append(setParts, "status = ?")
		args = append(args, *req.Status)
		
		// Set resolved_at if status is resolved
		if *req.Status == "resolved" {
			setParts = append(setParts, "resolved_at = NOW()")
		}
		// Set closed_at if status is closed
		if *req.Status == "closed" {
			setParts = append(setParts, "closed_at = NOW()")
		}
	}
	if req.PriorityID != nil {
		setParts = append(setParts, "priority_id = ?")
		args = append(args, *req.PriorityID)
	}
	if req.CategoryID != nil {
		setParts = append(setParts, "category_id = ?")
		args = append(args, *req.CategoryID)
	}
	if req.AssigneeID != nil {
		setParts = append(setParts, "assignee_id = ?")
		args = append(args, *req.AssigneeID)
	}
	if req.DueDate != nil {
		setParts = append(setParts, "due_date = ?")
		args = append(args, *req.DueDate)
	}
	if req.Tags != nil {
		setParts = append(setParts, "tags = ?")
		args = append(args, req.Tags)
	}
	if req.CustomFields != nil {
		setParts = append(setParts, "custom_fields = ?")
		args = append(args, req.CustomFields)
	}
	
	if len(setParts) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}
	
	setParts = append(setParts, "updated_at = NOW()")
	query += strings.Join(setParts, ", ") + " WHERE id = ?"
	args = append(args, id)
	
	_, err := db.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Ticket updated successfully"})
}

func deleteTicket(c *gin.Context) {
	id := c.Param("id")
	
	_, err := db.Exec("DELETE FROM tickets WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Ticket deleted successfully"})
}

// Comment Handlers
func getTicketComments(ticketID int64) ([]*TicketComment, error) {
	query := "SELECT * FROM ticket_comments WHERE ticket_id = ? ORDER BY created_at ASC"
	rows, err := db.Query(query, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var comments []*TicketComment
	for rows.Next() {
		comment := &TicketComment{}
		err := rows.Scan(&comment.ID, &comment.TicketID, &comment.UserID,
			&comment.Content, &comment.IsPublic, &comment.IsFromCustomer,
			&comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	
	return comments, nil
}

func createComment(c *gin.Context) {
	ticketID := c.Param("ticket_id")
	
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}
	
	isFromCustomer := false
	if req.IsFromCustomer != nil {
		isFromCustomer = *req.IsFromCustomer
	}
	
	query := `INSERT INTO ticket_comments (ticket_id, user_id, content, is_public, is_from_customer) 
		VALUES (?, ?, ?, ?, ?)`
	
	result, err := db.Exec(query, ticketID, req.UserID, req.Content, isPublic, isFromCustomer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	id, _ := result.LastInsertId()
	
	// Get created comment
	comment := &TicketComment{}
	selectQuery := "SELECT * FROM ticket_comments WHERE id = ?"
	db.QueryRow(selectQuery, id).Scan(&comment.ID, &comment.TicketID, &comment.UserID,
		&comment.Content, &comment.IsPublic, &comment.IsFromCustomer,
		&comment.CreatedAt, &comment.UpdatedAt)
	
	c.JSON(http.StatusCreated, gin.H{"comment": comment})
}

// Email Template Handlers
func getEmailTemplates(c *gin.Context) {
	templateType := c.Query("template_type")
	organizationID := c.Query("organization_id")
	isActive := c.Query("is_active")
	
	query := "SELECT * FROM email_templates WHERE 1=1"
	args := []interface{}{}
	
	if templateType != "" {
		query += " AND template_type = ?"
		args = append(args, templateType)
	}
	if organizationID != "" {
		query += " AND organization_id = ?"
		args = append(args, organizationID)
	}
	if isActive != "" {
		query += " AND is_active = ?"
		args = append(args, isActive == "true")
	}
	
	rows, err := db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	
	var templates []*EmailTemplate
	for rows.Next() {
		template := &EmailTemplate{}
		err := rows.Scan(&template.ID, &template.Name, &template.Subject,
			&template.Body, &template.TemplateType, &template.OrganizationID,
			&template.IsActive, &template.CreatedAt, &template.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		templates = append(templates, template)
	}
	
	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

func createEmailTemplate(c *gin.Context) {
	var req CreateEmailTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	
	query := `INSERT INTO email_templates (name, subject, body, template_type, organization_id, is_active) 
		VALUES (?, ?, ?, ?, ?, ?)`
	
	result, err := db.Exec(query, req.Name, req.Subject, req.Body,
		req.TemplateType, req.OrganizationID, isActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	id, _ := result.LastInsertId()
	
	// Get created template
	template := &EmailTemplate{}
	selectQuery := "SELECT * FROM email_templates WHERE id = ?"
	db.QueryRow(selectQuery, id).Scan(&template.ID, &template.Name, &template.Subject,
		&template.Body, &template.TemplateType, &template.OrganizationID,
		&template.IsActive, &template.CreatedAt, &template.UpdatedAt)
	
	c.JSON(http.StatusCreated, gin.H{"template": template})
}

// Dashboard Handler
func getDashboardStats(c *gin.Context) {
	organizationID := c.Query("organization_id")
	
	stats := &DashboardStats{
		TicketsByStatus:   make(map[string]int64),
		TicketsByPriority: make(map[string]int64),
	}
	
	// Base query condition
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	if organizationID != "" {
		whereClause += " AND organization_id = ?"
		args = append(args, organizationID)
	}
	
	// Total tickets
	db.QueryRow("SELECT COUNT(*) FROM tickets "+whereClause, args...).Scan(&stats.TotalTickets)
	
	// Open tickets
	openArgs := append(args, "open")
	db.QueryRow("SELECT COUNT(*) FROM tickets "+whereClause+" AND status = ?", openArgs...).Scan(&stats.OpenTickets)
	
	// Resolved tickets
	resolvedArgs := append(args, "resolved")
	db.QueryRow("SELECT COUNT(*) FROM tickets "+whereClause+" AND status = ?", resolvedArgs...).Scan(&stats.ResolvedTickets)
	
	// Tickets by status
	statusQuery := "SELECT status, COUNT(*) FROM tickets " + whereClause + " GROUP BY status"
	rows, _ := db.Query(statusQuery, args...)
	defer rows.Close()
	
	for rows.Next() {
		var status string
		var count int64
		rows.Scan(&status, &count)
		stats.TicketsByStatus[status] = count
	}
	
	c.JSON(http.StatusOK, stats)
}

// Utility functions
func generateTicketNumber() string {
	return "TKT-" + time.Now().Format("2006") + "-" + strconv.FormatInt(time.Now().UnixNano(), 10)[10:]
}

// Routes setup
func setupRoutes() *gin.Engine {
	r := gin.Default()
	
	api := r.Group("/api")
	{
		// Tickets
		api.GET("/tickets", getTickets)
		api.GET("/tickets/:id", getTicket)
		api.POST("/tickets", createTicket)
		api.PUT("/tickets/:id", updateTicket)
		api.DELETE("/tickets/:id", deleteTicket)
		
		// Comments
		api.POST("/tickets/:ticket_id/comments", createComment)
		
		// Email Templates
		api.GET("/email-templates", getEmailTemplates)
		api.POST("/email-templates", createEmailTemplate)
		
		// Dashboard
		api.GET("/dashboard/stats", getDashboardStats)
	}
	
	return r
}

func main() {
	initDB()
	defer db.Close()
	
	r := setupRoutes()
	r.Run(":6001")
}
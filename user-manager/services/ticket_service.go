package services

import (
	"database/sql"
	"strings"
	"ticket-service/config"
	"ticket-service/models"
	"ticket-service/utils"
)

type TicketService struct{}

func NewTicketService() *TicketService {
	return &TicketService{}
}

func (s *TicketService) GetTickets(page, limit int, filters map[string]string) ([]*models.Ticket, int64, error) {
	offset := (page - 1) * limit
	
	query := "SELECT * FROM tickets WHERE 1=1"
	args := []interface{}{}
	
	if status := filters["status"]; status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}
	if priorityID := filters["priority_id"]; priorityID != "" {
		query += " AND priority_id = ?"
		args = append(args, priorityID)
	}
	if assigneeID := filters["assignee_id"]; assigneeID != "" {
		query += " AND assignee_id = ?"
		args = append(args, assigneeID)
	}
	if organizationID := filters["organization_id"]; organizationID != "" {
		query += " AND organization_id = ?"
		args = append(args, organizationID)
	}
	
	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)
	
	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var tickets []*models.Ticket
	for rows.Next() {
		ticket := &models.Ticket{}
		err := rows.Scan(&ticket.ID, &ticket.TicketNumber, &ticket.Subject, &ticket.Description,
			&ticket.Status, &ticket.PriorityID, &ticket.CategoryID, &ticket.RequesterID,
			&ticket.RequesterEmail, &ticket.RequesterName, &ticket.AssigneeID,
			&ticket.OrganizationID, &ticket.CreatedAt, &ticket.UpdatedAt,
			&ticket.DueDate, &ticket.ResolvedAt, &ticket.ClosedAt,
			&ticket.Source, &ticket.Tags, &ticket.CustomFields)
		if err != nil {
			return nil, 0, err
		}
		tickets = append(tickets, ticket)
	}
	
	// Get total count
	var total int64
	countQuery := "SELECT COUNT(*) FROM tickets WHERE 1=1"
	countArgs := []interface{}{}
	
	if status := filters["status"]; status != "" {
		countQuery += " AND status = ?"
		countArgs = append(countArgs, status)
	}
	if priorityID := filters["priority_id"]; priorityID != "" {
		countQuery += " AND priority_id = ?"
		countArgs = append(countArgs, priorityID)
	}
	if assigneeID := filters["assignee_id"]; assigneeID != "" {
		countQuery += " AND assignee_id = ?"
		countArgs = append(countArgs, assigneeID)
	}
	if organizationID := filters["organization_id"]; organizationID != "" {
		countQuery += " AND organization_id = ?"
		countArgs = append(countArgs, organizationID)
	}
	
	config.DB.QueryRow(countQuery, countArgs...).Scan(&total)
	
	return tickets, total, nil
}

func (s *TicketService) GetTicketByID(id string) (*models.Ticket, []*models.TicketComment, error) {
	ticket := &models.Ticket{}
	query := "SELECT * FROM tickets WHERE id = ?"
	err := config.DB.QueryRow(query, id).Scan(&ticket.ID, &ticket.TicketNumber, &ticket.Subject,
		&ticket.Description, &ticket.Status, &ticket.PriorityID, &ticket.CategoryID,
		&ticket.RequesterID, &ticket.RequesterEmail, &ticket.RequesterName,
		&ticket.AssigneeID, &ticket.OrganizationID, &ticket.CreatedAt,
		&ticket.UpdatedAt, &ticket.DueDate, &ticket.ResolvedAt, &ticket.ClosedAt,
		&ticket.Source, &ticket.Tags, &ticket.CustomFields)
	
	if err != nil {
		return nil, nil, err
	}
	
	// Get comments
	comments, err := s.GetTicketComments(ticket.ID)
	if err != nil {
		return nil, nil, err
	}
	
	return ticket, comments, nil
}

func (s *TicketService) CreateTicket(req *models.CreateTicketRequest) (*models.Ticket, error) {
	// Generate ticket number
	ticketNumber := utils.GenerateTicketNumber()
	
	query := `INSERT INTO tickets (ticket_number, subject, description, priority_id, category_id, 
		requester_id, requester_email, requester_name, assignee_id, organization_id, 
		due_date, source, tags, custom_fields) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	result, err := config.DB.Exec(query, ticketNumber, req.Subject, req.Description,
		req.PriorityID, req.CategoryID, req.RequesterID, req.RequesterEmail,
		req.RequesterName, req.AssigneeID, req.OrganizationID, req.DueDate,
		req.Source, req.Tags, req.CustomFields)
	
	if err != nil {
		return nil, err
	}
	
	id, _ := result.LastInsertId()
	
	// Get created ticket
	ticket := &models.Ticket{}
	selectQuery := "SELECT * FROM tickets WHERE id = ?"
	config.DB.QueryRow(selectQuery, id).Scan(&ticket.ID, &ticket.TicketNumber, &ticket.Subject,
		&ticket.Description, &ticket.Status, &ticket.PriorityID, &ticket.CategoryID,
		&ticket.RequesterID, &ticket.RequesterEmail, &ticket.RequesterName,
		&ticket.AssigneeID, &ticket.OrganizationID, &ticket.CreatedAt,
		&ticket.UpdatedAt, &ticket.DueDate, &ticket.ResolvedAt, &ticket.ClosedAt,
		&ticket.Source, &ticket.Tags, &ticket.CustomFields)
	
	return ticket, nil
}

func (s *TicketService) UpdateTicket(id string, req *models.UpdateTicketRequest) error {
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
		return nil
	}
	
	setParts = append(setParts, "updated_at = NOW()")
	query += strings.Join(setParts, ", ") + " WHERE id = ?"
	args = append(args, id)
	
	_, err := config.DB.Exec(query, args...)
	return err
}

func (s *TicketService) DeleteTicket(id string) error {
	_, err := config.DB.Exec("DELETE FROM tickets WHERE id = ?", id)
	return err
}

func (s *TicketService) GetTicketComments(ticketID int64) ([]*models.TicketComment, error) {
	query := "SELECT * FROM ticket_comments WHERE ticket_id = ? ORDER BY created_at ASC"
	rows, err := config.DB.Query(query, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var comments []*models.TicketComment
	for rows.Next() {
		comment := &models.TicketComment{}
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

func (s *TicketService) CreateComment(ticketID string, req *models.CreateCommentRequest) (*models.TicketComment, error) {
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
	
	result, err := config.DB.Exec(query, ticketID, req.UserID, req.Content, isPublic, isFromCustomer)
	if err != nil {
		return nil, err
	}
	
	id, _ := result.LastInsertId()
	
	// Get created comment
	comment := &models.TicketComment{}
	selectQuery := "SELECT * FROM ticket_comments WHERE id = ?"
	config.DB.QueryRow(selectQuery, id).Scan(&comment.ID, &comment.TicketID, &comment.UserID,
		&comment.Content, &comment.IsPublic, &comment.IsFromCustomer,
		&comment.CreatedAt, &comment.UpdatedAt)
	
	return comment, nil
}

func (s *TicketService) GetDashboardStats(organizationID string) (*models.DashboardStats, error) {
	stats := &models.DashboardStats{
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
	config.DB.QueryRow("SELECT COUNT(*) FROM tickets "+whereClause, args...).Scan(&stats.TotalTickets)
	
	// Open tickets
	openArgs := append(args, "open")
	config.DB.QueryRow("SELECT COUNT(*) FROM tickets "+whereClause+" AND status = ?", openArgs...).Scan(&stats.OpenTickets)
	
	// Resolved tickets
	resolvedArgs := append(args, "resolved")
	config.DB.QueryRow("SELECT COUNT(*) FROM tickets "+whereClause+" AND status = ?", resolvedArgs...).Scan(&stats.ResolvedTickets)
	
	// Tickets by status
	statusQuery := "SELECT status, COUNT(*) FROM tickets " + whereClause + " GROUP BY status"
	rows, _ := config.DB.Query(statusQuery, args...)
	defer rows.Close()
	
	for rows.Next() {
		var status string
		var count int64
		rows.Scan(&status, &count)
		stats.TicketsByStatus[status] = count
	}
	
	return stats, nil
}

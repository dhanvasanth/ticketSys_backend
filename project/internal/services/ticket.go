package services

import (
    "errors"
    //"strconv"
    "time"
    "project/internal/models"
    "project/internal/repositories"
)

type TicketService interface {
    CreateTicket(userID uint, req *models.CreateTicketRequest) (*models.Ticket, error)
    GetTickets(userID uint, roleID uint, page, limit int) ([]*models.Ticket, int64, error)
    GetTicket(ticketID, userID uint, roleID uint) (*models.Ticket, error)
    UpdateTicket(ticketID, userID uint, req *models.UpdateTicketRequest) (*models.Ticket, error)
    AddComment(ticketID, userID uint, req *models.CreateCommentRequest) (*models.TicketComment, error)
}

type ticketService struct {
    ticketRepo repositories.TicketRepository
}

func NewTicketService(ticketRepo repositories.TicketRepository) TicketService {
    return &ticketService{ticketRepo: ticketRepo}
}

func (s *ticketService) CreateTicket(userID uint, req *models.CreateTicketRequest) (*models.Ticket, error) {
    ticketNumber, err := s.ticketRepo.GenerateTicketNumber()
    if err != nil {
        return nil, err
    }
    
    ticket := &models.Ticket{
        TicketNumber: ticketNumber,
        Subject:      req.Subject,
        Description:  req.Description,
        Status:       "open",
        Priority:     req.Priority,
        RequesterID:  userID,
        AssigneeID:   req.AssigneeID,
        Source:       req.Source,
        DueDate:      req.DueDate,
    }
    
    if ticket.Priority == "" {
        ticket.Priority = "medium"
    }
    if ticket.Source == "" {
        ticket.Source = "web"
    }
    
    if err := s.ticketRepo.Create(ticket); err != nil {
        return nil, err
    }
    
    return s.ticketRepo.GetByID(ticket.ID)
}

func (s *ticketService) GetTickets(userID uint, roleID uint, page, limit int) ([]*models.Ticket, int64, error) {
    offset := (page - 1) * limit
    
    // Check if user is admin/agent (can see all tickets) or regular user (only own tickets)
    // Admin role ID = 1, Agent role ID = 2
    if roleID == 1 || roleID == 2 {
        return s.ticketRepo.GetAll(limit, offset)
    }
    
    return s.ticketRepo.GetByUserID(userID, limit, offset)
}

func (s *ticketService) GetTicket(ticketID, userID uint, roleID uint) (*models.Ticket, error) {
    ticket, err := s.ticketRepo.GetByID(ticketID)
    if err != nil {
        return nil, err
    }
    
    // Check permissions: admin/agent can see all, users can see only their own
    if roleID != 1 && roleID != 2 && ticket.RequesterID != userID {
        return nil, errors.New("permission denied")
    }
    
    return ticket, nil
}

func (s *ticketService) UpdateTicket(ticketID, userID uint, req *models.UpdateTicketRequest) (*models.Ticket, error) {
    ticket, err := s.ticketRepo.GetByID(ticketID)
    if err != nil {
        return nil, err
    }
    
    // Update fields
    if req.Subject != nil {
        ticket.Subject = *req.Subject
    }
    if req.Description != nil {
        ticket.Description = *req.Description
    }
    if req.Status != nil {
        ticket.Status = *req.Status
        if *req.Status == "resolved" {
            now := time.Now()
            ticket.ResolvedAt = &now
        }
    }
    if req.Priority != nil {
        ticket.Priority = *req.Priority
    }
    if req.AssigneeID != nil {
        ticket.AssigneeID = req.AssigneeID
    }
    if req.DueDate != nil {
        ticket.DueDate = req.DueDate
    }
    
    if err := s.ticketRepo.Update(ticket); err != nil {
        return nil, err
    }
    
    return s.ticketRepo.GetByID(ticket.ID)
}

func (s *ticketService) AddComment(ticketID, userID uint, req *models.CreateCommentRequest) (*models.TicketComment, error) {
    comment := &models.TicketComment{
        TicketID: ticketID,
        UserID:   &userID,
        Content:  req.Content,
        IsPublic: true,
        IsFromCustomer: false,
    }
    
    if req.IsPublic != nil {
        comment.IsPublic = *req.IsPublic
    }
    if req.IsFromCustomer != nil {
        comment.IsFromCustomer = *req.IsFromCustomer
    }
    
    if err := s.ticketRepo.AddComment(comment); err != nil {
        return nil, err
    }
    
    return comment, nil
}
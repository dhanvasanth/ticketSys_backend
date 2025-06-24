package repositories

import (
    "ticket-service/internal/models"
    "gorm.io/gorm"
)

type TicketRepository interface {
    Create(ticket *models.Ticket) error
    GetByID(id uint) (*models.Ticket, error)
    GetByUserID(userID uint, limit, offset int) ([]*models.Ticket, int64, error)
    GetAll(limit, offset int) ([]*models.Ticket, int64, error)
    Update(ticket *models.Ticket) error
    AddComment(comment *models.TicketComment) error
    GetComments(ticketID uint) ([]*models.TicketComment, error)
    GenerateTicketNumber() (string, error)
}

type ticketRepository struct {
    db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
    return &ticketRepository{db: db}
}

func (r *ticketRepository) Create(ticket *models.Ticket) error {
    return r.db.Create(ticket).Error
}

func (r *ticketRepository) GetByID(id uint) (*models.Ticket, error) {
    var ticket models.Ticket
    err := r.db.Preload("Requester").Preload("Assignee").Preload("Comments.User").First(&ticket, id).Error
    if err != nil {
        return nil, err
    }
    return &ticket, nil
}

func (r *ticketRepository) GetByUserID(userID uint, limit, offset int) ([]*models.Ticket, int64, error) {
    var tickets []*models.Ticket
    var total int64
    
    query := r.db.Where("requester_id = ?", userID)
    query.Model(&models.Ticket{}).Count(&total)
    
    err := query.Preload("Requester").Preload("Assignee").
        Limit(limit).Offset(offset).
        Order("created_at DESC").
        Find(&tickets).Error
    
    return tickets, total, err
}

func (r *ticketRepository) GetAll(limit, offset int) ([]*models.Ticket, int64, error) {
    var tickets []*models.Ticket
    var total int64
    
    r.db.Model(&models.Ticket{}).Count(&total)
    
    err := r.db.Preload("Requester").Preload("Assignee").
        Limit(limit).Offset(offset).
        Order("created_at DESC").
        Find(&tickets).Error
    
    return tickets, total, err
}

func (r *ticketRepository) Update(ticket *models.Ticket) error {
    return r.db.Save(ticket).Error
}

func (r *ticketRepository) AddComment(comment *models.TicketComment) error {
    return r.db.Create(comment).Error
}

func (r *ticketRepository) GetComments(ticketID uint) ([]*models.TicketComment, error) {
    var comments []*models.TicketComment
    err := r.db.Preload("User").Where("ticket_id = ?", ticketID).Order("created_at ASC").Find(&comments).Error
    return comments, err
}

func (r *ticketRepository) GenerateTicketNumber() (string, error) {
    var count int64
    r.db.Model(&models.Ticket{}).Count(&count)
    return fmt.Sprintf("TKT-%d-%06d", time.Now().Year(), count+1), nil
}
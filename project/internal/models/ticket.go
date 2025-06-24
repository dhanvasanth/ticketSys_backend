package models

import (
    "encoding/json"
    "time"
)

type Ticket struct {
    BaseModel
    TicketNumber   string          `json:"ticket_number" gorm:"uniqueIndex;not null"`
    Subject        string          `json:"subject" gorm:"not null"`
    Description    string          `json:"description"`
    Status         string          `json:"status" gorm:"default:'open'"` // open, in_progress, resolved, closed
    Priority       string          `json:"priority" gorm:"default:'medium'"` // low, medium, high, critical
    RequesterID    uint            `json:"requester_id" gorm:"not null"`
    Requester      User            `json:"requester" gorm:"foreignKey:RequesterID"`
    AssigneeID     *uint           `json:"assignee_id"`
    Assignee       *User           `json:"assignee,omitempty" gorm:"foreignKey:AssigneeID"`
    Source         string          `json:"source" gorm:"default:'web'"` // web, email, phone, api
    Tags           json.RawMessage `json:"tags" gorm:"type:json"`
    DueDate        *time.Time      `json:"due_date"`
    ResolvedAt     *time.Time      `json:"resolved_at"`
    Comments       []TicketComment `json:"comments,omitempty"`
}

type TicketComment struct {
    BaseModel
    TicketID       uint   `json:"ticket_id" gorm:"not null"`
    UserID         *uint  `json:"user_id"`
    User           *User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
    Content        string `json:"content" gorm:"not null"`
    IsPublic       bool   `json:"is_public" gorm:"default:true"`
    IsFromCustomer bool   `json:"is_from_customer" gorm:"default:false"`
}

// Request DTOs
type CreateTicketRequest struct {
    Subject     string     `json:"subject" binding:"required"`
    Description string     `json:"description"`
    Priority    string     `json:"priority"`
    AssigneeID  *uint      `json:"assignee_id"`
    Source      string     `json:"source"`
    DueDate     *time.Time `json:"due_date"`
}

type UpdateTicketRequest struct {
    Subject     *string    `json:"subject"`
    Description *string    `json:"description"`
    Status      *string    `json:"status"`
    Priority    *string    `json:"priority"`
    AssigneeID  *uint      `json:"assignee_id"`
    DueDate     *time.Time `json:"due_date"`
}

type CreateCommentRequest struct {
    Content        string `json:"content" binding:"required"`
    IsPublic       *bool  `json:"is_public"`
    IsFromCustomer *bool  `json:"is_from_customer"`
}
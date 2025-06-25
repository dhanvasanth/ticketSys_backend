// project/internal/models/role_change_request.go
package models

import "time"

type RoleChangeRequest struct {
    BaseModel
    RequesterID     uint      `json:"requester_id" gorm:"not null"`
    Requester       User      `json:"requester" gorm:"foreignKey:RequesterID"`
    CurrentRoleID   uint      `json:"current_role_id" gorm:"not null"`
    CurrentRole     Role      `json:"current_role" gorm:"foreignKey:CurrentRoleID"`
    RequestedRoleID uint      `json:"requested_role_id" gorm:"not null"`
    RequestedRole   Role      `json:"requested_role" gorm:"foreignKey:RequestedRoleID"`
    Reason          string    `json:"reason" gorm:"not null"`
    Status          string    `json:"status" gorm:"default:'pending'"` // pending, approved, rejected
    AdminID         *uint     `json:"admin_id"`
    Admin           *User     `json:"admin,omitempty" gorm:"foreignKey:AdminID"`
    AdminNotes      string    `json:"admin_notes"`
    ProcessedAt     *time.Time `json:"processed_at"`
}

// Request DTOs
type CreateRoleChangeRequest struct {
    RequestedRoleID uint   `json:"requested_role_id" binding:"required"`
    Reason          string `json:"reason" binding:"required,min=10"`
}

type ProcessRoleChangeRequest struct {
    Status     string `json:"status" binding:"required,oneof=approved rejected"`
    AdminNotes string `json:"admin_notes"`
}
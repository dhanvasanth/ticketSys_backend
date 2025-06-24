package models

type Role struct {
    BaseModel
    Name        string `json:"name" gorm:"uniqueIndex;not null"` // admin, agent, user, customer
    Description string `json:"description"`
    Permissions string `json:"permissions" gorm:"type:json"` // JSON array of permissions
    IsActive    bool   `json:"is_active" gorm:"default:true"`
}

type Permission struct {
    ID          uint   `json:"id" gorm:"primaryKey"`
    Name        string `json:"name" gorm:"uniqueIndex;not null"` // create_ticket, assign_ticket, etc.
    Resource    string `json:"resource"`                         // ticket, user, comment
    Action      string `json:"action"`                           // create, read, update, delete
    Description string `json:"description"`
}
package models

type User struct {
    BaseModel
    Name         string `json:"name" gorm:"not null"`
    Email        string `json:"email" gorm:"uniqueIndex;not null"`
    Password     string `json:"-" gorm:"not null"`
    RoleID       uint   `json:"role_id" gorm:"not null"`
    Role         Role   `json:"role" gorm:"foreignKey:RoleID"`
    IsActive     bool   `json:"is_active" gorm:"default:true"`
    Organization string `json:"organization"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
    Name         string `json:"name" binding:"required"`
    Email        string `json:"email" binding:"required,email"`
    Password     string `json:"password" binding:"required,min=6"`
    Organization string `json:"organization"`
}

type UserResponse struct {
    ID           uint   `json:"id"`
    Name         string `json:"name"`
    Email        string `json:"email"`
    Role         Role   `json:"role"`
    Organization string `json:"organization"`
}
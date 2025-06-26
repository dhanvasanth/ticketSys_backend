// package models

// type User struct {
//     BaseModel
//     Name         string `json:"name" gorm:"not null"`
//     Email        string `json:"email" gorm:"uniqueIndex;not null"`
//     Password     string `json:"-" gorm:"not null"`
//     RoleID       uint   `json:"role_id" gorm:"not null"`
//     Role         Role   `json:"role" gorm:"foreignKey:RoleID"`
//     IsActive     bool   `json:"is_active" gorm:"default:true"`
//     Organization string `json:"organization"`
// }

// type LoginRequest struct {
//     Email    string `json:"email" binding:"required,email"`
//     Password string `json:"password" binding:"required,min=6"`
// }

// type RegisterRequest struct {
//     Name         string `json:"name" binding:"required"`
//     Email        string `json:"email" binding:"required,email"`
//     Password     string `json:"password" binding:"required,min=6"`
//     Organization string `json:"organization"`
// }

// type UserResponse struct {
//     ID           uint   `json:"id"`
//     Name         string `json:"name"`
//     Email        string `json:"email"`
//     Role         Role   `json:"role"`
//     Organization string `json:"organization"`
// }


// type VerifyLoginCodeRequest struct {
//     UserID uint   `json:"user_id" binding:"required"`
//     Email  string `json:"email" binding:"required,email"`
//     Code   string `json:"code" binding:"required"`
// }

// type ResendCodeRequest struct {
//     Email string `json:"email" binding:"required,email"`
// }

// type ChangePasswordRequest struct {
//     CurrentPassword string `json:"current_password" binding:"required"`
//     NewPassword     string `json:"new_password" binding:"required,min=6"`
// }

// type UpdateUserRequest struct {
//     FirstName string `json:"first_name"`
//     LastName  string `json:"last_name"`
// }

// project/internal/models/user.go
package models

import "time"

type User struct {
    BaseModel
    Name              string     `json:"name" gorm:"not null"`
    Email             string     `json:"email" gorm:"uniqueIndex;not null"`
    Password          string     `json:"-" gorm:"not null"`
    RoleID            uint       `json:"role_id" gorm:"not null"`
    Role              Role       `json:"role" gorm:"foreignKey:RoleID"`
    IsActive          bool       `json:"is_active" gorm:"default:true"`
    Organization      string     `json:"organization"`
    EmailVerified     bool       `json:"email_verified" gorm:"default:false"`
    EmailVerifiedAt   *time.Time `json:"email_verified_at"`
}

// Email Verification Code model
type EmailVerificationCode struct {
    BaseModel
    UserID    uint      `json:"user_id" gorm:"not null"`
    User      User      `json:"user" gorm:"foreignKey:UserID"`
    Email     string    `json:"email" gorm:"not null"`
    Code      string    `json:"code" gorm:"not null"`
    ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
    Used      bool      `json:"used" gorm:"default:false"`
    UsedAt    *time.Time `json:"used_at"`
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

type VerifyEmailRequest struct {
    Email string `json:"email" binding:"required,email"`
    Code  string `json:"code" binding:"required"`
}

type ResendVerificationRequest struct {
    Email string `json:"email" binding:"required,email"`
}

type UserResponse struct {
    ID            uint   `json:"id"`
    Name          string `json:"name"`
    Email         string `json:"email"`
    Role          Role   `json:"role"`
    Organization  string `json:"organization"`
    EmailVerified bool   `json:"email_verified"`
}

type UpdateUserRequest struct {
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
}

type ChangePasswordRequest struct {
    CurrentPassword string `json:"current_password" binding:"required"`
    NewPassword     string `json:"new_password" binding:"required,min=6"`
}
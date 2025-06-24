package models

type Users struct {
	ID       uint64         `gorm:"primaryKey" json:"id"`
	Email    string         `json:"email"`
	UserName string         `json:"userName"`
	Password string         `json:"-"` // Hide password in JSON
	Roles    []RelUserRoles `gorm:"foreignKey:UserID" json:"roles,omitempty"`
}

type Roles struct {
	ID   uint64 `gorm:"primaryKey" json:"id"`
	Role string `json:"role"`
}

type RelUserRoles struct {
	// Omit internal fields from JSON
	ID     uint64 `gorm:"primaryKey" json:"-"`
	UserID uint64 `json:"-"`
	RoleID uint64 `json:"-"`
	Role   Roles  `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

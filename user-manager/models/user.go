package models

type User struct {
	ID          uint64 `gorm:"primaryKey" json:"id"`
	Email       string `json:"email"`
	UserName string `json:"userName"`
	Password string `json:"phoneNumber"`
}

type Roles struct{
	ID			uint64 `gorm:"primaryKey" json:"id"`
	Role     	string `json:"role`
}


type rel_user_roles struct {
	ID 			uint64 `gorm:"primaryKey" json:"id"`
	UserID		uint64 `json:"userId"`
	RoleID   	uint64 `json:"roleId"`
}

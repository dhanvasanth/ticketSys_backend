package controllers

import (
	"errors"
	"user-manager/models"
	"user-manager/services"
)

func CreateUser(user *models.Users) error {
	if err := services.DB.Create(user).Error; err != nil {
		return err
	}

	var count int64
	if err := services.DB.Model(&models.Users{}).Count(&count).Error; err != nil {
		return err
	}

	var roleName string
	if count == 1 {
		roleName = "admin"
	} else {
		roleName = "customer"
	}

	var role models.Roles
	err := services.DB.Where("role = ?", roleName).First(&role).Error
	if err != nil {
		role = models.Roles{Role: roleName}
		if err := services.DB.Create(&role).Error; err != nil {
			return errors.New("failed to create role: " + roleName)
		}
	}

	relation := models.RelUserRoles{
		UserID: user.ID,
		RoleID: role.ID,
	}
	return services.DB.Create(&relation).Error
}

func GetAllUsers() ([]models.Users, error) {
	var users []models.Users
	err := services.DB.
		Preload("Roles.Role"). 
		Find(&users).Error
	return users, err
}

func GetUserByID(id uint64) (*models.Users, error) {
	var user models.Users
	err := services.DB.
		Preload("Roles.Role").
		First(&user, id).Error
	return &user, err
}

func UpdateUser(user *models.Users) error {
	return services.DB.Save(user).Error
}

func DeleteUser(id uint64) error {
	if err := services.DB.Where("user_id = ?", id).Delete(&models.RelUserRoles{}).Error; err != nil {
		return err
	}

	return services.DB.Delete(&models.Users{}, id).Error
}


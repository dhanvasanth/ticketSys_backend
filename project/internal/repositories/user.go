package repositories

import (
    //"errors"
    "project/internal/models"
    "gorm.io/gorm"
)

type UserRepository interface {
    Create(user *models.User) error
    GetByID(id uint) (*models.User, error)
    GetByEmail(email string) (*models.User, error)
    Update(user *models.User) error
    GetAll() ([]*models.User, error)
    GetRoles() ([]*models.Role, error)
    CreateRole(role *models.Role) error
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
    return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id uint) (*models.User, error) {
    var user models.User
    err := r.db.Preload("Role").First(&user, id).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
    var user models.User
    err := r.db.Preload("Role").Where("email = ?", email).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
    return r.db.Save(user).Error
}

func (r *userRepository) GetAll() ([]*models.User, error) {
    var users []*models.User
    err := r.db.Preload("Role").Find(&users).Error
    return users, err
}

func (r *userRepository) GetRoles() ([]*models.Role, error) {
    var roles []*models.Role
    err := r.db.Find(&roles).Error
    return roles, err
}

func (r *userRepository) CreateRole(role *models.Role) error {
    return r.db.Create(role).Error
}
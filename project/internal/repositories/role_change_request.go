// project/internal/repositories/role_change_request.go
package repositories

import (
    "project/internal/models"
    "gorm.io/gorm"
)

type RoleChangeRequestRepository interface {
    Create(req *models.RoleChangeRequest) error
    GetByID(id uint) (*models.RoleChangeRequest, error)
    GetByUserID(userID uint, limit, offset int) ([]*models.RoleChangeRequest, int64, error)
    GetAll(limit, offset int) ([]*models.RoleChangeRequest, int64, error)
    GetPending(limit, offset int) ([]*models.RoleChangeRequest, int64, error)
    Update(req *models.RoleChangeRequest) error
    HasPendingRequest(userID uint) (bool, error)
}

type roleChangeRequestRepository struct {
    db *gorm.DB
}

func NewRoleChangeRequestRepository(db *gorm.DB) RoleChangeRequestRepository {
    return &roleChangeRequestRepository{db: db}
}

func (r *roleChangeRequestRepository) Create(req *models.RoleChangeRequest) error {
    return r.db.Create(req).Error
}

func (r *roleChangeRequestRepository) GetByID(id uint) (*models.RoleChangeRequest, error) {
    var req models.RoleChangeRequest
    err := r.db.Preload("Requester").
        Preload("CurrentRole").
        Preload("RequestedRole").
        Preload("Admin").
        First(&req, id).Error
    if err != nil {
        return nil, err
    }
    return &req, nil
}

func (r *roleChangeRequestRepository) GetByUserID(userID uint, limit, offset int) ([]*models.RoleChangeRequest, int64, error) {
    var requests []*models.RoleChangeRequest
    var total int64
    
    query := r.db.Where("requester_id = ?", userID)
    query.Model(&models.RoleChangeRequest{}).Count(&total)
    
    err := query.Preload("Requester").
        Preload("CurrentRole").
        Preload("RequestedRole").
        Preload("Admin").
        Limit(limit).Offset(offset).
        Order("created_at DESC").
        Find(&requests).Error
    
    return requests, total, err
}

func (r *roleChangeRequestRepository) GetAll(limit, offset int) ([]*models.RoleChangeRequest, int64, error) {
    var requests []*models.RoleChangeRequest
    var total int64
    
    r.db.Model(&models.RoleChangeRequest{}).Count(&total)
    
    err := r.db.Preload("Requester").
        Preload("CurrentRole").
        Preload("RequestedRole").
        Preload("Admin").
        Limit(limit).Offset(offset).
        Order("created_at DESC").
        Find(&requests).Error
    
    return requests, total, err
}

func (r *roleChangeRequestRepository) GetPending(limit, offset int) ([]*models.RoleChangeRequest, int64, error) {
    var requests []*models.RoleChangeRequest
    var total int64
    
    query := r.db.Where("status = ?", "pending")
    query.Model(&models.RoleChangeRequest{}).Count(&total)
    
    err := query.Preload("Requester").
        Preload("CurrentRole").
        Preload("RequestedRole").
        Preload("Admin").
        Limit(limit).Offset(offset).
        Order("created_at ASC").
        Find(&requests).Error
    
    return requests, total, err
}

func (r *roleChangeRequestRepository) Update(req *models.RoleChangeRequest) error {
    return r.db.Save(req).Error
}

func (r *roleChangeRequestRepository) HasPendingRequest(userID uint) (bool, error) {
    var count int64
    err := r.db.Model(&models.RoleChangeRequest{}).
        Where("requester_id = ? AND status = ?", userID, "pending").
        Count(&count).Error
    return count > 0, err
}
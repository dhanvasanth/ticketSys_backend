// project/internal/services/role_change_request.go
package services

import (
    "errors"
    "time"
    "project/internal/models"
    "project/internal/repositories"
)

type RoleChangeRequestService interface {
    CreateRequest(userID uint, req *models.CreateRoleChangeRequest) (*models.RoleChangeRequest, error)
    GetUserRequests(userID uint, page, limit int) ([]*models.RoleChangeRequest, int64, error)
    GetAllRequests(page, limit int) ([]*models.RoleChangeRequest, int64, error)
    GetPendingRequests(page, limit int) ([]*models.RoleChangeRequest, int64, error)
    GetRequestByID(id uint) (*models.RoleChangeRequest, error)
    ProcessRequest(requestID, adminID uint, req *models.ProcessRoleChangeRequest) (*models.RoleChangeRequest, error)
}

type roleChangeRequestService struct {
    roleChangeRepo repositories.RoleChangeRequestRepository
    userRepo       repositories.UserRepository
}

func NewRoleChangeRequestService(
    roleChangeRepo repositories.RoleChangeRequestRepository,
    userRepo repositories.UserRepository,
) RoleChangeRequestService {
    return &roleChangeRequestService{
        roleChangeRepo: roleChangeRepo,
        userRepo:       userRepo,
    }
}

func (s *roleChangeRequestService) CreateRequest(userID uint, req *models.CreateRoleChangeRequest) (*models.RoleChangeRequest, error) {
    // Check if user has pending request
    hasPending, err := s.roleChangeRepo.HasPendingRequest(userID)
    if err != nil {
        return nil, err
    }
    if hasPending {
        return nil, errors.New("you already have a pending role change request")
    }
    
    // Get current user to get current role
    user, err := s.userRepo.GetByID(userID)
    if err != nil {
        return nil, err
    }
    
    // Check if requested role is different from current role
    if user.RoleID == req.RequestedRoleID {
        return nil, errors.New("requested role is the same as current role")
    }
    
    // Validate requested role exists
    roles, err := s.userRepo.GetRoles()
    if err != nil {
        return nil, err
    }
    
    roleExists := false
    for _, role := range roles {
        if role.ID == req.RequestedRoleID {
            roleExists = true
            break
        }
    }
    
    if !roleExists {
        return nil, errors.New("requested role does not exist")
    }
    
    // Create role change request
    roleChangeReq := &models.RoleChangeRequest{
        RequesterID:     userID,
        CurrentRoleID:   user.RoleID,
        RequestedRoleID: req.RequestedRoleID,
        Reason:          req.Reason,
        Status:          "pending",
    }
    
    if err := s.roleChangeRepo.Create(roleChangeReq); err != nil {
        return nil, err
    }
    
    return s.roleChangeRepo.GetByID(roleChangeReq.ID)
}

func (s *roleChangeRequestService) GetUserRequests(userID uint, page, limit int) ([]*models.RoleChangeRequest, int64, error) {
    offset := (page - 1) * limit
    return s.roleChangeRepo.GetByUserID(userID, limit, offset)
}

func (s *roleChangeRequestService) GetAllRequests(page, limit int) ([]*models.RoleChangeRequest, int64, error) {
    offset := (page - 1) * limit
    return s.roleChangeRepo.GetAll(limit, offset)
}

func (s *roleChangeRequestService) GetPendingRequests(page, limit int) ([]*models.RoleChangeRequest, int64, error) {
    offset := (page - 1) * limit
    return s.roleChangeRepo.GetPending(limit, offset)
}

func (s *roleChangeRequestService) GetRequestByID(id uint) (*models.RoleChangeRequest, error) {
    return s.roleChangeRepo.GetByID(id)
}

func (s *roleChangeRequestService) ProcessRequest(requestID, adminID uint, req *models.ProcessRoleChangeRequest) (*models.RoleChangeRequest, error) {
    // Get the request
    roleChangeReq, err := s.roleChangeRepo.GetByID(requestID)
    if err != nil {
        return nil, err
    }
    
    // Check if request is still pending
    if roleChangeReq.Status != "pending" {
        return nil, errors.New("request has already been processed")
    }
    
    // Update request status
    roleChangeReq.Status = req.Status
    roleChangeReq.AdminID = &adminID
    roleChangeReq.AdminNotes = req.AdminNotes
    now := time.Now()
    roleChangeReq.ProcessedAt = &now
    
    // If approved, update user's role
    if req.Status == "approved" {
        user, err := s.userRepo.GetByID(roleChangeReq.RequesterID)
        if err != nil {
            return nil, err
        }
        
        user.RoleID = roleChangeReq.RequestedRoleID
        if err := s.userRepo.Update(user); err != nil {
            return nil, err
        }
    }
    
    // Save the request
    if err := s.roleChangeRepo.Update(roleChangeReq); err != nil {
        return nil, err
    }
    
    return s.roleChangeRepo.GetByID(roleChangeReq.ID)
}
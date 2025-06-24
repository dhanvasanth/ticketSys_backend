package services

import (
    "project/internal/models"
    "project/internal/repositories"
)

type UserService interface {
    GetProfile(userID uint) (*models.UserResponse, error)
    UpdateProfile(userID uint, req map[string]interface{}) (*models.UserResponse, error)
    GetAllUsers() ([]*models.UserResponse, error)
    GetRoles() ([]*models.Role, error)
    CreateRole(req *models.Role) error
}

type userService struct {
    userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
    return &userService{userRepo: userRepo}
}

func (s *userService) GetProfile(userID uint) (*models.UserResponse, error) {
    user, err := s.userRepo.GetByID(userID)
    if err != nil {
        return nil, err
    }
    
    return &models.UserResponse{
        ID:           user.ID,
        Name:         user.Name,
        Email:        user.Email,
        Role:         user.Role,
        Organization: user.Organization,
    }, nil
}

func (s *userService) UpdateProfile(userID uint, req map[string]interface{}) (*models.UserResponse, error) {
    user, err := s.userRepo.GetByID(userID)
    if err != nil {
        return nil, err
    }
    
    // Update allowed fields
    if name, ok := req["name"]; ok {
        user.Name = name.(string)
    }
    if org, ok := req["organization"]; ok {
        user.Organization = org.(string)
    }
    
    if err := s.userRepo.Update(user); err != nil {
        return nil, err
    }
    
    return s.GetProfile(userID)
}

func (s *userService) GetAllUsers() ([]*models.UserResponse, error) {
    users, err := s.userRepo.GetAll()
    if err != nil {
        return nil, err
    }
    
    var responses []*models.UserResponse
    for _, user := range users {
        responses = append(responses, &models.UserResponse{
            ID:           user.ID,
            Name:         user.Name,
            Email:        user.Email,
            Role:         user.Role,
            Organization: user.Organization,
        })
    }
    
    return responses, nil
}

func (s *userService) GetRoles() ([]*models.Role, error) {
    return s.userRepo.GetRoles()
}

func (s *userService) CreateRole(req *models.Role) error {
    return s.userRepo.CreateRole(req)
}
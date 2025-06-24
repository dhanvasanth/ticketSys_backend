package services

import (
    "encoding/json"
    "errors"
    "time"
    "ticket-service/internal/config"
    "ticket-service/internal/models"
    "ticket-service/internal/repositories"
    
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

type AuthService interface {
    Register(req *models.RegisterRequest) (*models.UserResponse, error)
    Login(req *models.LoginRequest) (string, *models.UserResponse, error)
}

type authService struct {
    userRepo repositories.UserRepository
    jwtCfg   *config.JWTConfig
}

func NewAuthService(userRepo repositories.UserRepository, jwtCfg *config.JWTConfig) AuthService {
    return &authService{
        userRepo: userRepo,
        jwtCfg:   jwtCfg,
    }
}

func (s *authService) Register(req *models.RegisterRequest) (*models.UserResponse, error) {
    // Check if user already exists
    if _, err := s.userRepo.GetByEmail(req.Email); err == nil {
        return nil, errors.New("user already exists")
    }
    
    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    
    // Get default role (user)
    roles, err := s.userRepo.GetRoles()
    if err != nil {
        return nil, err
    }
    
    var defaultRoleID uint = 3 // Default to "user" role
    for _, role := range roles {
        if role.Name == "user" {
            defaultRoleID = role.ID
            break
        }
    }
    
    // Create user
    user := &models.User{
        Name:         req.Name,
        Email:        req.Email,
        Password:     string(hashedPassword),
        RoleID:       defaultRoleID,
        Organization: req.Organization,
        IsActive:     true,
    }
    
    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }
    
    // Get user with role
    createdUser, err := s.userRepo.GetByID(user.ID)
    if err != nil {
        return nil, err
    }
    
    return &models.UserResponse{
        ID:           createdUser.ID,
        Name:         createdUser.Name,
        Email:        createdUser.Email,
        Role:         createdUser.Role,
        Organization: createdUser.Organization,
    }, nil
}

func (s *authService) Login(req *models.LoginRequest) (string, *models.UserResponse, error) {
    // Get user by email
    user, err := s.userRepo.GetByEmail(req.Email)
    if err != nil {
        return "", nil, errors.New("invalid credentials")
    }
    
    // Check password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        return "", nil, errors.New("invalid credentials")
    }
    
    // Check if user is active
    if !user.IsActive {
        return "", nil, errors.New("account is deactivated")
    }
    
    // Parse permissions
    var permissions []string
    if err := json.Unmarshal([]byte(user.Role.Permissions), &permissions); err != nil {
        permissions = []string{}
    }
    
    // Generate JWT token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id":     user.ID,
        "role_id":     user.RoleID,
        "role_name":   user.Role.Name,
        "permissions": permissions,
        "exp":         time.Now().Add(time.Hour * time.Duration(s.jwtCfg.ExpiresHours)).Unix(),
    })
    
    tokenString, err := token.SignedString([]byte(s.jwtCfg.Secret))
    if err != nil {
        return "", nil, err
    }
    
    userResponse := &models.UserResponse{
        ID:           user.ID,
        Name:         user.Name,
        Email:        user.Email,
        Role:         user.Role,
        Organization: user.Organization,
    }
    
    return tokenString, userResponse, nil
}
package services

import (
    "encoding/json"
    "errors"
    "time"
    "project/internal/config"
    "project/internal/models"
    "project/internal/repositories"
    
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

// // project/internal/services/auth.go
// package services

// import (
//     "crypto/rand"
//     "encoding/json"
//     "errors"
//     "fmt"
//     "math/big"
//     "time"
//     "project/internal/config"
//     "project/internal/models"
//     "project/internal/repositories"
   
//     "github.com/golang-jwt/jwt/v5"
//     "golang.org/x/crypto/bcrypt"
// )

// type AuthService interface {
//     Register(req *models.RegisterRequest) (*models.UserResponse, error)
//     Login(req *models.LoginRequest) (string, *models.UserResponse, error)
//     VerifyEmail(req *models.VerifyEmailRequest) (*models.UserResponse, error)
//     ResendVerificationCode(req *models.ResendVerificationRequest) error
// }

// type authService struct {
//     userRepo     repositories.UserRepository
//     jwtCfg       *config.JWTConfig
//     emailService EmailService
// }

// func NewAuthService(userRepo repositories.UserRepository, jwtCfg *config.JWTConfig, emailService EmailService) AuthService {
//     return &authService{
//         userRepo:     userRepo,
//         jwtCfg:       jwtCfg,
//         emailService: emailService,
//     }
// }

// func (s *authService) Register(req *models.RegisterRequest) (*models.UserResponse, error) {
//     // Check if user already exists
//     if _, err := s.userRepo.GetByEmail(req.Email); err == nil {
//         return nil, errors.New("user already exists")
//     }
   
//     // Hash password
//     hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
//     if err != nil {
//         return nil, err
//     }
   
//     // Get default role (user)
//     roles, err := s.userRepo.GetRoles()
//     if err != nil {
//         return nil, err
//     }
   
//     var defaultRoleID uint = 3 // Default to "user" role
//     for _, role := range roles {
//         if role.Name == "user" {
//             defaultRoleID = role.ID
//             break
//         }
//     }
   
//     // Create user (email not verified initially)
//     user := &models.User{
//         Name:          req.Name,
//         Email:         req.Email,
//         Password:      string(hashedPassword),
//         RoleID:        defaultRoleID,
//         Organization:  req.Organization,
//         IsActive:      true,
//         EmailVerified: false,
//     }
   
//     if err := s.userRepo.Create(user); err != nil {
//         return nil, err
//     }
   
//     // Send verification email
//     if err := s.sendVerificationEmail(user.ID, user.Email); err != nil {
//         // Log error but don't fail registration
//         fmt.Printf("Failed to send verification email: %v\n", err)
//     }
   
//     // Get user with role
//     createdUser, err := s.userRepo.GetByID(user.ID)
//     if err != nil {
//         return nil, err
//     }
   
//     return &models.UserResponse{
//         ID:            createdUser.ID,
//         Name:          createdUser.Name,
//         Email:         createdUser.Email,
//         Role:          createdUser.Role,
//         Organization:  createdUser.Organization,
//         EmailVerified: createdUser.EmailVerified,
//     }, nil
// }

// func (s *authService) Login(req *models.LoginRequest) (string, *models.UserResponse, error) {
//     // Get user by email
//     user, err := s.userRepo.GetByEmail(req.Email)
//     if err != nil {
//         return "", nil, errors.New("invalid credentials")
//     }
   
//     // Check password
//     if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
//         return "", nil, errors.New("invalid credentials")
//     }
   
//     // Check if user is active
//     if !user.IsActive {
//         return "", nil, errors.New("account is deactivated")
//     }
   
//     // Note: We're allowing login without email verification
//     // You can uncomment this if you want to require email verification for login
//     // if !user.EmailVerified {
//     //     return "", nil, errors.New("email not verified. Please check your email for verification code")
//     // }
   
//     // Parse permissions
//     var permissions []string
//     if err := json.Unmarshal([]byte(user.Role.Permissions), &permissions); err != nil {
//         permissions = []string{}
//     }
   
//     // Generate JWT token
//     token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
//         "user_id":        user.ID,
//         "role_id":        user.RoleID,
//         "role_name":      user.Role.Name,
//         "permissions":    permissions,
//         "email_verified": user.EmailVerified,
//         "exp":            time.Now().Add(time.Hour * time.Duration(s.jwtCfg.ExpiresHours)).Unix(),
//     })
   
//     tokenString, err := token.SignedString([]byte(s.jwtCfg.Secret))
//     if err != nil {
//         return "", nil, err
//     }
   
//     userResponse := &models.UserResponse{
//         ID:            user.ID,
//         Name:          user.Name,
//         Email:         user.Email,
//         Role:          user.Role,
//         Organization:  user.Organization,
//         EmailVerified: user.EmailVerified,
//     }
   
//     return tokenString, userResponse, nil
// }

// func (s *authService) VerifyEmail(req *models.VerifyEmailRequest) (*models.UserResponse, error) {
//     // Get user by email
//     user, err := s.userRepo.GetByEmail(req.Email)
//     if err != nil {
//         return nil, errors.New("user not found")
//     }
    
//     // Check if already verified
//     if user.EmailVerified {
//         return nil, errors.New("email already verified")
//     }
    
//     // Get verification code
//     verificationCode, err := s.userRepo.GetVerificationCode(user.ID, req.Code)
//     if err != nil {
//         return nil, errors.New("invalid verification code")
//     }
    
//     // Check if code is expired
//     if time.Now().After(verificationCode.ExpiresAt) {
//         return nil, errors.New("verification code expired")
//     }
    
//     // Check if code is already used
//     if verificationCode.Used {
//         return nil, errors.New("verification code already used")
//     }
    
//     // Mark code as used
//     now := time.Now()
//     verificationCode.Used = true
//     verificationCode.UsedAt = &now
//     if err := s.userRepo.UpdateVerificationCode(verificationCode); err != nil {
//         return nil, err
//     }
    
//     // Mark user as verified
//     user.EmailVerified = true
//     user.EmailVerifiedAt = &now
//     if err := s.userRepo.Update(user); err != nil {
//         return nil, err
//     }
    
//     return &models.UserResponse{
//         ID:            user.ID,
//         Name:          user.Name,
//         Email:         user.Email,
//         Role:          user.Role,
//         Organization:  user.Organization,
//         EmailVerified: user.EmailVerified,
//     }, nil
// }

// func (s *authService) ResendVerificationCode(req *models.ResendVerificationRequest) error {
//     // Get user by email
//     user, err := s.userRepo.GetByEmail(req.Email)
//     if err != nil {
//         return errors.New("user not found")
//     }
    
//     // Check if already verified
//     if user.EmailVerified {
//         return errors.New("email already verified")
//     }
    
//     // Send new verification email
//     return s.sendVerificationEmail(user.ID, user.Email)
// }

// func (s *authService) sendVerificationEmail(userID uint, email string) error {
//     // Generate verification code
//     code, err := s.generateVerificationCode()
//     if err != nil {
//         return err
//     }
    
//     // Save verification code to database
//     verificationCode := &models.EmailVerificationCode{
//         UserID:    userID,
//         Email:     email,
//         Code:      code,
//         ExpiresAt: time.Now().Add(15 * time.Minute), // 15 minutes expiry
//         Used:      false,
//     }
    
//     if err := s.userRepo.CreateVerificationCode(verificationCode); err != nil {
//         return err
//     }
    
//     // Send email
//     return s.emailService.SendVerificationEmail(email, code)
// }

// func (s *authService) generateVerificationCode() (string, error) {
//     // Generate 6-digit random code
//     max := big.NewInt(999999)
//     min := big.NewInt(100000)
//     n, err := rand.Int(rand.Reader, max.Sub(max, min).Add(max, big.NewInt(1)))
//     if err != nil {
//         return "", err
//     }
//     return fmt.Sprintf("%06d", n.Add(n, min).Int64()), nil
// }


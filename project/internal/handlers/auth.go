package handlers

import (
    "net/http"
    "project/internal/models"
    "project/internal/services"
    
    "github.com/gin-gonic/gin"
)

type AuthHandler struct {
    authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req models.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    user, err := h.authService.Register(&req)
    if err != nil {
        c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, gin.H{
        "message": "User registered successfully",
        "user":    user,
    })
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req models.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    token, user, err := h.authService.Login(&req)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message": "Login successful",
        "token":   token,
        "user":    user,
    })
}





// // project/internal/handlers/auth.go
// package handlers

// import (
//     "net/http"
//     "project/internal/models"
//     "project/internal/services"
    
//     "github.com/gin-gonic/gin"
// )

// type AuthHandler struct {
//     authService services.AuthService
// }

// func NewAuthHandler(authService services.AuthService) *AuthHandler {
//     return &AuthHandler{authService: authService}
// }

// func (h *AuthHandler) Register(c *gin.Context) {
//     var req models.RegisterRequest
//     if err := c.ShouldBindJSON(&req); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }
    
//     user, err := h.authService.Register(&req)
//     if err != nil {
//         c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
//         return
//     }
    
//     c.JSON(http.StatusCreated, gin.H{
//         "message": "User registered successfully. Please check your email for verification code.",
//         "user":    user,
//     })
// }

// func (h *AuthHandler) Login(c *gin.Context) {
//     var req models.LoginRequest
//     if err := c.ShouldBindJSON(&req); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }
    
//     token, user, err := h.authService.Login(&req)
//     if err != nil {
//         c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
//         return
//     }
    
//     c.JSON(http.StatusOK, gin.H{
//         "message": "Login successful",
//         "token":   token,
//         "user":    user,
//     })
// }

// func (h *AuthHandler) VerifyEmail(c *gin.Context) {
//     var req models.VerifyEmailRequest
//     if err := c.ShouldBindJSON(&req); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }
    
//     user, err := h.authService.VerifyEmail(&req)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }
    
//     c.JSON(http.StatusOK, gin.H{
//         "message": "Email verified successfully",
//         "user":    user,
//     })
// }

// func (h *AuthHandler) ResendVerificationCode(c *gin.Context) {
//     var req models.ResendVerificationRequest
//     if err := c.ShouldBindJSON(&req); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }
    
//     err := h.authService.ResendVerificationCode(&req)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }
    
//     c.JSON(http.StatusOK, gin.H{
//         "message": "Verification code sent successfully",
//     })
// }


















































// package handlers

// import (
//     "net/http"
//     "strconv"
//     "project/internal/models"
//     "project/internal/services"
//     "project/internal/utils"
    
//     "github.com/gin-gonic/gin"
//     "go.uber.org/zap"
// )

// type AuthHandler struct {
//     authService       services.AuthService
//     emailClientService services.EmailClientService
//     logger            *zap.Logger
// }

// func NewAuthHandler(authService services.AuthService, emailClientService services.EmailClientService, logger *zap.Logger) *AuthHandler {
//     return &AuthHandler{
//         authService:       authService,
//         emailClientService: emailClientService,
//         logger:            logger,
//     }
// }

// func (h *AuthHandler) Register(c *gin.Context) {
//     var req models.RegisterRequest
    
//     if err := c.ShouldBindJSON(&req); err != nil {
//         h.logger.Error("Invalid registration request", zap.Error(err))
//         utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
//         return
//     }
    
//     user, err := h.authService.Register(&req)
//     if err != nil {
//         h.logger.Error("Registration failed", zap.Error(err))
//         utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
//         return
//     }
    
//     h.logger.Info("User registered successfully", zap.String("email", user.Email))
//     utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", user)
// }

// func (h *AuthHandler) Login(c *gin.Context) {
//     var req models.LoginRequest
    
//     if err := c.ShouldBindJSON(&req); err != nil {
//         h.logger.Error("Invalid login request", zap.Error(err))
//         utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
//         return
//     }
    
//     tempToken, user, err := h.authService.Login(&req)
//     if err != nil {
//         h.logger.Error("Login failed", zap.Error(err))
//         utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
//         return
//     }
    
//     h.logger.Info("Login initiated, verification code sent", zap.String("email", user.Email))
    
//     response := map[string]interface{}{
//         "message":    "Verification code sent to your email",
//         "temp_token": tempToken,
//         "user":       user,
//         "requires_verification": true,
//     }
    
//     utils.SuccessResponse(c, http.StatusOK, "Verification code sent", response)
// }

// func (h *AuthHandler) VerifyLoginCode(c *gin.Context) {
//     var req models.VerifyLoginCodeRequest
    
//     if err := c.ShouldBindJSON(&req); err != nil {
//         h.logger.Error("Invalid verification request", zap.Error(err))
//         utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
//         return
//     }
    
//     // Verify the code through auth service
//     if err := h.authService.VerifyLoginCode(req.UserID, req.Email, req.Code); err != nil {
//         h.logger.Error("Code verification failed", 
//             zap.Uint("user_id", req.UserID),
//             zap.String("email", req.Email),
//             zap.Error(err))
//         utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
//         return
//     }
    
//     // Generate final JWT token
//     user, err := h.authService.GetUserByID(req.UserID)
//     if err != nil {
//         h.logger.Error("Failed to get user after verification", zap.Error(err))
//         utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to complete login")
//         return
//     }
    
//     token, err := h.authService.GenerateJWTToken(user.ID, user.Email, user.Role.Name)
//     if err != nil {
//         h.logger.Error("Failed to generate JWT token", zap.Error(err))
//         utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
//         return
//     }
    
//     h.logger.Info("Login completed successfully", zap.String("email", user.Email))
    
//     response := map[string]interface{}{
//         "token": token,
//         "user":  user,
//     }
    
//     utils.SuccessResponse(c, http.StatusOK, "Login successful", response)
// }

// func (h *AuthHandler) ResendVerificationCode(c *gin.Context) {
//     userIDStr := c.Param("user_id")
//     userID, err := strconv.ParseUint(userIDStr, 10, 32)
//     if err != nil {
//         utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
//         return
//     }
    
//     var req models.ResendCodeRequest
//     if err := c.ShouldBindJSON(&req); err != nil {
//         h.logger.Error("Invalid resend request", zap.Error(err))
//         utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
//         return
//     }
    
//     if err := h.emailClientService.SendVerificationCode(uint(userID), req.Email); err != nil {
//         h.logger.Error("Failed to resend verification code", zap.Error(err))
//         utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to resend verification code")
//         return
//     }
    
//     h.logger.Info("Verification code resent", 
//         zap.Uint("user_id", uint(userID)), 
//         zap.String("email", req.Email))
    
//     utils.SuccessResponse(c, http.StatusOK, "Verification code resent successfully", nil)
// }

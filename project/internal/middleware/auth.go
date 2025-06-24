package middleware

import (
    "net/http"
    "strings"
    "ticket-service/internal/config"
    
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(cfg *config.JWTConfig) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }
        
        tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
        
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(cfg.Secret), nil
        })
        
        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            c.Set("user_id", uint(claims["user_id"].(float64)))
            c.Set("role_id", uint(claims["role_id"].(float64)))
            c.Set("role_name", claims["role_name"].(string))
            c.Set("permissions", claims["permissions"])
        }
        
        c.Next()
    }
}

// Permission-based middleware
func RequirePermission(permission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        permissions, exists := c.Get("permissions")
        if !exists {
            c.JSON(http.StatusForbidden, gin.H{"error": "No permissions found"})
            c.Abort()
            return
        }
        
        // Check if user has the required permission or admin privileges
        permList, ok := permissions.([]interface{})
        if !ok {
            c.JSON(http.StatusForbidden, gin.H{"error": "Invalid permissions format"})
            c.Abort()
            return
        }
        
        // Admin has all permissions
        for _, perm := range permList {
            if perm.(string) == "*" || perm.(string) == permission {
                c.Next()
                return
            }
        }
        
        c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
        c.Abort()
    }
}
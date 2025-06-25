// package middleware

// import (
//     "net/http"
//     "strings"
//     "project/internal/config"
  
//     "github.com/gin-gonic/gin"
//     "github.com/golang-jwt/jwt/v5"
// )

// func AuthMiddleware(cfg *config.JWTConfig) gin.HandlerFunc {
//     return func(c *gin.Context) {
//         authHeader := c.GetHeader("Authorization")
//         if authHeader == "" {
//             c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
//             c.Abort()
//             return
//         }
        
//         tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
        
//         token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//             return []byte(cfg.Secret), nil
//         })
        
//         if err != nil || !token.Valid {
//             c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
//             c.Abort()
//             return
//         }
        
//         if claims, ok := token.Claims.(jwt.MapClaims); ok {
//             c.Set("user_id", uint(claims["user_id"].(float64)))
//             c.Set("role_id", uint(claims["role_id"].(float64)))
//             c.Set("role_name", claims["role_name"].(string))
//             c.Set("permissions", claims["permissions"])
//         }
        
//         c.Next()
//     }
// }

// // Permission-based middleware
// func RequirePermission(permission string) gin.HandlerFunc {
//     return func(c *gin.Context) {
//         permissions, exists := c.Get("permissions")
//         if !exists {
//             c.JSON(http.StatusForbidden, gin.H{"error": "No permissions found"})
//             c.Abort()
//             return
//         }
        
//         // Check if user has the required permission or admin privileges
//         permList, ok := permissions.([]interface{})
//         if !ok {
//             c.JSON(http.StatusForbidden, gin.H{"error": "Invalid permissions format"})
//             c.Abort()
//             return
//         }
        
//         // Admin has all permissions
//         for _, perm := range permList {
//             if perm.(string) == "*" || perm.(string) == permission {
//                 c.Next()
//                 return
//             }
//         }
        
//         c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
//         c.Abort()
//     }
// }


// func RequirePermission(permission string) gin.HandlerFunc {
//     return func(c *gin.Context) {
//         permissions, exists := c.Get("permissions")

//         if !exists {
//             c.JSON(http.StatusForbidden, gin.H{"error": "No permissions found"})
//             c.Abort()
//             return
//         }
//         if exists {
//             fmt.Printf("User permissions: %+v (type: %T)\n", permissions, permissions)
//         }
//         // Handle nil permissions
//         if permissions == nil {
//             c.JSON(http.StatusForbidden, gin.H{"error": "No permissions assigned"})
//             c.Abort()
//             return
//         }
        
//         // Check if user has the required permission or admin privileges
//         permList, ok := permissions.([]interface{})
//         if !ok {
//             c.JSON(http.StatusForbidden, gin.H{"error": "Invalid permissions format"})
//             c.Abort()
//             return
//         }
        
//         // Check if permissions list is empty
//         if len(permList) == 0 {
//             c.JSON(http.StatusForbidden, gin.H{"error": "No permissions assigned"})
//             c.Abort()
//             return
//         }
        
//         // Check permissions with safe type assertion
//         for _, perm := range permList {
//             if permStr, ok := perm.(string); ok {
//                 if permStr == "*" || permStr == permission {
//                     c.Next()
//                     return
//                 }
//             }
//         }
        
//         c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
//         c.Abort()
//     }
// }


package middleware

import (
    "net/http"
    "strings"
    "project/internal/config"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "go.uber.org/zap"
)

func AuthMiddleware(cfg *config.JWTConfig, logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            logger.Warn("Missing authorization header", 
                zap.String("path", c.Request.URL.Path),
                zap.String("method", c.Request.Method),
                zap.String("client_ip", c.ClientIP()))
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }
        
        tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
        
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(cfg.Secret), nil
        })
        
        if err != nil || !token.Valid {
            logger.Warn("Invalid token provided",
                zap.String("path", c.Request.URL.Path),
                zap.String("method", c.Request.Method),
                zap.String("client_ip", c.ClientIP()),
                zap.Error(err))
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            userID := uint(claims["user_id"].(float64))
            roleID := uint(claims["role_id"].(float64))
            roleName := claims["role_name"].(string)
            permissions := claims["permissions"]
            
            c.Set("user_id", userID)
            c.Set("role_id", roleID)
            c.Set("role_name", roleName)
            c.Set("permissions", permissions)
            
            logger.Info("User authenticated successfully",
                zap.Uint("user_id", userID),
                zap.Uint("role_id", roleID),
                zap.String("role_name", roleName),
                zap.String("path", c.Request.URL.Path),
                zap.String("method", c.Request.Method),
                zap.Any("permissions", permissions))
        }
        
        c.Next()
    }
}

// Permission-based middleware
func RequirePermission(permission string, logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        permissions, exists := c.Get("permissions")
        if !exists {
            logger.Warn("No permissions found in context",
                zap.String("required_permission", permission),
                zap.String("path", c.Request.URL.Path),
                zap.String("method", c.Request.Method))
            c.JSON(http.StatusForbidden, gin.H{"error": "No permissions found"})
            c.Abort()
            return
        }
        
        // Check if user has the required permission or admin privileges
        permList, ok := permissions.([]interface{})
        if !ok {
            logger.Error("Invalid permissions format",
                zap.String("required_permission", permission),
                zap.String("path", c.Request.URL.Path),
                zap.String("method", c.Request.Method),
                zap.Any("permissions", permissions))
            c.JSON(http.StatusForbidden, gin.H{"error": "Invalid permissions format"})
            c.Abort()
            return
        }
        
        // Admin has all permissions
        for _, perm := range permList {
            if perm.(string) == "*" || perm.(string) == permission {
                userID, _ := c.Get("user_id")
                roleName, _ := c.Get("role_name")
                
                logger.Info("Permission granted",
                    zap.Any("user_id", userID),
                    zap.Any("role_name", roleName),
                    zap.String("required_permission", permission),
                    zap.String("granted_permission", perm.(string)),
                    zap.String("path", c.Request.URL.Path),
                    zap.String("method", c.Request.Method),
                    zap.Any("user_permissions", permissions))
                
                c.Next()
                return
            }
        }
        
        userID, _ := c.Get("user_id")
        roleName, _ := c.Get("role_name")
        
        logger.Warn("Permission denied",
            zap.Any("user_id", userID),
            zap.Any("role_name", roleName),
            zap.String("required_permission", permission),
            zap.String("path", c.Request.URL.Path),
            zap.String("method", c.Request.Method),
            zap.Any("user_permissions", permissions))
        
        c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
        c.Abort()
    }
}




func RequirePermissionWithOwnership(generalPerm, ownPerm string, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions, exists := c.Get("permissions")
		if !exists {
			logger.Warn("No permissions found in context",
				zap.String("required_general_permission", generalPerm),
				zap.String("required_own_permission", ownPerm),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method))
			c.JSON(http.StatusForbidden, gin.H{"error": "No permissions found"})
			c.Abort()
			return
		}

		permList, ok := permissions.([]interface{})
		if !ok {
			logger.Error("Invalid permissions format",
				zap.String("required_general_permission", generalPerm),
				zap.String("required_own_permission", ownPerm),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.Any("permissions", permissions))
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid permissions format"})
			c.Abort()
			return
		}

		userID, _ := c.Get("userid")
		roleName, _ := c.Get("role_name")

		// Check if user has admin privileges or general permission
		for _, perm := range permList {
			permStr := perm.(string)
			if permStr == "*" || permStr == generalPerm {
				logger.Info("General permission granted",
					zap.Any("user_id", userID),
					zap.Any("role_name", roleName),
					zap.String("granted_permission", permStr),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method))
				c.Next()
				return
			}
		}

		// Check if user has ownership-based permission
		for _, perm := range permList {
			permStr := perm.(string)
			if permStr == ownPerm {
				// Set flag to indicate ownership validation is required
				c.Set("require_ownership_validation", true)
				c.Set("validated_user_id", userID)
				
				logger.Info("Ownership-based permission granted",
					zap.Any("user_id", userID),
					zap.Any("role_name", roleName),
					zap.String("granted_permission", permStr),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method))
				c.Next()
				return
			}
		}

		logger.Warn("Permission denied",
			zap.Any("user_id", userID),
			zap.Any("role_name", roleName),
			zap.String("required_general_permission", generalPerm),
			zap.String("required_own_permission", ownPerm),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.Any("user_permissions", permissions))

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}

// ValidateOwnership middleware to validate if user owns the resource
func ValidateOwnership(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if ownership validation is required
		requireValidation, exists := c.Get("require_ownership_validation")
		if !exists || !requireValidation.(bool) {
			c.Next()
			return
		}

		validatedUserID, _ := c.Get("validated_user_id")
		
		// Get ticket ID from URL parameter
		ticketIDStr := c.Param("id")
		if ticketIDStr == "" {
			// For POST requests, you might need to get ticket_id from request body
			var requestBody struct {
				TicketID int `json:"ticket_id"`
			}
			if err := c.ShouldBindJSON(&requestBody); err == nil {
				ticketIDStr = strconv.Itoa(requestBody.TicketID)
			}
		}

		if ticketIDStr == "" {
			logger.Error("Ticket ID not found for ownership validation",
				zap.Any("user_id", validatedUserID),
				zap.String("path", c.Request.URL.Path))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket ID required"})
			c.Abort()
			return
		}

		ticketID, err := strconv.Atoi(ticketIDStr)
		if err != nil {
			logger.Error("Invalid ticket ID format",
				zap.String("ticket_id", ticketIDStr),
				zap.Any("user_id", validatedUserID))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
			c.Abort()
			return
		}

		// Here you would typically check the database to verify ownership
		// This is a placeholder - replace with your actual database query
		if !isTicketOwnedByUser(ticketID, validatedUserID) {
			logger.Warn("Ownership validation failed",
				zap.Int("ticket_id", ticketID),
				zap.Any("user_id", validatedUserID),
				zap.String("path", c.Request.URL.Path))
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: not your ticket"})
			c.Abort()
			return
		}

		logger.Info("Ownership validation successful",
			zap.Int("ticket_id", ticketID),
			zap.Any("user_id", validatedUserID),
			zap.String("path", c.Request.URL.Path))

		c.Next()
	}
}

// Placeholder function - replace with your actual database query
func isTicketOwnedByUser(ticketID int, userID interface{}) bool {
	// Example implementation:
	// var ticket Ticket
	// if err := db.Where("id = ? AND user_id = ?", ticketID, userID).First(&ticket).Error; err != nil {
	//     return false
	// }
	// return true
	
	// For now, returning true as placeholder
	return true
}
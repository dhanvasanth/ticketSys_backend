package utils

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

func SuccessResponse(c *gin.Context, code int, message string, data interface{}) {
    c.JSON(code, Response{
        Success: true,
        Message: message,
        Data:    data,
    })
}

func ErrorResponse(c *gin.Context, code int, message string) {
    c.JSON(code, Response{
        Success: false,
        Error:   message,
    })
}

func ValidationError(c *gin.Context, err error) {
    ErrorResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error())
}

func InternalError(c *gin.Context, err error) {
    ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
    // Log the actual error for debugging
    // log.Printf("Internal error: %v", err)
}

func NotFoundError(c *gin.Context, resource string) {
    ErrorResponse(c, http.StatusNotFound, resource+" not found")
}

func UnauthorizedError(c *gin.Context, message string) {
    ErrorResponse(c, http.StatusUnauthorized, message)
}

func ForbiddenError(c *gin.Context, message string) {
    ErrorResponse(c, http.StatusForbidden, message)
}
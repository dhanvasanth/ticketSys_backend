// project/internal/handlers/role_change_request.go
package handlers

import (
    "net/http"
    "strconv"
    "project/internal/models"
    "project/internal/services"
    
    "github.com/gin-gonic/gin"
)

type RoleChangeRequestHandler struct {
    roleChangeService services.RoleChangeRequestService
}

func NewRoleChangeRequestHandler(roleChangeService services.RoleChangeRequestService) *RoleChangeRequestHandler {
    return &RoleChangeRequestHandler{roleChangeService: roleChangeService}
}

func (h *RoleChangeRequestHandler) CreateRequest(c *gin.Context) {
    userID := c.GetUint("user_id")
    
    var req models.CreateRoleChangeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    request, err := h.roleChangeService.CreateRequest(userID, &req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, gin.H{
        "message": "Role change request created successfully",
        "request": request,
    })
}

func (h *RoleChangeRequestHandler) GetMyRequests(c *gin.Context) {
    userID := c.GetUint("user_id")
    
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    
    requests, total, err := h.roleChangeService.GetUserRequests(userID, page, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "requests": requests,
        "total":    total,
        "page":     page,
        "limit":    limit,
    })
}

func (h *RoleChangeRequestHandler) GetAllRequests(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    status := c.Query("status") // pending, approved, rejected, or empty for all
    
    var requests []*models.RoleChangeRequest
    var total int64
    var err error
    
    if status == "pending" {
        requests, total, err = h.roleChangeService.GetPendingRequests(page, limit)
    } else {
        requests, total, err = h.roleChangeService.GetAllRequests(page, limit)
    }
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "requests": requests,
        "total":    total,
        "page":     page,
        "limit":    limit,
        "status":   status,
    })
}

func (h *RoleChangeRequestHandler) GetRequest(c *gin.Context) {
    requestID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
        return
    }
    
    request, err := h.roleChangeService.GetRequestByID(uint(requestID))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
        return
    }
    
    // Check if user can view this request
    userID := c.GetUint("user_id")
    roleName := c.GetString("role_name")
    
    // Admin can see all, user can see only their own
    if roleName != "admin" && request.RequesterID != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"request": request})
}

func (h *RoleChangeRequestHandler) ProcessRequest(c *gin.Context) {
    adminID := c.GetUint("user_id")
    
    requestID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
        return
    }
    
    var req models.ProcessRoleChangeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    request, err := h.roleChangeService.ProcessRequest(uint(requestID), adminID, &req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    message := "Role change request rejected"
    if req.Status == "approved" {
        message = "Role change request approved and user role updated"
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message": message,
        "request": request,
    })
}
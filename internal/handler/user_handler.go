package handler

import (
    "net/http"
    "strconv"
    "BAZ/internal/middleware"
    "BAZ/internal/usecase"
    "github.com/gin-gonic/gin"
)

type UserHandler struct {
    userUsecase *usecase.UserUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
    return &UserHandler{userUsecase: userUsecase}
}

func (h *UserHandler) GetOwnProfile(c *gin.Context) {
    userID := middleware.GetUserID(c)
    user, err := h.userUsecase.GetOwnProfile(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "user not found"})
        return
    }

    c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetPublicProfile(c *gin.Context) {
    userIDStr := c.Param("userId")
    userID, err := strconv.ParseInt(userIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid user id"})
        return
    }

    user, err := h.userUsecase.GetPublicProfile(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "user not found"})
        return
    }

    c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteOwnAccount(c *gin.Context) {
    var req struct {
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "password required"})
        return
    }

    userID := middleware.GetUserID(c)
    err := h.userUsecase.DeleteOwnAccount(c.Request.Context(), userID, req.Password)
    if err != nil {
        if err.Error() == "invalid password" {
            c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": err.Error()})
            return
        }
        c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": err.Error()})
        return
    }

    c.JSON(http.StatusNoContent, nil)
}
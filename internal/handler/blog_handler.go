package handler

import (
    "net/http"
    "strconv"
    "BAZ/internal/domain"
    "BAZ/internal/middleware"
    "BAZ/internal/usecase"
    "github.com/gin-gonic/gin"
)

type BlogHandler struct {
    blogUsecase *usecase.BlogUsecase
}

func NewBlogHandler(blogUsecase *usecase.BlogUsecase) *BlogHandler {
    return &BlogHandler{blogUsecase: blogUsecase}
}

func (h *BlogHandler) Create(c *gin.Context) {
    var req domain.CreateBlogRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
        return
    }
    
    userID := middleware.GetUserID(c)
    blog, err := h.blogUsecase.Create(c.Request.Context(), userID, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, blog)
}

func (h *BlogHandler) GetByID(c *gin.Context) {
    idStr := c.Param("blogId")
    blogID, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid blog id"})
        return
    }
    
    blog, err := h.blogUsecase.GetByID(c.Request.Context(), blogID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "blog not found"})
        return
    }
    
    c.JSON(http.StatusOK, blog)
}

func (h *BlogHandler) GetAll(c *gin.Context) { // Parameter - der Gin Context (enthält Request + Response)
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    
    blogs, total, err := h.blogUsecase.GetAll(c.Request.Context(), page, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "page":  page,
        "limit": limit,
        "total": total,
        "data":  blogs,
    })
}

func (h *BlogHandler) Update(c *gin.Context) {
    idStr := c.Param("blogId")
    blogID, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid blog id"})
        return
    }
    
    var req domain.UpdateBlogRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
        return
    }
    
    userID := middleware.GetUserID(c)
    blog, err := h.blogUsecase.Update(c.Request.Context(), blogID, userID, &req)
    if err != nil {
        if err.Error() == "you are not the author of this blog" {
            c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": err.Error()})
            return
        }
        c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, blog)
}

func (h *BlogHandler) Delete(c *gin.Context) {
    idStr := c.Param("blogId")
    blogID, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid blog id"})
        return
    }
    
    userID := middleware.GetUserID(c)
    err = h.blogUsecase.Delete(c.Request.Context(), blogID, userID)
    if err != nil {
        if err.Error() == "you are not the author of this blog" {
            c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": err.Error()})
            return
        }
        c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": err.Error()})
        return
    }
    
    c.JSON(http.StatusNoContent, nil)
}
package handler

import (
    "net/http"
    "strconv"
    "BAZ/internal/domain"
    "BAZ/internal/middleware"
    "BAZ/internal/usecase"
    "github.com/gin-gonic/gin"
)

type CommentHandler struct {
    commentUsecase *usecase.CommentUsecase
}

func NewCommentHandler(commentUsecase *usecase.CommentUsecase) *CommentHandler {
    return &CommentHandler{commentUsecase: commentUsecase}
}

func (h *CommentHandler) Create(c *gin.Context) {
    blogIDStr := c.Param("blogId")
    blogID, err := strconv.ParseInt(blogIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid blog id"})
        return
    }

    var req domain.CreateCommentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
        return
    }

    userID := middleware.GetUserID(c)
    comment, err := h.commentUsecase.Create(c.Request.Context(), blogID, userID, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, comment)
}

func (h *CommentHandler) GetByID(c *gin.Context) {
    commentIDStr := c.Param("commentId")
    commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid comment id"})
        return
    }

    comment, err := h.commentUsecase.GetByID(c.Request.Context(), commentID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "comment not found"})
        return
    }

    c.JSON(http.StatusOK, comment)
}

func (h *CommentHandler) GetByBlogID(c *gin.Context) {
    blogIDStr := c.Param("blogId")
    blogID, err := strconv.ParseInt(blogIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid blog id"})
        return
    }

    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

    comments, total, err := h.commentUsecase.GetByBlogID(c.Request.Context(), blogID, page, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "page":  page,
        "limit": limit,
        "total": total,
        "data":  comments,
    })
}

func (h *CommentHandler) Update(c *gin.Context) {
    commentIDStr := c.Param("commentId")
    commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid comment id"})
        return
    }

    var req domain.UpdateCommentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
        return
    }

    userID := middleware.GetUserID(c)
    comment, err := h.commentUsecase.Update(c.Request.Context(), commentID, userID, &req)
    if err != nil {
        if err.Error() == "you are not the author of this comment" {
            c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": err.Error()})
            return
        }
        c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": err.Error()})
        return
    }

    c.JSON(http.StatusOK, comment)
}

func (h *CommentHandler) Delete(c *gin.Context) {
    commentIDStr := c.Param("commentId")
    commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid comment id"})
        return
    }

    userID := middleware.GetUserID(c)
    err = h.commentUsecase.Delete(c.Request.Context(), commentID, userID)
    if err != nil {
        if err.Error() == "you are not the author of this comment" {
            c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": err.Error()})
            return
        }
        c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": err.Error()})
        return
    }

    c.JSON(http.StatusNoContent, nil)
}
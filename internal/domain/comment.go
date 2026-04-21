package domain

import "time"

type Comment struct {
	CommentID int64     `json:"commentId"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	BlogID    int64     `json:"blogId"`
	UserID    int64     `json:"userId"`
}

type CreateCommentRequest struct {
	Body string `json:"body" binding:"required"`
}

type UpdateCommentRequest struct {
	Body string `json:"body" binding:"required"`
}

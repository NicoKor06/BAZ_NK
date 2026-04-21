package domain

import "time"

type Blog struct {
	BlogID    int64     `json:"blogId"`
	Headline  string    `json:"headline"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	UserID    int64     `json:"userId"`
}

type CreateBlogRequest struct {
	Headline string `json:"headline" binding:"required"`
	Body     string `json:"body" binding:"required"`
}

type UpdateBlogRequest struct {
	Headline string `json:"headline"`
	Body     string `json:"body"`
}

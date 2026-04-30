package domain

import "time"

type User struct {
	UserID     int64     `json:"userId"`
	Username   string    `json:"username"`
	Firstname  string    `json:"firstname"`
	Lastname   string    `json:"lastname"`
	Email      string    `json:"email"`
	Password   string    `json:"-"`
	Birthday   time.Time `json:"birthday"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	LastOnline time.Time `json:"lastOnline"`
}

type UserPublic struct {
	UserID     int64     `json:"userId"`
	Username   string    `json:"username"`
	Firstname  string    `json:"firstname"`
	Lastname   string    `json:"lastname"`
	Role       string    `json:"role"`
	LastOnline time.Time `json:"lastOnline"`
}

type RegisterRequest struct {
	Username  string    `json:"username" binding:"required"`
	Firstname string    `json:"firstname" binding:"required"`
	Lastname  string    `json:"lastname" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	Password  string    `json:"password" binding:"required,min=6"`
	Birthday  time.Time `json:"birthday"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

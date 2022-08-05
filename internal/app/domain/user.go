package domain

import "time"

type User struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Occupation     string    `json:"occupation"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"password_hash"`
	AvatarFileName string    `json:"avatar_file_name"`
	Role           string    `json:"role"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"update_at"`
}

type UserParam struct {
	Name       string `json:"name" binding:"required"`
	Occupation string `json:"occupation" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required,gte=6"`
}

type UserResponse struct {
	Name       string `json:"name"`
	Occupation string `json:"occupation"`
	Email      string `json:"email"`
}

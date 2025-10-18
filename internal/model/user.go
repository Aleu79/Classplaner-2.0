package model

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username" binding:"required"`
	RoleID    int       `json:"role_id" binding:"required"`
	Role      Role      `json:"role,omitempty"`
	FirstName string    `json:"first_name" binding:"required"`
	LastName  string    `json:"last_name" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	Password  string    `json:"password,omitempty" binding:"required"`
	Addresses []Address `json:"addresses,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}

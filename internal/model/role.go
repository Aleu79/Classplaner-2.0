package model

import "time"

type Role struct {
	ID          int       `json:"id"`
	Name        string    `json:"name" binding:"required" sql:"type:VARCHAR(250)"`
	Description string    `json:"description,omitempty" sql:"type:VARCHAR(250)"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at,omitempty"`
}

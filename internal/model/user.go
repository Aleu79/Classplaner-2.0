package model

import "time"

type User struct {
	ID        int       `json:"id" sql:"primary_key;auto_increment"`
	Username  string    `json:"username" binding:"required" sql:"type:VARCHAR(100)"`
	RoleID    int       `json:"role_id" binding:"required"`
	Role      Role      `json:"role,omitempty" sql:"foreign_key:RoleID;references:role"`
	FirstName string    `json:"first_name" binding:"required" sql:"type:VARCHAR(100)"`
	LastName  string    `json:"last_name" binding:"required" sql:"type:VARCHAR(100)"`
	Email     string    `json:"email" binding:"required,email" sql:"type:VARCHAR(250)"`
	Password  string    `json:"password,omitempty" binding:"required"  sql:"type:VARCHAR(500)"`
	Addresses []Address `json:"addresses,omitempty" sql:"foreign_key:ID;references:address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}

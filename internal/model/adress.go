package model

import "time"

type Address struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id" binding:"required"`
	Name       string    `json:"name,omitempty"`
	IsPrimary  bool      `json:"is_primary"`
	CityID     int       `json:"city_id,omitempty"`
	ProvinceID int       `json:"province_id,omitempty"`
	Address1   string    `json:"address1,omitempty"`
	Address2   string    `json:"address2,omitempty"`
	Phone      string    `json:"phone,omitempty"`
	Email      string    `json:"email,omitempty"`
	PostCode   string    `json:"post_code,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at,omitempty"`
}

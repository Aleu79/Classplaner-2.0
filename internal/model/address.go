package model

import "time"

type Address struct {
	ID         int       `json:"id" sql:"primary_key;auto_increment"`
	UserID     int       `json:"user_id" binding:"required"`
	Name       string    `json:"name,omitempty" sql:"type:VARCHAR(250)"`
	IsPrimary  bool      `json:"is_primary"`
	CityID     int       `json:"city_id,omitempty"`
	ProvinceID int       `json:"province_id,omitempty"`
	Address1   string    `json:"address1,omitempty" sql:"type:VARCHAR(250)"`
	Address2   string    `json:"address2,omitempty" sql:"type:VARCHAR(250)"`
	Phone      string    `json:"phone,omitempty" sql:"type:VARCHAR(250)"`
	Email      string    `json:"email,omitempty" sql:"type:VARCHAR(250)"`
	PostCode   string    `json:"post_code,omitempty" sql:"type:VARCHAR(250)"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at,omitempty"`
}

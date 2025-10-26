package model

import "database/sql"

type Comment struct {
	ID         int64          `json:"id"`
	Task       int64          `json:"id_task"`
	Text       string         `json:"text"`
	UserName   string         `json:"userName"`
	Time       string         `json:"time"`
	User_photo sql.NullString `json:"user_photo"`
}

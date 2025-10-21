package model

import "database/sql"

type Comment struct {
	ID         int            `json:"id"`
	Task       int            `json:"id_task"`
	Text       string         `json:"text"`
	UserName   string         `json:"userName"`
	Time       string         `json:"time"`
	User_photo sql.NullString `json:"user_photo"`
}

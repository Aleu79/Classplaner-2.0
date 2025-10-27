package model

import "database/sql"

type Comment struct {
	ID         int64          `json:"id" sql:"primary_key;auto_increment"`
	Task       Tasks          `json:"id_task" sql:"type:VARCHAR(100)" sql:"foreign_key:id;references:id"` //forange
	Text       string         `json:"text" sql:"type:VARCHAR(500)"`
	UserName   string         `json:"userName" sql:"type:VARCHAR(100)"`
	Time       string         `json:"time" sql:"type:timestamp"`
	User_photo sql.NullString `json:"user_photo" sql:"type:VARCHAR(300)"`
}

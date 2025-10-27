package model

type Submission struct {
	ID           int    `json:"id_submission" form:"id_submission" sql:"primary_key;auto_increment"`
	ID_user      User   `json:"id_user" form:"id_user" sql:"foreign_key:id;references:id"`
	ID_task      Tasks  `json:"id_task" form:"id_task" sql:"foreign_key:id;references:id"`
	File         string `json:"submission_file" form:"submission_file" sql:"type:VARCHAR(500)"`
	Comment      string `json:"submission_comment" form:"submission_comment" sql:"type:VARCHAR(500)"`
	Date         string `json:"submission_date" form:"submission_date" sql:"type:timestamp"`
	Calification string `json:"calification" form:"calification" sql:"type:FLOAT"`
	Feedback     string `json:"feedback" form:"feedback" sql:"type:VARCHAR(500)"`
	Username     string `json:"user_name" sql:"type:VARCHAR(100)"`
	Lastname     string `json:"user_lastname" sql:"type:VARCHAR(100)"`
	Alias        string `json:"user_alias" sql:"type:VARCHAR(100)"`
	Photo        string `json:"user_photo" sql:"type:VARCHAR(500)"`
}

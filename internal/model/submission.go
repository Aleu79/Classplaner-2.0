package model

type Submission struct {
	ID           int    `json:"id_submission" form:"id_submission"`
	ID_user      int    `json:"id_user" form:"id_user"`
	ID_task      int    `json:"id_task" form:"id_task"`
	File         string `json:"submission_file" form:"submission_file"`
	Comment      string `json:"submission_comment" form:"submission_comment"`
	Date         string `json:"submission_date" form:"submission_date"`
	Calification string `json:"calification" form:"calification"`
	Feedback     string `json:"feedback" form:"feedback"`
	Username     string `json:"user_name"`
	Lastname     string `json:"user_lastname"`
	Alias        string `json:"user_alias"`
	Photo        string `json:"user_photo"`
}

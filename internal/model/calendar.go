package model

type Calendar struct {
	Id        int    `sql:"primary_key;auto_increment"`
	Title     string `json:"title" sql:"type:VARCHAR(100)"`
	Desc      string `json:"description" sql:"type:VARCHAR(400)"`
	ID_task   int    `json:"id_task" sql:"foreign_key:id;references:id"`
	Created   string `json:"created_on" sql:"type:timestamp"`
	Deliver   string `json:"deliver_until" sql:"type:timestamp"`
	ClassName string `json:"class_name" sql:"type:VARCHAR(100)"`
	Curso     string `json:"class_curso" sql:"type:VARCHAR(10)"`
}

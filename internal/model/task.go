package model

type Tasks struct {
	ID          int     `json:"id_task" sql:"primary_key;auto_increment"`
	Clase       Classes `json:"id_class" sql:"foreign_key:id;references:id"`
	Titulo      string  `json:"title" sql:"type:VARCHAR(200)"`
	Description string  `json:"description" sql:"type:VARCHAR(1000)"`
	Creado      string  `json:"created_on" sql:"type:timestamp"`
	Limite      string  `json:"deliver_until" sql:"type:timestamp"`
}

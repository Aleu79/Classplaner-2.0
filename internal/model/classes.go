package model

type Classes struct {
	ID       int    `json:"id_class" sql:"primary_key;auto_increment"`
	Name     string `json:"class_name" sql:"type:VARCHAR(100)"`
	Profesor int16  `json:"class_profesor" sql:"type:VARCHAR(100)"` //forange
	Curso    string `json:"class_curso" sql:"type:VARCHAR(100)"`
	Color    string `json:"class_color" sql:"type:VARCHAR(100)"`
	Token    string `json:"class_token" sql:"type:VARCHAR(100)"`
}

type UserClass struct {
	Name     string `json:"user_name" sql:"type:VARCHAR(100)"`
	LastName string `json:"user_lastname" sql:"type:VARCHAR(100)"`
	Photo    string `json:"user_photo" sql:"type:VARCHAR(300)"`
	Type     string `json:"user_type" sql:"type:VARCHAR(100)"`
}

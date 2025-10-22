package model

type Classes struct {
	ID       int16  `json:"id_class"`
	Name     string `json:"class_name"`
	Profesor int16  `json:"class_profesor"`
	Curso    string `json:"class_curso"`
	Color    string `json:"class_color"`
	Token    string `json:"class_token"`
}

type UserClass struct {
	Name     string `json:"user_name"`
	LastName string `json:"user_lastname"`
	Photo    string `json:"user_photo"`
	Type     string `json:"user_type"`
}

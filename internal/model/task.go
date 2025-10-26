package model

type Tasks struct {
	ID          int    `json:"id_task"`
	Clase       int    `json:"id_class"`
	Titulo      string `json:"title"`
	Description string `json:"description"`
	Creado      string `json:"created_on"`
	Limite      string `json:"deliver_until"`
}

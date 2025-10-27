package model

type TaskState string

const (
	StatePending   TaskState = "pendiente"
	StateSubmitted TaskState = "entregada"
	StateLate      TaskState = "atrasada"
)

type Tasks struct {
<<<<<<< HEAD
	ID          int     `json:"id_task" sql:"primary_key;auto_increment"`
	Clase       Classes `json:"id_class" sql:"foreign_key:id;references:id"`
	Titulo      string  `json:"title" sql:"type:VARCHAR(200)"`
	Description string  `json:"description" sql:"type:VARCHAR(1000)"`
	Creado      string  `json:"created_on" sql:"type:timestamp"`
	Limite      string  `json:"deliver_until" sql:"type:timestamp"`
=======
	ID          int       `json:"id_task"`
	Clase       int       `json:"id_class"`
	Titulo      string    `json:"title"`
	Estado      TaskState `json:"estado"`
	Description string    `json:"description"`
	Creado      string    `json:"created_on"`
	Limite      string    `json:"deliver_until"`
>>>>>>> fe019da
}

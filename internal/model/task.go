package model

type TaskState string

const (
	StatePending   TaskState = "pendiente"
	StateSubmitted TaskState = "entregada"
	StateLate      TaskState = "atrasada"
)

type Tasks struct {
	ID          int       `json:"id_task"`
	Clase       int       `json:"id_class"`
	Titulo      string    `json:"title"`
	Estado      TaskState `json:"estado"`
	Description string    `json:"description"`
	Creado      string    `json:"created_on"`
	Limite      string    `json:"deliver_until"`
}

package model

type Calendar struct {
	Title     string `json:"title"`
	Desc      string `json:"description"`
	IDtask    int    `json:"id_task"`
	Created   string `json:"created_on"`
	Deliver   string `json:"deliver_until"`
	ClassName string `json:"class_name"`
	Curso     string `json:"class_curso"`
}

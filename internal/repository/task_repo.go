package repository

import "database/sql"

type TasksRepository interface {
}

type TasksSQL struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TasksRepository {
	return &TasksSQL{db: db}
}

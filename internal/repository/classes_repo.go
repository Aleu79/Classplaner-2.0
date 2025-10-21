package repository

import "database/sql"

type ClassesRepository interface {
}

type ClassesSQL struct {
	db *sql.DB
}

func NewClassesRepository(db *sql.DB) ClassesRepository {
	return &ClassesSQL{db: db}
}

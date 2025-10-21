package repository

import "database/sql"

type SubmissionRepository interface {
}

type SubmissionSQL struct {
	db *sql.DB
}

func NewSubmissionRepository(db *sql.DB) SubmissionRepository {
	return &SubmissionSQL{db: db}
}

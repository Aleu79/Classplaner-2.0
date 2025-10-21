package repository

import "database/sql"

type CommentRepository interface {
}

type CommentSQL struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return CommentSQL{db: db}
}

package repository

import (
	"classplanner/internal/model"
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *model.Comment, userID int64) error
	GetByTaskID(ctx context.Context, taskID int64, limit, offset uint64) ([]*model.Comment, error)
}

type CommentSQL struct {
	db *sql.DB
	sb sq.StatementBuilderType
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &CommentSQL{
		db: db,
		sb: sq.StatementBuilder.PlaceholderFormat(sq.Question),
	}
}

// Crear comentario
func (r *CommentSQL) Create(ctx context.Context, comment *model.Comment, userID int64) error {
	query, args, err := r.sb.
		Insert("comments").
		Columns("id_user", "id_task", "comment").
		Values(userID, comment.Task, comment.Text).
		ToSql()
	if err != nil {
		return fmt.Errorf("error construyendo query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error insertando comentario: %w", err)
	}

	return nil
}

// Obtener comentarios por ID de tarea
func (r *CommentSQL) GetByTaskID(ctx context.Context, taskID int64, limit, offset uint64) ([]*model.Comment, error) {
	builder := r.sb.
		Select(
			"c.id_comment",
			"c.id_task",
			"c.comment",
			"c.created_on",
			"u.user_name",
			"u.user_photo",
		).
		From("comments c").
		Join("users u ON c.id_user = u.id_user").
		Where(sq.Eq{"c.id_task": taskID}).
		OrderBy("c.created_on DESC")

	if limit > 0 {
		builder = builder.Limit(limit)
	}
	if offset > 0 {
		builder = builder.Offset(offset)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error construyendo query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando query: %w", err)
	}
	defer rows.Close()

	var comments []*model.Comment

	for rows.Next() {
		var c model.Comment
		var createdOn []byte

		if err := rows.Scan(
			&c.ID,
			&c.Task,
			&c.Text,
			&createdOn,
			&c.UserName,
			&c.User_photo,
		); err != nil {
			return nil, fmt.Errorf("error leyendo fila: %w", err)
		}

		c.Time = string(createdOn)
		comments = append(comments, &c)
	}

	return comments, nil
}

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
	GetByTaskID(ctx context.Context, taskID int64, limit, offset int) ([]*model.Comment, error)
	GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*model.Comment, error)
	Delete(ctx context.Context, commentID, userID int64, isAdmin bool) error
	Update(ctx context.Context, commentID, userID int64, newText string) error
	CountByTask(ctx context.Context, taskID int64) (int64, error)
	Exists(ctx context.Context, commentID int64) (bool, error)
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
func (r *CommentSQL) GetByTaskID(ctx context.Context, taskID int64, limit, offset int) ([]*model.Comment, error) {
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
		builder = builder.Limit(uint64(limit))
	}
	if offset > 0 {
		builder = builder.Offset(uint64(offset))
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
		var createdOn string

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

		c.Time = createdOn
		comments = append(comments, &c)
	}

	return comments, nil
}

// Obtener comentarios por ID de usuario
func (r *CommentSQL) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*model.Comment, error) {
	builder := r.sb.
		Select("id_comment", "id_task", "comment", "created_on").
		From("comments").
		Where(sq.Eq{"id_user": userID}).
		OrderBy("created_on DESC")

	if limit > 0 {
		builder = builder.Limit(uint64(limit))
	}
	if offset > 0 {
		builder = builder.Offset(uint64(offset))
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
		if err := rows.Scan(&c.ID, &c.Task, &c.Text, &c.Time); err != nil {
			return nil, fmt.Errorf("error escaneando fila: %w", err)
		}
		comments = append(comments, &c)
	}

	return comments, nil
}

// Eliminar comentario (el usuario solo puede eliminar el suyo si no es admin)
func (r *CommentSQL) Delete(ctx context.Context, commentID, userID int64, isAdmin bool) error {
	builder := r.sb.
		Delete("comments").
		Where(sq.Eq{"id_comment": commentID})

	if !isAdmin {
		builder = builder.Where(sq.Eq{"id_user": userID})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("error construyendo delete: %w", err)
	}

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error ejecutando delete: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no se encontró el comentario o no tienes permisos")
	}

	return nil
}

// Actualizar texto del comentario
func (r *CommentSQL) Update(ctx context.Context, commentID, userID int64, newText string) error {
	query, args, err := r.sb.
		Update("comments").
		Set("comment", newText).
		Where(sq.Eq{"id_comment": commentID, "id_user": userID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("error construyendo update: %w", err)
	}

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error ejecutando update: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no se encontró el comentario o no tienes permisos")
	}

	return nil
}

// Contar comentarios por tarea
func (r *CommentSQL) CountByTask(ctx context.Context, taskID int64) (int64, error) {
	query, args, err := r.sb.
		Select("COUNT(*)").
		From("comments").
		Where(sq.Eq{"id_task": taskID}).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("error construyendo query: %w", err)
	}

	var count int64
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error ejecutando query: %w", err)
	}

	return count, nil
}

// Verificar si un comentario existe
func (r *CommentSQL) Exists(ctx context.Context, commentID int64) (bool, error) {
	query, args, err := r.sb.
		Select("1").
		From("comments").
		Where(sq.Eq{"id_comment": commentID}).
		Limit(1).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("error construyendo query: %w", err)
	}

	var exists int
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("error ejecutando query: %w", err)
	}

	return true, nil
}

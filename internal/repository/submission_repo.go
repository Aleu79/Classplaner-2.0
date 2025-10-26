package repository

import (
	"classplanner/internal/model"
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type SubmissionRepository interface {
	Create(ctx context.Context, s *model.Submission) error
	GetByUserAndTask(ctx context.Context, userID, taskID int64) (*model.Submission, error)
	GetByTask(ctx context.Context, taskID int64, limit, offset uint64) ([]*model.Submission, error)
	Update(ctx context.Context, s *model.Submission) error
}

type SubmissionSQL struct {
	db *sql.DB
	sb sq.StatementBuilderType
}

func NewSubmissionRepository(db *sql.DB) SubmissionRepository {
	return &SubmissionSQL{
		db: db,
		sb: sq.StatementBuilder.PlaceholderFormat(sq.Question),
	}
}

// Crear una nueva entrega
func (r *SubmissionSQL) Create(ctx context.Context, s *model.Submission) error {
	query, args, err := r.sb.
		Insert("submissions").
		Columns("id_user", "id_task", "submission_file", "submission_comment", "submission_date").
		Values(s.ID_user, s.ID_task, s.File, s.Comment, time.Now()).
		ToSql()
	if err != nil {
		return fmt.Errorf("error construyendo query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error insertando submission: %w", err)
	}
	return nil
}

// Obtener una entrega específica (por usuario y tarea)
func (r *SubmissionSQL) GetByUserAndTask(ctx context.Context, userID, taskID int64) (*model.Submission, error) {
	query, args, err := r.sb.
		Select("submission_file", "submission_comment", "submission_date", "calification", "feedback").
		From("submissions").
		Where(sq.Eq{"id_user": userID, "id_task": taskID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error construyendo query: %w", err)
	}

	row := r.db.QueryRowContext(ctx, query, args...)

	var s model.Submission
	var calif, feedback sql.NullString
	if err := row.Scan(&s.File, &s.Comment, &s.Date, &calif, &feedback); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error ejecutando query: %w", err)
	}

	if calif.Valid {
		s.Calification = calif.String
	}
	if feedback.Valid {
		s.Feedback = feedback.String
	}

	return &s, nil
}

// Obtener todas las entregas de una tarea (con paginación)
func (r *SubmissionSQL) GetByTask(ctx context.Context, taskID int64, limit, offset uint64) ([]*model.Submission, error) {
	builder := r.sb.
		Select("s.id_submission", "s.id_user", "s.id_task", "s.submission_file", "s.submission_comment", "s.submission_date", "s.calification", "s.feedback", "u.user_name", "u.user_lastname", "u.user_alias", "u.user_photo").
		From("submissions s").
		Join("users u ON s.id_user = u.id_user").
		Where(sq.Eq{"s.id_task": taskID}).
		OrderBy("s.submission_date DESC")

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

	var subs []*model.Submission
	for rows.Next() {
		var s model.Submission
		var date []byte

		if err := rows.Scan(
			&s.ID,
			&s.ID_user,
			&s.ID_task,
			&s.File,
			&s.Comment,
			&date,
			&s.Calification,
			&s.Feedback,
			&s.Username,
			&s.Lastname,
			&s.Alias,
			&s.Photo,
		); err != nil {
			return nil, fmt.Errorf("error escaneando fila: %w", err)
		}

		s.Date = string(date)
		subs = append(subs, &s)
	}

	return subs, nil
}

// Actualizar una entrega existente
func (r *SubmissionSQL) Update(ctx context.Context, s *model.Submission) error {
	query, args, err := r.sb.
		Update("submissions").
		SetMap(map[string]interface{}{
			"submission_comment": s.Comment,
			"submission_file":    s.File,
			"calification":       s.Calification,
			"feedback":           s.Feedback,
		}).
		Where(sq.Eq{"id_submission": s.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("error construyendo query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error actualizando submission: %w", err)
	}
	return nil
}

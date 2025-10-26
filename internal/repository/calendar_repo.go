package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"classplanner/internal/model"

	sq "github.com/Masterminds/squirrel"
)

type CalendarRepositoy interface {
	CalendarByUser(ctx context.Context, userID int, userType string, limit, offset uint64) ([]*model.Calendar, error)
	CalendarByUserAndToken(ctx context.Context, userID int, userType, classToken string, limit, offset uint64) ([]*model.Calendar, error)
}

type CalendarSQL struct {
	db *sql.DB
	sb sq.StatementBuilderType
}

func NewCalendarRepositoy(db *sql.DB) CalendarRepositoy {
	return &CalendarSQL{
		db: db,
		sb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// CalendarByUser obtiene todas las tareas de un usuario según su tipo con paginación
func (r *CalendarSQL) CalendarByUser(ctx context.Context, userID int, userType string, limit, offset uint64) ([]*model.Calendar, error) {
	builder := r.sb.
		Select("tasks.title", "tasks.description", "tasks.id_task", "tasks.created_on", "tasks.deliver_until", "classes.class_name", "classes.class_curso").
		From("tasks").
		Join("classes ON classes.id_class = tasks.id_class").
		Limit(limit).
		Offset(offset)

	switch userType {
	case "docente":
		builder = builder.Where(sq.Eq{"classes.class_profesor": userID})
	case "alumno":
		builder = builder.
			Join("class_users ON class_users.id_class = classes.id_class").
			Join("users ON users.id_user = class_users.id_user").
			Where(sq.Eq{"users.id_user": userID})
	default:
		return nil, fmt.Errorf("CalendarByUser: invalid user type: %s", userType)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("CalendarByUser build query error: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("CalendarByUser query execution error: %w", err)
	}
	defer rows.Close()

	var calendars []*model.Calendar
	for rows.Next() {
		c := &model.Calendar{}
		var created, deliver time.Time

		if err := rows.Scan(&c.Title, &c.Desc, &c.IDtask, &created, &deliver, &c.ClassName, &c.Curso); err != nil {
			return nil, fmt.Errorf("CalendarByUser scan error: %w", err)
		}

		c.Created = created.Format(time.RFC3339)
		c.Deliver = deliver.Format(time.RFC3339)
		calendars = append(calendars, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("CalendarByUser rows iteration error: %w", err)
	}

	return calendars, nil
}

// CalendarByUserAndToken obtiene las tareas de un usuario para una clase específica con paginación
func (r *CalendarSQL) CalendarByUserAndToken(ctx context.Context, userID int, userType, token string, limit, offset uint64) ([]*model.Calendar, error) {
	builder := r.sb.
		Select("tasks.title", "tasks.description", "tasks.id_task", "tasks.created_on", "tasks.deliver_until", "classes.class_name", "classes.class_curso").
		From("tasks").
		Join("classes ON classes.id_class = tasks.id_class").
		Limit(limit).
		Offset(offset)

	switch userType {
	case "docente":
		builder = builder.Where(sq.Eq{"classes.class_profesor": userID})
	case "alumno":
		builder = builder.
			Join("class_users ON class_users.id_class = classes.id_class").
			Join("users ON users.id_user = class_users.id_user").
			Where(sq.Eq{"users.id_user": userID})
	default:
		return nil, fmt.Errorf("CalendarByUserAndToken: invalid user type: %s", userType)
	}

	// Filtrar por token de clase
	builder = builder.Where(sq.Eq{"classes.class_token": token})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("CalendarByUserAndToken build query error: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("CalendarByUserAndToken query execution error: %w", err)
	}
	defer rows.Close()

	var calendars []*model.Calendar
	for rows.Next() {
		c := &model.Calendar{}
		var created, deliver time.Time

		if err := rows.Scan(&c.Title, &c.Desc, &c.IDtask, &created, &deliver, &c.ClassName, &c.Curso); err != nil {
			continue
		}

		c.Created = created.Format(time.RFC3339)
		c.Deliver = deliver.Format(time.RFC3339)
		calendars = append(calendars, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("CalendarByUserAndToken rows iteration error: %w", err)
	}

	return calendars, nil
}

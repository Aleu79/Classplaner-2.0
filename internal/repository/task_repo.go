package repository

import (
	"database/sql"
	"fmt"
	"time"

	"classplanner/internal/model"

	sq "github.com/Masterminds/squirrel"
)

type TasksRepository interface {
	Create(task *model.Tasks) error
	GetByUser(userID int, userType string) ([]*model.Tasks, error)
}

type TasksSQL struct {
	db *sql.DB
	sb sq.StatementBuilderType
}

func NewTaskRepository(db *sql.DB) *TasksSQL {
	return &TasksSQL{
		db: db,
		sb: sq.StatementBuilder.PlaceholderFormat(sq.Question),
	}
}

// Crear nueva tarea
func (r *TasksSQL) Create(task *model.Tasks) error {
	query := r.sb.Insert("tasks").
		Columns("id_class", "title", "description", "created_on", "deliver_until").
		Values(task.Clase, task.Titulo, task.Description, time.Now(), task.Limite)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("error al construir query: %w", err)
	}

	_, err = r.db.Exec(sqlStr, args...)
	if err != nil {
		return fmt.Errorf("error al insertar tarea: %w", err)
	}

	return nil
}

// Obtener tareas según tipo de usuario (alumno o docente)
func (r *TasksSQL) GetByUser(userID int, userType string) ([]*model.Tasks, error) {
	var builder sq.SelectBuilder

	if userType == "alumno" {
		builder = r.sb.
			Select("t.id_task", "t.id_class", "t.title", "t.description", "t.created_on", "t.deliver_until").
			From("tasks t").
			Join("classes c ON c.id_class = t.id_class").
			Join("class_users cu ON cu.id_class = c.id_class").
			Join("users u ON u.id_user = cu.id_user").
			Where(sq.Eq{"u.id_user": userID})
	} else if userType == "docente" {
		builder = r.sb.
			Select("t.id_task", "t.id_class", "t.title", "t.description", "t.created_on", "t.deliver_until").
			From("tasks t").
			Join("classes c ON c.id_class = t.id_class").
			Where(sq.Eq{"c.class_profesor": userID})
	} else {
		return nil, fmt.Errorf("tipo de usuario no válido: %s", userType)
	}

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error al construir query: %w", err)
	}

	rows, err := r.db.Query(sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar query: %w", err)
	}
	defer rows.Close()

	var tasks []*model.Tasks
	for rows.Next() {
		t := &model.Tasks{}
		var createdOn, deliverUntil sql.NullString

		if err := rows.Scan(&t.ID, &t.Clase, &t.Titulo, &t.Description, &createdOn, &deliverUntil); err != nil {
			fmt.Println("error al escanear fila:", err)
			continue
		}

		t.Creado = createdOn.String
		if deliverUntil.Valid {
			t.Limite = deliverUntil.String
		} else {
			t.Limite = ""
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

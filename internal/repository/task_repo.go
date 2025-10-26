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
	GetByClass(classID, limit, offset int) ([]*model.Tasks, error)
	GetByID(taskID int) (*model.Tasks, error)
	Update(task *model.Tasks) error
	Delete(taskID int) error
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

// Inserta una nueva tarea en la base de datos
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

// Obtiene todas las tareas asociadas a un usuario según su tipo ("alumno" o "docente")
func (r *TasksSQL) GetByUser(userID int, userType string) ([]*model.Tasks, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("userID inválido")
	}

	var builder sq.SelectBuilder

	switch userType {
	case "alumno":
		builder = r.sb.
			Select("t.id_task", "t.id_class", "t.title", "t.description", "t.created_on", "t.deliver_until").
			From("tasks t").
			Join("classes c ON c.id_class = t.id_class").
			Join("class_users cu ON cu.id_class = c.id_class").
			Join("users u ON u.id_user = cu.id_user").
			Where(sq.Eq{"u.id_user": userID})

	case "docente":
		builder = r.sb.
			Select("t.id_task", "t.id_class", "t.title", "t.description", "t.created_on", "t.deliver_until").
			From("tasks t").
			Join("classes c ON c.id_class = t.id_class").
			Where(sq.Eq{"c.class_profesor": userID})

	default:
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

// Obtiene todas las tareas de una clase específica, con soporte de paginación (limit y offset)
func (r *TasksSQL) GetByClass(classID int, limit, offset int) ([]*model.Tasks, error) {
	if classID <= 0 {
		return nil, fmt.Errorf("classID inválido")
	}

	builder := r.sb.Select("id_task", "id_class", "title", "description", "created_on", "deliver_until").
		From("tasks").
		Where(sq.Eq{"id_class": classID}).
		OrderBy("created_on DESC")

	if limit > 0 {
		builder = builder.Limit(uint64(limit))
	}
	if offset > 0 {
		builder = builder.Offset(uint64(offset))
	}

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error construyendo query: %w", err)
	}

	rows, err := r.db.Query(sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando query: %w", err)
	}
	defer rows.Close()

	var tasks []*model.Tasks
	for rows.Next() {
		t := &model.Tasks{}
		var createdOn, deliverUntil sql.NullString
		if err := rows.Scan(&t.ID, &t.Clase, &t.Titulo, &t.Description, &createdOn, &deliverUntil); err != nil {
			return nil, fmt.Errorf("error escaneando fila: %w", err)
		}
		t.Creado = createdOn.String
		if deliverUntil.Valid {
			t.Limite = deliverUntil.String
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

// Obtiene una tarea específica por su ID
func (r *TasksSQL) GetByID(taskID int) (*model.Tasks, error) {
	if taskID <= 0 {
		return nil, fmt.Errorf("taskID inválido")
	}

	builder := r.sb.Select("id_task", "id_class", "title", "description", "created_on", "deliver_until").
		From("tasks").
		Where(sq.Eq{"id_task": taskID})

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error construyendo query: %w", err)
	}

	row := r.db.QueryRow(sqlStr, args...)
	t := &model.Tasks{}
	var createdOn, deliverUntil sql.NullString
	if err := row.Scan(&t.ID, &t.Clase, &t.Titulo, &t.Description, &createdOn, &deliverUntil); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error escaneando fila: %w", err)
	}

	t.Creado = createdOn.String
	if deliverUntil.Valid {
		t.Limite = deliverUntil.String
	}

	return t, nil
}

// Update modifica una tarea existente
func (r *TasksSQL) Update(task *model.Tasks) error {
	if task == nil || task.ID <= 0 {
		return fmt.Errorf("task inválido")
	}

	builder := r.sb.Update("tasks").
		SetMap(map[string]interface{}{
			"title":         task.Titulo,
			"description":   task.Description,
			"deliver_until": task.Limite,
		}).
		Where(sq.Eq{"id_task": task.ID})

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("error construyendo query: %w", err)
	}

	if _, err := r.db.Exec(sqlStr, args...); err != nil {
		return fmt.Errorf("error actualizando tarea: %w", err)
	}

	return nil
}

// Delete elimina una tarea por su ID
func (r *TasksSQL) Delete(taskID int) error {
	if taskID <= 0 {
		return fmt.Errorf("taskID inválido")
	}

	builder := r.sb.Delete("tasks").
		Where(sq.Eq{"id_task": taskID})

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("error construyendo query: %w", err)
	}

	if _, err := r.db.Exec(sqlStr, args...); err != nil {
		return fmt.Errorf("error eliminando tarea: %w", err)
	}

	return nil
}

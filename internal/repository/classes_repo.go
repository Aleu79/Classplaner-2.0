package repository

import (
	"context"
	"database/sql"
	"fmt"

	"classplanner/internal/model"

	sq "github.com/Masterminds/squirrel"
)

type ClassesRepository interface {
	CreateClass(ctx context.Context, class *model.Classes) error
	ClassesByTeacher(ctx context.Context, teacherID int) ([]*model.Classes, error)
	ClassesByStudent(ctx context.Context, studentID int) ([]*model.Classes, error)
	JoinClass(ctx context.Context, userID int, token string) error
	UsersFromClass(ctx context.Context, classID int) ([]*model.UserClass, error)
}

type ClassesSQL struct {
	db *sql.DB
	sb sq.StatementBuilderType
}

func NewClassesRepository(db *sql.DB) ClassesRepository {
	return &ClassesSQL{
		db: db,
		sb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Crear clase
func (r *ClassesSQL) CreateClass(ctx context.Context, class *model.Classes) error {
	query, args, err := r.sb.
		Insert("classes").
		Columns("class_name", "class_profesor", "class_curso", "class_color", "class_token").
		Values(class.Name, class.Profesor, class.Curso, class.Color, class.Token).
		ToSql()
	if err != nil {
		return fmt.Errorf("CreateClass build query error: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("CreateClass execution error: %w", err)
	}
	return nil
}

// Clases por profesor
func (r *ClassesSQL) ClassesByTeacher(ctx context.Context, teacherID int) ([]*model.Classes, error) {
	query, args, err := r.sb.
		Select("id_class", "class_name", "class_profesor", "class_curso", "class_color", "class_token").
		From("classes").
		Where(sq.Eq{"class_profesor": teacherID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ClassesByTeacher build query error: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ClassesByTeacher query error: %w", err)
	}
	defer rows.Close()

	var classes []*model.Classes
	for rows.Next() {
		c := &model.Classes{}
		if err := rows.Scan(&c.ID, &c.Name, &c.Profesor, &c.Curso, &c.Color, &c.Token); err != nil {
			return nil, fmt.Errorf("ClassesByTeacher scan error: %w", err)
		}
		classes = append(classes, c)
	}

	return classes, nil
}

// Clases por alumno
func (r *ClassesSQL) ClassesByStudent(ctx context.Context, studentID int) ([]*model.Classes, error) {
	query, args, err := r.sb.
		Select("c.id_class", "c.class_name", "c.class_profesor", "c.class_curso", "c.class_color", "c.class_token").
		From("class_users cu").
		Join("classes c ON cu.id_class = c.id_class").
		Where(sq.Eq{"cu.id_user": studentID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ClassesByStudent build query error: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ClassesByStudent query error: %w", err)
	}
	defer rows.Close()

	var classes []*model.Classes
	for rows.Next() {
		c := &model.Classes{}
		if err := rows.Scan(&c.ID, &c.Name, &c.Profesor, &c.Curso, &c.Color, &c.Token); err != nil {
			return nil, fmt.Errorf("ClassesByStudent scan error: %w", err)
		}
		classes = append(classes, c)
	}

	return classes, nil
}

// Unirse a clase
func (r *ClassesSQL) JoinClass(ctx context.Context, userID int, token string) error {
	var classID int
	err := r.db.QueryRowContext(ctx, "SELECT id_class FROM classes WHERE class_token = $1", token).Scan(&classID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("JoinClass: class not found")
	}
	if err != nil {
		return fmt.Errorf("JoinClass query error: %w", err)
	}

	query, args, err := r.sb.
		Insert("class_users").
		Columns("id_user", "id_class").
		Values(userID, classID).
		ToSql()
	if err != nil {
		return fmt.Errorf("JoinClass build insert query error: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("JoinClass execution insert error: %w", err)
	}

	return nil
}

// Usuarios de una clase
func (r *ClassesSQL) UsersFromClass(ctx context.Context, classID int) ([]*model.UserClass, error) {
	studentsBuilder := r.sb.
		Select("u.user_name", "u.user_lastname", "u.user_photo", "u.user_type").
		From("class_users cu").
		Join("users u ON u.id_user = cu.id_user").
		Where(sq.Eq{"cu.id_class": classID})

	profBuilder := r.sb.
		Select("u.user_name", "u.user_lastname", "u.user_photo", "u.user_type").
		From("classes c").
		Join("users u ON u.id_user = c.class_profesor").
		Where(sq.Eq{"c.id_class": classID})

	q := fmt.Sprintf("(%s) UNION (%s)", mustToSql(studentsBuilder), mustToSql(profBuilder))

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("UsersFromClass query error: %w", err)
	}
	defer rows.Close()

	var users []*model.UserClass
	for rows.Next() {
		u := &model.UserClass{}
		var photo sql.NullString
		if err := rows.Scan(&u.Name, &u.LastName, &photo, &u.Type); err != nil {
			return nil, fmt.Errorf("UsersFromClass scan error: %w", err)
		}
		u.Photo = photo.String
		users = append(users, u)
	}

	return users, nil
}

// helper para convertir un builder a SQL y panic si falla (solo para UNION)
func mustToSql(builder sq.SelectBuilder) string {
	query, _, err := builder.ToSql()
	if err != nil {
		panic(fmt.Sprintf("mustToSql error: %v", err))
	}
	return query
}

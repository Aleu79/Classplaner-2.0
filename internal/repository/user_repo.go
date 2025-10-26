package repository

import (
	"context"
	"database/sql"
	"fmt"

	"classplanner/internal/model"

	sq "github.com/Masterminds/squirrel"
)

type UserRepository interface {
	GetAll() ([]*model.User, error)
	GetByID(ctx context.Context, id int) (*model.User, error)
	GetByEmailOrUser(ctx context.Context, user string) (*model.User, error)
	SearchByUserOrEmail(ctx context.Context, user string) ([]*model.User, error)
	Exists(ctx context.Context, id int) (bool, error)
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	Update(ctx context.Context, id int, user *model.User) (*model.User, error)
	Delete(ctx context.Context, id int) error
}

type userSQL struct {
	db *sql.DB
	sb sq.StatementBuilderType
}

func NewUserRepository(db *sql.DB) UserRepository {
	sb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &userSQL{db: db, sb: sb}
}

// Obtener todos los usuarios
func (r *userSQL) GetAll() ([]*model.User, error) {
	query, args, err := r.sb.
		Select("id", "username", "role_id", "first_name", "last_name", "email", "created_at", "updated_at").
		From("users").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetAll build query error: %w", err)
	}

	rows, err := r.db.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, fmt.Errorf("GetAll query error: %w", err)
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		u := &model.User{}
		if err := rows.Scan(&u.ID, &u.Username, &u.RoleID, &u.FirstName, &u.LastName, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("GetAll scan error: %w", err)
		}
		users = append(users, u)
	}

	return users, nil
}

// Obtener usuario por ID
func (r *userSQL) GetByID(ctx context.Context, id int) (*model.User, error) {
	query, args, err := r.sb.
		Select("id", "username", "role_id", "first_name", "last_name", "email", "password", "created_at", "updated_at").
		From("users").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetByID build query error: %w", err)
	}

	row := r.db.QueryRowContext(ctx, query, args...)
	u := &model.User{}
	if err := row.Scan(&u.ID, &u.Username, &u.RoleID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("GetByID scan error: %w", err)
	}
	return u, nil
}

// Obtener usuario por username o email
func (r *userSQL) GetByEmailOrUser(ctx context.Context, user string) (*model.User, error) {
	query, args, err := r.sb.
		Select("id", "username", "role_id", "first_name", "last_name", "email", "password", "created_at", "updated_at").
		From("users").
		Where(sq.Or{sq.Eq{"username": user}, sq.Eq{"email": user}}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetByEmailOrUser build query error: %w", err)
	}

	row := r.db.QueryRowContext(ctx, query, args...)
	u := &model.User{}
	if err := row.Scan(&u.ID, &u.Username, &u.RoleID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("GetByEmailOrUser scan error: %w", err)
	}
	return u, nil
}

// Buscar usuarios por nombre o email (ILIKE)
func (r *userSQL) SearchByUserOrEmail(ctx context.Context, user string) ([]*model.User, error) {
	pattern := "%" + user + "%"
	query, args, err := r.sb.
		Select("id", "username", "role_id", "first_name", "last_name", "email", "created_at", "updated_at").
		From("users").
		Where(sq.Or{sq.ILike{"username": pattern}, sq.ILike{"email": pattern}}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("SearchByUserOrEmail build query error: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("SearchByUserOrEmail query error: %w", err)
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		u := &model.User{}
		if err := rows.Scan(&u.ID, &u.Username, &u.RoleID, &u.FirstName, &u.LastName, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("SearchByUserOrEmail scan error: %w", err)
		}
		users = append(users, u)
	}

	return users, nil
}

// Verifica si existe un usuario por ID
func (r *userSQL) Exists(ctx context.Context, id int) (bool, error) {
	query, args, err := r.sb.Select("1").From("users").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return false, fmt.Errorf("Exists build query error: %w", err)
	}

	var exists int
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("Exists query error: %w", err)
	}
	return true, nil
}

// Crear usuario
func (r *userSQL) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	query, args, err := r.sb.
		Insert("users").
		Columns("username", "email", "password", "role_id", "first_name", "last_name", "created_at", "updated_at").
		Values(user.Username, user.Email, user.Password, user.RoleID, user.FirstName, user.LastName, sq.Expr("NOW()"), sq.Expr("NOW()")).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("CreateUser build query error: %w", err)
	}

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("CreateUser query error: %w", err)
	}
	return user, nil
}

// Actualizar usuario
func (r *userSQL) Update(ctx context.Context, id int, user *model.User) (*model.User, error) {
	query, args, err := r.sb.
		Update("users").
		SetMap(map[string]interface{}{
			"username":   user.Username,
			"email":      user.Email,
			"role_id":    user.RoleID,
			"password":   user.Password,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"updated_at": sq.Expr("NOW()"),
		}).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("Update build query error: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("Update exec error: %w", err)
	}

	user.ID = id
	return user, nil
}

// Eliminar usuario
func (r *userSQL) Delete(ctx context.Context, id int) error {
	query, args, err := r.sb.Delete("users").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("Delete build query error: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("Delete exec error: %w", err)
	}
	return nil
}

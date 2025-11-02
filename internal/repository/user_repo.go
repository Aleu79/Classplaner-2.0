package repository

import (
	"context"
	"database/sql"

	"classplanner/internal/model"
	"classplanner/pkg/errorsper"

	sq "github.com/Masterminds/squirrel"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]*model.User, *errorsper.AppError)
	GetByID(ctx context.Context, id int) (*model.User, *errorsper.AppError)
	GetByEmailOrUser(ctx context.Context, user string) (*model.User, *errorsper.AppError)
	SearchByUserOrEmail(ctx context.Context, user string) ([]*model.User, *errorsper.AppError)
	Exists(ctx context.Context, id int) (bool, *errorsper.AppError)
	CreateUser(ctx context.Context, user *model.User) (*model.User, *errorsper.AppError)
	Update(ctx context.Context, id int, user *model.User) (*model.User, *errorsper.AppError)
	Delete(ctx context.Context, id int) *errorsper.AppError
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
func (r *userSQL) GetAll(ctx context.Context) ([]*model.User, *errorsper.AppError) {
	query, args, err := r.sb.Select("id", "username", "role_id", "first_name", "last_name", "email", "created_at", "updated_at").
		From("users").ToSql()
	if err != nil {
		return nil, errorsper.ErrInternal(err, "UserRepository.GetAll ToSql")
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errorsper.ErrInternal(err, "UserRepository.GetAll Query")
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		u := &model.User{}
		if err := rows.Scan(&u.ID, &u.Username, &u.RoleID, &u.FirstName, &u.LastName, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, errorsper.ErrInternal(err, "UserRepository.GetAll Scan")
		}
		users = append(users, u)
	}
	return users, nil
}

// Obtener usuario por ID
func (r *userSQL) GetByID(ctx context.Context, id int) (*model.User, *errorsper.AppError) {
	query, args, err := r.sb.Select("id", "username", "role_id", "first_name", "last_name", "email", "password", "created_at", "updated_at").
		From("users").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, errorsper.ErrInternal(err, "UserRepository.GetByID ToSql")
	}

	row := r.db.QueryRowContext(ctx, query, args...)
	u := &model.User{}
	if err := row.Scan(&u.ID, &u.Username, &u.RoleID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errorsper.ErrNotFound("Usuario no encontrado", "UserRepository.GetByID")
		}
		return nil, errorsper.ErrInternal(err, "UserRepository.GetByID Scan")
	}
	return u, nil
}

// Obtener usuario por username o email
func (r *userSQL) GetByEmailOrUser(ctx context.Context, user string) (*model.User, *errorsper.AppError) {
	query, args, err := r.sb.Select("id", "username", "role_id", "first_name", "last_name", "email", "password", "created_at", "updated_at").
		From("users").Where(sq.Or{sq.Eq{"username": user}, sq.Eq{"email": user}}).ToSql()
	if err != nil {
		return nil, errorsper.ErrInternal(err, "UserRepository.GetByEmailOrUser ToSql")
	}

	row := r.db.QueryRowContext(ctx, query, args...)
	u := &model.User{}
	if err := row.Scan(&u.ID, &u.Username, &u.RoleID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // usuario no encontrado no es error crítico
		}
		return nil, errorsper.ErrInternal(err, "UserRepository.GetByEmailOrUser Scan")
	}
	return u, nil
}

// Buscar usuarios por nombre o email
func (r *userSQL) SearchByUserOrEmail(ctx context.Context, user string) ([]*model.User, *errorsper.AppError) {
	pattern := "%" + user + "%"
	query, args, err := r.sb.Select("id", "username", "role_id", "first_name", "last_name", "email", "created_at", "updated_at").
		From("users").Where(sq.Or{sq.ILike{"username": pattern}, sq.ILike{"email": pattern}}).ToSql()
	if err != nil {
		return nil, errorsper.ErrInternal(err, "UserRepository.SearchByUserOrEmail ToSql")
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errorsper.ErrInternal(err, "UserRepository.SearchByUserOrEmail Query")
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		u := &model.User{}
		if err := rows.Scan(&u.ID, &u.Username, &u.RoleID, &u.FirstName, &u.LastName, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, errorsper.ErrInternal(err, "UserRepository.SearchByUserOrEmail Scan")
		}
		users = append(users, u)
	}
	return users, nil
}

// Verifica si existe un usuario
func (r *userSQL) Exists(ctx context.Context, id int) (bool, *errorsper.AppError) {
	query, args, err := r.sb.Select("1").From("users").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return false, errorsper.ErrInternal(err, "UserRepository.Exists ToSql")
	}

	var exists int
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, errorsper.ErrInternal(err, "UserRepository.Exists Scan")
	}
	return true, nil
}

// Crear usuario
func (r *userSQL) CreateUser(ctx context.Context, user *model.User) (*model.User, *errorsper.AppError) {
	query, args, err := r.sb.Insert("users").
		Columns("username", "email", "password", "role_id", "first_name", "last_name", "created_at", "updated_at").
		Values(user.Username, user.Email, user.Password, user.RoleID, user.FirstName, user.LastName, sq.Expr("NOW()"), sq.Expr("NOW()")).
		Suffix("RETURNING id").ToSql()
	if err != nil {
		return nil, errorsper.ErrInternal(err, "UserRepository.CreateUser ToSql")
	}

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&user.ID)
	if err != nil {
		return nil, errorsper.ErrInternal(err, "UserRepository.CreateUser Exec")
	}
	return user, nil
}

// Actualizar usuario
func (r *userSQL) Update(ctx context.Context, id int, user *model.User) (*model.User, *errorsper.AppError) {
	query, args, err := r.sb.Update("users").
		SetMap(map[string]interface{}{
			"username":   user.Username,
			"email":      user.Email,
			"role_id":    user.RoleID,
			"password":   user.Password,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"updated_at": sq.Expr("NOW()"),
		}).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, errorsper.ErrInternal(err, "UserRepository.Update ToSql")
	}

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, errorsper.ErrInternal(err, "UserRepository.Update Exec")
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return nil, errorsper.ErrNotFound("Usuario no encontrado", "UserRepository.Update")
	}

	user.ID = id
	return user, nil
}

// Eliminar usuario
func (r *userSQL) Delete(ctx context.Context, id int) *errorsper.AppError {
	query, args, err := r.sb.Delete("users").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return errorsper.ErrInternal(err, "UserRepository.Delete ToSql")
	}

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return errorsper.ErrInternal(err, "UserRepository.Delete Exec")
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errorsper.ErrNotFound("Usuario no encontrado", "UserRepository.Delete")
	}

	return nil
}

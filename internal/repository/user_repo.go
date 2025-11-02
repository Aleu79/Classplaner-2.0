package repository

import (
	"context"
	"database/sql"

	"classplanner/internal/infrastructure/logger"
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
	db     *sql.DB
	sb     sq.StatementBuilderType
	logger *logger.Logger
}

func NewUserRepository(db *sql.DB, logger *logger.Logger) UserRepository {
	sb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &userSQL{db: db, sb: sb, logger: logger}
}

// GetAll obtiene todos los usuarios
func (r *userSQL) GetAll(ctx context.Context) ([]*model.User, *errorsper.AppError) {
	ctx = logger.EnsureContext(ctx)
	r.logger.Info(ctx, "UserRepository.GetAll - iniciando consulta")

	query, args, err := r.sb.Select("u.id", "u.username", "u.role_id", "r.name AS role", "u.first_name", "u.last_name",
		"u.email", "u.phone", "u.created_at", "u.updated_at").
		From("users u").
		LeftJoin("roles r ON r.id = u.role_id").
		ToSql()
	if err != nil {
		r.logger.Error(ctx, "GetAll - error generando query: %v", err)
		return nil, errorsper.ErrInternal(err, "UserRepository.GetAll ToSql")
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error(ctx, "GetAll - error ejecutando query: %v", err)
		return nil, errorsper.ErrInternal(err, "UserRepository.GetAll Query")
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		u := &model.User{}
		if err := rows.Scan(&u.ID, &u.Username, &u.RoleID, &u.Role, &u.FirstName, &u.LastName,
			&u.Email, &u.Phone, &u.CreatedAt, &u.UpdatedAt); err != nil {
			r.logger.Error(ctx, "GetAll - error scanneando fila: %v", err)
			return nil, errorsper.ErrInternal(err, "UserRepository.GetAll Scan")
		}
		users = append(users, u)
	}

	r.logger.Info(ctx, "GetAll - %d usuarios obtenidos", len(users))
	return users, nil
}

// Obtener usuario por ID
func (r *userSQL) GetByID(ctx context.Context, id int) (*model.User, *errorsper.AppError) {
	ctx = logger.EnsureContext(ctx)
	r.logger.Info(ctx, "GetByID - id=%d", id)

	query, args, err := r.sb.Select("u.id", "u.username", "u.role_id", "r.name AS role", "u.first_name", "u.last_name",
		"u.email", "u.password", "u.phone", "u.created_at", "u.updated_at").
		From("users u").
		LeftJoin("roles r ON r.id = u.role_id").
		Where(sq.Eq{"u.id": id}).
		ToSql()
	if err != nil {
		r.logger.Error(ctx, "GetByID - error ToSql: %v", err)
		return nil, errorsper.ErrInternal(err, "UserRepository.GetByID ToSql")
	}

	row := r.db.QueryRowContext(ctx, query, args...)
	u := &model.User{}
	if err := row.Scan(&u.ID, &u.Username, &u.RoleID, &u.Role, &u.FirstName, &u.LastName,
		&u.Email, &u.Password, &u.Phone, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn(ctx, "GetByID - user not found id=%d", id)
			return nil, errorsper.ErrNotFound("Usuario no encontrado", "UserRepository.GetByID")
		}
		r.logger.Error(ctx, "GetByID - scan error: %v", err)
		return nil, errorsper.ErrInternal(err, "UserRepository.GetByID Scan")
	}

	r.logger.Info(ctx, "GetByID - éxito id=%d", id)
	return u, nil
}

// Obtener usuario por username o email
func (r *userSQL) GetByEmailOrUser(ctx context.Context, user string) (*model.User, *errorsper.AppError) {
	ctx = logger.EnsureContext(ctx)
	r.logger.Info(ctx, "GetByEmailOrUser - user=%s", user)

	query, args, err := r.sb.Select("id", "username", "role_id", "first_name", "last_name",
		"email", "password", "phone", "created_at", "updated_at").
		From("users").
		Where(sq.Or{sq.Eq{"username": user}, sq.Eq{"email": user}}).
		ToSql()
	if err != nil {
		r.logger.Error(ctx, "GetByEmailOrUser - error ToSql: %v", err)
		return nil, errorsper.ErrInternal(err, "UserRepository.GetByEmailOrUser ToSql")
	}

	row := r.db.QueryRowContext(ctx, query, args...)
	u := &model.User{}
	if err := row.Scan(&u.ID, &u.Username, &u.RoleID, &u.FirstName, &u.LastName,
		&u.Email, &u.Password, &u.Phone, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn(ctx, "GetByEmailOrUser - user not found: %s", user)
			return nil, nil // usuario no encontrado no es un error critico
		}
		r.logger.Error(ctx, "GetByEmailOrUser - scan error: %v", err)
		return nil, errorsper.ErrInternal(err, "UserRepository.GetByEmailOrUser Scan")
	}

	r.logger.Info(ctx, "GetByEmailOrUser - éxito user=%s", user)
	return u, nil
}

// Buscar usuarios por nombre o email
func (r *userSQL) SearchByUserOrEmail(ctx context.Context, user string) ([]*model.User, *errorsper.AppError) {
	ctx = logger.EnsureContext(ctx)
	r.logger.Info(ctx, "SearchByUserOrEmail - user=%s", user)

	pattern := "%" + user + "%"
	query, args, err := r.sb.Select("id", "username", "role_id", "first_name", "last_name",
		"email", "phone", "created_at", "updated_at").
		From("users").
		Where(sq.Or{sq.ILike{"username": pattern}, sq.ILike{"email": pattern}}).
		ToSql()
	if err != nil {
		r.logger.Error(ctx, "SearchByUserOrEmail - error ToSql: %v", err)
		return nil, errorsper.ErrInternal(err, "UserRepository.SearchByUserOrEmail ToSql")
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error(ctx, "SearchByUserOrEmail - error query: %v", err)
		return nil, errorsper.ErrInternal(err, "UserRepository.SearchByUserOrEmail Query")
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		u := &model.User{}
		if err := rows.Scan(&u.ID, &u.Username, &u.RoleID, &u.FirstName, &u.LastName,
			&u.Email, &u.Phone, &u.CreatedAt, &u.UpdatedAt); err != nil {
			r.logger.Error(ctx, "SearchByUserOrEmail - scan error: %v", err)
			return nil, errorsper.ErrInternal(err, "UserRepository.SearchByUserOrEmail Scan")
		}
		users = append(users, u)
	}

	r.logger.Info(ctx, "SearchByUserOrEmail - %d resultados encontrados", len(users))
	return users, nil
}

// Verifica si existe un usuario
func (r *userSQL) Exists(ctx context.Context, id int) (bool, *errorsper.AppError) {
	ctx = logger.EnsureContext(ctx)
	r.logger.Debug(ctx, "Exists - id=%d", id)

	query, args, err := r.sb.Select("1").From("users").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		r.logger.Error(ctx, "Exists - error ToSql: %v", err)
		return false, errorsper.ErrInternal(err, "UserRepository.Exists ToSql")
	}

	var exists int
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		r.logger.Debug(ctx, "Exists - no existe id=%d", id)
		return false, nil
	}
	if err != nil {
		r.logger.Error(ctx, "Exists - scan error: %v", err)
		return false, errorsper.ErrInternal(err, "UserRepository.Exists Scan")
	}

	r.logger.Debug(ctx, "Exists - existe id=%d", id)
	return true, nil
}

// Crear usuario
func (r *userSQL) CreateUser(ctx context.Context, user *model.User) (*model.User, *errorsper.AppError) {
	ctx = logger.EnsureContext(ctx)
	r.logger.Info(ctx, "CreateUser - username=%s email=%s", user.Username, user.Email)

	query, args, err := r.sb.Insert("users").
		Columns("username", "email", "password", "role_id", "first_name", "last_name", "phone", "created_at", "updated_at").
		Values(user.Username, user.Email, user.Password, user.RoleID, user.FirstName, user.LastName, user.Phone, sq.Expr("NOW()"), sq.Expr("NOW()")).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		r.logger.Error(ctx, "CreateUser - error ToSql: %v", err)
		return nil, errorsper.ErrInternal(err, "UserRepository.CreateUser ToSql")
	}

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&user.ID)
	if err != nil {
		r.logger.Error(ctx, "CreateUser - exec error: %v", err)
		return nil, errorsper.ErrInternal(err, "UserRepository.CreateUser Exec")
	}

	r.logger.Info(ctx, "CreateUser - éxito id=%d", user.ID)
	return user, nil
}

// Actualizar usuario
func (r *userSQL) Update(ctx context.Context, id int, user *model.User) (*model.User, *errorsper.AppError) {
	ctx = logger.EnsureContext(ctx)
	r.logger.Info(ctx, "Update - id=%d username=%s", id, user.Username)

	query, args, err := r.sb.Update("users").
		SetMap(map[string]interface{}{
			"username":   user.Username,
			"email":      user.Email,
			"role_id":    user.RoleID,
			"password":   user.Password,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"phone":      user.Phone,
			"updated_at": sq.Expr("NOW()"),
		}).Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		r.logger.Error(ctx, "Update - error ToSql: %v", err)
		return nil, errorsper.ErrInternal(err, "UserRepository.Update ToSql")
	}

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		r.logger.Error(ctx, "Update - exec error: %v", err)
		return nil, errorsper.ErrInternal(err, "UserRepository.Update Exec")
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		r.logger.Warn(ctx, "Update - usuario no encontrado id=%d", id)
		return nil, errorsper.ErrNotFound("Usuario no encontrado", "UserRepository.Update")
	}

	user.ID = id
	r.logger.Info(ctx, "Update - éxito id=%d", id)
	return user, nil
}

// Eliminar usuario
func (r *userSQL) Delete(ctx context.Context, id int) *errorsper.AppError {
	ctx = logger.EnsureContext(ctx)
	r.logger.Info(ctx, "Delete - id=%d", id)

	query, args, err := r.sb.Delete("users").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		r.logger.Error(ctx, "Delete - error ToSql: %v", err)
		return errorsper.ErrInternal(err, "UserRepository.Delete ToSql")
	}

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		r.logger.Error(ctx, "Delete - exec error: %v", err)
		return errorsper.ErrInternal(err, "UserRepository.Delete Exec")
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		r.logger.Warn(ctx, "Delete - usuario no encontrado id=%d", id)
		return errorsper.ErrNotFound("Usuario no encontrado", "UserRepository.Delete")
	}

	r.logger.Info(ctx, "Delete - éxito id=%d", id)
	return nil
}

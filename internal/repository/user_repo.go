package repository

import (
	"database/sql"

	"classplanner-2.0/internal/model"
)

type UserRepository interface {
	GetAll() ([]*model.User, error)
	GetByID(id int) (*model.User, error)
	GetByEmailOrUser(user string) (*model.User, error)
	SearchByUserOrEmail(user string) ([]*model.User, error)
	Exists(id int) (bool, error)
	CreateUser(user *model.User) (*model.User, error)
	Update(id int, user *model.User) (*model.User, error)
	Delete(id int) error
}

type userSQL struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userSQL{db: db}
}

// Obtener todos los usuarios
func (r *userSQL) GetAll() ([]*model.User, error) {
	q := `SELECT id, username, role_id, first_name, last_name, email, created_at, updated_at 
	      FROM users`
	rows, err := r.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		u := &model.User{}
		if err := rows.Scan(&u.ID, &u.Username, &u.RoleID, &u.FirstName, &u.LastName, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// Obtener usuario por ID
func (r *userSQL) GetByID(id int) (*model.User, error) {
	q := `SELECT id, username, role_id, first_name, last_name, email, password, created_at, updated_at
	      FROM users WHERE id = $1`
	row := r.db.QueryRow(q, id)

	u := &model.User{}
	if err := row.Scan(&u.ID, &u.Username, &u.RoleID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

// Obtener usuario por username o email
func (r *userSQL) GetByEmailOrUser(user string) (*model.User, error) {
	q := `SELECT id, username, role_id, first_name, last_name, email, password, created_at, updated_at
	      FROM users WHERE username = $1 OR email = $1`
	row := r.db.QueryRow(q, user)

	u := &model.User{}
	if err := row.Scan(&u.ID, &u.Username, &u.RoleID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

// Buscar usuarios por nombre o email (LIKE)
func (r *userSQL) SearchByUserOrEmail(user string) ([]*model.User, error) {
	q := `SELECT id, username, role_id, first_name, last_name, email, created_at, updated_at 
	      FROM users WHERE username ILIKE $1 OR email ILIKE $1`
	rows, err := r.db.Query(q, "%"+user+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		u := &model.User{}
		if err := rows.Scan(&u.ID, &u.Username, &u.RoleID, &u.FirstName, &u.LastName, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// Verifica si existe un usuario por ID
func (r *userSQL) Exists(id int) (bool, error) {
	q := `SELECT 1 FROM users WHERE id = $1`
	row := r.db.QueryRow(q, id)
	var exists int
	err := row.Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// Crear usuario
func (r *userSQL) CreateUser(user *model.User) (*model.User, error) {
	q := `INSERT INTO users (username, email, password, role_id, first_name, last_name, created_at, updated_at)
	      VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	      RETURNING id`
	err := r.db.QueryRow(q, user.Username, user.Email, user.Password, user.RoleID, user.FirstName, user.LastName).Scan(&user.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Actualizar usuario
func (r *userSQL) Update(id int, user *model.User) (*model.User, error) {
	q := `UPDATE users 
	      SET username = $1, email = $2, role_id = $3, password = $4, first_name = $5, last_name = $6, updated_at = NOW() 
	      WHERE id = $7`
	_, err := r.db.Exec(q, user.Username, user.Email, user.RoleID, user.Password, user.FirstName, user.LastName, id)
	if err != nil {
		return nil, err
	}
	user.ID = id
	return user, nil
}

// Eliminar usuario
func (r *userSQL) Delete(id int) error {
	q := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(q, id)
	return err
}

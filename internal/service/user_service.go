package service

import (
	"classplanner/internal/model"
	"classplanner/internal/repository"
	"classplanner/internal/security"
	"context"
	"errors"
	"time"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Registro de usuario
func (s *UserService) Register(ctx context.Context, u *model.User) (*model.User, error) {
	// Verificar si el username ya existe
	existing, err := s.repo.GetByEmailOrUser(ctx, u.Username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("username ya existe")
	}

	// Verificar si el email ya está registrado
	existing, err = s.repo.GetByEmailOrUser(ctx, u.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email ya registrado")
	}

	// Hashear la contraseña
	hashed, err := security.HashPassword(u.Password)
	if err != nil {
		return nil, err
	}

	u.Password = hashed
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	return s.repo.CreateUser(ctx, u)
}

// Login de usuario
func (s *UserService) Login(ctx context.Context, userOrEmail, password string) (*model.User, error) {
	u, err := s.repo.GetByEmailOrUser(ctx, userOrEmail)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("usuario no encontrado")
	}

	if !security.CheckPasswordHash(password, u.Password) {
		return nil, errors.New("contraseña incorrecta")
	}

	return u, nil
}

// Obtener todos los usuarios
func (s *UserService) GetAll() ([]*model.User, error) {
	return s.repo.GetAll()
}

// Obtener usuario por ID
func (s *UserService) GetByID(ctx context.Context, id int) (*model.User, error) {
	return s.repo.GetByID(ctx, id)
}

// Actualizar usuario
func (s *UserService) Update(ctx context.Context, id int, u *model.User) (*model.User, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("usuario no encontrado")
	}

	if u.Password != "" {
		hashed, err := security.HashPassword(u.Password)
		if err != nil {
			return nil, err
		}
		u.Password = hashed
	}

	u.UpdatedAt = time.Now()
	return s.repo.Update(ctx, id, u)
}

// Logout de usuario (stateless)
func (s *UserService) Logout(_ context.Context, _ *model.User) error {
	// En un backend stateless (JWT o tokens), el logout se maneja del lado del cliente.
	return nil
}

// Eliminar usuario
func (s *UserService) Delete(ctx context.Context, id int) error {
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("usuario no encontrado")
	}
	return s.repo.Delete(ctx, id)
}

// Buscar usuarios por nombre o email
func (s *UserService) SearchByUserOrEmail(ctx context.Context, query string) ([]*model.User, error) {
	return s.repo.SearchByUserOrEmail(ctx, query)
}

// Verifica si un usuario existe por ID
func (s *UserService) Exists(ctx context.Context, id int) (bool, error) {
	return s.repo.Exists(ctx, id)
}

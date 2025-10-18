package service

import (
	"classplanner/internal/model"
	"classplanner/internal/repository"
	"classplanner/internal/security"
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
func (s *UserService) Register(u *model.User) (*model.User, error) {

	existing, _ := s.repo.GetByEmailOrUser(u.Username)
	if existing != nil {
		return nil, errors.New("username ya existe")
	}
	existing, _ = s.repo.GetByEmailOrUser(u.Email)
	if existing != nil {
		return nil, errors.New("email ya registrado")
	}

	hashed, err := security.HashPassword(u.Password)
	if err != nil {
		return nil, err
	}
	u.Password = hashed
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	return s.repo.CreateUser(u)
}

// Login de usuario
func (s *UserService) Login(userOrEmail, password string) (*model.User, error) {
	u, err := s.repo.GetByEmailOrUser(userOrEmail)
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
func (s *UserService) GetByID(id int) (*model.User, error) {
	return s.repo.GetByID(id)
}

// Actualizar usuario
func (s *UserService) Update(id int, u *model.User) (*model.User, error) {
	existing, err := s.repo.GetByID(id)
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
	return s.repo.Update(id, u)
}

// Logout de usuario
func (s *UserService) Logout(u *model.User) error {
	// En un backend stateless (sin sesiones), el logout se maneja en el cliente eliminando el token
	// Aquí solo devolvemos nil como confirmación.
	return nil
}

// Eliminar usuario
func (s *UserService) Delete(id int) error {
	exists, err := s.repo.Exists(id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("usuario no encontrado")
	}
	return s.repo.Delete(id)
}

// Buscar usuarios por nombre o email
func (s *UserService) SearchByUserOrEmail(query string) ([]*model.User, error) {
	return s.repo.SearchByUserOrEmail(query)
}

// Verifica si un usuario existe por ID
func (s *UserService) Exists(id int) (bool, error) {
	return s.repo.Exists(id)
}

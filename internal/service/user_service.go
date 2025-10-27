package service

import (
	"classplanner/internal/infrastructure/logger"
	"classplanner/internal/model"
	"classplanner/internal/repository"
	"classplanner/internal/security"
	"context"
	"errors"
	"time"
)

type UserService struct {
	repo   repository.UserRepository
	logger *logger.Logger
}

func NewUserService(repo repository.UserRepository, l *logger.Logger) *UserService {
	return &UserService{repo: repo, logger: l}
}

// Registro de usuario
func (s *UserService) Register(ctx context.Context, u *model.User) (*model.User, error) {
	ctx = logger.EnsureContext(ctx)
	s.logger.Info(ctx, "Intentando registrar usuario: %s", u.Username)

	existing, err := s.repo.GetByEmailOrUser(ctx, u.Username)
	if err != nil {
		s.logger.Error(ctx, "Error al verificar username existente: %v", err)
		return nil, err
	}
	if existing != nil {
		s.logger.Warn(ctx, "Intento de registro con username ya existente: %s", u.Username)
		return nil, errors.New("username ya existe")
	}

	existing, err = s.repo.GetByEmailOrUser(ctx, u.Email)
	if err != nil {
		s.logger.Error(ctx, "Error al verificar email existente: %v", err)
		return nil, err
	}
	if existing != nil {
		s.logger.Warn(ctx, "Intento de registro con email ya registrado: %s", u.Email)
		return nil, errors.New("email ya registrado")
	}

	hashed, err := security.HashPassword(u.Password)
	if err != nil {
		s.logger.Error(ctx, "Error al hashear contraseña: %v", err)
		return nil, err
	}

	u.Password = hashed
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	created, err := s.repo.CreateUser(ctx, u)
	if err != nil {
		s.logger.Error(ctx, "Error al crear usuario: %v", err)
		return nil, err
	}

	s.logger.Info(ctx, "Usuario registrado correctamente: %s", created.Username)
	return created, nil
}

// Login de usuario
func (s *UserService) Login(ctx context.Context, userOrEmail, password string) (*model.User, error) {
	ctx = logger.EnsureContext(ctx)
	s.logger.Info(ctx, "Intentando login para: %s", userOrEmail)

	u, err := s.repo.GetByEmailOrUser(ctx, userOrEmail)
	if err != nil {
		s.logger.Error(ctx, "Error al obtener usuario en login: %v", err)
		return nil, err
	}
	if u == nil {
		s.logger.Warn(ctx, "Usuario no encontrado en login: %s", userOrEmail)
		return nil, errors.New("usuario no encontrado")
	}

	if !security.CheckPasswordHash(password, u.Password) {
		s.logger.Warn(ctx, "Contraseña incorrecta para usuario: %s", userOrEmail)
		return nil, errors.New("contraseña incorrecta")
	}

	s.logger.Info(ctx, "Login exitoso para usuario: %s", u.Username)
	return u, nil
}

// Obtener todos los usuarios
func (s *UserService) GetAll(ctx context.Context) ([]*model.User, error) {
	ctx = logger.EnsureContext(ctx)
	users, err := s.repo.GetAll()
	if err != nil {
		s.logger.Error(ctx, "Error al obtener todos los usuarios: %v", err)
		return nil, err
	}
	s.logger.Info(ctx, "Usuarios obtenidos correctamente (%d encontrados)", len(users))
	return users, nil
}

// Obtener usuario por ID
func (s *UserService) GetByID(ctx context.Context, id int) (*model.User, error) {
	ctx = logger.EnsureContext(ctx)
	s.logger.Debug(ctx, "Buscando usuario por ID: %d", id)

	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "Error al obtener usuario por ID %d: %v", id, err)
		return nil, err
	}
	if u == nil {
		s.logger.Warn(ctx, "Usuario no encontrado con ID: %d", id)
		return nil, errors.New("usuario no encontrado")
	}
	return u, nil
}

// Actualizar usuario
func (s *UserService) Update(ctx context.Context, id int, u *model.User) (*model.User, error) {
	ctx = logger.EnsureContext(ctx)
	s.logger.Info(ctx, "Intentando actualizar usuario ID: %d", id)

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "Error al obtener usuario antes de actualizar: %v", err)
		return nil, err
	}
	if existing == nil {
		s.logger.Warn(ctx, "Intento de actualizar usuario inexistente ID: %d", id)
		return nil, errors.New("usuario no encontrado")
	}

	if u.Password != "" {
		hashed, err := security.HashPassword(u.Password)
		if err != nil {
			s.logger.Error(ctx, "Error al hashear nueva contraseña para usuario ID %d: %v", id, err)
			return nil, err
		}
		u.Password = hashed
	}

	u.UpdatedAt = time.Now()
	updated, err := s.repo.Update(ctx, id, u)
	if err != nil {
		s.logger.Error(ctx, "Error al actualizar usuario ID %d: %v", id, err)
		return nil, err
	}

	s.logger.Info(ctx, "Usuario actualizado correctamente ID: %d", id)
	return updated, nil
}

// Eliminar usuario
func (s *UserService) Delete(ctx context.Context, id int) error {
	ctx = logger.EnsureContext(ctx)
	s.logger.Info(ctx, "Intentando eliminar usuario ID: %d", id)

	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "Error al verificar existencia del usuario ID %d: %v", id, err)
		return err
	}
	if !exists {
		s.logger.Warn(ctx, "Intento de eliminar usuario inexistente ID: %d", id)
		return errors.New("usuario no encontrado")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error(ctx, "Error al eliminar usuario ID %d: %v", id, err)
		return err
	}

	s.logger.Info(ctx, "Usuario eliminado correctamente ID: %d", id)
	return nil
}

// Buscar usuarios por nombre o email
func (s *UserService) SearchByUserOrEmail(ctx context.Context, query string) ([]*model.User, error) {
	ctx = logger.EnsureContext(ctx)
	s.logger.Debug(ctx, "Buscando usuarios con query: %s", query)

	users, err := s.repo.SearchByUserOrEmail(ctx, query)
	if err != nil {
		s.logger.Error(ctx, "Error en búsqueda de usuarios: %v", err)
		return nil, err
	}

	s.logger.Info(ctx, "Se encontraron %d usuarios para query: %s", len(users), query)
	return users, nil
}

// Verifica si un usuario existe por ID
func (s *UserService) Exists(ctx context.Context, id int) (bool, error) {
	ctx = logger.EnsureContext(ctx)
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "Error al verificar existencia de usuario ID %d: %v", id, err)
		return false, err
	}
	s.logger.Debug(ctx, "Existencia usuario ID %d: %t", id, exists)
	return exists, nil
}

// Logout (stateless)
func (s *UserService) Logout(ctx context.Context, _ *model.User) error {
	ctx = logger.EnsureContext(ctx)
	s.logger.Info(ctx, "Logout de usuario (stateless)")
	return nil
}

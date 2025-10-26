package service

import (
	"classplanner/internal/model"
	"classplanner/internal/repository"
	"context"
	"errors"
	"fmt"
)

type ClassesService struct {
	repo repository.ClassesRepository
}

func NewClassRepository(repo repository.ClassesRepository) *ClassesService {
	return &ClassesService{repo: repo}
}

// Crear clase
func (s *ClassesService) CreateClasses(ctx context.Context, class *model.Classes) (*model.Classes, error) {
	exist, err := s.repo.Exists(ctx, int(class.ID))
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, errors.New("la clase ya existe o los datos ingresados no son correctos")
	}

	created, err := s.repo.CreateClass(ctx, class)
	if err != nil {
		return nil, err
	}
	return created, nil
}

// Obtener clases por profesor
func (s *ClassesService) GetClassByTeacher(ctx context.Context, teacherID int) ([]*model.Classes, error) {
	classes, err := s.repo.ClassesByTeacher(ctx, teacherID)
	if err != nil {
		return nil, err
	}
	if len(classes) == 0 {
		return nil, fmt.Errorf("no se encontraron clases para el profesor con ID %d", teacherID)
	}

	return classes, nil
}

// Obtener clases por alumno
func (s *ClassesService) GetClassByStudent(ctx context.Context, studentID int) ([]*model.Classes, error) {
	classes, err := s.repo.ClassesByStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}
	if len(classes) == 0 {
		return nil, fmt.Errorf("no se encontraron clases para el alumno con ID %d", studentID)
	}
	return classes, nil
}

// Unirse a clase por token
func (s *ClassesService) JoinClass(ctx context.Context, userID int, token string) error {
	if token == "" {
		return errors.New("el token no puede estar vacío")
	}

	err := s.repo.JoinClass(ctx, userID, token)
	if err != nil {
		return fmt.Errorf("error al unirse a la clase: %w", err)
	}

	return nil
}

// Obtener usuarios de una clase (profesor + alumnos)
func (s *ClassesService) GetUsersFromClass(ctx context.Context, classID int) ([]*model.UserClass, error) {
	users, err := s.repo.UsersFromClass(ctx, classID)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("no se encontraron usuarios para la clase con ID %d", classID)
	}

	return users, nil
}

package service

import (
	"classplanner/internal/model"
	"classplanner/internal/repository"
	"context"
	"errors"
	"fmt"
	"time"
)

type SubmissionService struct {
	repo repository.SubmissionRepository
}

func NewSubmissionService(repo repository.SubmissionRepository) *SubmissionService {
	return &SubmissionService{repo: repo}
}

// Crear una nueva entrega
func (s *SubmissionService) CreateSubmission(ctx context.Context, sub *model.Submission) error {
	if sub == nil {
		return errors.New("submission no puede ser nil")
	}
	if sub.ID_user <= 0 {
		return errors.New("usuario inválido")
	}
	if sub.ID_task <= 0 {
		return errors.New("tarea inválida")
	}
	if sub.File == "" && sub.Comment == "" {
		return errors.New("el submission debe tener un archivo o un comentario")
	}
	sub.Date = time.Now().Format(time.RFC3339)
	return s.repo.Create(ctx, sub)
}

// Obtener entrega por usuario y tarea
func (s *SubmissionService) GetByUserAndTask(ctx context.Context, userID, taskID int64) (*model.Submission, error) {
	if userID <= 0 || taskID <= 0 {
		return nil, errors.New("usuario o tarea inválidos")
	}
	sub, err := s.repo.GetByUserAndTask(ctx, userID, taskID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo submission: %w", err)
	}
	if sub == nil {
		return nil, fmt.Errorf("no se encontró submission para el usuario %d en la tarea %d", userID, taskID)
	}
	return sub, nil
}

// Obtener todas las entregas de una tarea (paginadas)
func (s *SubmissionService) GetByTask(ctx context.Context, taskID int64, limit, offset uint64) ([]*model.Submission, error) {
	if taskID <= 0 {
		return nil, errors.New("tarea inválida")
	}
	if limit == 0 {
		limit = 50
	}
	return s.repo.GetByTask(ctx, taskID, limit, offset)
}

// Actualizar entrega existente
func (s *SubmissionService) UpdateSubmission(ctx context.Context, sub *model.Submission) error {
	if sub == nil {
		return errors.New("submission no puede ser nil")
	}
	if sub.ID <= 0 {
		return errors.New("submission inválido")
	}
	if sub.File == "" && sub.Comment == "" && sub.Calification == "" && sub.Feedback == "" {
		return errors.New("no hay campos para actualizar")
	}
	return s.repo.Update(ctx, sub)
}

// Obtener entregas de una tarea dentro de un rango de fechas
func (s *SubmissionService) GetByTaskAndDate(ctx context.Context, taskID int64, from, to time.Time) ([]*model.Submission, error) {
	if taskID <= 0 {
		return nil, errors.New("tarea inválida")
	}
	if from.IsZero() || to.IsZero() {
		return nil, errors.New("fechas inválidas")
	}
	if from.After(to) {
		return nil, errors.New("la fecha 'from' no puede ser posterior a 'to'")
	}
	return s.repo.GetByTaskAndDate(ctx, taskID, from, to)
}

// Obtener últimas N entregas de un usuario
func (s *SubmissionService) GetLastSubmissionsByUser(ctx context.Context, userID int64, limit int) ([]*model.Submission, error) {
	if userID <= 0 {
		return nil, errors.New("usuario inválido")
	}
	if limit <= 0 {
		limit = 10
	}
	return s.repo.GetLastSubmissionsByUser(ctx, userID, limit)
}

// Contar entregas de una tarea
func (s *SubmissionService) CountByTask(ctx context.Context, taskID int64) (int64, error) {
	if taskID <= 0 {
		return 0, errors.New("tarea inválida")
	}
	count, err := s.repo.CountByTask(ctx, taskID)
	if err != nil {
		return 0, fmt.Errorf("error contando submissions: %w", err)
	}
	return count, nil
}

package service

import (
	"classplanner/internal/model"
	"classplanner/internal/repository"
	"context"
	"fmt"
)

type CommentService struct {
	repo repository.CommentRepository
}

func NewCommentService(repo repository.CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

// Crear un comentario
func (s *CommentService) CreateComment(ctx context.Context, comment *model.Comment, userID int64) (*model.Comment, error) {
	if comment == nil {
		return nil, fmt.Errorf("comentario no puede ser nil")
	}
	if err := s.repo.Create(ctx, comment, userID); err != nil {
		return nil, fmt.Errorf("error creando comentario: %w", err)
	}
	return comment, nil
}

// Obtener comentarios por ID de tarea
func (s *CommentService) GetCommentsByTask(ctx context.Context, taskID int64, limit, offset int) ([]*model.Comment, error) {
	comments, err := s.repo.GetByTaskID(ctx, taskID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo comentarios por tarea: %w", err)
	}
	return comments, nil
}

// Obtener comentarios por ID de usuario
func (s *CommentService) GetCommentsByUser(ctx context.Context, userID int64, limit, offset int) ([]*model.Comment, error) {
	comments, err := s.repo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo comentarios por usuario: %w", err)
	}
	return comments, nil
}

// Eliminar comentario (si es admin o dueño)
func (s *CommentService) DeleteComment(ctx context.Context, commentID, userID int64, isAdmin bool) error {
	exists, err := s.repo.Exists(ctx, commentID)
	if err != nil {
		return fmt.Errorf("error verificando existencia del comentario: %w", err)
	}
	if !exists {
		return fmt.Errorf("comentario no encontrado")
	}

	if err := s.repo.Delete(ctx, commentID, userID, isAdmin); err != nil {
		return fmt.Errorf("error eliminando comentario: %w", err)
	}
	return nil
}

// Actualizar comentario (si es admin o dueño)
func (s *CommentService) UpdateComment(ctx context.Context, commentID, userID int64, newText string) error {
	exists, err := s.repo.Exists(ctx, commentID)
	if err != nil {
		return fmt.Errorf("error verificando existencia del comentario: %w", err)
	}
	if !exists {
		return fmt.Errorf("comentario no encontrado")
	}

	if newText == "" {
		return fmt.Errorf("el texto del comentario no puede estar vacío")
	}

	if err := s.repo.Update(ctx, commentID, userID, newText); err != nil {
		return fmt.Errorf("error actualizando comentario: %w", err)
	}
	return nil
}

// Contar comentarios de una tarea
func (s *CommentService) CountCommentsByTask(ctx context.Context, taskID int64) (int64, error) {
	count, err := s.repo.CountByTask(ctx, taskID)
	if err != nil {
		return 0, fmt.Errorf("error contando comentarios de la tarea: %w", err)
	}
	return count, nil
}

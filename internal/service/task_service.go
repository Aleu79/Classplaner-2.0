package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"classplanner/internal/model"
	"classplanner/internal/repository"
)

type TasksService struct {
	taskRepo       repository.TasksRepository
	submissionRepo *SubmissionService // para chequear entregas
}

func NewTasksService(taskRepo repository.TasksRepository, submissionRepo *SubmissionService) *TasksService {
	return &TasksService{
		taskRepo:       taskRepo,
		submissionRepo: submissionRepo,
	}
}

// Valida y crea una nueva tarea
func (s *TasksService) CreateTask(task *model.Tasks) error {
	if task == nil {
		return fmt.Errorf("task no puede ser nil")
	}
	if task.Clase <= 0 {
		return fmt.Errorf("id de clase inválido")
	}
	if task.Titulo == "" {
		return fmt.Errorf("el título no puede estar vacío")
	}
	if task.Limite != "" {
		if _, err := time.Parse("2006-01-02", task.Limite); err != nil {
			return fmt.Errorf("formato de fecha límite inválido: %w", err)
		}
	}

	if err := s.taskRepo.Create(task); err != nil {
		log.Printf("Error creando tarea: %v", err)
		return err
	}

	return nil
}

// Obtiene las tareas de un usuario y asigna el estado según fecha y entrega
func (s *TasksService) GetTasksByUser(ctx context.Context, userID int, userType string) ([]*model.Tasks, error) {
	tasks, err := s.taskRepo.GetByUser(userID, userType)
	if err != nil {
		log.Printf("Error obteniendo tareas por usuario: %v", err)
		return nil, err
	}

	now := time.Now()
	for _, t := range tasks {
		t.Estado = model.StatePending

		// Chequea si la tarea fue entregada
		if s.submissionRepo != nil {
			sub, _ := s.submissionRepo.GetByUserAndTask(ctx, int64(userID), int64(t.ID))
			if sub != nil {
				t.Estado = model.StateSubmitted
				continue
			}
		}

		// Si no fue entregada, chequea si está atrasada
		if t.Limite != "" {
			if deadline, err := time.Parse("2006-01-02", t.Limite); err == nil && deadline.Before(now) {
				t.Estado = model.StateLate
			}
		}
	}

	return tasks, nil
}

// Obtiene las tareas de una clase con estado y filtros opcionales
func (s *TasksService) GetTasksByClass(ctx context.Context, classID, limit, offset int, from, to string, userID int) ([]*model.Tasks, error) {
	tasks, err := s.taskRepo.GetByClass(classID, limit, offset)
	if err != nil {
		log.Printf("Error obteniendo tareas por clase: %v", err)
		return nil, err
	}

	now := time.Now()
	filtered := make([]*model.Tasks, 0, len(tasks))
	var fromDate, toDate time.Time

	if from != "" {
		if fromDate, err = time.Parse("2006-01-02", from); err != nil {
			return nil, fmt.Errorf("formato de fecha 'from' inválido: %w", err)
		}
	}
	if to != "" {
		if toDate, err = time.Parse("2006-01-02", to); err != nil {
			return nil, fmt.Errorf("formato de fecha 'to' inválido: %w", err)
		}
	}

	for _, t := range tasks {
		t.Estado = model.StatePending

		// Chequea si la tarea fue entregada por este usuario
		if s.submissionRepo != nil && userID > 0 {
			sub, _ := s.submissionRepo.GetByUserAndTask(ctx, int64(userID), int64(t.ID))
			if sub != nil {
				t.Estado = model.StateSubmitted
				filtered = append(filtered, t)
				continue
			}
		}

		// Si no fue entregada, chequea si está atrasada
		if t.Limite != "" {
			if deadline, err := time.Parse("2006-01-02", t.Limite); err == nil && deadline.Before(now) {
				t.Estado = model.StateLate
			}
		}

		// Filtro por rango de fechas
		taskDate, err := time.Parse("2006-01-02", t.Limite)
		if err == nil {
			if (from == "" || !taskDate.Before(fromDate)) && (to == "" || !taskDate.After(toDate)) {
				filtered = append(filtered, t)
			}
		} else {
			filtered = append(filtered, t)
		}
	}

	return filtered, nil
}

// Obtiene una tarea específica con su estado
func (s *TasksService) GetTaskByID(ctx context.Context, taskID, userID int) (*model.Tasks, error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		log.Printf("Error obteniendo tarea por ID: %v", err)
		return nil, err
	}
	if task == nil {
		return nil, fmt.Errorf("tarea no encontrada")
	}

	task.Estado = model.StatePending

	if s.submissionRepo != nil && userID > 0 {
		sub, _ := s.submissionRepo.GetByUserAndTask(ctx, int64(userID), int64(task.ID))
		if sub != nil {
			task.Estado = model.StateSubmitted
			return task, nil
		}
	}

	if task.Limite != "" {
		if deadline, err := time.Parse("2006-01-02", task.Limite); err == nil && deadline.Before(time.Now()) {
			task.Estado = model.StateLate
		}
	}

	return task, nil
}

// Actualiza una tarea existente
func (s *TasksService) UpdateTask(task *model.Tasks) error {
	if task == nil || task.ID <= 0 {
		return fmt.Errorf("task inválido")
	}
	if task.Titulo == "" {
		return fmt.Errorf("el título no puede estar vacío")
	}
	if task.Limite != "" {
		if _, err := time.Parse("2006-01-02", task.Limite); err != nil {
			return fmt.Errorf("formato de fecha límite inválido: %w", err)
		}
	}

	if err := s.taskRepo.Update(task); err != nil {
		log.Printf("Error actualizando tarea: %v", err)
		return err
	}

	return nil
}

// Elimina una tarea por ID
func (s *TasksService) DeleteTask(taskID int) error {
	if taskID <= 0 {
		return fmt.Errorf("taskID inválido")
	}

	if err := s.taskRepo.Delete(taskID); err != nil {
		log.Printf("Error eliminando tarea: %v", err)
		return err
	}

	return nil
}

// Genera un reporte rápido de tareas por estado
func (s *TasksService) ReportTasksByState(tasks []*model.Tasks) map[model.TaskState]int {
	report := map[model.TaskState]int{
		model.StatePending:   0,
		model.StateSubmitted: 0,
		model.StateLate:      0,
	}

	for _, t := range tasks {
		report[t.Estado]++
	}

	return report
}

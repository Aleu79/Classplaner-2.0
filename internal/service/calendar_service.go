package service

import (
	"classplanner/internal/model"
	"classplanner/internal/repository"
	"context"
	"fmt"
)

type CalendarService struct {
	repo repository.CalendarRepositoy
}

func NewCalendarService(repo repository.CalendarRepositoy) *CalendarService {
	return &CalendarService{repo: repo}
}

// GetCalendarByUser recupera todas las tareas de un usuario (con paginación opcional)
func (s *CalendarService) GetCalendarByUser(ctx context.Context, userID int, userType string, limit, offset uint64) ([]*model.Calendar, error) {
	calendars, err := s.repo.CalendarByUser(ctx, userID, userType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("CalendarService: GetCalendarByUser failed: %w", err)
	}
	return calendars, nil
}

// GetCalendarByUserAndToken recupera tareas para un usuario filtradas por token de clase
func (s *CalendarService) GetCalendarByUserAndToken(ctx context.Context, userID int, userType, classToken string, limit, offset uint64) ([]*model.Calendar, error) {
	calendars, err := s.repo.CalendarByUserAndToken(ctx, userID, userType, classToken, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("CalendarService: GetCalendarByUserAndToken failed: %w", err)
	}
	return calendars, nil
}

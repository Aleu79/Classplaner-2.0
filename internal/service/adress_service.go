package service

import (
	"classplanner/internal/model"
	"classplanner/internal/repository"
	"context"
	"errors"
	"time"
)

type AddressService struct {
	repo repository.AddressRepository
}

func NewAddressService(repo repository.AddressRepository) *AddressService {
	return &AddressService{repo: repo}
}

// Obtener todas las direcciones de un usuario
func (s *AddressService) GetByUserID(ctx context.Context, userID int) ([]*model.Address, error) {
	addresses, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

// Crear nueva dirección
func (s *AddressService) CreateAddress(ctx context.Context, address *model.Address) (*model.Address, error) {
	address.CreatedAt = time.Now()
	address.UpdatedAt = time.Now()
	return s.repo.CreateAddress(ctx, address)
}

// Actualizar dirección existente
func (s *AddressService) UpdateAddress(ctx context.Context, id int, address *model.Address) (*model.Address, error) {

	existing, err := s.repo.GetByUserID(ctx, address.UserID)
	if err != nil {
		return nil, err
	}

	found := false
	for _, a := range existing {
		if a.ID == id {
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New("dirección no encontrada para este usuario")
	}

	address.UpdatedAt = time.Now()
	return s.repo.UpdateAddress(ctx, id, address)
}

// Eliminar dirección (lógicamente)
func (s *AddressService) DeleteAddress(ctx context.Context, id int) error {
	return s.repo.DeleteAddress(ctx, id)
}

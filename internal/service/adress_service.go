package service

import (
	"classplanner/internal/model"
	"classplanner/internal/repository"
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
func (s *AddressService) GetByUserID(userID int) ([]*model.Address, error) {
	addresses, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

// Crear nueva dirección
func (s *AddressService) CreateAddress(address *model.Address) (*model.Address, error) {
	address.CreatedAt = time.Now()
	address.UpdatedAt = time.Now()
	return s.repo.CreateAddress(address)
}

// Actualizar dirección existente
func (s *AddressService) UpdateAddress(id int, address *model.Address) (*model.Address, error) {

	existing, err := s.repo.GetByUserID(address.UserID)
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
	return s.repo.UpdateAddress(id, address)
}

// Eliminar dirección (lógicamente)
func (s *AddressService) DeleteAddress(id int) error {
	return s.repo.DeleteAddress(id)
}

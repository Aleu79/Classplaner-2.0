package address

import (
	"classplanner/internal/model"
	"classplanner/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AddressHandler struct {
	service *service.AddressService
}

func NewAddressHandler(s *service.AddressService) *AddressHandler {
	return &AddressHandler{service: s}
}

// Obtener direcciones de un usuario
func (h *AddressHandler) GetByUserID(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("user_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID de usuario inválido"})
	}

	addresses, err := h.service.GetByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(addresses)
}

// Crear dirección
func (h *AddressHandler) Create(c *fiber.Ctx) error {
	address := new(model.Address)
	if err := c.BodyParser(address); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	newAddress, err := h.service.CreateAddress(address)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(newAddress)
}

// Actualizar dirección
func (h *AddressHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	address := new(model.Address)
	if err := c.BodyParser(address); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	updated, err := h.service.UpdateAddress(id, address)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(updated)
}

// Eliminar dirección
func (h *AddressHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	if err := h.service.DeleteAddress(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "dirección eliminada"})
}

package users

import (
	"classplanner/internal/model"
	"classplanner/internal/service"
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// Registro
func (h *UserHandler) Register(c *fiber.Ctx) error {
	ctx := context.Background()

	u := new(model.User)
	if err := c.BodyParser(u); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := h.service.Register(ctx, u)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// Login
func (h *UserHandler) Login(c *fiber.Ctx) error {
	ctx := context.Background()

	req := struct {
		UserOrEmail string `json:"user_or_email"`
		Password    string `json:"password"`
	}{}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := h.service.Login(ctx, req.UserOrEmail, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}

// Logout
func (h *UserHandler) Logout(c *fiber.Ctx) error {
	ctx := context.Background()

	err := h.service.Logout(ctx, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "logout exitoso",
	})
}

// Actualizar usuario
func (h *UserHandler) Update(c *fiber.Ctx) error {
	ctx := context.Background()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	u := new(model.User)
	if err := c.BodyParser(u); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	updatedUser, err := h.service.Update(ctx, id, u)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(updatedUser)
}

// Verificar si un usuario existe
func (h *UserHandler) Exists(c *fiber.Ctx) error {
	ctx := context.Background()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	exists, err := h.service.Exists(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"exists": exists})
}

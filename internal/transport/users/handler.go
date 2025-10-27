package users

import (
	"classplanner/internal/infrastructure/logger"
	"classplanner/internal/model"
	"classplanner/internal/service"
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service *service.UserService
	logger  *logger.Logger
}

// Ahora el constructor recibe logger también
func NewUserHandler(s *service.UserService, l *logger.Logger) *UserHandler {
	return &UserHandler{
		service: s,
		logger:  l,
	}
}

// Helper para crear contexto con requestID y userID
func (h *UserHandler) GetCtx(c *fiber.Ctx) context.Context {
	ctx := context.Background()
	if reqID := c.Locals("requestid"); reqID != nil {
		if rid, ok := reqID.(string); ok {
			ctx = logger.WithReqID(ctx, rid)
		}
	}
	if uid := c.Locals("userID"); uid != nil {
		if id, ok := uid.(int); ok {
			ctx = logger.WithUserID(ctx, id)
		}
	}
	return ctx
}

// Registro
func (h *UserHandler) Register(c *fiber.Ctx) error {
	ctx := h.GetCtx(c)

	u := new(model.User)
	if err := c.BodyParser(u); err != nil {
		h.logger.Warn(ctx, "BodyParser error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := h.service.Register(ctx, u)
	if err != nil {
		h.logger.Error(ctx, "Error registrando usuario: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	h.logger.Info(ctx, "Usuario registrado ID=%d", user.ID)
	return c.Status(fiber.StatusCreated).JSON(user)
}

// Login
func (h *UserHandler) Login(c *fiber.Ctx) error {
	ctx := h.GetCtx(c)

	req := struct {
		UserOrEmail string `json:"user_or_email"`
		Password    string `json:"password"`
	}{}

	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn(ctx, "BodyParser error login: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := h.service.Login(ctx, req.UserOrEmail, req.Password)
	if err != nil {
		h.logger.Warn(ctx, "Login fallido: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	h.logger.Info(ctx, "Login exitoso ID=%d", user.ID)
	return c.JSON(user)
}

// Logout
func (h *UserHandler) Logout(c *fiber.Ctx) error {
	ctx := h.GetCtx(c)

	err := h.service.Logout(ctx, nil)
	if err != nil {
		h.logger.Error(ctx, "Error en logout: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	h.logger.Info(ctx, "Logout exitoso")
	return c.JSON(fiber.Map{"message": "logout exitoso"})
}

// Actualizar usuario
func (h *UserHandler) Update(c *fiber.Ctx) error {
	ctx := h.GetCtx(c)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		h.logger.Warn(ctx, "ID inválido: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	u := new(model.User)
	if err := c.BodyParser(u); err != nil {
		h.logger.Warn(ctx, "BodyParser error update: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	updatedUser, err := h.service.Update(ctx, id, u)
	if err != nil {
		h.logger.Error(ctx, "Error actualizando usuario ID=%d: %v", id, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	h.logger.Info(ctx, "Usuario actualizado ID=%d", id)
	return c.JSON(updatedUser)
}

// Verificar si un usuario existe
func (h *UserHandler) Exists(c *fiber.Ctx) error {
	ctx := h.GetCtx(c)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		h.logger.Warn(ctx, "ID inválido exists: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	exists, err := h.service.Exists(ctx, id)
	if err != nil {
		h.logger.Error(ctx, "Error checking exists ID=%d: %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	h.logger.Info(ctx, "Exists ID=%d: %v", id, exists)
	return c.JSON(fiber.Map{"exists": exists})
}

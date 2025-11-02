package users

import (
	"classplanner/internal/infrastructure/logger"
	"classplanner/internal/model"
	"classplanner/internal/service"
	"classplanner/pkg/response"
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service *service.UserService
	logger  *logger.Logger
}

func NewUserHandler(s *service.UserService, l *logger.Logger) *UserHandler {
	return &UserHandler{
		service: s,
		logger:  l,
	}
}

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

// Register
func (h *UserHandler) Register(c *fiber.Ctx) error {
	ctx := h.GetCtx(c)

	u := new(model.User)
	if err := c.BodyParser(u); err != nil {
		h.logger.Warn(ctx, "Error parsing body: %v", err)
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	user, err := h.service.Register(ctx, u)
	if err != nil {
		h.logger.Error(ctx, "Error registering user: %v", err)
		return response.BadRequest(c, "No se pudo registrar el usuario", err.Error())
	}

	h.logger.Info(ctx, "User registered ID=%d", user.ID)
	return response.Created(c, "Usuario registrado correctamente", user)
}

// Login
func (h *UserHandler) Login(c *fiber.Ctx) error {
	ctx := h.GetCtx(c)

	req := struct {
		UserOrEmail string `json:"user_or_email"`
		Password    string `json:"password"`
	}{}

	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn(ctx, "Error parsing login body: %v", err)
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	user, err := h.service.Login(ctx, req.UserOrEmail, req.Password)
	if err != nil {
		h.logger.Warn(ctx, "Login failed: %v", err)
		return response.Unauthorized(c, "Usuario o contraseña incorrectos")
	}

	h.logger.Info(ctx, "Login successful ID=%d", user.ID)
	return response.Success(c, "Login exitoso", user)
}

// Logout
func (h *UserHandler) Logout(c *fiber.Ctx) error {
	ctx := h.GetCtx(c)

	if err := h.service.Logout(ctx, nil); err != nil {
		h.logger.Error(ctx, "Logout error: %v", err)
		return response.InternalError(c, "Error en logout", err.Error())
	}

	h.logger.Info(ctx, "Logout successful")
	return response.Success(c, "Logout exitoso", nil)
}

// Update
func (h *UserHandler) Update(c *fiber.Ctx) error {
	ctx := h.GetCtx(c)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		h.logger.Warn(ctx, "Invalid ID: %v", err)
		return response.BadRequest(c, "ID inválido", err.Error())
	}

	u := new(model.User)
	if err := c.BodyParser(u); err != nil {
		h.logger.Warn(ctx, "Error parsing body for update: %v", err)
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	updatedUser, err := h.service.Update(ctx, id, u)
	if err != nil {
		h.logger.Error(ctx, "Error updating user ID=%d: %v", id, err)
		return response.BadRequest(c, "No se pudo actualizar el usuario", err.Error())
	}

	h.logger.Info(ctx, "User updated ID=%d", id)
	return response.Success(c, "Usuario actualizado correctamente", updatedUser)
}

// Exists
func (h *UserHandler) Exists(c *fiber.Ctx) error {
	ctx := h.GetCtx(c)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		h.logger.Warn(ctx, "Invalid ID for exists: %v", err)
		return response.BadRequest(c, "ID inválido", err.Error())
	}

	exists, err := h.service.Exists(ctx, id)
	if err != nil {
		h.logger.Error(ctx, "Error checking exists ID=%d: %v", id, err)
		return response.InternalError(c, "Error verificando existencia de usuario", err.Error())
	}

	h.logger.Info(ctx, "Exists check ID=%d: %v", id, exists)
	return response.Success(c, "Verificación de existencia exitosa", fiber.Map{"exists": exists})
}

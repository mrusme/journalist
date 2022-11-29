package users

import (
	"context"
	// "github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	// "github.com/mrusme/journalist/ent/user"
	// "github.com/mrusme/journalist/ent"
)

type UserListResponse struct {
	Success bool             `json:"success"`
	Users   *[]UserShowModel `json:"users"`
	Message string           `json:"message"`
}

// List godoc
// @Summary      List users
// @Description  Get all users
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  UserListResponse
// @Failure      400  {object}  UserListResponse
// @Failure      404  {object}  UserListResponse
// @Failure      500  {object}  UserListResponse
// @Router       /users [get]
// @security     BasicAuth
func (h *handler) List(ctx *fiber.Ctx) error {
	var err error

	role := ctx.Locals("role").(string)

	if role != "admin" {
		h.logger.Debug(
			"User not allowed to list users",
			zap.Error(err),
		)
		return ctx.
			Status(fiber.StatusForbidden).
			JSON(UserListResponse{
				Success: false,
				Users:   nil,
				Message: "Only admins are allowed to list users",
			})
	}

	dbUsers, err := h.entClient.User.
		Query().
		All(context.Background())
	if err != nil {
		return ctx.
			Status(fiber.StatusInternalServerError).
			JSON(UserListResponse{
				Success: false,
				Users:   nil,
				Message: err.Error(),
			})
	}

	showUsers := make([]UserShowModel, len(dbUsers))

	for i, dbUser := range dbUsers {
		showUsers[i] = UserShowModel{
			ID:       dbUser.ID.String(),
			Username: dbUser.Username,
			Role:     dbUser.Role,
		}
	}

	return ctx.
		Status(fiber.StatusOK).
		JSON(UserListResponse{
			Success: true,
			Users:   &showUsers,
			Message: "",
		})
}

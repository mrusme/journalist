package users

import (
	"context"
	// "github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	// "github.com/mrusme/journalist/ent/user"
	// "github.com/mrusme/journalist/ent"
)

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
      JSON(&fiber.Map{
        "success": false,
        "users": nil,
        "message": "Only admins are allowed to list users",
      })
  }

  dbUsers, err := h.entClient.User.
    Query().
    All(context.Background())
  if err != nil {
    return ctx.
      Status(fiber.StatusInternalServerError).
      JSON(&fiber.Map{
        "success": false,
        "users": nil,
        "message": err.Error(),
      })
  }

  showUsers := make([]UserShowModel, len(dbUsers))

  for i, dbUser := range dbUsers {
    showUsers[i] = UserShowModel{
      ID: dbUser.ID.String(),
      Username: dbUser.Username,
      Role: dbUser.Role,
    }
  }

  return ctx.
    Status(fiber.StatusOK).
    JSON(&fiber.Map{
      "success": true,
      "users": showUsers,
      "message": "",
    })
}


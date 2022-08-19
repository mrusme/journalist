package users

import (
	"context"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	"github.com/mrusme/journalist/ent/user"
	// "github.com/mrusme/journalist/ent"
)

func (h *handler) Show(ctx *fiber.Ctx) error {
  var err error

  param_id := ctx.Params("id")
  id, err := uuid.Parse(param_id)
  if err != nil {
    return ctx.
      Status(fiber.StatusBadRequest).
      JSON(&fiber.Map{
        "success": false,
        "user": nil,
        "message": err.Error(),
      })
  }

  user_id := ctx.Locals("user_id").(string)
  role := ctx.Locals("role").(string)

  if param_id != user_id && role != "admin" {
    return ctx.
      Status(fiber.StatusForbidden).
      JSON(&fiber.Map{
        "success": false,
        "user": nil,
        "message": "Only admins are allowed to see other users",
      })
  }

  dbUser, err := h.EntClient.User.
    Query().
    Where(
      user.ID(id),
    ).
    Only(context.Background())
  if err != nil {
    return ctx.
      Status(fiber.StatusInternalServerError).
      JSON(&fiber.Map{
        "success": false,
        "user": nil,
        "message": err.Error(),
      })
  }

  showUser := UserShowModel{
    ID: dbUser.ID.String(),
    Username: dbUser.Username,
    Role: dbUser.Role,
  }

  return ctx.
    Status(fiber.StatusOK).
    JSON(&fiber.Map{
      "success": true,
      "user": showUser,
      "message": "",
    })
}


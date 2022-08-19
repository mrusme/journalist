package users

import (
  "context"
	"github.com/gofiber/fiber/v2"
	// "github.com/mrusme/journalist/ent"
	"github.com/mrusme/journalist/ent/user"
)

func (h *handler) List(ctx *fiber.Ctx) error {
  username := ctx.Locals("username").(string)

  u, err := h.EntClient.User.
    Query().
    Where(user.Username(username)).
    Only(context.Background())
  if err != nil {
    return ctx.
      Status(fiber.StatusUnauthorized).
      JSON(&fiber.Map{
        "success": false,
        "users": nil,
        "message": err.Error(),
      })
  }

  if u.Role != "admin" {
    return ctx.
      Status(fiber.StatusForbidden).
      JSON(&fiber.Map{
        "success": false,
        "users": nil,
        "message": "Only admins are allowed to list users",
      })
  }

  list, err := h.EntClient.User.
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

  return ctx.
    Status(fiber.StatusOK).
    JSON(&fiber.Map{
      "success": true,
      "users": list,
      "message": "",
    })
}


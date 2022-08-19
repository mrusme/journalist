package users

import (
  "log"
	"context"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	"github.com/mrusme/journalist/ent/user"
	// "github.com/mrusme/journalist/ent"
)

func (h *handler) List(ctx *fiber.Ctx) error {
  role := ctx.Locals("role").(string)

  if role != "admin" {
    return ctx.
      Status(fiber.StatusForbidden).
      JSON(&fiber.Map{
        "success": false,
        "users": nil,
        "message": "Only admins are allowed to list users",
      })
  }

  var list[]struct {
    ID                uuid.UUID     `json:"id"`
    Username          string        `json:"username"`
    Role              string        `json:"role"`
  }
  log.Println("ID: %s", user.FieldID)
  err := h.EntClient.User.
    Query().
    Select(
      user.FieldID,
      user.FieldUsername,
      user.FieldRole,
    ).
    Scan(context.Background(), &list)
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


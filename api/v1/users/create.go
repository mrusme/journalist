package users

import (
	"context"
	// "github.com/google/uuid"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"
	// "github.com/mrusme/journalist/ent/user"
	// "github.com/mrusme/journalist/ent"
)

func (h *handler) Create(ctx *fiber.Ctx) error {
  var err error

  role := ctx.Locals("role").(string)

  if role != "admin" {
    h.logger.Debug(
      "User not allowed to create users",
      zap.Error(err),
    )
    return ctx.
      Status(fiber.StatusForbidden).
      JSON(&fiber.Map{
        "success": false,
        "user": nil,
        "message": "Only admins are allowed to create users",
      })
  }

  createUser := new(UserCreateModel)
  if err = ctx.BodyParser(createUser); err != nil {
    h.logger.Debug(
      "Body parsing failed",
      zap.Error(err),
    )
    return ctx.
      Status(fiber.StatusInternalServerError).
      JSON(&fiber.Map{
        "success": false,
        "user": nil,
        "message": err.Error(),
      })
  }

  validate := validator.New()
  if err = validate.Struct(*createUser); err != nil {
    h.logger.Debug(
      "Validation failed",
      zap.Error(err),
    )
    return ctx.
      Status(fiber.StatusBadRequest).
      JSON(&fiber.Map{
        "success": false,
        "user": nil,
        "message": err.Error(),
      })
  }

  dbUser, err := h.entClient.User.
    Create().
    SetUsername(createUser.Username).
    SetPassword(createUser.Password).
    SetRole(createUser.Role).
    Save(context.Background())

  if err != nil {
    h.logger.Debug(
      "Could not create user",
      zap.Error(err),
    )
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



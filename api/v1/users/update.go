package users

import (
	"context"
	"github.com/google/uuid"
  "github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v2"
	// "github.com/mrusme/journalist/ent/user"
	// "github.com/mrusme/journalist/ent"
)

func (h *handler) Update(ctx *fiber.Ctx) error {
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
        "message": "Only admins are allowed to update other users",
      })
  }

  updateUser := new(UserUpdateModel)
  if err = ctx.BodyParser(updateUser); err != nil {
    return ctx.
      Status(fiber.StatusInternalServerError).
      JSON(&fiber.Map{
        "success": false,
        "user": nil,
        "message": err.Error(),
      })
  }

  validate := validator.New()
  if err = validate.Struct(*updateUser); err != nil {
    return ctx.
      Status(fiber.StatusBadRequest).
      JSON(&fiber.Map{
        "success": false,
        "user": nil,
        "message": err.Error(),
      })
  }

  dbUserTmp := h.EntClient.User.
    UpdateOneID(id)

  if updateUser.Role != "" {
    if role == "admin" {
      dbUserTmp = dbUserTmp.SetRole(updateUser.Role)
    } else {
      return ctx.
        Status(fiber.StatusForbidden).
        JSON(&fiber.Map{
          "success": false,
          "user": nil,
          "message": "Only admins are allowed to update roles",
        })
    }
  }

  if updateUser.Password != "" {
    dbUserTmp = dbUserTmp.
      SetPassword(updateUser.Password)
  }

  dbUser, err := dbUserTmp.Save(context.Background())

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




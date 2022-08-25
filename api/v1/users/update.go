package users

import (
	"context"
	"github.com/google/uuid"
  "github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v2"
	// "github.com/mrusme/journalist/ent/user"
	// "github.com/mrusme/journalist/ent"
)

type UserUpdateResponse struct {
  Success           bool            `json:"success"`
  User              *UserShowModel  `json:"user"`
  Message           string          `json:"message"`
}

// Update godoc
// @Summary      Update a user
// @Description  Change an existing user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string true "User ID"
// @Param        user body      UserUpdateModel true "Change user"
// @Success      200  {object}  UserUpdateResponse
// @Failure      400  {object}  UserUpdateResponse
// @Failure      404  {object}  UserUpdateResponse
// @Failure      500  {object}  UserUpdateResponse
// @Router       /users/{id} [put]
func (h *handler) Update(ctx *fiber.Ctx) error {
  var err error

  param_id := ctx.Params("id")
  id, err := uuid.Parse(param_id)
  if err != nil {
    return ctx.
      Status(fiber.StatusBadRequest).
      JSON(UserUpdateResponse{
        Success: false,
        User: nil,
        Message: err.Error(),
      })
  }

  user_id := ctx.Locals("user_id").(string)
  role := ctx.Locals("role").(string)

  if param_id != user_id && role != "admin" {
    return ctx.
      Status(fiber.StatusForbidden).
      JSON(UserUpdateResponse{
        Success: false,
        User: nil,
        Message: "Only admins are allowed to update other users",
      })
  }

  updateUser := new(UserUpdateModel)
  if err = ctx.BodyParser(updateUser); err != nil {
    return ctx.
      Status(fiber.StatusInternalServerError).
      JSON(UserUpdateResponse{
        Success: false,
        User: nil,
        Message: err.Error(),
      })
  }

  validate := validator.New()
  if err = validate.Struct(*updateUser); err != nil {
    return ctx.
      Status(fiber.StatusBadRequest).
      JSON(UserUpdateResponse{
        Success: false,
        User: nil,
        Message: err.Error(),
      })
  }

  dbUserTmp := h.entClient.User.
    UpdateOneID(id)

  if updateUser.Role != "" {
    if role == "admin" {
      dbUserTmp = dbUserTmp.SetRole(updateUser.Role)
    } else {
      return ctx.
        Status(fiber.StatusForbidden).
        JSON(UserUpdateResponse{
          Success: false,
          User: nil,
          Message: "Only admins are allowed to update roles",
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
      JSON(UserUpdateResponse{
        Success: false,
        User: nil,
        Message: err.Error(),
      })
  }

  showUser := UserShowModel{
    ID: dbUser.ID.String(),
    Username: dbUser.Username,
    Role: dbUser.Role,
  }

  return ctx.
    Status(fiber.StatusOK).
    JSON(UserUpdateResponse{
      Success: true,
      User: &showUser,
      Message: "",
    })
}


package users

import (
  "context"
  "github.com/google/uuid"

  "github.com/gofiber/fiber/v2"
  "github.com/mrusme/journalist/ent/user"
  // "github.com/mrusme/journalist/ent"
)

type UserShowResponse struct {
  Success           bool           `json:"success"`
  User              *UserShowModel `json:"user"`
  Message           string         `json:"message"`
}

// Show godoc
// @Summary      Show a user
// @Description  Get user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string true "User ID"
// @Success      200  {object}  UserShowResponse
// @Failure      400  {object}  UserShowResponse
// @Failure      404  {object}  UserShowResponse
// @Failure      500  {object}  UserShowResponse
// @Router       /users/{id} [get]
// @security     BasicAuth
func (h *handler) Show(ctx *fiber.Ctx) error {
  var err error

  param_id := ctx.Params("id")
  id, err := uuid.Parse(param_id)
  if err != nil {
    return ctx.
      Status(fiber.StatusBadRequest).
      JSON(UserShowResponse{
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
      JSON(UserShowResponse{
        Success: false,
        User: nil,
        Message: "Only admins are allowed to see other users",
      })
  }

  dbUser, err := h.entClient.User.
    Query().
    Where(
      user.ID(id),
    ).
    Only(context.Background())
  if err != nil {
    return ctx.
      Status(fiber.StatusInternalServerError).
      JSON(UserShowResponse{
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
    JSON(UserShowResponse{
      Success: true,
      User: &showUser,
      Message: "",
    })
}


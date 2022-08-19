package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrusme/journalist/ent"
)

type handler struct {
  EntClient *ent.Client
}

type UserShowModel struct {
  ID                string        `json:"id"`
  Username          string        `json:"username"`
  Role              string        `json:"role"`
}

type UserCreateModel struct {
  Username          string        `json:"username",validate:"required,alphanum,max=32"`
  Password          string        `json:"password",validate:"required"`
  Role              string        `json:"role",validate:"required"`
}

type UserUpdateModel struct {
  Password          string        `json:"password",validate:""`
  Role              string        `json:"role",validate:""`
}

func Register(fiberRouter *fiber.Router, entClient *ent.Client) () {
  endpoint := new(handler)
  endpoint.EntClient = entClient

  usersRouter := (*fiberRouter).Group("/users")
  usersRouter.Get("/", endpoint.List)
  usersRouter.Get("/:id", endpoint.Show)
  usersRouter.Post("/", endpoint.Create)
  usersRouter.Put("/:id", endpoint.Update)
  // usersRouter.Delete("/:id", endpoint.Destroy)
}


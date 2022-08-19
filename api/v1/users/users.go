package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrusme/journalist/ent"
)

type handler struct {
  EntClient *ent.Client
}

func Register(fiberRouter *fiber.Router, entClient *ent.Client) () {
  endpoint := new(handler)
  endpoint.EntClient = entClient

  usersRouter := (*fiberRouter).Group("/users")
  usersRouter.Get("/", endpoint.List)
  /* usersRouter.Get("/:id", endpoint.Show)
  usersRouter.Post("/", endpoint.Create)
  usersRouter.Put("/:id", endpoint.Update)
  usersRouter.Delete("/:id", endpoint.Destroy) */
}


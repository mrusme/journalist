package subscriptions

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

  subscriptionsRouter := (*fiberRouter).Group("/subscriptions")
  subscriptionsRouter.Get("/", endpoint.List)
  // subscriptionsRouter.Get("/:id", endpoint.Show)
  // subscriptionsRouter.Post("/", endpoint.Create)
  // subscriptionsRouter.Put("/:id", endpoint.Update)
  // subscriptionsRouter.Delete("/:id", endpoint.Destroy)
}


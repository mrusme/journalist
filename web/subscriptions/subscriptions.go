package subscriptions

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrusme/journalist/ent"
	"github.com/mrusme/journalist/journalistd"
)

type handler struct {
  config    *journalistd.Config
  EntClient *ent.Client
}

func Register(
  config *journalistd.Config,
  fiberRouter *fiber.Router,
  entClient *ent.Client,
) () {
  endpoint := new(handler)
  endpoint.config = config
  endpoint.EntClient = entClient

  subscriptionsRouter := (*fiberRouter).Group("/subscriptions")
  subscriptionsRouter.Get("/", endpoint.List)
  // subscriptionsRouter.Get("/:id", endpoint.Show)
  // subscriptionsRouter.Post("/", endpoint.Create)
  // subscriptionsRouter.Put("/:id", endpoint.Update)
  // subscriptionsRouter.Delete("/:id", endpoint.Destroy)
}


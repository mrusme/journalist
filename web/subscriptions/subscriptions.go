package subscriptions

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrusme/journalist/ent"
	"github.com/mrusme/journalist/journalistd"
	"go.uber.org/zap"
)

type handler struct {
  config    *journalistd.Config
  entClient *ent.Client
  logger    *zap.Logger
}

func Register(
  config *journalistd.Config,
  fiberRouter *fiber.Router,
  entClient *ent.Client,
  logger *zap.Logger,
) () {
  endpoint := new(handler)
  endpoint.config = config
  endpoint.entClient = entClient
  endpoint.logger = logger

  subscriptionsRouter := (*fiberRouter).Group("/subscriptions")
  subscriptionsRouter.Get("/", endpoint.List)
  // subscriptionsRouter.Get("/:id", endpoint.Show)
  // subscriptionsRouter.Post("/", endpoint.Create)
  // subscriptionsRouter.Put("/:id", endpoint.Update)
  // subscriptionsRouter.Delete("/:id", endpoint.Destroy)
}


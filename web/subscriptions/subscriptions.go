package subscriptions

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrusme/journalist/ent"
	"github.com/mrusme/journalist/lib"
	"go.uber.org/zap"
)

type handler struct {
  jctx      *lib.JournalistContext

  config    *lib.Config
  entClient *ent.Client
  logger    *zap.Logger
}

func Register(
  jctx *lib.JournalistContext,
  fiberRouter *fiber.Router,
) () {
  endpoint := new(handler)
  endpoint.jctx = jctx
  endpoint.config = endpoint.jctx.Config
  endpoint.entClient = endpoint.jctx.EntClient
  endpoint.logger = endpoint.jctx.Logger

  subscriptionsRouter := (*fiberRouter).Group("/subscriptions")
  subscriptionsRouter.Get("/", endpoint.List)
  // subscriptionsRouter.Get("/:id", endpoint.Show)
  // subscriptionsRouter.Post("/", endpoint.Create)
  // subscriptionsRouter.Put("/:id", endpoint.Update)
  // subscriptionsRouter.Delete("/:id", endpoint.Destroy)
}


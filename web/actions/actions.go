package actions

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

  actionsRouter := (*fiberRouter).Group("/actions")
  actionsRouter.Get("/read/:id", endpoint.Read)
  actionsRouter.Get("/read_older/:id", endpoint.ReadOlder)
  actionsRouter.Get("/read_newer/:id", endpoint.ReadNewer)
  // actionsRouter.Get("/read_all/:id", endpoint.ReadAll)
}



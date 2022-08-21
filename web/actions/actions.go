package actions

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

  actionsRouter := (*fiberRouter).Group("/actions")
  actionsRouter.Get("/read/:id", endpoint.Read)
  // actionsRouter.Get("/read_older/:id", endpoint.ReadOlder)
  // actionsRouter.Get("/read_newer/:id", endpoint.ReadNewer)
  // actionsRouter.Get("/read_all/:id", endpoint.ReadAll)
}



package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrusme/journalist/api/v1/feeds"
	"github.com/mrusme/journalist/api/v1/tokens"
	"github.com/mrusme/journalist/api/v1/users"
	"github.com/mrusme/journalist/ent"
	"github.com/mrusme/journalist/journalistd"
	"go.uber.org/zap"
)

func Register(
  config *journalistd.Config,
  fiberRouter *fiber.Router,
  entClient *ent.Client,
  logger *zap.Logger,
) () {
  v1 := (*fiberRouter).Group("/v1")

  users.Register(
    config,
    &v1,
    entClient,
    logger,
  )

  tokens.Register(
    config,
    &v1,
    entClient,
    logger,
  )

  feeds.Register(
    config,
    &v1,
    entClient,
    logger,
  )
}

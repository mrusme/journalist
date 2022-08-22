package tokens

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

type TokenShowModel struct {
  ID                string        `json:"id"`
  Type              string        `json:"type"`
  Name              string        `json:"tokenname"`
  Token             string        `json:"token"`
}

type TokenCreateModel struct {
  Name              string        `json:"name",validate:"required,alphanum,max=32"`
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

  tokensRouter := (*fiberRouter).Group("/tokens")
  // tokensRouter.Get("/", endpoint.List)
  // tokensRouter.Get("/:id", endpoint.Show)
  tokensRouter.Post("/", endpoint.Create)
  // tokensRouter.Put("/:id", endpoint.Update)
  // tokensRouter.Delete("/:id", endpoint.Destroy)
}



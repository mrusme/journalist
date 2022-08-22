package tokens

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
  jctx *lib.JournalistContext,
  fiberRouter *fiber.Router,
) () {
  endpoint := new(handler)
  endpoint.jctx = jctx
  endpoint.config = endpoint.jctx.Config
  endpoint.entClient = endpoint.jctx.EntClient
  endpoint.logger = endpoint.jctx.Logger

  tokensRouter := (*fiberRouter).Group("/tokens")
  // tokensRouter.Get("/", endpoint.List)
  // tokensRouter.Get("/:id", endpoint.Show)
  tokensRouter.Post("/", endpoint.Create)
  // tokensRouter.Put("/:id", endpoint.Update)
  // tokensRouter.Delete("/:id", endpoint.Destroy)
}



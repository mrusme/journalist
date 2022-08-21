package tokens

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrusme/journalist/ent"
)

type handler struct {
  EntClient *ent.Client
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


func Register(fiberRouter *fiber.Router, entClient *ent.Client) () {
  endpoint := new(handler)
  endpoint.EntClient = entClient

  tokensRouter := (*fiberRouter).Group("/tokens")
  // tokensRouter.Get("/", endpoint.List)
  // tokensRouter.Get("/:id", endpoint.Show)
  tokensRouter.Post("/", endpoint.Create)
  // tokensRouter.Put("/:id", endpoint.Update)
  // tokensRouter.Delete("/:id", endpoint.Destroy)
}



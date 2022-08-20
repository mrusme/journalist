package subscriptions

import (
  "log"
	// "context"
	// "github.com/google/uuid"

	"context"

	"github.com/gofiber/fiber/v2"
	// "github.com/mrusme/journalist/ent/user"
	// "github.com/mrusme/journalist/ent"
)

func (h *handler) List(ctx *fiber.Ctx) error {
  dbItems, err := h.EntClient.Item.Query().All(context.Background())
  if err != nil {
    ctx.SendStatus(fiber.StatusInternalServerError)
    return err
  }

  err = ctx.Render("views/subscriptions.list", fiber.Map{
    "Title": "Hello World",
    "Items": dbItems,
  })
  ctx.Set("Content-type", "text/xml; charset=utf-8")
  return err
}



package subscriptions

import (
	// "context"
	// "github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	// "github.com/mrusme/journalist/ent/user"
	// "github.com/mrusme/journalist/ent"
)

func (h *handler) List(ctx *fiber.Ctx) error {
  err := ctx.Render("views/subscriptions.list", fiber.Map{
    "Title": "Hello World",
  })
  ctx.Set("Content-type", "text/xml; charset=utf-8")
  return err
}



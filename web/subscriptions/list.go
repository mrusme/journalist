package subscriptions

import (
	// "log"
	// "context"
	"github.com/google/uuid"
	// "github.com/mrusme/journalist/ent"
	"github.com/mrusme/journalist/ent/item"
	// "github.com/mrusme/journalist/ent/read"
	"github.com/mrusme/journalist/ent/subscription"
	"github.com/mrusme/journalist/ent/user"

	"context"

	"github.com/gofiber/fiber/v2"
	// "github.com/mrusme/journalist/ent/user"
	// "github.com/mrusme/journalist/ent"
)

func (h *handler) List(ctx *fiber.Ctx) error {
  sessionUserId := ctx.Locals("user_id").(string)
  myId, err := uuid.Parse(sessionUserId)
  if err != nil {
    ctx.SendStatus(fiber.StatusInternalServerError)
    return err
  }

  dbItems, err := h.EntClient.Subscription.
    Query().
    Where(subscription.UserID(myId)).
    QueryFeed().
    QueryItems().
    Where(
      item.Not(
        item.HasReadByUsersWith(
          user.ID(myId),
        ),
      ),
    ).
    All(context.Background())
  if err != nil {
    ctx.SendStatus(fiber.StatusInternalServerError)
    return err
  }

  err = ctx.Render("views/subscriptions.list", fiber.Map{
    "Title": "Subscriptions",
    "Items": dbItems,
  })
  ctx.Set("Content-type", "text/xml; charset=utf-8")
  return err
}


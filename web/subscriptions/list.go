package subscriptions

import (
  "fmt"
  "time"
	"github.com/google/uuid"
	"github.com/mrusme/journalist/ent/item"
	"github.com/mrusme/journalist/ent/subscription"
	"github.com/mrusme/journalist/ent/user"

	"context"

	"github.com/gofiber/fiber/v2"
)

func (h *handler) List(ctx *fiber.Ctx) error {
  qat := ctx.Query("qat")
  sessionUsername := ctx.Locals("username").(string)
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
    "Config": h.config,
    "QAT": qat,

    "Title": "Subscriptions",
    "Link": fmt.Sprintf(
      "%s/subscriptions?qat=%s",
      h.config.Server.Endpoint.Web,
      qat,
    ),
    "Description": fmt.Sprintf(
      "%s' subscriptions",
      sessionUsername,
    ),
    "Generator": "Journalist",
    "Language": "en-us",
    "LastBuildDate": time.Now(),

    "Items": dbItems,
  })
  ctx.Set("Content-type", "text/xml; charset=utf-8")
  return err
}


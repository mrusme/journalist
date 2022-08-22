package feeds

import (
	"context"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	"github.com/mrusme/journalist/ent/feed"
	// "github.com/mrusme/journalist/ent"
)

func (h *handler) Show(ctx *fiber.Ctx) error {
  var err error

  param_id := ctx.Params("id")
  id, err := uuid.Parse(param_id)
  if err != nil {
    return ctx.
      Status(fiber.StatusBadRequest).
      JSON(&fiber.Map{
        "success": false,
        "feed": nil,
        "message": err.Error(),
      })
  }

  feed_id := ctx.Locals("feed_id").(string)
  role := ctx.Locals("role").(string)

  if param_id != feed_id && role != "admin" {
    return ctx.
      Status(fiber.StatusForbidden).
      JSON(&fiber.Map{
        "success": false,
        "feed": nil,
        "message": "Only admins are allowed to see other feeds",
      })
  }

  dbFeed, err := h.entClient.Feed.
    Query().
    Where(
      feed.ID(id),
    ).
    Only(context.Background())
  if err != nil {
    return ctx.
      Status(fiber.StatusInternalServerError).
      JSON(&fiber.Map{
        "success": false,
        "feed": nil,
        "message": err.Error(),
      })
  }

  showFeed := FeedShowModel{
    ID: dbFeed.ID.String(),
    Name: dbFeed.FeedTitle,
    URL: dbFeed.FeedFeedLink,
    Group: "*",
  }

  return ctx.
    Status(fiber.StatusOK).
    JSON(&fiber.Map{
      "success": true,
      "feed": showFeed,
      "message": "",
    })
}



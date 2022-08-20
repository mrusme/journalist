package feeds

import (
	"context"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	"github.com/mrusme/journalist/ent/user"
	// "github.com/mrusme/journalist/ent"
)

func (h *handler) List(ctx *fiber.Ctx) error {
  var showFeeds []FeedShowModel

  role := ctx.Locals("role").(string)

  if role == "admin" {
    dbFeeds, err := h.EntClient.Feed.
      Query().
      All(context.Background())
    if err != nil {
      return ctx.
        Status(fiber.StatusInternalServerError).
        JSON(&fiber.Map{
          "success": false,
          "feeds": nil,
          "message": err.Error(),
        })
    }

    showFeeds = make([]FeedShowModel, len(dbFeeds))

    for i, dbFeed := range dbFeeds {
      showFeeds[i] = FeedShowModel{
        ID: dbFeed.ID.String(),
        Name: dbFeed.FeedTitle,
        URL: dbFeed.FeedFeedLink,
        Group: "*",
      }
    }
  } else {
    sessionUserId := ctx.Locals("user_id").(string)
    myId, err := uuid.Parse(sessionUserId)
    if err != nil {
      return ctx.
        Status(fiber.StatusInternalServerError).
        JSON(&fiber.Map{
          "success": false,
          "feed": nil,
          "message": err.Error(),
        })
    }

    dbUser, err := h.EntClient.User.
      Query().
      WithSubscribedFeeds().
      WithSubscriptions().
      Where(
        user.ID(myId),
      ).
      Only(context.Background())

    for i, feed := range dbUser.Edges.SubscribedFeeds {
      showFeeds = append(showFeeds, FeedShowModel{
        ID: feed.ID.String(),
        Name: dbUser.Edges.Subscriptions[i].Name,
        URL: feed.URL,
        Group: dbUser.Edges.Subscriptions[i].Group,
      })
    }
  }

  return ctx.
    Status(fiber.StatusOK).
    JSON(&fiber.Map{
      "success": true,
      "feeds": showFeeds,
      "message": "",
    })
}


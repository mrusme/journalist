package feeds

import (
	"context"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	"github.com/mrusme/journalist/ent/feed"
	// "github.com/mrusme/journalist/ent"
	"go.uber.org/zap"
)

type FeedShowResponse struct {
  Success           bool           `json:"success"`
  Feed              *FeedShowModel `json:"feed"`
  Message           string         `json:"message"`
}

// Show godoc
// @Summary      Show a feed
// @Description  Get feed by ID
// @Tags         feeds
// @Accept       json
// @Produce      json
// @Param        id   path      string true "Feed ID"
// @Success      200  {object}  FeedShowResponse
// @Failure      400  {object}  FeedShowResponse
// @Failure      404  {object}  FeedShowResponse
// @Failure      500  {object}  FeedShowResponse
// @Router       /feeds/{id} [get]
// @security     BasicAuth
func (h *handler) Show(ctx *fiber.Ctx) error {
  var err error

  param_id := ctx.Params("id")
  id, err := uuid.Parse(param_id)
  if err != nil {
    h.logger.Debug(
      "Could not parse user ID",
      zap.Error(err),
    )
    return ctx.
      Status(fiber.StatusBadRequest).
      JSON(FeedShowResponse{
        Success: false,
        Feed: nil,
        Message: err.Error(),
      })
  }

  feed_id := ctx.Locals("feed_id").(string)
  role := ctx.Locals("role").(string)

  if param_id != feed_id && role != "admin" {
    h.logger.Debug(
      "User not allowed to see other feeds",
      zap.Error(err),
    )
    return ctx.
      Status(fiber.StatusForbidden).
      JSON(FeedShowResponse{
        Success: false,
        Feed: nil,
        Message: "Only admins are allowed to see other feeds",
      })
  }

  dbFeed, err := h.entClient.Feed.
    Query().
    Where(
      feed.ID(id),
    ).
    Only(context.Background())
  if err != nil {
    h.logger.Debug(
      "Could not query feed",
      zap.String("feedID", param_id),
      zap.Error(err),
    )
    return ctx.
      Status(fiber.StatusInternalServerError).
      JSON(FeedShowResponse{
        Success: false,
        Feed: nil,
        Message: err.Error(),
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
    JSON(FeedShowResponse{
      Success: true,
      Feed: &showFeed,
      Message: "",
    })
}



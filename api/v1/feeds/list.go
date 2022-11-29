package feeds

import (
	"context"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	"github.com/mrusme/journalist/ent/user"

	// "github.com/mrusme/journalist/ent"
	"go.uber.org/zap"
)

type FeedListResponse struct {
	Success bool             `json:"success"`
	Feeds   *[]FeedShowModel `json:"feeds"`
	Message string           `json:"message"`
}

// List godoc
// @Summary      List feeds
// @Description  Get all feeds
// @Tags         feeds
// @Accept       json
// @Produce      json
// @Success      200  {object}  FeedListResponse
// @Failure      400  {object}  FeedListResponse
// @Failure      404  {object}  FeedListResponse
// @Failure      500  {object}  FeedListResponse
// @Router       /feeds [get]
// @security     BasicAuth
func (h *handler) List(ctx *fiber.Ctx) error {
	var showFeeds []FeedShowModel

	role := ctx.Locals("role").(string)

	if role == "admin" {
		dbFeeds, err := h.entClient.Feed.
			Query().
			All(context.Background())
		if err != nil {
			h.logger.Debug(
				"Could not query all feeds",
				zap.Error(err),
			)
			return ctx.
				Status(fiber.StatusInternalServerError).
				JSON(FeedListResponse{
					Success: false,
					Feeds:   nil,
					Message: err.Error(),
				})
		}

		showFeeds = make([]FeedShowModel, len(dbFeeds))

		for i, dbFeed := range dbFeeds {
			showFeeds[i] = FeedShowModel{
				ID:    dbFeed.ID.String(),
				Name:  dbFeed.FeedTitle,
				URL:   dbFeed.FeedFeedLink,
				Group: "*",
			}
		}
	} else {
		sessionUserId := ctx.Locals("user_id").(string)
		myId, err := uuid.Parse(sessionUserId)
		if err != nil {
			h.logger.Debug(
				"Could not parse user ID",
				zap.Error(err),
			)
			return ctx.
				Status(fiber.StatusInternalServerError).
				JSON(FeedListResponse{
					Success: false,
					Feeds:   nil,
					Message: err.Error(),
				})
		}

		dbUser, err := h.entClient.User.
			Query().
			WithSubscribedFeeds().
			WithSubscriptions().
			Where(
				user.ID(myId),
			).
			Only(context.Background())

		for i, feed := range dbUser.Edges.SubscribedFeeds {
			showFeeds = append(showFeeds, FeedShowModel{
				ID:    feed.ID.String(),
				Name:  dbUser.Edges.Subscriptions[i].Name,
				URL:   feed.URL,
				Group: dbUser.Edges.Subscriptions[i].Group,
			})
		}
	}

	return ctx.
		Status(fiber.StatusOK).
		JSON(FeedListResponse{
			Success: true,
			Feeds:   &showFeeds,
			Message: "",
		})
}

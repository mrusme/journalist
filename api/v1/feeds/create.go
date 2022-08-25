package feeds

import (
  // "strings"
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	// "github.com/mrusme/journalist/ent/user"
	// "github.com/mrusme/journalist/ent"

	"github.com/mrusme/journalist/crawler"
	"github.com/mrusme/journalist/rss"

	"go.uber.org/zap"
)

type FeedCreateResponse struct {
  Success           bool           `json:"success"`
  Feed              *FeedShowModel `json:"feed"`
  Message           string         `json:"message"`
}

// Create godoc
// @Summary      Create a feed
// @Description  Add a new feed
// @Tags         feeds
// @Accept       json
// @Produce      json
// @Param        feed body      FeedCreateModel  true "Add feed"
// @Success      200  {object}  FeedCreateResponse
// @Failure      400  {object}  FeedCreateResponse
// @Failure      404  {object}  FeedCreateResponse
// @Failure      500  {object}  FeedCreateResponse
// @Router       /feeds [post]
func (h *handler) Create(ctx *fiber.Ctx) error {
  var err error

  // sessionId := ctx.Locals("user_id").(string)
  // sessionRole := ctx.Locals("role").(string)

  createFeed := new(FeedCreateModel)
  if err = ctx.BodyParser(createFeed); err != nil {
    h.logger.Debug(
      "Body parsing failed",
      zap.Error(err),
    )
    return ctx.
      Status(fiber.StatusInternalServerError).
      JSON(FeedCreateResponse{
        Success: false,
        Feed: nil,
        Message: err.Error(),
      })
  }

  validate := validator.New()
  if err = validate.Struct(*createFeed); err != nil {
    h.logger.Debug(
      "Validation failed",
      zap.Error(err),
    )
    return ctx.
      Status(fiber.StatusBadRequest).
      JSON(FeedCreateResponse{
        Success: false,
        Feed: nil,
        Message: err.Error(),
      })
  }

  crwlr := crawler.New(h.logger)
  defer crwlr.Close()

  crwlr.SetLocation(createFeed.URL)

  if createFeed.Username != "" && createFeed.Password != "" {
    crwlr.SetBasicAuth(createFeed.Username, createFeed.Password)
  }

  _, feedLink, err := crwlr.GetFeedLink()
  if err != nil {
    h.logger.Debug(
      "Could not get feed link",
      zap.Error(err),
    )
    return ctx.
      Status(fiber.StatusBadRequest).
      JSON(FeedCreateResponse{
        Success: false,
        Feed: nil,
        Message: err.Error(),
      })
  }

  rc, errr := rss.NewClient(
    feedLink,
    createFeed.Username,
    createFeed.Password,
    false,
    h.logger,
  )
  if len(errr) > 0 {
    h.logger.Debug(
      "Could not fetch feed",
      zap.Error(err),
    )
    return ctx.
      Status(fiber.StatusInternalServerError).
      JSON(FeedCreateResponse{
        Success: false,
        Feed: nil,
        Message: err.Error(),
      })
  }

  dbFeedTmp := h.entClient.Feed.
    Create()

  dbFeedTmp = rc.SetFeed(
    feedLink,
    createFeed.Username,
    createFeed.Password,
    dbFeedTmp,
  )
  feedId, err := dbFeedTmp.
    OnConflict().
    UpdateNewValues().
    ID(context.Background())
  if err != nil {
    h.logger.Debug(
      "Could not upsert feed",
      zap.Error(err),
    )
    return ctx.
      Status(fiber.StatusInternalServerError).
      JSON(FeedCreateResponse{
        Success: false,
        Feed: nil,
        Message: err.Error(),
      })
  }

  sessionUserId := ctx.Locals("user_id").(string)
  myId, err := uuid.Parse(sessionUserId)
  if err != nil {
    h.logger.Debug(
      "Could not parse user ID",
      zap.Error(err),
    )
    return ctx.
      Status(fiber.StatusInternalServerError).
      JSON(FeedCreateResponse{
        Success: false,
        Feed: nil,
        Message: err.Error(),
      })
  }

  dbSubscriptionTmp := h.entClient.Subscription.
    Create().
    SetUserID(myId).
    SetFeedID(feedId)

  if createFeed.Name != "" {
    dbSubscriptionTmp = dbSubscriptionTmp.
      SetName(createFeed.Name)
  }

  if createFeed.Group != "" {
    dbSubscriptionTmp = dbSubscriptionTmp.
      SetGroup(createFeed.Group)
  }

  dbSubscription, err := dbSubscriptionTmp.
    Save(context.Background())
  if err != nil {
    h.logger.Debug(
      "Could not add feed subscription",
      zap.Error(err),
    )
    return ctx.
      Status(fiber.StatusInternalServerError).
      JSON(FeedCreateResponse{
        Success: false,
        Feed: nil,
        Message: err.Error(),
      })
  }

  showFeed := FeedShowModel{
    ID: feedId.String(),
    Name: dbSubscription.Name,
    URL: createFeed.URL,
    Group: dbSubscription.Group,
  }

  return ctx.
    Status(fiber.StatusOK).
    JSON(FeedCreateResponse{
      Success: true,
      Feed: &showFeed,
      Message: "",
    })
}




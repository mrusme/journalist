package feeds

import (
	"context"
	// "github.com/google/uuid"
  "github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v2"
	// "github.com/mrusme/journalist/ent/feed"
	// "github.com/mrusme/journalist/ent"

	"github.com/mrusme/journalist/crawler"
)

func (h *handler) Create(ctx *fiber.Ctx) error {
  var err error

  // sessionId := ctx.Locals("user_id").(string)
  // sessionRole := ctx.Locals("role").(string)

  createFeed := new(FeedCreateModel)
  if err = ctx.BodyParser(createFeed); err != nil {
    return ctx.
      Status(fiber.StatusInternalServerError).
      JSON(&fiber.Map{
        "success": false,
        "feed": nil,
        "message": err.Error(),
      })
  }

  validate := validator.New()
  if err = validate.Struct(*createFeed); err != nil {
    return ctx.
      Status(fiber.StatusBadRequest).
      JSON(&fiber.Map{
        "success": false,
        "feed": nil,
        "message": err.Error(),
      })
  }

  dbFeedTmp := h.EntClient.Feed.
    Create()

  crwlr := crawler.New()
  defer crwlr.Close()

  crwlr.SetLocation(createFeed.URL)

  if createFeed.Username != "" && createFeed.Password != "" {
    crwlr.SetBasicAuth(createFeed.Username, createFeed.Password)

    dbFeedTmp = dbFeedTmp.
      SetUsername(createFeed.Username).
      SetPassword(createFeed.Password)
  }

  feedType, feedLink, err := crwlr.GetFeedLink()
  if err != nil {
    return ctx.
      Status(fiber.StatusBadRequest).
      JSON(&fiber.Map{
        "success": false,
        "feed": nil,
        "message": err.Error(),
      })
  }

  dbFeed, err := dbFeedTmp.
    SetURL(createFeed.URL).
    OnConflict().
    Ignore().
    Save(context.Background())
  if err != nil {
    return ctx.
      Status(fiber.StatusInternalServerError).
      JSON(&fiber.Map{
        "success": false,
        "feed": nil,
        "message": err.Error(),
      })
  }

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

  h.EntClient.User.UpdateOneID(myId).AddSubscription()

  showFeed := FeedShowModel{
    ID: dbFeed.ID.String(),
    URL: dbFeed.URL,
  }

  return ctx.
    Status(fiber.StatusOK).
    JSON(&fiber.Map{
      "success": true,
      "feed": showFeed,
      "message": "",
    })
}




package feeds

import (
	"context"
	// "github.com/google/uuid"
  "github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v2"
	// "github.com/mrusme/journalist/ent/feed"
	// "github.com/mrusme/journalist/ent"
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

  // if createFeed.Name != "" {
    // dbFeedTmp = dbFeedTmp.
      // SetName(createFeed.Name)
  // }

  dbFeed, err := dbFeedTmp.
    SetURL(createFeed.URL).
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




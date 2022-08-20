package feeds

import (
  "strings"
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	// "github.com/mrusme/journalist/ent/user"
	// "github.com/mrusme/journalist/ent"

	"github.com/mrusme/journalist/crawler"
	"github.com/mrusme/journalist/rss"
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

  _, feedLink, err := crwlr.GetFeedLink()
  if err != nil {
    return ctx.
      Status(fiber.StatusBadRequest).
      JSON(&fiber.Map{
        "success": false,
        "feed": nil,
        "message": err.Error(),
      })
  }

  rssClient, err := rss.NewClient(feedLink)

  if rssClient.Feed.Author != nil {
    dbFeedTmp.
      SetFeedAuthorName(rssClient.Feed.Author.Name).
      SetFeedAuthorEmail(rssClient.Feed.Author.Email)
  }
  if rssClient.Feed.Image != nil {
    dbFeedTmp.
      SetFeedImageTitle(rssClient.Feed.Image.Title).
      SetFeedImageURL(rssClient.Feed.Image.URL)
  }

  feedId, err := dbFeedTmp.
    SetURL(feedLink).
    SetFeedTitle(rssClient.Feed.Title).
    SetFeedDescription(rssClient.Feed.Description).
    SetFeedLink(rssClient.Feed.Link).
    SetFeedFeedLink(rssClient.Feed.FeedLink).
    SetFeedUpdated(rssClient.Feed.Updated).
    SetFeedPublished(rssClient.Feed.Published).
    SetFeedLanguage(rssClient.Feed.Language).
    SetFeedCopyright(rssClient.Feed.Copyright).
    SetFeedGenerator(rssClient.Feed.Generator).
    SetFeedCopyright(rssClient.Feed.Copyright).
    SetFeedCategories(strings.Join(rssClient.Feed.Categories, ", ")).
    OnConflict().
    // Ignore().
    UpdateNewValues().
    ID(context.Background())
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

  dbSubscriptionTmp := h.EntClient.Subscription.
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
    return ctx.
      Status(fiber.StatusInternalServerError).
      JSON(&fiber.Map{
        "success": false,
        "feed": nil,
        "message": err.Error(),
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
    JSON(&fiber.Map{
      "success": true,
      "feed": showFeed,
      "message": "",
    })
}




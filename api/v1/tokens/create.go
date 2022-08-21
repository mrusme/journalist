package tokens

import (
	"context"
	"time"
  "fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/mrusme/journalist/rss"

	"github.com/gofiber/fiber/v2"
	// "github.com/mrusme/journalist/ent/token"
	// "github.com/mrusme/journalist/ent"
)

func (h *handler) Create(ctx *fiber.Ctx) error {
  var err error

  createToken := new(TokenCreateModel)
  if err = ctx.BodyParser(createToken); err != nil {
    return ctx.
      Status(fiber.StatusInternalServerError).
      JSON(&fiber.Map{
        "success": false,
        "token": nil,
        "message": err.Error(),
      })
  }

  validate := validator.New()
  if err = validate.Struct(*createToken); err != nil {
    return ctx.
      Status(fiber.StatusBadRequest).
      JSON(&fiber.Map{
        "success": false,
        "token": nil,
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

  // TODO: Move GenerateGUID to a common helper, rename
  token := rss.GenerateGUID(
    fmt.Sprintf(
      "%s-%d",
      sessionUserId,
      time.Now().UnixNano(),
    ),
  )

  dbToken, err := h.EntClient.Token.
    Create().
    SetType("qat").
    SetName(createToken.Name).
    SetToken(token).
    Save(context.Background())

  if err != nil {
    return ctx.
      Status(fiber.StatusInternalServerError).
      JSON(&fiber.Map{
        "success": false,
        "token": nil,
        "message": err.Error(),
      })
  }

  _, err = h.EntClient.User.
    UpdateOneID(myId).
    AddTokenIDs(dbToken.ID).
    Save(context.Background())

  if err != nil {
    return ctx.
      Status(fiber.StatusInternalServerError).
      JSON(&fiber.Map{
        "success": false,
        "token": nil,
        "message": err.Error(),
      })
  }

  showToken := TokenShowModel{
    ID: dbToken.ID.String(),
    Type: dbToken.Type,
    Name: dbToken.Name,
    Token: dbToken.Token,
  }

  return ctx.
    Status(fiber.StatusOK).
    JSON(&fiber.Map{
      "success": true,
      "token": showToken,
      "message": "",
    })
}



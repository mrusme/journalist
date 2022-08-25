package tokens

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/mrusme/journalist/rss"
	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"
	// "github.com/mrusme/journalist/ent/token"
	// "github.com/mrusme/journalist/ent"
)

type TokenCreateResponse struct {
  Success           bool            `json:"success"`
  Token             *TokenShowModel `json:"token"`
  Message           string          `json:"message"`
}

// Create godoc
// @Summary      Create a token
// @Description  Add a new token
// @Tags         tokens
// @Accept       json
// @Produce      json
// @Param        token body     TokenCreateModel true "Add token"
// @Success      200  {object}  TokenCreateResponse
// @Failure      400  {object}  TokenCreateResponse
// @Failure      404  {object}  TokenCreateResponse
// @Failure      500  {object}  TokenCreateResponse
// @Router       /tokens [post]
func (h *handler) Create(ctx *fiber.Ctx) error {
  var err error

  createToken := new(TokenCreateModel)
  if err = ctx.BodyParser(createToken); err != nil {
    h.logger.Debug(
      "Body parsing failed",
      zap.Error(err),
    )
    return ctx.
      Status(fiber.StatusInternalServerError).
      JSON(TokenCreateResponse{
        Success: false,
        Token: nil,
        Message: err.Error(),
      })
  }

  validate := validator.New()
  if err = validate.Struct(*createToken); err != nil {
    h.logger.Debug(
      "Validation failed",
      zap.Error(err),
    )
    return ctx.
      Status(fiber.StatusBadRequest).
      JSON(TokenCreateResponse{
        Success: false,
        Token: nil,
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
      JSON(TokenCreateResponse{
        Success: false,
        Token: nil,
        Message: err.Error(),
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

  dbToken, err := h.entClient.Token.
    Create().
    SetType("qat").
    SetName(createToken.Name).
    SetToken(token).
    Save(context.Background())

  if err != nil {
    h.logger.Debug(
      "Could create token",
      zap.Error(err),
    )
    return ctx.
      Status(fiber.StatusInternalServerError).
      JSON(TokenCreateResponse{
        Success: false,
        Token: nil,
        Message: err.Error(),
      })
  }

  _, err = h.entClient.User.
    UpdateOneID(myId).
    AddTokenIDs(dbToken.ID).
    Save(context.Background())

  if err != nil {
    h.logger.Debug(
      "Could not add new token to user",
      zap.Error(err),
    )
    return ctx.
      Status(fiber.StatusInternalServerError).
      JSON(TokenCreateResponse{
        Success: false,
        Token: nil,
        Message: err.Error(),
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
    JSON(TokenCreateResponse{
      Success: true,
      Token: &showToken,
      Message: "",
    })
}




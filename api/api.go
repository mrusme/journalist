package api

import (
	"encoding/base64"
	"strings"

	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/mrusme/journalist/api/v1"
	"github.com/mrusme/journalist/ent"
	"github.com/mrusme/journalist/ent/user"
	"github.com/mrusme/journalist/lib"
)

// @title Journalist API
// @version 1.0
// @description The Journalist REST API v1

// @contact.name Marius
// @contact.url https://xn--gckvb8fzb.com
// @contact.email marius@xn--gckvb8fzb.com

// @license.name GPL-3.0
// @license.url https://github.com/mrusme/journalist/blob/master/LICENSE

// @host localhost:8000
// @BasePath /api/v1
// @accept json
// @produce json
// @schemes http
// @securityDefinitions.basic  BasicAuth
func Register(
  jctx *lib.JournalistContext,
  fiberApp *fiber.App,
) () {
  api := fiberApp.Group("/api")
  api.Use(cors.New())
  api.Use(authorizer(jctx.EntClient))

  v1.Register(
    jctx,
    &api,
  )
}

// TODO: Move to `middlewares`
func authorizer(entClient *ent.Client) fiber.Handler {
  return func (ctx *fiber.Ctx) error {
    auth := ctx.Get(fiber.HeaderAuthorization)

    if len(auth) <= 6 || strings.ToLower(auth[:5]) != "basic" {
      return ctx.SendStatus(fiber.StatusUnauthorized)
    }

    raw, err := base64.StdEncoding.DecodeString(auth[6:])
    if err != nil {
      return ctx.SendStatus(fiber.StatusUnauthorized)
    }

    creds := utils.UnsafeString(raw)

    index := strings.Index(creds, ":")
    if index == -1 {
      return ctx.SendStatus(fiber.StatusUnauthorized)
    }

    username := creds[:index]
    password := creds[index+1:]

    u, err := entClient.User.
      Query().
      Where(user.Username(username)).
      Only(context.Background())
    if err != nil {
      return ctx.SendStatus(fiber.StatusUnauthorized)
    }

    if u.Password != password {
      return ctx.SendStatus(fiber.StatusUnauthorized)
    }

    ctx.Locals("user_id", u.ID.String())
    ctx.Locals("username", u.Username)
    // ctx.Locals("password", u.Password)
    ctx.Locals("role", u.Role)
    return ctx.Next()
  }
}

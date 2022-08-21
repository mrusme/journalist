package web

import (
  "encoding/base64"
  "strings"
  // "log"
  "context"
  "github.com/gofiber/fiber/v2"
  "github.com/gofiber/fiber/v2/utils"
  "github.com/mrusme/journalist/journalistd"
  "github.com/mrusme/journalist/ent"
  "github.com/mrusme/journalist/ent/user"
  "github.com/mrusme/journalist/web/subscriptions"
)

func Register(config *journalistd.Config, fiberApp *fiber.App, entClient *ent.Client) () {
  web := fiberApp.Group("/web")
  // web.Use(authorizer(entClient))

  subscriptions.Register(&web, entClient)
}

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


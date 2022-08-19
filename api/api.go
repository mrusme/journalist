package api

import (
  "log"
  "context"
	"github.com/gofiber/fiber/v2"
  "github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/mrusme/journalist/ent"
	"github.com/mrusme/journalist/ent/user"
  "github.com/mrusme/journalist/api/v1"
)

func Register(fiberApp *fiber.App, entClient *ent.Client) () {
  bamw := basicauth.New(basicauth.Config{
      Realm: "Forbidden",
      Authorizer: func(username, password string) bool {
        _, err := entClient.User.
          Query().
          Where(user.Username(username)).
          Only(context.Background())
        if err == nil {
          return true
        }

        log.Printf("%v\n", err)
        return false
      },
      Unauthorized: func(c *fiber.Ctx) error {
          return c.SendStatus(fiber.StatusUnauthorized)
      },
      ContextUsername: "username",
      ContextPassword: "password",
  })
  api := fiberApp.Group("/api", bamw)

  v1.Register(&api, entClient)
}

package web

import (
	// "log"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/mrusme/journalist/ent"
	"github.com/mrusme/journalist/ent/token"
	"github.com/mrusme/journalist/ent/user"
	"github.com/mrusme/journalist/lib"
	"github.com/mrusme/journalist/web/actions"
	"github.com/mrusme/journalist/web/subscriptions"
)

func Register(
	jctx *lib.JournalistContext,
	fiberApp *fiber.App,
) {
	web := fiberApp.Group("/web")
	web.Use(authorizer(jctx.EntClient))

	actions.Register(
		jctx,
		&web,
	)

	subscriptions.Register(
		jctx,
		&web,
	)
}

// TODO: Move to `middlewares`
func authorizer(entClient *ent.Client) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		qat := ctx.Query("qat")
		if qat == "" {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}

		u, err := entClient.User.
			Query().
			WithTokens().
			Where(
				user.HasTokensWith(
					token.Token(qat),
				),
			).
			Only(context.Background())
		if err != nil {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}

		if u == nil {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}

		ctx.Locals("user_id", u.ID.String())
		ctx.Locals("username", u.Username)
		// ctx.Locals("password", u.Password)
		ctx.Locals("role", u.Role)
		return ctx.Next()
	}
}

package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrusme/journalist/api/v1/feeds"
	"github.com/mrusme/journalist/api/v1/tokens"
	"github.com/mrusme/journalist/api/v1/users"
	"github.com/mrusme/journalist/lib"
)

func Register(
  jctx *lib.JournalistContext,
  fiberRouter *fiber.Router,
) () {
  v1 := (*fiberRouter).Group("/v1")

  users.Register(
    jctx,
    &v1,
  )

  tokens.Register(
    jctx,
    &v1,
  )

  feeds.Register(
    jctx,
    &v1,
  )
}

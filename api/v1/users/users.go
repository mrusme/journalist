package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrusme/journalist/ent"
	"github.com/mrusme/journalist/lib"
	"go.uber.org/zap"
)

type handler struct {
  jctx      *lib.JournalistContext
  config    *lib.Config
  entClient *ent.Client
  logger    *zap.Logger
}

type UserShowModel struct {
  ID                string        `json:"id"`
  Username          string        `json:"username"`
  Role              string        `json:"role"`
}

type UserCreateModel struct {
  Username          string        `json:"username",validate:"required,alphanum,max=32"`
  Password          string        `json:"password",validate:"required"`
  Role              string        `json:"role",validate:"required"`
}

type UserUpdateModel struct {
  Password          string        `json:"password",validate:"min=5"`
  Role              string        `json:"role",validate:""`
}

func Register(
  jctx *lib.JournalistContext,
  fiberRouter *fiber.Router,
) () {
  endpoint := new(handler)
  endpoint.jctx = jctx
  endpoint.config = endpoint.jctx.Config
  endpoint.entClient = endpoint.jctx.EntClient
  endpoint.logger = endpoint.jctx.Logger

  usersRouter := (*fiberRouter).Group("/users")
  usersRouter.Get("/", endpoint.List)
  usersRouter.Get("/:id", endpoint.Show)
  usersRouter.Post("/", endpoint.Create)
  usersRouter.Put("/:id", endpoint.Update)
  // usersRouter.Delete("/:id", endpoint.Destroy)
}


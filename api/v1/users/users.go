package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrusme/journalist/ent"
	"github.com/mrusme/journalist/journalistd"
	"go.uber.org/zap"
)

type handler struct {
  config    *journalistd.Config
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
  config *journalistd.Config,
  fiberRouter *fiber.Router,
  entClient *ent.Client,
  logger *zap.Logger,
) () {
  endpoint := new(handler)
  endpoint.config = config
  endpoint.entClient = entClient
  endpoint.logger = logger

  usersRouter := (*fiberRouter).Group("/users")
  usersRouter.Get("/", endpoint.List)
  usersRouter.Get("/:id", endpoint.Show)
  usersRouter.Post("/", endpoint.Create)
  usersRouter.Put("/:id", endpoint.Update)
  // usersRouter.Delete("/:id", endpoint.Destroy)
}


package feeds

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrusme/journalist/ent"
	"github.com/mrusme/journalist/lib"
	"go.uber.org/zap"
)

type handler struct {
	jctx *lib.JournalistContext

	config    *lib.Config
	entClient *ent.Client
	logger    *zap.Logger
}

type FeedShowModel struct {
	ID    string `json:"id"`
	Name  string `json:"name,omitempty" validate:"omitempty,max=32"`
	URL   string `json:"url"`
	Group string `json:"group,omitempty" validate:"omitempty,max=32"`
}

type FeedCreateModel struct {
	Name     string `json:"name,omitempty" validate:"omitempty,max=32"`
	URL      string `json:"url" validate:"required,url"`
	Username string `json:"username,omitempty" validate:"omitempty,required_with=password"`
	Password string `json:"password,omitempty" validate:"omitempty,required_with=username"`
	Group    string `json:"group,omitempty" validate:"omitempty,max=32"`
}

/* type FeedUpdateModel struct {
  Password          string        `json:"password,omitempty" validate:"omitempty,min=5"`
} */

func Register(
	jctx *lib.JournalistContext,
	fiberRouter *fiber.Router,
) {
	endpoint := new(handler)
	endpoint.jctx = jctx
	endpoint.config = endpoint.jctx.Config
	endpoint.entClient = endpoint.jctx.EntClient
	endpoint.logger = endpoint.jctx.Logger

	feedsRouter := (*fiberRouter).Group("/feeds")
	feedsRouter.Get("/", endpoint.List)
	feedsRouter.Get("/:id", endpoint.Show)
	feedsRouter.Post("/", endpoint.Create)
	// feedsRouter.Put("/:id", endpoint.Update)
	// feedsRouter.Delete("/:id", endpoint.Destroy)
}

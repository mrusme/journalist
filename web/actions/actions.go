package actions

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

func Register(
	jctx *lib.JournalistContext,
	fiberRouter *fiber.Router,
) {
	endpoint := new(handler)
	endpoint.jctx = jctx
	endpoint.config = endpoint.jctx.Config
	endpoint.entClient = endpoint.jctx.EntClient
	endpoint.logger = endpoint.jctx.Logger

	actionsRouter := (*fiberRouter).Group("/actions")
	actionsRouter.Get("/read/:id", endpoint.Read)
	actionsRouter.Get("/read_older/:id", endpoint.ReadOlder)
	actionsRouter.Get("/read_newer/:id", endpoint.ReadNewer)
	// actionsRouter.Get("/read_all/:id", endpoint.ReadAll)
}

func (h *handler) resp(ctx *fiber.Ctx, content fiber.Map) error {
	err := ctx.Render("views/actions", content)
	ctx.Set("Content-type", "text/html; charset=utf-8")
	return err
}

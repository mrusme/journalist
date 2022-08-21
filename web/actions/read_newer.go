package actions

import (
	"github.com/google/uuid"
	"github.com/mrusme/journalist/ent/item"

	"github.com/gofiber/fiber/v2"
)

func (h *handler) ReadNewer(ctx *fiber.Ctx) error {
  id := ctx.Params("id")
  group := ctx.Query("group")

  sessionUserId := ctx.Locals("user_id").(string)
  myId, err := uuid.Parse(sessionUserId)
  if err != nil {
    ctx.SendStatus(fiber.StatusInternalServerError)
    return err
  }

  err = h.readWithItemCondition(
    myId,
    group,
    id,
    item.ItemPublishedGT,
  )
  if err != nil {
    ctx.SendStatus(fiber.StatusInternalServerError)
    return err
  }

  // err = ctx.Render("views/actions.read", fiber.Map{
  // })
  // ctx.Set("Content-type", "text/html; charset=utf-8")
  // return err
  return ctx.SendStatus(fiber.StatusNoContent)
}




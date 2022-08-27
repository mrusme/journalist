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
    h.resp(ctx, fiber.Map{
      "Success": false,
      "Title": "Error",
      "Message": err.Error(),
    })
    return err
  }

  err = h.readWithItemCondition(
    myId,
    group,
    id,
    item.ItemPublishedGT,
  )
  if err != nil {
    h.resp(ctx, fiber.Map{
      "Success": false,
      "Title": "Error",
      "Message": err.Error(),
    })
    return err
  }

  return h.resp(ctx, fiber.Map{
    "Success": true,
    "Title": "Marked as read",
    "Message": "Item was marked as read!",
  })
}




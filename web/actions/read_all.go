package actions

import (
	"github.com/google/uuid"
	"github.com/mrusme/journalist/ent"
	"github.com/mrusme/journalist/ent/item"
	"github.com/mrusme/journalist/ent/user"
	"github.com/mrusme/journalist/ent/subscription"

	"context"

	"github.com/gofiber/fiber/v2"
)

func (h *handler) ReadAll(ctx *fiber.Ctx) error {
  // id := ctx.Params("id")
  group := ctx.Query("group")

  sessionUserId := ctx.Locals("user_id").(string)
  myId, err := uuid.Parse(sessionUserId)
  if err != nil {
    ctx.SendStatus(fiber.StatusInternalServerError)
    return err
  }

  dbTmp := h.entClient.User.
    Query().
    Where(
      user.ID(myId),
    ).
    QuerySubscriptions()

  if group != "" {
    dbTmp = dbTmp.
      Where(subscription.Group(group))
  }

  dbItems, err := dbTmp.
    QueryFeed().
    QueryItems().
    Where(
      item.Not(
        item.HasReadByUsersWith(user.ID(myId)),
      ),
    ).
    All(context.Background())

  if err != nil {
    ctx.SendStatus(fiber.StatusInternalServerError)
    return err
  }

  bulkReads := make([]*ent.ReadCreate, len(dbItems))
  for i, item := range dbItems {
    bulkReads[i] = h.entClient.Read.
      Create().
      SetUserID(myId).
      SetItemID(item.ID)
  }
  err = h.entClient.Read.
    CreateBulk(bulkReads...).
    OnConflict().
    Ignore().
    Exec(context.Background())

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


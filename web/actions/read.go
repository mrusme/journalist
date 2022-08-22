package actions

import (
  "time"

	"github.com/google/uuid"
	"github.com/mrusme/journalist/ent"
	"github.com/mrusme/journalist/ent/predicate"
	"github.com/mrusme/journalist/ent/item"
	"github.com/mrusme/journalist/ent/user"
	"github.com/mrusme/journalist/ent/subscription"

	"context"

	"github.com/gofiber/fiber/v2"
)

func (h *handler) Read(ctx *fiber.Ctx) error {
  id := ctx.Params("id")
  // qat := ctx.Query("qat")
  // group := ctx.Query("group")

  // sessionUsername := ctx.Locals("username").(string)
  sessionUserId := ctx.Locals("user_id").(string)
  myId, err := uuid.Parse(sessionUserId)
  if err != nil {
    ctx.SendStatus(fiber.StatusInternalServerError)
    return err
  }

  dbUser, err := h.EntClient.User.
  Query().
  Where(
    user.ID(myId),
  ).
  Only(context.Background())

  dbItem, err := h.EntClient.Item.
  Query().
  Where(
    item.ItemGUID(id),
  ).
  Only(context.Background())
  if err != nil {
    ctx.SendStatus(fiber.StatusInternalServerError)
    return err
  }

  err = dbUser.
    Update().
    AddReadItemIDs(dbItem.ID).
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

func (h *handler) readWithItemCondition(
  userId uuid.UUID,
  group string,
  itemGUID string,
  itemCondition func(v time.Time) predicate.Item,
) error {
  dbItem, err := h.EntClient.Item.
    Query().
    Where(
      item.ItemGUID(itemGUID),
    ).
    Only(context.Background())
  if err != nil {
    return err
  }

  dbTmp := h.EntClient.User.
    Query().
    Where(
      user.ID(userId),
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
      item.Or(
        item.ID(dbItem.ID),
        itemCondition(dbItem.ItemPublished),
      ),
    ).
    All(context.Background())

  if err != nil {
    return err
  }

  bulkReads := make([]*ent.ReadCreate, len(dbItems))
  for i, item := range dbItems {
    bulkReads[i] = h.EntClient.Read.
      Create().
      SetUserID(userId).
      SetItemID(item.ID)
  }
  err = h.EntClient.Read.
    CreateBulk(bulkReads...).
    OnConflict().
    Ignore().
    Exec(context.Background())

  if err != nil {
    return err
  }

  return nil
}

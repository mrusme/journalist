package actions

import (
	"time"

	"github.com/google/uuid"
	"github.com/mrusme/journalist/ent"
	"github.com/mrusme/journalist/ent/item"
	"github.com/mrusme/journalist/ent/predicate"
	"github.com/mrusme/journalist/ent/subscription"
	"github.com/mrusme/journalist/ent/user"

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
		h.resp(ctx, fiber.Map{
			"Success": false,
			"Title":   "Error",
			"Message": err.Error(),
		})
		return err
	}

	dbUser, err := h.entClient.User.
		Query().
		Where(
			user.ID(myId),
		).
		Only(context.Background())

	dbItem, err := h.entClient.Item.
		Query().
		Where(
			item.ItemGUID(id),
		).
		Only(context.Background())
	if err != nil {
		h.resp(ctx, fiber.Map{
			"Success": false,
			"Title":   "Error",
			"Message": err.Error(),
		})
		return err
	}

	err = dbUser.
		Update().
		AddReadItemIDs(dbItem.ID).
		Exec(context.Background())
	if err != nil {
		h.resp(ctx, fiber.Map{
			"Success": false,
			"Title":   "Error",
			"Message": err.Error(),
		})
		return err
	}

	return h.resp(ctx, fiber.Map{
		"Success": true,
		"Title":   "Marked as read",
		"Message": "Item was marked as read!",
	})
}

func (h *handler) readWithItemCondition(
	userId uuid.UUID,
	group string,
	itemGUID string,
	itemCondition func(v time.Time) predicate.Item,
) error {
	dbItem, err := h.entClient.Item.
		Query().
		Where(
			item.ItemGUID(itemGUID),
		).
		Only(context.Background())
	if err != nil {
		return err
	}

	dbTmp := h.entClient.User.
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
		bulkReads[i] = h.entClient.Read.
			Create().
			SetUserID(userId).
			SetItemID(item.ID)
	}
	err = h.entClient.Read.
		CreateBulk(bulkReads...).
		OnConflict().
		Ignore().
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

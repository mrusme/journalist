package actions

import (
	"github.com/google/uuid"
	"github.com/mrusme/journalist/ent"
	"github.com/mrusme/journalist/ent/item"
	"github.com/mrusme/journalist/ent/subscription"
	"github.com/mrusme/journalist/ent/user"

	"context"

	"github.com/gofiber/fiber/v2"
)

func (h *handler) ReadAll(ctx *fiber.Ctx) error {
	// id := ctx.Params("id")
	group := ctx.Query("group")

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
		h.resp(ctx, fiber.Map{
			"Success": false,
			"Title":   "Error",
			"Message": err.Error(),
		})
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

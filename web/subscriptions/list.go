package subscriptions

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mrusme/journalist/ent"
	"github.com/mrusme/journalist/ent/item"
	"github.com/mrusme/journalist/ent/subscription"
	"github.com/mrusme/journalist/ent/user"
	"github.com/mrusme/journalist/journalistd"

	"context"

	"github.com/gofiber/fiber/v2"
)

func (h *handler) List(ctx *fiber.Ctx) error {
	qat := ctx.Query("qat")
	group := ctx.Query("group")
	sessionUsername := ctx.Locals("username").(string)
	sessionUserId := ctx.Locals("user_id").(string)
	myId, err := uuid.Parse(sessionUserId)
	if err != nil {
		ctx.SendStatus(fiber.StatusInternalServerError)
		return err
	}

	dbItemsTmp := h.entClient.Subscription.
		Query()

	if group == "" {
		dbItemsTmp = dbItemsTmp.Where(
			subscription.UserID(myId),
		)
	} else {
		dbItemsTmp = dbItemsTmp.Where(
			subscription.UserID(myId),
			subscription.Group(group),
		)
	}

	dbItems, err := dbItemsTmp.
		QueryFeed().
		QueryItems().
		Where(
			item.Not(
				item.HasReadByUsersWith(
					user.ID(myId),
				),
			),
		).
		Order(
			ent.Desc(item.FieldItemPublished),
		).
		All(context.Background())
	if err != nil {
		ctx.SendStatus(fiber.StatusInternalServerError)
		return err
	}

	err = ctx.Render("views/subscriptions.list", fiber.Map{
		"Config": h.config,
		"Token": fiber.Map{
			"Type":  "qat",
			"Token": qat,
		},
		"Group": group,

		"Title":         h.tmplTitle(group),
		"Link":          h.tmplLink("qat", qat, group),
		"Description":   h.tmplDescription(sessionUsername, group),
		"Generator":     h.tmplGenerator(),
		"Language":      "en-us",
		"LastBuildDate": time.Now(),

		"Items": dbItems,
	})
	ctx.Set("Content-type", "text/xml; charset=utf-8")
	return err
}

func (h *handler) tmplTitle(group string) string {
	var title string = "Subscriptions"
	if group != "" {
		title = group
	}

	return title
}

func (h *handler) tmplDescription(
	username string,
	group string,
) string {
	var description string = ""

	if username[len(username)-1] == 's' {
		description = fmt.Sprintf(
			"%s' subscriptions",
			username,
		)
	} else {
		description = fmt.Sprintf(
			"%s's subscriptions",
			username,
		)
	}

	if group != "" {
		description = fmt.Sprintf(
			"%s in %s",
			description,
			group,
		)
	}

	return description
}

func (h *handler) tmplLink(
	tokenType string,
	token string,
	group string,
) string {
	return fmt.Sprintf(
		"%s/subscriptions?group=%s&%s=%s",
		h.config.Server.Endpoint.Web,
		group,
		tokenType,
		token,
	)
}

func (h *handler) tmplGenerator() string {
	return fmt.Sprintf(
		"Journalist %s",
		journalistd.Version(),
	)
}

package feeds

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrusme/journalist/ent"
)

type handler struct {
  EntClient *ent.Client
}

type FeedShowModel struct {
  ID                string        `json:"id"`
  Name              string        `json:"name",validate:"optional,alphanum,max=32"`
  URL               string        `json:"url"`
  Group             string        `json:"group",validate:"optional,alphanum,max=32"`
}

type FeedCreateModel struct {
  Name              string        `json:"name",validate:"optional,alphanum,max=32"`
  URL               string        `json:"url",validate:"required,url"`
  Username          string        `json:"username",validate:"optional,required_with=password"`
  Password          string        `json:"password",validate:"optional,required_with=username"`
  Group             string        `json:"group",validate:"optional,alphanum,max=32"`
}

type FeedUpdateModel struct {
  Password          string        `json:"password",validate:"min=5"`
  Role              string        `json:"role",validate:""`
}

func Register(fiberRouter *fiber.Router, entClient *ent.Client) () {
  endpoint := new(handler)
  endpoint.EntClient = entClient

  feedsRouter := (*fiberRouter).Group("/feeds")
  // feedsRouter.Get("/", endpoint.List)
  // feedsRouter.Get("/:id", endpoint.Show)
  feedsRouter.Post("/", endpoint.Create)
  // feedsRouter.Put("/:id", endpoint.Update)
  // feedsRouter.Delete("/:id", endpoint.Destroy)
}



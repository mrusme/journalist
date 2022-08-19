package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrusme/journalist/ent"
  "github.com/mrusme/journalist/api/v1/users"
)
func Register(fiberRouter *fiber.Router, entClient *ent.Client) () {
  v1 := (*fiberRouter).Group("/v1")
  users.Register(&v1, entClient)
}

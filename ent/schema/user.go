package schema

import (
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	validate := validator.New()

	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		// StorageKey("oid"),
		field.String("username").
			Validate(func(s string) error {
				return validate.Var(s, "required,alphanum,max=32")
			}).
			Unique(),
		field.String("password").
			Validate(func(s string) error {
				return validate.Var(s, "required,min=5")
			}).
			Sensitive(),
		field.String("role").
			Default("user").
			Match(regexp.MustCompile("^(admin|user)$")),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.Time("deleted_at").
			Default(nil).
			Optional().
			Nillable(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tokens", Token.Type),
		edge.To("subscribed_feeds", Feed.Type).
			Through("subscriptions", Subscription.Type),
		edge.To("read_items", Item.Type).
			Through("reads", Read.Type),
	}
}

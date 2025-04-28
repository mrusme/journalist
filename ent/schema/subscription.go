package schema

import (
	"time"

	"github.com/go-playground/validator/v10"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Subscription holds the schema definition for the Subscription entity.
type Subscription struct {
	ent.Schema
}

// Fields of the Subscription.
func (Subscription) Fields() []ent.Field {
	validate := validator.New()

	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		// StorageKey("oid"),

		field.UUID("user_id", uuid.UUID{}),
		field.UUID("feed_id", uuid.UUID{}),

		field.String("name").
			Validate(func(s string) error {
				return validate.Var(s, "required")
			}),
		field.String("group").
			Validate(func(s string) error {
				return validate.Var(s, "required,alphanum,max=32")
			}),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the Subscription.
func (Subscription) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).
			Unique().
			Required().
			Field("user_id"),
		edge.To("feed", Feed.Type).
			Unique().
			Required().
			Field("feed_id"),
	}
}

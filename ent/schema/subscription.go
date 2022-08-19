package schema

import (
	// "regexp"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/edge"
	"github.com/google/uuid"
)

// Subscription holds the schema definition for the Subscription entity.
type Subscription struct {
	ent.Schema
}

// Fields of the Subscription.
func (Subscription) Fields() []ent.Field {
  return []ent.Field{
    field.UUID("id", uuid.UUID{}).
      Default(uuid.New).
      StorageKey("oid"),
    field.UUID("user_id", uuid.UUID{}),
    field.UUID("feed_id", uuid.UUID{}),
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

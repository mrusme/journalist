package schema

import (
	// "regexp"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/edge"
	"github.com/google/uuid"
)

// Feed holds the schema definition for the Feed entity.
type Feed struct {
	ent.Schema
}

// Fields of the Feed.
func (Feed) Fields() []ent.Field {
  return []ent.Field{
    field.UUID("id", uuid.UUID{}).
      Default(uuid.New).
      StorageKey("oid"),
  }
}

// Edges of the Feed.
func (Feed) Edges() []ent.Edge {
  return []ent.Edge{
    edge.To("items", Item.Type),
    edge.From("subscribed_users", User.Type).
      Ref("subscribed_feeds").
      Through("subscriptions", Subscription.Type),
  }
}

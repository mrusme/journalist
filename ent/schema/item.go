package schema

import (
	// "regexp"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/edge"
	"github.com/google/uuid"
)

// Item holds the schema definition for the Item entity.
type Item struct {
	ent.Schema
}

// Fields of the Item.
func (Item) Fields() []ent.Field {
  return []ent.Field{
    field.UUID("id", uuid.UUID{}).
      Default(uuid.New).
      StorageKey("oid"),
  }
}

// Edges of the Item.
func (Item) Edges() []ent.Edge {
  return []ent.Edge{
    edge.From("feed", Feed.Type).
      Ref("items").
      Unique(),
    edge.From("read_by_users", User.Type).
      Ref("read_items").
      Through("reads", Read.Type),
  }
}

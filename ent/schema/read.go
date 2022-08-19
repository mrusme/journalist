package schema

import (
	// "regexp"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/edge"
	"github.com/google/uuid"
)

// Read holds the schema definition for the Read entity.
type Read struct {
	ent.Schema
}

// Fields of the Read.
func (Read) Fields() []ent.Field {
  return []ent.Field{
    field.UUID("id", uuid.UUID{}).
      Default(uuid.New).
      StorageKey("oid"),
    field.UUID("user_id", uuid.UUID{}),
    field.UUID("item_id", uuid.UUID{}),
  }
}

// Edges of the Read.
func (Read) Edges() []ent.Edge {
  return []ent.Edge{
    edge.To("user", User.Type).
      Unique().
      Required().
      Field("user_id"),
    edge.To("item", Item.Type).
      Unique().
      Required().
      Field("item_id"),
  }
}

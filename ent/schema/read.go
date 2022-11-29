package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
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
			Default(uuid.New),
		// StorageKey("oid"),

		field.UUID("user_id", uuid.UUID{}),
		field.UUID("item_id", uuid.UUID{}),

		field.Time("created_at").
			Default(time.Now),
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

package schema

import (
	"regexp"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
  return []ent.Field{
    field.UUID("id", uuid.UUID{}).
      Default(uuid.New).
      StorageKey("oid"),
    field.String("username").
      Match(regexp.MustCompile("^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$")).
      Unique(),
    field.String("password").
      Sensitive(),
  }
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}

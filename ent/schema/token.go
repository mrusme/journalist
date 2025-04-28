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

// Token holds the schema definition for the Token entity.
type Token struct {
	ent.Schema
}

// Fields of the Token.
func (Token) Fields() []ent.Field {
	validate := validator.New()

	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		// StorageKey("oid"),
		field.String("type").
			Default("qat").
			Match(regexp.MustCompile("^(qat|jwt)$")),
		field.String("name").
			Validate(func(s string) error {
				return validate.Var(s, "required,alphanum,max=32")
			}),
		field.String("token").
			Unique().
			Sensitive(),
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

// Edges of the Token.
func (Token) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("tokens").
			Unique(),
	}
}
